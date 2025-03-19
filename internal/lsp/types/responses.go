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
	JsonRpc string `json:"jsonrpc"`
	Id      any    `json:"id"`
	Result  any    `json:"result"`
	Error   any    `json:"error"`
}

type Location struct {
	Uri      string `json:"uri"`
	LocRange Range  `json:"range"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type CompletionItemLabelDetails struct {
	Detail      string `json:"detail"`
	Description string `json:"description"`
}

type CompletionItem struct {
	Label string `json:"label"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type SemanticTokens struct {
	Data []uint `json:"data"`
}
