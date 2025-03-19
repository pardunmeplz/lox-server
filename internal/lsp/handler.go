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
	var requestObj lsp.JsonRpcRequest

	if err := json.Unmarshal([]byte(msg), &requestObj); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	if requestObj.Method == "" {
		return nil, nil
	}

	responseObj, err := processRequest(requestObj)
	if err != nil {
		return nil, fmt.Errorf("invalid Request: %v", err)
	}
	if responseObj == nil {
		return nil, nil
	}

	response, err := json.Marshal(*responseObj)
	if err != nil {
		return nil, fmt.Errorf("invalid Response: %v", err)
	}

	return response, nil
}

func processRequest(request lsp.JsonRpcRequest) (*lsp.JsonRpcResponse, error) {
	switch request.Id.(type) {
	case string:
		id, err := strconv.Atoi(request.Id.(string))
		if err == nil {
			serverState.idCount = id
		}
	case int:
		serverState.idCount = request.Id.(int)
	}

	switch request.Method {
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
		var params lsp.DidOpenTextDocumentParams
		err := getRequestValues(&params, request)
		if err != nil {
			return nil, nil
		}

		serverState.documents[params.TextDocument.Uri] = &DocumentService{Uri: params.TextDocument.Uri}
		serverState.documents[params.TextDocument.Uri].Initialize()
		go (func() {
			serverState.documents[params.TextDocument.Uri].ParseCode(params.TextDocument.Text, params.TextDocument.Version)
		})()
		return nil, nil
	case "textDocument/didClose":
		var params lsp.DidCloseTextDocumentParams
		err := getRequestValues(&params, request)
		if err != nil {
			return nil, nil
		}

		delete(serverState.documents, params.TextDocument.Uri)
		return nil, nil
	case "textDocument/didChange":
		var params lsp.DidChangeTextDocumentParams
		err := getRequestValues(&params, request)
		if err != nil {
			return nil, nil
		}

		go (func() {
			serverState.documents[params.TextDocument.Uri].ParseCode(params.ContentChanges[0].Text, params.TextDocument.Version)
		})()
		return nil, nil
	case "textDocument/definition":
		return protocolDefinition(request), nil
	case "textDocument/references":
		return protocolReferences(request), nil
	case "textDocument/formatting":
		return protocolFormatting(request), nil
	case "textDocument/completion":
		return protocolCompletion(request), nil
	case "textDocument/semanticTokens/full":
		return protocolSemanticTokens(request), nil
	case "/cancelRequest":
		return nil, nil

	}

	return nil, fmt.Errorf("Invalid Method: %v", request.Method)
}

func getRequestValues[T any](document *T, request lsp.JsonRpcRequest) error {
	params, err := json.Marshal(request.Params)
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
