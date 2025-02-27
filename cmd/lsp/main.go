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
        for (var x = 12; false;x = x + 3){
           testing(12, 55, brandon); 
        }
        `)
}
