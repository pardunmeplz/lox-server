package lsp

import lsp "lox-server/internal/lsp/types"

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
			"capabilities": map[string]any{},
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
