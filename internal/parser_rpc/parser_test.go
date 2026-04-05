package parser_rpc

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func initParser() ParserState {
	st := ParserState{PARSER_STATUS_HEADER, JsonRpcRequest{}}
	return st
}

func TestPositiveFlow(t *testing.T) {
	parser := initParser()
	req := []byte("Content-Length: 128\r\n\r\n{\"jsonrpc\": \"2.0\",\"id\": 0,\"result\": {\"capabilities\": {\"textDocumentSync\": 1,\"completionProvider\": { \"resolveProvider\": true }}}}")

	consumed, err := ParseJsonRpcRequest(req, &parser)
	if err != nil {
		t.Fatal(err)
	}

	if consumed != len(req) {
		t.Fatal(fmt.Sprintf("invalid consumed value %d instead of %d", consumed, len(req)))
	}

	if parser.Request.Headers["Content-Length"] != "128" {
		t.Fatal(fmt.Sprintf("invalid content length value %s", parser.Request.Headers["Content-Length"]))
	}

	expectedContent := []byte("{\"jsonrpc\": \"2.0\",\"id\": 0,\"result\": {\"capabilities\": {\"textDocumentSync\": 1,\"completionProvider\": { \"resolveProvider\": true }}}}")
	if !bytes.Equal(parser.Request.ContentBytes, expectedContent) {
		t.Fatal(fmt.Sprintf("invalid content %s", string(parser.Request.ContentBytes)))
	}
}

func TestContentLength(t *testing.T) {
	// missin content length check
	parser := initParser()
	req := []byte("\r\n\r\n{\"jsonrpc\": \"2.0\",\"id\": 0,\"result\": {\"capabilities\": {\"textDocumentSync\": 1,\"completionProvider\": { \"resolveProvider\": true }}}}")

	_, err := ParseJsonRpcRequest(req, &parser)
	if !strings.Contains(err.Error(), MISSING_CONTENT_LENGTH) {
		t.Fatal("Error expected, got %w", err)
	}

	// invalid content length check
	parser = initParser()
	req = []byte("Content-Length: 12A\r\n\r\n{\"jsonrpc\": \"2.0\",\"id\": 0,\"result\": {\"capabilities\": {\"textDocumentSync\": 1,\"completionProvider\": { \"resolveProvider\": true }}}}")

	_, err = ParseJsonRpcRequest(req, &parser)
	if !strings.Contains(err.Error(), "Invalid Content-Length") {
		t.Fatal("Error expected, got %w", err)
	}

	// content length smaller than incoming request bytes, check if content length is honored correctly
	parser = initParser()
	req = []byte("Content-Length: 128\r\n\r\n{\"jsonrpc\": \"2.0\",\"id\": 0,\"result\": {\"capabilities\": {\"textDocumentSync\": 1,\"completionProvider\": { \"resolveProvider\": true }}}}bogus value should be ignored")

	_, err = ParseJsonRpcRequest(req, &parser)
	if err != nil {
		t.Fatal(err)
	}

	expectedContent := []byte("{\"jsonrpc\": \"2.0\",\"id\": 0,\"result\": {\"capabilities\": {\"textDocumentSync\": 1,\"completionProvider\": { \"resolveProvider\": true }}}}")
	if !bytes.Equal(parser.Request.ContentBytes, expectedContent) {
		t.Fatal(fmt.Sprintf("invalid content %s", string(parser.Request.ContentBytes)))
	}
}
