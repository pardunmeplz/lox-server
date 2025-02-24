package lox

import (
	"encoding/json"
	"fmt"
)

func ParseCode(code string) error {
	var scanner Scanner
	var parser Parser
	tokens, _, err := scanner.Scan(code)
	if err != nil {
		return err
	}
	fmt.Println(tokens)
	ast, _ := parser.Parse(tokens)
	printable, err := (json.Marshal(ast))
	if err != nil {
		return err
	}
	fmt.Println(string(printable))
	return nil
}

func FindErrors(code string) ([]CompileError, error) {
	var scanner Scanner
	_, codeErrors, err := scanner.Scan(code)
	if err != nil {
		return nil, err
	}
	return codeErrors, nil

}
