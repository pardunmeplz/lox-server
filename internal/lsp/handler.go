package lsp

import (
	"encoding/json"
	"fmt"
	"lox-server/internal/lox"
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
		var document lsp.DidOpenTextDocumentParams
		params, err := json.Marshal(request["params"])
		if err != nil {
			return nil, fmt.Errorf("Marshal failed : %v", err)
		}

		json.Unmarshal(params, &document)
		go checkForErrors(&document.TextDocument.Text, document.TextDocument.Uri, document.TextDocument.Version)
		return nil, nil
	case "textDocument/didClose":
		return nil, nil
	case "textDocument/didChange":
		return nil, nil
	}

	return nil, fmt.Errorf("Invalid Method: %v", request["method"])
}

func checkForErrors(code *string, uri string, version int) {
	errors := lox.FindErrors(*code)
	if errors == nil {
		return
	}
	responseObj := lsp.PublishDiagnosticParams{Uri: uri, Version: version, Diagnostics: lsp.Diagnostic{
		Severity: 1,
		ErrRange: lsp.Range{
			Start: lsp.Position{Line: 0, Character: 0},
			End:   lsp.Position{Line: 1, Character: 0},
		},
		Message: "message here",
	}}
	response, err := json.Marshal(responseObj)
	if err != nil {
		serverState.logger.Print(fmt.Sprintf("invalid Response: %v\n", err))
		return
	}
	if err := writeMessage(serverState.writer, response); err != nil {
		serverState.logger.Print(fmt.Sprintf("Error writing response: %v\n", err))
	}
}
