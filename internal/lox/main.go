package lox

import (
	"encoding/json"
	"fmt"
)

func PrintParse(code string) error {
	var scanner Scanner
	var parser Parser
	tokens, _, err := scanner.Scan(code)
	if err != nil {
		return err
	}
	fmt.Println(tokens)

	ast, _, errorList := parser.Parse(tokens)
	printable, err := (json.Marshal(ast))
	if err != nil {
		return err
	}
	fmt.Println(errorList)
	fmt.Println(string(printable))
	return nil
}

func ParseCode(code string) ([]CompileError, []Node, error) {
	var scanner Scanner
	var parser Parser
	tokens, scanErrors, err := scanner.Scan(code)
	if err != nil {
		return nil, nil, err
	}

	_, identifiers, parseErrors := parser.Parse(tokens)
	return append(parseErrors, scanErrors...), identifiers, nil
}

func FindErrors(code string) ([]CompileError, error) {
	var scanner Scanner
	var parser Parser
	tokens, codeErrors, err := scanner.Scan(code)
	if err != nil {
		return nil, err
	}

	_, _, parseErrors := parser.Parse(tokens)

	return append(parseErrors, codeErrors...), nil

}
