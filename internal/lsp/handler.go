package lsp

import (
	"encoding/json"
	"fmt"
)

func handleRequest(msg string) ([]byte, error) {
	var requestObj map[string]any

	if err := json.Unmarshal([]byte(msg), &requestObj); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}
	if _, exists := requestObj["id"]; exists == false {
		return nil, nil
	}

	responseObj, err := processRequest(requestObj)

	if err != nil {
		return nil, fmt.Errorf("invalid Request: %v", err)
	}

	response, err := json.Marshal(responseObj)

	if err != nil {
		return nil, fmt.Errorf("invalid Request: %v", err)
	}

	header := []byte(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(response)))

	return append(header, response...), nil
}

func processRequest(request map[string]any) (map[string]any, error) {

	switch request["method"] {
	case "initialize":
		return protocolInitialize(request)
	}

	return nil, fmt.Errorf("Invalid Method: %v", request["method"])

}
