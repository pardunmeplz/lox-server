package lsp

import "lox-server/internal/lox"

var classContextKeywords []string = []string{
	//keywords
	"this",
	"super.",
}

var classSnippets []string = []string{
	//snippets
	"name() {}",
}

var commonKeywords []string = []string{
	//snippets
	"fun name() {}",
	"class Name {init(){ }}",
	"for (var i = 0; i < ; i = i + 1) { }",
	//keywords
	"if",
	"true",
	"false",
	"nil",
	"else",
	"for",
	"while",
	"fun",
	"class",
	"var",
	"and",
	"or",
	"print",
	// native functions
	"clock",
}

func getKeywords(scopeContext int, classContext int, functionContext int) []string {
	if scopeContext == lox.CLASS_CONTEXT {
		return classSnippets
	}
	keywords := commonKeywords

	if functionContext != lox.GLOBAL_CONTEXT {
		keywords = append(keywords, "return")
	}

	if classContext == lox.CLASS_CONTEXT {
		keywords = append(keywords, classContextKeywords...)
	}

	return keywords

}
