package lsp

import (
	"encoding/json"
	"fmt"
)

func handleRequest(message []byte) (any, error) {
	var request map[string]any
	if err := json.Unmarshal(message, &request); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	response := map[string]any{
		"jsonrpc": "2.0",
		"id":      request["id"],
		"result":  fmt.Sprintf("received method: %v", request["method"]),
	}

	return response, nil
}
