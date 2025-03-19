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
	ScopeTable map[lox.ScopeRange][]lox.Token
	Errors     []lox.CompileError
	Uri        string
	Mutex      sync.Mutex
	EOF        lox.Token
	IsError    bool
}

func (loxService *DocumentService) Initialize() {
	loxService.AST = make([]lox.Node, 0)
	loxService.References = make([]lox.Node, 0)
	loxService.SymbolMap = make(map[lox.Token][]lox.Token)
	loxService.Errors = make([]lox.CompileError, 0)
}

func (loxService *DocumentService) ParseCode(code string, version int) {
	tokens, ast, compileErrors, references, symbolMap, scopeTable, err := lox.ParseCode(code)
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
	loxService.ScopeTable = scopeTable
	loxService.EOF = tokens[len(tokens)-1]
	loxService.IsError = false

	for _, error := range compileErrors {
		loxService.IsError = error.Source < lox.ERROR_RESOLVER || loxService.IsError
	}

	responseObj := diagnosticNotification(compileErrors, loxService.Uri, version)
	response, err := json.Marshal(responseObj)
	sendNotification(response)
}

func (loxService *DocumentService) GetCompletion(position lsp.Position) []lsp.CompletionItem {
	items := make([]lsp.CompletionItem, 0)
	var scope *lox.ScopeRange = nil

	for scopeRange := range loxService.ScopeTable {
		inScope := scopeRange.ScopeContext == lox.GLOBAL_CONTEXT || (scopeRange.StartLine <= int(position.Line) &&
			scopeRange.EndLine >= int(position.Line))
		if !inScope {
			continue
		}
		if scope == nil { // deepest scopes come first
			scope = &scopeRange
		}

		for _, definition := range loxService.ScopeTable[scopeRange] {
			label, ok := definition.Value.(string)
			if !ok {
				continue
			}
			items = append(items, lsp.CompletionItem{
				Label: label,
			})

		}

	}
	if scope.ScopeContext == lox.CLASS_CONTEXT {
		items = make([]lsp.CompletionItem, 0)
	}
	keywords := getKeywords(scope.ScopeContext, scope.ClassContext, scope.FunctionContext)

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

func (loxService *DocumentService) GetSemanticTokens() []uint {
	// line, character, length, tokentype, tokenModifier(none)
	response := []uint{}
	var lastToken *lox.Token = &lox.Token{Character: 0, Length: 0, Line: 0}

	for _, token := range loxService.Tokens {
		switch token.TokenType {
		case lox.FOR, lox.AND, lox.FUN, lox.VAR, lox.WHILE, lox.IF, lox.ELSE, lox.THIS, lox.SUPER, lox.CLASS, lox.OR, lox.PRINT, lox.RETURN:
			if token.Line == lastToken.Line {
				response = append(response, 0, uint(token.Character)-uint(lastToken.Character), uint(token.Length), 2, 0)
			} else {
				response = append(response, uint(token.Line)-uint(lastToken.Line), uint(token.Character), uint(token.Length), 2, 0)
			}
			lastToken = &token
		case lox.IDENTIFIER:
			if token.Line == lastToken.Line {
				response = append(response, 0, uint(token.Character)-uint(lastToken.Character), uint(token.Length), 0, 0)
			} else {
				response = append(response, uint(token.Line)-uint(lastToken.Line), uint(token.Character), uint(token.Length), 0, 0)
			}
			lastToken = &token
		}
	}

	return response
}
