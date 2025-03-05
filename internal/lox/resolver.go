package lox

type Resolver struct {
	Ast         []Node
	SymbolTable map[string]Token
}

func (resolver *Resolver) initialize(ast []Node) {
	resolver.Ast = ast
	resolver.SymbolTable = make(map[string]Token)
}

func (resolver *Resolver) Resolve(ast []Node) map[string]Token {
	resolver.initialize(ast)
	for _, node := range ast {
		node.accept(resolver)
	}
	return resolver.SymbolTable
}
