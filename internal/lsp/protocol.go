package lsp

import (
	"encoding/json"
	"fmt"
	"lox-server/internal/lox"
	lsp "lox-server/internal/lsp/types"
)

func initializeCheck(request lsp.JsonRpcRequest) *lsp.JsonRpcResponse {
	if !serverState.initialized {
		response := lsp.JsonRpcResponse{
			JsonRpc: "2.0",
			Id:      request.Id,
			Error:   lsp.InvalidRequest,
		}
		return &response
	}
	return nil
}

func shutdownCheck(request lsp.JsonRpcRequest) *lsp.JsonRpcResponse {
	if !serverState.initialized {
		response := lsp.JsonRpcResponse{
			JsonRpc: "2.0",
			Id:      request.Id,
			Error:   lsp.InvalidRequest,
		}
		return &response
	}
	return nil
}

func protocolInitialize(request lsp.JsonRpcRequest) (*lsp.JsonRpcResponse, error) {
	shutCheck := shutdownCheck(request)
	if shutCheck != nil {
		serverState.initialized = false
		return shutCheck, nil
	}

	responseObj := lsp.JsonRpcResponse{
		JsonRpc: "2.0",
		Id:      request.Id,
		Result: map[string]any{
			"capabilities": map[string]any{
				"textDocumentSync": map[string]any{
					"openClose": true,
					"change":    1,
				},
				"definitionProvider": true,
				"referencesProvider": true,
			},
			"serverInfo": map[string]any{
				"name":    "LoxServer",
				"version": "0.1.0",
			},
		},
	}

	return &responseObj, nil

}

func protocolShutdown(request lsp.JsonRpcRequest) *lsp.JsonRpcResponse {

	initialCheck := initializeCheck(request)
	if initialCheck != nil {
		serverState.shutdown = false
		return initialCheck
	}
	responseObj := lsp.JsonRpcResponse{
		JsonRpc: "2.0",
		Id:      request.Id,
		Result:  nil,
	}
	return &responseObj
}

func protocolReferences(request lsp.JsonRpcRequest) *lsp.JsonRpcResponse {

	responseObj := lsp.JsonRpcResponse{
		JsonRpc: "2.0",
		Id:      request.Id,
		Result:  nil,
	}

	requestjson, err := json.Marshal(request.Params)
	var requestObj lsp.ReferenceParams

	if err != nil {
		return &responseObj
	}

	err = json.Unmarshal(requestjson, &requestObj)

	if err != nil {
		return &responseObj
	}

	document, ok := serverState.documents[requestObj.TextDocument.Uri]
	if !ok {
		serverState.logger.Print(fmt.Sprintf("Get Reference Error: URI %s not found", requestObj.TextDocument.Uri))
		return &responseObj
	}

	references := document.GetReferences(requestObj.Position, requestObj.Context.IncludeDeclaration)

	if references == nil {
		return &responseObj
	}

	test, err := json.Marshal(references)

	serverState.logger.Print(fmt.Sprintf("References: %s", test))

	responseParams := make([]lsp.Location, 0, 4)
	for _, reference := range references {
		responseParams = append(responseParams, lsp.Location{
			Uri: requestObj.TextDocument.Uri,
			LocRange: lsp.Range{
				Start: reference,
				End:   reference,
			},
		})
	}

	responseObj.Result = responseParams
	return &responseObj
}

func protocolDefinition(request lsp.JsonRpcRequest) *lsp.JsonRpcResponse {

	responseObj := lsp.JsonRpcResponse{
		JsonRpc: "2.0",
		Id:      request.Id,
		Result:  nil,
	}

	requestjson, err := json.Marshal(request.Params)
	var requestObj lsp.DefinitionParams

	if err != nil {
		return &responseObj
	}

	err = json.Unmarshal(requestjson, &requestObj)

	if err != nil {
		return &responseObj
	}

	definition, _ := serverState.documents[requestObj.TextDocument.Uri].GetDefinition(requestObj.Position)

	responseObj.Result = lsp.Location{
		Uri: requestObj.TextDocument.Uri,
		LocRange: lsp.Range{
			Start: definition,
			End:   definition,
		},
	}

	return &responseObj
}

func diagnosticNotification(parseErrors []lox.CompileError, uri string, version int) lsp.JsonRpcNotification {

	diagnostic := []lsp.Diagnostic{}
	for _, e := range parseErrors {
		diagnostic = append(diagnostic, lsp.Diagnostic{
			Severity: e.Severity,
			Message:  e.Message,
			ErrRange: lsp.Range{
				Start: lsp.Position{
					Line:      uint(e.Line),
					Character: uint(e.Char),
				},
				End: lsp.Position{
					Line:      uint(e.Line),
					Character: uint(e.Char),
				},
			},
		})
	}

	result := lsp.PublishDiagnosticParams{Uri: uri, Version: version, Diagnostics: diagnostic}

	responseObj := lsp.JsonRpcNotification{
		JsonRpc: "2.0",
		Params:  result,
		Method:  "textDocument/publishDiagnostics",
	}

	return responseObj
}

func register(id int) lsp.JsonRpcRequest {
	requestObj := lsp.JsonRpcRequest{
		JsonRpc: "2.0",
		Id:      id,
		Method:  "client/registerCapability",
		Params: lsp.RegistrationParams{
			Registrations: []lsp.Registration{{
				Id:     "definition",
				Method: "textDocument/definition",
				RegisterOptions: lsp.DefinitionRegistrationOptions{
					TextDocumentRegistrationOptions: lsp.TextDocumentRegistrationOptions{
						DocumentSelector: []lsp.DocumentFilter{
							{
								Pattern: "**/*.lox",
							},
						},
					},
				},
			}},
		},
	}

	return requestObj

}
