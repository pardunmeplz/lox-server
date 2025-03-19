package lsp

type TextDocumentIdentifier struct {
	Uri string `json:"uri"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

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
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type RegistrationParams struct {
	Registrations []Registration `json:"registrations"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type DefinitionParams struct {
	TextDocumentPositionParams `json:",inline"`
}

type ReferenceContext struct {
	IncludeDeclaration bool `json:"includeDeclaration"`
}

type ReferenceParams struct {
	TextDocumentPositionParams `json:",inline"`
	Context                    ReferenceContext `json:"context"`
}

type CompletionParams struct {
	TextDocumentPositionParams `json:",inline"`
}

type SemanticTokensParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

type FormattingOptions struct {
	TabSize                int  `json:"tabSize"`
	InsertSpaces           bool `json:"insertSpaces"`
	TrimTrailingWhiteSpace bool `json:"trimTrailingWhiteSpace"`
	InsertFinalNewLine     bool `json:"insertFinalNewLine"`
	TrimFinalNewLines      bool `json:"trimFinalNewLines"`
}

type Registration struct {
	Id              string `json:"id"`
	Method          string `json:"method"`
	RegisterOptions any    `json:"registerOptions"`
}

type DocumentFilter struct {
	Language string `json:"language"`
	Scheme   string `json:"scheme"`
	Pattern  string `json:"pattern"`
}

type TextDocumentRegistrationOptions struct {
	DocumentSelector []DocumentFilter `json:"documentSelector"`
}

type DefinitionRegistrationOptions struct {
	TextDocumentRegistrationOptions `json:",inline"`
}

type JsonRpcNotification struct {
	JsonRpc string `json:"jsonrpc"`
	Method  any    `json:"method"`
	Params  any    `json:"params"`
}

type Id interface {
	string | int
}

type JsonRpcRequest struct {
	JsonRpc string `json:"jsonrpc"`
	Id      any    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}
