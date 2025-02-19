package lox

import "fmt"

func ScanLoxCode(code string) {
	fmt.Println(scan(code))
}

func FindErrors(code string) ([]CompileError, error) {
	_, codeErrors, err := scan(code)
	if err != nil {
		return nil, err
	}
	return codeErrors, nil

}
