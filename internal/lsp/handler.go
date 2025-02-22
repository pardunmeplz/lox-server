package lsp

import (
	"bytes"
	"encoding/json"
	"fmt"
	lsp "lox-server/internal/lsp/types"
	"strconv"
	"syscall"
)

func split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})

	if !found {
		return 0, nil, nil
	}
	contentLength, err := strconv.Atoi(string(header[len("Content-Length: "):]))
	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	bodyStart := len(header) + 4
	totalLength := len(header) + 4 + contentLength

	return totalLength, data[bodyStart:totalLength], nil
}

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
		serverState.notificationChannel <- request
		return nil, nil
	case "textDocument/didClose":
		return nil, nil
	case "textDocument/didChange":
		serverState.notificationChannel <- request
		return nil, nil
	}

	return nil, fmt.Errorf("Invalid Method: %v", request["method"])
}

func processNotification(request map[string]any) []byte {
	switch request["method"] {
	case "textDocument/didOpen":
		var document lsp.DidOpenTextDocumentParams
		err := getRequestValues(&document, request)
		if err != nil {
			return nil
		}
		responseObj, err := diagnosticNotification(document.TextDocument.Text, document.TextDocument.Uri, document.TextDocument.Version)
		if err != nil {
			serverState.logger.Print(fmt.Sprintf("Parse Error: %v\n", err))
			return nil
		}
		response, err := json.Marshal(responseObj)
		if err != nil {
			serverState.logger.Print(fmt.Sprintf("invalid Response: %v\n", err))
			return nil
		}
		return response

	case "textDocument/didChange":
		var document lsp.DidChangeTextDocumentParams
		err := getRequestValues(&document, request)
		if err != nil {
			return nil
		}
		responseObj, err := diagnosticNotification(document.ContentChanges[0].Text, document.TextDocument.Uri, document.TextDocument.Version)
		if err != nil {
			serverState.logger.Print(fmt.Sprintf("Parse Error: %v\n", err))
			return nil
		}
		response, err := json.Marshal(responseObj)
		if err != nil {
			serverState.logger.Print(fmt.Sprintf("invalid Response: %v\n", err))
			return nil
		}
		return response

	}
	return nil

}

func getRequestValues[T any](document *T, request map[string]any) error {
	params, err := json.Marshal(request["params"])
	if err != nil {
		serverState.logger.Print(fmt.Sprintf("Marshal failed : %v", err))
		return err
	}
	err = json.Unmarshal(params, &document)
	if err != nil {
		serverState.logger.Print(fmt.Sprintf("Params Unmarshal failed : %v", err))
		return err
	}
	return nil
}
