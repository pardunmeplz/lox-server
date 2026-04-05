package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	lspt "lox-server/internal/lsp/types"
	rpc "lox-server/internal/parser_rpc"
	"sync"
)

/*
Router: just a map of string to function that will run the function based on exact string match of incoming request's
*/
type Server struct {
	Router map[string]func(*rpc.JsonRpcRequest) *lspt.JsonRpcResponse
	// you need to set up an bufio.NewReader(os.Stdin) in main
	Reader *bufio.Reader
	buffer []byte
	// you need to set up an bufio.NewWriter(os.Stdout) in main
	Writer   *bufio.Writer
	WriterMu sync.Mutex
}

const (
	BUFFER_SIZE      = 1024
	BUFFER_THRESHOLD = 128
	BUFFER_INCREMENT = 256
)

/*
Starts an application loop that will keep running listening for requests
and responding as per the server router

todo: Bubble all errors up to the server loop
set up a log to put out the error in an error log and fast fail, exit server on error
*/
func StartServer(server *Server) {
	server.buffer = make([]byte, BUFFER_SIZE)
	bufferIndex := 0
	rpcParser := rpc.ParserState{ParserStatus: rpc.PARSER_STATUS_HEADER, Request: rpc.JsonRpcRequest{}}

	for {

		// grow buffer when running out of space
		if len(server.buffer)-bufferIndex < BUFFER_THRESHOLD {
			server.buffer = append(server.buffer, make([]byte, BUFFER_INCREMENT)...)
		}

		// read from io
		bytes, err := server.Reader.Read(server.buffer[bufferIndex:])
		if err != nil {
			if err == io.EOF {
				// suggests clent shut down
				return
			}
			// log error
			return
		}
		bufferIndex += bytes

		// parse current buffer
		consumed, err := rpc.ParseJsonRpcRequest(server.buffer[:bufferIndex], &rpcParser)
		if err != nil {
			// log error
			return
		}

		// discard consumed bytes
		bufferIndex -= consumed
		server.buffer = server.buffer[consumed:]

		// reset parser and process request if parsing done
		if rpcParser.ParserStatus == rpc.PARSER_STATUS_DONE {
			// todo: this could cause responses to be sent out of order
			// the spec does have some leeway for response order for parallel execution
			// but it is still expected to for the most part to return responses in order
			go func(request rpc.JsonRpcRequest) {
				respBytes, err := routeRequest(server, &request)
				if err != nil {
					// handle server failure ( log and respond with appropriate error response or fast fail )
					return
				}
				if respBytes == nil {
					// nil response means router doesn't want to send a response via the loop
					return
				}
				WriteMessage(server, respBytes)
			}(*&rpcParser.Request)
			rpcParser.ParserStatus = rpc.PARSER_STATUS_HEADER
			rpcParser.Request = rpc.JsonRpcRequest{}
		}

	}
}

func routeRequest(server *Server, req *rpc.JsonRpcRequest) ([]byte, error) {

	var route, found = server.Router[req.Content.Method]
	if !found {
		// todo: log missing method ( consider if error should be logged and ignored or fast failed or check if LSP spec has a "not supported" response appropriate for this )
		return nil, nil
	}
	var resp = route(req)
	if resp == nil {
		return nil, nil
	}

	responseBytes, err := json.Marshal(&resp)
	if err != nil {
		return nil, fmt.Errorf("invalid Response: %v", err)
	}
	return responseBytes, nil
}

func WriteMessage(server *Server, response []byte) error {
	header := []byte(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(response)))
	server.WriterMu.Lock()
	defer server.WriterMu.Unlock()

	_, err := server.Writer.Write(append(header, response...))
	if err != nil {
		return err
	}
	server.Writer.Flush()
	return err
}
