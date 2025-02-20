package lsp

import (
	"lox-server/internal/lox"
	lsp "lox-server/internal/lsp/types"
)

func initializeCheck(request map[string]any) map[string]any {
	if !serverState.initialized {
		return map[string]any{
			"jsonrpc": "2.0",
			"id":      request["id"],
			"error":   lsp.InvalidRequest,
		}
	}
	return nil
}

func shutdownCheck(request map[string]any) map[string]any {
	if !serverState.initialized {
		return map[string]any{
			"jsonrpc": "2.0",
			"id":      request["id"],
			"error":   lsp.InvalidRequest,
		}
	}
	return nil
}

func protocolInitialize(request map[string]any) (map[string]any, error) {
	shutCheck := shutdownCheck(request)
	if shutCheck != nil {
		serverState.initialized = false
		return shutCheck, nil
	}

	responseObj := map[string]any{
		"jsonrpc": "2.0",
		"id":      request["id"],
		"result": map[string]any{
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
	return responseObj, nil

}

func protocolShutdown(request map[string]any) map[string]any {
	initialCheck := initializeCheck(request)
	if initialCheck != nil {
		serverState.shutdown = false
		return initialCheck
	}
	return map[string]any{
		"jsonrpc": "2.0",
		"id":      request["id"],
		"result":  nil,
	}
}

func diagnosticNotification(code string, uri string, version int) (lsp.JsonRpcNotification, bool, error) {

	parseErrors, err := lox.FindErrors(code)
	if err != nil {
		return lsp.JsonRpcNotification{}, false, err
	}
	if parseErrors == nil || len(parseErrors) == 0 {
		return lsp.JsonRpcNotification{}, false, nil
	}

	diagnostic := []lsp.Diagnostic{}
	for _, e := range parseErrors {
		diagnostic = append(diagnostic, lsp.Diagnostic{
			Severity: 1,
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
	return responseObj, true, nil

}
