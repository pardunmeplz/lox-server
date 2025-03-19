package main

import (
	"lox-server/internal/lox"
	"lox-server/internal/lsp"
)

func main() {
	startServer()
	//testLanguage()
}

func startServer() {
	lsp.StartServer()
}

func testLanguage() {
	var test = `



while (true) {
}
`
	lox.PrintParse(test)
}
