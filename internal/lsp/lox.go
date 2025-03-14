package lsp

import (
	"encoding/json"
	"lox-server/internal/lox"
	lsp "lox-server/internal/lsp/types"
	"sync"
)

/* document level logic like language features and state are handled here*/

type DocumentService struct {
	AST         []lox.Node
	Definitions []lox.Node
	References  map[lox.Token][]lox.Token
	Errors      []lox.CompileError
	Uri         string
	Mutex       sync.Mutex
	EOF         lox.Token
}

func (loxService *DocumentService) Initialize() {
	loxService.AST = make([]lox.Node, 0)
	loxService.Definitions = make([]lox.Node, 0)
	loxService.References = make(map[lox.Token][]lox.Token)
	loxService.Errors = make([]lox.CompileError, 0)
}

func (loxService *DocumentService) ParseCode(code string, version int) {
	defer loxService.Mutex.Unlock()
	loxService.Mutex.Lock()

	tokens, ast, compileErrors, definitions, references, err := lox.ParseCode(code)
	if err != nil {
		return
	}
	loxService.AST = ast
	loxService.Definitions = definitions
	loxService.Errors = compileErrors
	loxService.References = references
	loxService.EOF = tokens[len(tokens)-1]
	responseObj := diagnosticNotification(compileErrors, loxService.Uri, version)
	response, err := json.Marshal(responseObj)
	sendNotification(response)
}

func (loxService *DocumentService) GetDefinition(position lsp.Position) (lsp.Position, bool) {
	for _, definable := range loxService.Definitions {
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
	for definition := range loxService.References {
		name, ok := definition.Value.(string)
		if !ok {
			continue
		}
		atCursor := definition.Line == int(position.Line) &&
			definition.Character <= int(position.Character) &&
			definition.Character+len(name) >= int(position.Character)

		if atCursor {
			response := make([]lsp.Position, 0, 4)
			for _, reference := range loxService.References[definition] {
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
