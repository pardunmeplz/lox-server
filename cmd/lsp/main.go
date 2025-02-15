package main

//import "lox-server/internal/lsp"
import "lox-server/internal/lox"

func main() {
	// disabling server to test scanning and parsing first
	//lsp.StartServer()
	lox.ScanLoxCode(`
        +-
        */,
        // ++testing comments
        " test error"
        1232.2
        .
        `)

}
