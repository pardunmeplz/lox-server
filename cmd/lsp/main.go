package main

//import "lox-server/internal/lsp"
import "lox-server/internal/lox"

func main() {
	// disabling server to test scanning and parsing first
	//lsp.StartServer()
	lox.ScanLoxCode(`
        var myTest = 121.5;
        var myName = "second test"
        print myTest + 25
        `)

}
