package parser

type JsonRpcRequest struct {
	Headers      map[string]string
	Content      *JsonRpcRequestContent
	ContentBytes []byte
}

type JsonRpcRequestContent struct {
	JsonRpc string `json:"jsonrpc"`
	Id      any    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}
