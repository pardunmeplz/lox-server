package parser

import (
	"fmt"
	"testing"
)

func initParser() ParserState {
	st := ParserState{PARSER_STATUS_HEADER, JsonRpcRequest{make(map[string]string), nil}}
	return st
}

func TestPositiveFlow(t *testing.T) {
	parser := initParser()
	req := []byte("Content-Length: 129\r\n\r\n{\"jsonrpc\": \"2.0\",\"id\": 0,\"result\": {\"capabilities\": {\"textDocumentSync\": 1,\"completionProvider\": { \"resolveProvider\": true }}}}")

	consumed, err := parse(req, &parser)
	if err != nil {
		t.Fatal(err)
	}

	if consumed != len(req) {
		t.Log(parser.Request)
		t.Fatal(fmt.Sprintf("invalid consumed value %d instead of %d", consumed, len(req)))
	}

}
