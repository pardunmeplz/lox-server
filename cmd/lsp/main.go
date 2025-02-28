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
        fun testing(){
            return 12;
        }
        for (var x = 12; false;x = x + 3){
           testing(12, 55, brandon); 
        }
        `)
}
