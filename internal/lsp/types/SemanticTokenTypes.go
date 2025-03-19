package lsp

const (
	namespace = "namespace"
	/**
	 * Represents a generic type. Acts as a fallback for types which
	 * can"t be mapped to a specific type like class or enum.
	 */
	Type          = "type"
	Class         = "class"
	Enum          = "enum"
	Interface     = "interface"
	Struct        = "struct"
	TypeParameter = "typeParameter"
	Parameter     = "parameter"
	Variable      = "variable"
	Property      = "property"
	EnumMember    = "enumMember"
	Event         = "event"
	Function      = "function"
	Method        = "method"
	Macro         = "macro"
	Keyword       = "keyword"
	Modifier      = "modifier"
	Comment       = "comment"
	String        = "string"
	Number        = "number"
	Regexp        = "regexp"
	Operator      = "operator"
)

type SemanticTokensLegend struct {
	TokenTypes     []string `json:"tokenTypes"`
	TokenModifiers []string `json:"tokenModifiers"`
}

var Legend = SemanticTokensLegend{
	TokenTypes:     []string{Variable, Method, Keyword, Comment, Number, String, Operator},
	TokenModifiers: []string{},
}
