package lox

import "fmt"

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

func (parser *Parser) Parse(input []token) Expr {
	parser.initialize(input)
	return parser.expression()
}

func (parser *Parser) expression() Expr {
	return parser.logicalAnd()
}

func (parser *Parser) logicalOr() Expr {
	expr := parser.logicalAnd()

	for token := parser.peekParser(); token.tokenType == OR; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.logicalAnd()
		expr = &Binary{Left: expr, Right: right, Operation: token.tokenType}
	}

	return expr
}

func (parser *Parser) logicalAnd() Expr {
	expr := parser.equality()

	for token := parser.peekParser(); token.tokenType == AND; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.equality()
		expr = &Binary{Left: expr, Right: right, Operation: token.tokenType}
	}

	return expr
}

func (parser *Parser) equality() Expr {
	expr := parser.comparison()

	for token := parser.peekParser(); token.tokenType == EQUALEQUAL || token.tokenType == BANGEQUAL; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.comparison()
		expr = &Binary{Left: expr, Right: right, Operation: token.tokenType}
	}

	return expr
}

func (parser *Parser) comparison() Expr {
	expr := parser.term()

	for token := parser.peekParser(); token.tokenType == GREATER || token.tokenType == GREATEREQUAL ||
		token.tokenType == LESS || token.tokenType == LESSEQUAL; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.term()
		expr = &Binary{Left: expr, Right: right, Operation: token.tokenType}
	}

	return expr
}

func (parser *Parser) term() Expr {
	expr := parser.factor()

	for token := parser.peekParser(); token.tokenType == PLUS || token.tokenType == MINUS; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.factor()
		expr = &Binary{Left: expr, Right: right, Operation: token.tokenType}
	}

	return expr
}

func (parser *Parser) factor() Expr {
	expr := parser.unary()

	for token := parser.peekParser(); token.tokenType == STAR || token.tokenType == SLASH; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.unary()
		expr = &Binary{Left: expr, Right: right, Operation: token.tokenType}
	}

	return expr
}

func (parser *Parser) unary() Expr {
	if token := parser.peekParser(); token.tokenType == MINUS || token.tokenType == BANG {
		parser.advanceParser()
		return &Unary{Expression: parser.unary(), Operation: token.tokenType}
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
	case PARANLEFT:
		expr := parser.expression()
		parser.consume(PARANRIGHT, fmt.Sprintf("Expected ')' at line %d character %d", currToken.line, currToken.character))
		return &Group{Expression: expr}
	}
	return &Primary{}
}

func (parser *Parser) advanceParser() {
	parser.currentToken++

}

func (parser *Parser) match(tokenType int) bool {
	if tokenType == parser.peekParser().tokenType {
		parser.advanceParser()
		return true
	}
	return false
}

func (parser *Parser) consume(tokenType int, message string) {
	if parser.match(tokenType) {
		return
	}
	// todo: add error here

}

func (parser *Parser) peekParser() token {
	return parser.tokenList[parser.currentToken]
}
