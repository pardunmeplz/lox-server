package lsp

type TextDocumentItem struct {
	Uri        string `json:"uri"`
	LanguageId string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type TextDocumentContentChangeEvent struct {
	TextRange   Range  `json:"range"`
	RangeLength *uint  `json:"rangeLength"`
	Text        string `json:"text"`
}

type DidOpenTextDocumentParams struct {
	TextDocument   TextDocumentItem               `json:"textDocument"`
	ContentChanges TextDocumentContentChangeEvent `json:"contentChanges"`
}

type JsonRpcNotification struct {
	JsonRpc string `json:"jsonRpc"`
	Method  any    `json:"method"`
	Params  any    `json:"params"`
}
