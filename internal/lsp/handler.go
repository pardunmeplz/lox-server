package lsp

import (
	"encoding/json"
	"fmt"
	"syscall"
)

func handleRequest(msg string) ([]byte, error) {
	var requestObj map[string]any

	if err := json.Unmarshal([]byte(msg), &requestObj); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	responseObj, err := processRequest(requestObj)
	if err != nil {
		return nil, fmt.Errorf("invalid Request: %v", err)
	}
	if responseObj == nil {
		return nil, nil
	}

	response, err := json.Marshal(responseObj)
	if err != nil {
		return nil, fmt.Errorf("invalid Request: %v", err)
	}

	header := []byte(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(response)))

	return append(header, response...), nil
}

func processRequest(request map[string]any) (map[string]any, error) {

	switch request["method"] {
	case "initialize":
		serverState.initialized = true
		return protocolInitialize(request)
	case "shutdown":
		serverState.shutdown = true
		return protocolShutdown(request), nil
	case "exit":
		if serverState.shutdown {
			syscall.Exit(0)
		} else {
			syscall.Exit(1)
		}
	case "initialized":
		return nil, nil
	case "textDocument/didOpen":
		return nil, nil
	case "textDocument/didClose":
		return nil, nil
	case "textDocument/didChange":
		return nil, nil
	}

	return nil, fmt.Errorf("Invalid Method: %v", request["method"])

}
