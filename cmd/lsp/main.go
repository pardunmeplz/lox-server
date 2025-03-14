package main

import (
	"lox-server/internal/lox"
	"lox-server/internal/lsp"
)

func main() {
	//startServer()
	testLanguage()
}

func startServer() {
	lsp.StartServer()
}

func testLanguage() {
	//lox.PrintParse("{\r\n    var quacker;\r\nprint \"quacker\";\r\n}")
	lox.PrintParse("{\r\n    var quacker;\r\n    print \"quacker\";\r\n}\r\n\r\nvar\r\n\r\nprint quacker;\r\n")
}
