package parser

type JsonRpcRequest struct {
	Headers map[string]string
	Content []byte
}
