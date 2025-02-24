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
        if(12)
        {
            print 12+55;
        } else {
        15;
        }
        `)
}
