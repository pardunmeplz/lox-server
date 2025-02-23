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
	lox.ParseCode(`
        -2*5+2*4
        `)
}
