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

if (true) {    
    
}

fun testfunc(test,newtest,mytest) {    
    
}
`
	lox.PrintParse(test)
}
