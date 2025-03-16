package lsp

import (
	"encoding/json"
	"lox-server/internal/lox"
	lsp "lox-server/internal/lsp/types"
	"sync"
)

/* document level logic like language features and state are handled here*/

type DocumentService struct {
	AST        []lox.Node
	Tokens     []lox.Token
	References []lox.Node
	SymbolMap  map[lox.Token][]lox.Token
	Errors     []lox.CompileError
	Uri        string
	Mutex      sync.Mutex
	EOF        lox.Token
}

var keywords []string = []string{
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
	"this",
	"super",
	"return",
}

func (loxService *DocumentService) Initialize() {
	loxService.AST = make([]lox.Node, 0)
	loxService.References = make([]lox.Node, 0)
	loxService.SymbolMap = make(map[lox.Token][]lox.Token)
	loxService.Errors = make([]lox.CompileError, 0)
}

func (loxService *DocumentService) ParseCode(code string, version int) {
	tokens, ast, compileErrors, references, symbolMap, err := lox.ParseCode(code)
	if err != nil {
		return
	}

	defer loxService.Mutex.Unlock()
	loxService.Mutex.Lock()

	loxService.AST = ast
	loxService.Tokens = tokens
	loxService.References = references
	loxService.Errors = compileErrors
	loxService.SymbolMap = symbolMap
	loxService.EOF = tokens[len(tokens)-1]

	responseObj := diagnosticNotification(compileErrors, loxService.Uri, version)
	response, err := json.Marshal(responseObj)
	sendNotification(response)
}

func (loxService *DocumentService) GetCompletion(position lsp.Position) []lsp.CompletionItem {
	items := make([]lsp.CompletionItem, 0, 0)
	for definition := range loxService.SymbolMap {
		label, ok := definition.Value.(string)
		if !ok {
			continue
		}
		items = append(items, lsp.CompletionItem{
			Label: label,
		})
	}

	for _, label := range keywords {
		items = append(items, lsp.CompletionItem{
			Label: label,
		})
	}
	return items

}

func (loxService *DocumentService) GetToken(position lsp.Position) lox.Token {
	var currToken lox.Token
	for _, token := range loxService.Tokens {
		crossedCursor := token.Line > int(position.Line) ||
			(token.Line == int(position.Line) && token.Character > int(position.Character))

		if crossedCursor {
			return currToken
		}
		currToken = token
	}
	return currToken
}

func (loxService *DocumentService) GetDefinition(position lsp.Position) (lsp.Position, bool) {
	for _, definable := range loxService.References {
		switch definable.(type) {
		case *lox.Variable:
			variable := definable.(*lox.Variable)
			name, ok := variable.Identifier.Value.(string)
			if !ok {
				continue
			}
			atCursor := variable.Identifier.Line == int(position.Line) &&
				variable.Identifier.Character <= int(position.Character) &&
				variable.Identifier.Character+len(name) >= int(position.Character)

			if atCursor {
				return lsp.Position{
					Line:      uint(variable.Definition.Line),
					Character: uint(variable.Definition.Character),
				}, true
			}
		default:
			continue
		}
	}
	return position, false

}

func (loxService *DocumentService) GetFormattedCode() string {
	var formatter lox.Formatter
	return formatter.Format(loxService.AST)
}

func (loxService *DocumentService) GetReferences(position lsp.Position, addDefinition bool) []lsp.Position {
	// check if cursor is on a definition
	for definition := range loxService.SymbolMap {
		name, ok := definition.Value.(string)
		if !ok {
			continue
		}
		atCursor := definition.Line == int(position.Line) &&
			definition.Character <= int(position.Character) &&
			definition.Character+len(name) >= int(position.Character)

		if atCursor {
			response := make([]lsp.Position, 0, 4)
			for _, reference := range loxService.SymbolMap[definition] {
				response = append(response, lsp.Position{
					Line:      uint(reference.Line),
					Character: uint(reference.Character),
				})

			}

			if addDefinition {
				response = append(response, lsp.Position{
					Line:      uint(definition.Line),
					Character: uint(definition.Character),
				})
			}
			return response
		}
	}

	// check if cusor is on a reference
	definition, found := loxService.GetDefinition(position)
	if found {
		return loxService.GetReferences(definition, addDefinition)
	}

	return nil

}
