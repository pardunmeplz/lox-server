package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type ParserState struct {
	ParserStatus ParserStatus
	Request      JsonRpcRequest
}

type ParserStatus int

const (
	HEADER_CONTENT_LENGTH = "Content-Length"
)

const (
	PARSER_STATUS_HEADER ParserStatus = iota
	PARSER_STATUS_CONTENT
	PARSER_STATUS_DONE
)

const (
	MALFORMED_REQUEST      = "Malformed Request Line"
	MISSING_CONTENT_LENGTH = "Missing Content-Length header"
	INVALID_CONTENT_LENGTH = "Invalid Content-Length value %s"
)

func parseJsonRpcRequest(data []byte, parserState *ParserState) (int, error) {
	totalConsumed := 0
	for parserState.ParserStatus == PARSER_STATUS_HEADER {
		consumed, err := parseHeader(data, parserState)
		if err != nil {
			return 0, err
		}
		if consumed == 0 {
			break
		}
		totalConsumed += consumed
		data = data[consumed:]
	}
	if parserState.ParserStatus == PARSER_STATUS_CONTENT {
		consumed, err := parseContent(data, parserState)
		if err != nil {
			return totalConsumed, err
		}
		totalConsumed += consumed
	}

	return totalConsumed, nil

}

func parseHeader(data []byte, parserState *ParserState) (int, error) {
	// get bytes upto the terminator of the first header
	index := bytes.Index(data, []byte{'\r', '\n'})

	// no header found yet
	if index == -1 {
		return 0, nil
	}
	// if terminator is at start of line, the message is now moving on to the content part
	if index == 0 {
		parserState.ParserStatus = PARSER_STATUS_CONTENT
		return 2, nil
	}
	// get key and value of header using ": " separator
	key, value, found := bytes.Cut(data[:index], []byte{':', ' '})
	if !found {
		return 0, fmt.Errorf(MALFORMED_REQUEST)
	}
	if parserState.Request.Headers == nil {
		parserState.Request.Headers = make(map[string]string)
	}
	parserState.Request.Headers[string(key)] = string(value)

	// +2 because we will consider the /r/n as consumed as well
	return index + 2, nil
}

func parseContent(data []byte, parserState *ParserState) (int, error) {
	lengthStr, found := parserState.Request.Headers[HEADER_CONTENT_LENGTH]

	// content length is mandatory as per spec
	if !found {
		return 0, fmt.Errorf(MISSING_CONTENT_LENGTH)
	}

	// if its not an integer its invalid
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return 0, fmt.Errorf(INVALID_CONTENT_LENGTH, lengthStr)
	}

	// return if entire content is not yet streamed
	if length > len(data) {
		return 0, nil
	}

	// parse the base json rpc message that is expected to always be the same for all requests
	parserState.Request.ContentBytes = data[:length]
	parserState.Request.Content = &JsonRpcRequestContent{}

	err = json.Unmarshal(parserState.Request.ContentBytes, parserState.Request.Content)
	if err != nil {
		return 0, err
	}

	parserState.ParserStatus = PARSER_STATUS_DONE
	return length, nil
}
