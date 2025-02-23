package lox

import "fmt"

func ScanCode(code string) {
	var scanner Scanner
	fmt.Println(scanner.Scan(code))
}

func FindErrors(code string) ([]CompileError, error) {
	var scanner Scanner
	_, codeErrors, err := scanner.Scan(code)
	if err != nil {
		return nil, err
	}
	return codeErrors, nil

}
