package lsp

import (
	"encoding/json"
	"fmt"
	lsp "lox-server/internal/lsp/types"
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
		return nil, fmt.Errorf("invalid Response: %v", err)
	}

	return response, nil
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
		go checkForErrors(request)
		return nil, nil
	case "textDocument/didClose":
		return nil, nil
	case "textDocument/didChange":
		return nil, nil
	}

	return nil, fmt.Errorf("Invalid Method: %v", request["method"])
}

func checkForErrors(request map[string]any) {
	// get values
	var document lsp.DidOpenTextDocumentParams
	params, err := json.Marshal(request["params"])
	if err != nil {
		serverState.logger.Print(fmt.Sprintf("Marshal failed : %v", err))
		return
	}
	json.Unmarshal(params, &document)

	// gen notificaion
	responseObj, isError := diagnosticNotification(document.TextDocument.Text, document.TextDocument.Uri, document.TextDocument.Version)
	if !isError {
		return
	}

	// send notification
	response, err := json.Marshal(responseObj)
	if err != nil {
		serverState.logger.Print(fmt.Sprintf("invalid Response: %v\n", err))
		return
	}
	if err := writeMessage(response); err != nil {
		serverState.logger.Print(fmt.Sprintf("Error writing response: %v\n", err))
	}
	serverState.logger.Print(string(response))
}
