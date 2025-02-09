package lsp

func protocolInitialize(request map[string]any) (any, error) {
	responseObj := map[string]any{
		"jsonrpc": "2.0",
		"id":      request["id"],
		"result":  "received method: initialize",
	}

	return responseObj, nil

}
