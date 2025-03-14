package lox

import (
	"encoding/json"
	"fmt"
)

func PrintParse(code string) error {
	var scanner Scanner
	var parser Parser
	var formatter Formatter
	tokens, _, err := scanner.Scan(code)
	if err != nil {
		return err
	}
	fmt.Println(tokens)

	ast, _, _, errorList := parser.Parse(tokens)
	formatCode := formatter.Format(ast)
	printable, err := (json.Marshal(ast))
	if err != nil {
		return err
	}
	fmt.Println(errorList)
	fmt.Println(string(printable))
	fmt.Println(formatCode)
	return nil
}

func ParseCode(code string) ([]CompileError, []Node, map[Token][]Token, error) {
	var scanner Scanner
	var parser Parser
	tokens, scanErrors, err := scanner.Scan(code)
	if err != nil {
		return nil, nil, nil, err
	}

	_, identifiers, references, parseErrors := parser.Parse(tokens)
	return append(parseErrors, scanErrors...), identifiers, references, nil
}

func FindErrors(code string) ([]CompileError, error) {
	var scanner Scanner
	var parser Parser
	tokens, codeErrors, err := scanner.Scan(code)
	if err != nil {
		return nil, err
	}

	_, _, _, parseErrors := parser.Parse(tokens)

	return append(parseErrors, codeErrors...), nil

}
