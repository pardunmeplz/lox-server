package lsp

func protocolInitialize(request map[string]any) (map[string]any, error) {
	responseObj := map[string]any{
		"jsonrpc":      "2.0",
		"id":           request["id"],
		"capabilities": map[string]any{},
		"serverInfo": map[string]any{
			"name":    "LoxServer",
			"version": "0.1.0",
		},
	}
	return responseObj, nil

}
