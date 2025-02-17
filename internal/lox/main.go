package lox

import "fmt"

func ScanLoxCode(code string) {
	fmt.Println(scan(code))
}

func FindErrors(code string) []error {
	_, err := scan(code)
	if err == nil {
		return nil
	}
	return []error{err}

}
