package lsp

type Diagnostic struct {
	/*
	   1 = Error
	   2 = Warning
	   3 = Info
	   4 = Hint
	*/
	Severity int    `json:"severity"`
	ErrRange Range  `json:"range"`
	Message  string `json:"message"`
}

type PublishDiagnosticParams struct {
	Uri         string       `json:"uri"`
	Version     int          `json:"version"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type JsonRpcResponse struct {
	JsonRpc string `json:"jsonRpc"`
	Id      any    `json:"id"`
	Result  any    `json:"result"`
	Error   any    `json:"error"`
}
