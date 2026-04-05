package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	lsp "lox-server/internal/lsp/types"
	"os"
	"sync"
)

/*
Router: just a map of string to function that will run the function based on exact string match of incoming request's
*/
type Server struct {
	Router   map[string]func(*lsp.JsonRpcRequest) *lsp.JsonRpcResponse
	Reader   bufio.Reader
	buffer   []byte
	Writer   *os.File
	WriterMu sync.Mutex
}

/*
Starts an application loop that will keep running listening for requests
and responding as per the server router
*/
func StartServer(server *Server) {
	for {
		line, err := server.Reader.ReadBytes()
	}
}

func routeRequest(server *Server, msg []byte) ([]byte, error) {
	var requestObj lsp.JsonRpcRequest

	if err := json.Unmarshal([]byte(msg), &requestObj); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	if requestObj.Method == "" {
		return nil, nil
	}

	var route = server.Router[requestObj.Method]
	if route == nil {
		return nil, nil
	}
	var resp = route(&requestObj)

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
	return err
}
