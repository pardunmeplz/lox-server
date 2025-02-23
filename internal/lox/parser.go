package lox

/*
   program        → declaration* EOF ;

   declaration    → varDeclaration | statement | funcDecl | classDecl ;

   funcDecl       → "fun" function;
   function       → IDENTIFIER "(" parameters? ")" block;
   parameters     → IDENTIFIER ( "," IDENTIFIER )*;

   classDecl      → "class" IDENTIFIER ( "<" IDENTIFIER )? "{" function* "}" ;

   varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;

   statement      → exprStmt | ifStmt | whileStmt | forStmt | returnStmt |  printStmt | block;
   ifStmt         → "if" "(" expression ")" statement
                      (else statement)?;

   returnStmt     → "return" expression? ";" ;

   whileStmt      → "while" "(" expression ")" statement;
   forStmt        → "for" "(" varDecl | exprStmt | ";" expression? ";" expression? ")" statement;

   block          → "{" declaration "}";
   exprStmt       → expression ";" ;
   printStmt      → "print" expression ";" ;

   expression     → assignment;
   assignment     → (call ".")? IDENTIFIER "=" assignment | logicalOr;
   logicalOr      → logicalAnd ( "or" logicalAnd)*;
   logicalAnd     → equality ( "and" equality)*;
   equality       → comparison ( ( "!=" | "==" ) comparison )* ;
   comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
   term           → factor ( ( "-" | "+" ) factor )* ;
   factor         → unary ( ( "/" | "*" ) unary )* ;
   unary          → ( "!" | "-" ) unary | primary ;
   call           → primary ( "(" arguments? ")" )* | getExpression;
   getExpression  → primary ( "." IDENTIFIER )*;
   arguments      → expression ( "," expression )*;
   primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | "super" "." IDENTIFIER ;
*/

type Parser struct {
	tokenList    []token
	currentToken int
}

func (parser *Parser) initialize(input []token) {
	parser.tokenList = input
	parser.currentToken = 0
}

func (parser *Parser) Parse(input []token) {
	parser.initialize(input)
}

func expression() {
}

func (parser *Parser) term() {

}

func (parser *Parser) factor() Expr {
	left := parser.unary()
	token := parser.peekParser()

	switch token.tokenType {
	case STAR:
		right := parser.factor()
		return &Binary{Left: left, Right: right, Operation: '*'}
	case SLASH:
		right := parser.factor()
		return &Binary{Left: left, Right: right, Operation: '/'}
	}
	return left

}

func (parser *Parser) unary() Expr {
	token := parser.peekParser()

	switch token.tokenType {
	case BANG:
		expr := parser.unary()
		return &Unary{Expression: expr, Operation: '!'}
	case MINUS:
		expr := parser.unary()
		return &Unary{Expression: expr, Operation: '-'}
	}
	return parser.primary()
}

func (parser *Parser) primary() Expr {

	currToken := parser.peekParser()
	parser.advanceParser()

	switch currToken.tokenType {
	case STRING:
		return &Primary{ValType: "string", Value: currToken.value}
	case NUMBER:
		return &Primary{ValType: "number", Value: currToken.value}
	case TRUE:
		return &Primary{ValType: "boolean", Value: true}
	case FALSE:
		return &Primary{ValType: "boolean", Value: true}
	case NIL:
		return &Primary{ValType: "nil", Value: nil}
	}
	return &Primary{}
}

func (parser *Parser) advanceParser() {
	parser.currentToken++
}

func (parser *Parser) peekParser() token {
	return parser.tokenList[parser.currentToken]
}
