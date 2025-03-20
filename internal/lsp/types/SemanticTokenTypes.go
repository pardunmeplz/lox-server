package lsp

const (
	/**
	 * Represents a generic type. Acts as a fallback for types which
	 * can"t be mapped to a specific type like class or enum.
	 */
	Type      = "type"
	Parameter = "parameter"
	Variable  = "variable"
	Property  = "property"
	Function  = "function"
	Method    = "method"
	Keyword   = "keyword"
	Modifier  = "modifier"
	Comment   = "comment"
	String    = "string"
	Number    = "number"
	Operator  = "operator"
)

type SemanticTokensLegend struct {
	TokenTypes     []string `json:"tokenTypes"`
	TokenModifiers []string `json:"tokenModifiers"`
}

var Legend = SemanticTokensLegend{
	TokenTypes:     []string{Variable, Method, Keyword, Type, Comment, Number, String, Operator, Parameter, Property},
	TokenModifiers: []string{},
}
