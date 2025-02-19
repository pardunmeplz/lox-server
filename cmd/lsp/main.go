package main

import (
	"fmt"
	"lox-server/internal/lox"
	"lox-server/internal/lsp"
)

func main() {
	startServer()
}

func startServer() {
	lsp.StartServer()
}

func testLanguage() {
	errors, err := lox.FindErrors(`
class testing{
  honk(){
    print this.name + "says Honkkkk";
  }
}

class subTesting < testing{
  honk(){
    print this.name + "says Haaank";
    super.honk();
  }
}

var test = subTesting();
^
test.field = 22;
test.name = "Quacker
        `)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println(errors)
}
