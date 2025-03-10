package lsp

import (
	"encoding/json"
	"lox-server/internal/lox"
	"sync"
)

/* document level logic like language features and state are handled here*/

type DocumentService struct {
	Definables []lox.Node
	Errors     []lox.CompileError
	Uri        string
	Mutex      sync.Mutex
}

func (loxService *DocumentService) ParseCode(code string, version int) {
	compileErrors, definables, err := lox.ParseCode(code)
	if err != nil {
		return
	}
	loxService.Definables = definables
	loxService.Errors = compileErrors

	responseObj := diagnosticNotification(compileErrors, loxService.Uri, version)
	response, err := json.Marshal(responseObj)
	sendNotification(response)
}

// func processNotification(request) []byte {
// 	switch request.Method {
// 	case "textDocument/didOpen":
// 		var document lsp.DidOpenTextDocumentParams
// 		err := getRequestValues(&document, request)
// 		if err != nil {
// 			return nil
// 		}

// 	case "textDocument/didChange":
// 		var document lsp.DidChangeTextDocumentParams
// 		err := getRequestValues(&document, request)
// 		if err != nil {
// 			return nil
// 		}
// 		responseObj, err := diagnosticNotification(document.ContentChanges[0].Text, document.TextDocument.Uri, document.TextDocument.Version)
// 		if err != nil {
// 			serverState.logger.Print(fmt.Sprintf("Parse Error: %v\n", err))
// 			return nil
// 		}
// 		response, err := json.Marshal(responseObj)
// 		if err != nil {
// 			serverState.logger.Print(fmt.Sprintf("invalid Response: %v\n", err))
// 			return nil
// 		}
// 		return response

// 	}
// 	return nil

// }
