package lsp

import (
	"lox-server/internal/lox"
	lsp "lox-server/internal/lsp/types"
)

func initializeCheck(request map[string]any) *lsp.JsonRpcResponse {
	if !serverState.initialized {
		response := lsp.JsonRpcResponse{
			JsonRpc: "2.0",
			Id:      request["id"],
			Error:   lsp.InvalidRequest,
		}
		return &response
	}
	return nil
}

func shutdownCheck(request map[string]any) *lsp.JsonRpcResponse {
	if !serverState.initialized {
		response := lsp.JsonRpcResponse{
			JsonRpc: "2.0",
			Id:      request["id"],
			Error:   lsp.InvalidRequest,
		}
		return &response
	}
	return nil
}

func protocolInitialize(request map[string]any) (*lsp.JsonRpcResponse, error) {
	shutCheck := shutdownCheck(request)
	if shutCheck != nil {
		serverState.initialized = false
		return shutCheck, nil
	}

	responseObj := lsp.JsonRpcResponse{
		JsonRpc: "2.0",
		Id:      request["id"],
		Result: map[string]any{
			"capabilities": map[string]any{
				"textDocumentSync": map[string]any{
					"openClose": true,
					"change":    1,
				},
			},
			"serverInfo": map[string]any{
				"name":    "LoxServer",
				"version": "0.1.0",
			}},
	}

	return &responseObj, nil

}

func protocolShutdown(request map[string]any) *lsp.JsonRpcResponse {

	initialCheck := initializeCheck(request)
	if initialCheck != nil {
		serverState.shutdown = false
		return initialCheck
	}
	responseObj := lsp.JsonRpcResponse{
		JsonRpc: "2.0",
		Id:      request["id"],
		Result:  nil,
	}
	return &responseObj
}

func diagnosticNotification(code string, uri string, version int) (lsp.JsonRpcNotification, error) {

	parseErrors, err := lox.FindErrors(code)
	if err != nil || parseErrors == nil {
		parseErrors = make([]lox.CompileError, 0)
	}

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

	return responseObj, err
}
