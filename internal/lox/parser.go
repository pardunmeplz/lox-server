package lox

import (
	"fmt"
)

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
                      ("else" statement)?;

    returnStmt     → "return" expression? ";" ;

    whileStmt      → "while" "(" expression ")" statement;
   forStmt        → "for" "(" varDecl | exprStmt | ";" expression? ";" expression? ")" statement;

    block          → "{" declaration* "}";
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
	tokenList    []Token
	errorList    []CompileError
	currentToken int
}

func (parser *Parser) initialize(input []Token) {
	parser.tokenList = input
	parser.currentToken = 0
	parser.errorList = make([]CompileError, 0)
}

func (parser *Parser) Parse(input []Token) ([]Node, []CompileError) {
	parser.initialize(input)
	program := make([]Node, 0)
	for token := parser.peekParser(); token.TokenType != EOF; token = parser.peekParser() {
		program = append(program, parser.declaration())
	}
	return program, parser.errorList
}

func (parser *Parser) declaration() Node {
	switch {
	case parser.match(VAR):
		return parser.varDeclaration()
	default:
		return parser.statement()
	}
}

func (parser *Parser) varDeclaration() Node {
	identifier := parser.peekParser()
	parser.consume(IDENTIFIER, "Expected identifier after var declaration")
	var value Node = &Primary{ValType: "nil", Value: nil}
	if parser.match(EQUAL) {
		value = parser.expression()
	}
	parser.consume(SEMICOLON, "Expected ; at end of statement")

	return &VarDecl{Identifier: identifier, Value: value}

}

func (parser *Parser) statement() Node {
	switch {
	case parser.match(PRINT):
		expr := parser.expression()
		parser.consume(SEMICOLON, "Expected ; at end of statement")
		return &PrintStmt{Expr: expr}

	case parser.match(RETURN):
		if parser.match(SEMICOLON) {
			return &ReturnStmt{Expr: &Primary{ValType: "nil", Value: nil}}
		}
		expr := parser.expression()
		parser.consume(SEMICOLON, "Expected ; at end of statement")
		return &PrintStmt{Expr: expr}

	case parser.match(BRACELEFT):
		return parser.block()

	case parser.match(IF):
		return parser.ifStmt()

	case parser.match(WHILE):
		return parser.whileStmt()

	case parser.match(FOR):
		return parser.forStmt()

	default:
		return parser.exprStmt()
	}
}

func (parser *Parser) exprStmt() Node {
	expr := parser.expression()
	parser.consume(SEMICOLON, "Expected ; at end of statement")
	return &ExpressionStmt{Expr: expr}
}

func (parser *Parser) block() Node {
	body := make([]Node, 0)
	for token := parser.peekParser(); token.TokenType != EOF && token.TokenType != BRACERIGHT; token = parser.peekParser() {
		body = append(body, parser.declaration())
	}
	parser.consume(BRACERIGHT, "Expected '}' at end of block")
	return &BlockStmt{Body: body}
}

func (parser *Parser) ifStmt() Node {
	parser.consume(PARANLEFT, "Expected '(' after if")
	condition := parser.expression()
	parser.consume(PARANRIGHT, "Expected ')' after condition")

	thenBranch := parser.statement()
	var elseBranch Node = nil

	if parser.match(ELSE) {
		elseBranch = parser.statement()
	}

	return &IfStmt{Condition: condition, Then: thenBranch, Else: elseBranch}
}

func (parser *Parser) whileStmt() Node {
	parser.consume(PARANLEFT, "Expected '(' after while")
	condition := parser.expression()
	parser.consume(PARANRIGHT, "Expected ')' after condition")

	body := parser.statement()

	return &WhileStmt{Condition: condition, Then: body}
}

func (parser *Parser) forStmt() Node {
	parser.consume(PARANLEFT, "Expected '(' after for")

	var initializer Node = nil
	if !parser.match(SEMICOLON) {
		if parser.match(VAR) {
			initializer = parser.varDeclaration()
		} else {
			initializer = parser.exprStmt()
		}
	}

	var condition Node = &Primary{ValType: "boolean", Value: true}
	if !parser.match(SEMICOLON) {
		condition = parser.expression()
		parser.consume(SEMICOLON, "Expected ; after condition")
	}

	var assignment Node = nil
	if parser.peekParser().TokenType != PARANRIGHT {
		expr := parser.expression()
		assignment = &ExpressionStmt{Expr: expr}
	}

	parser.consume(PARANRIGHT, "Expected ')' before body")

	loop := &WhileStmt{Condition: condition, Then: parser.statement()}

	if assignment != nil {
		loop.Then = &BlockStmt{Body: []Node{loop.Then, assignment}}
	}

	if initializer == nil {
		return loop
	}

	return &BlockStmt{Body: []Node{initializer, loop}}
}

func (parser *Parser) expression() Node {
	return parser.logicalAnd()
}

func (parser *Parser) logicalOr() Node {
	expr := parser.logicalAnd()

	for token := parser.peekParser(); token.TokenType == OR; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.logicalAnd()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) logicalAnd() Node {
	expr := parser.equality()

	for token := parser.peekParser(); token.TokenType == AND; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.equality()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) equality() Node {
	expr := parser.comparison()

	for token := parser.peekParser(); token.TokenType == EQUALEQUAL || token.TokenType == BANGEQUAL; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.comparison()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) comparison() Node {
	expr := parser.term()

	for token := parser.peekParser(); token.TokenType == GREATER || token.TokenType == GREATEREQUAL ||
		token.TokenType == LESS || token.TokenType == LESSEQUAL; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.term()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) term() Node {
	expr := parser.factor()

	for token := parser.peekParser(); token.TokenType == PLUS || token.TokenType == MINUS; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.factor()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) factor() Node {
	expr := parser.unary()

	for token := parser.peekParser(); token.TokenType == STAR || token.TokenType == SLASH; token = parser.peekParser() {
		parser.advanceParser()
		right := parser.unary()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) unary() Node {
	if token := parser.peekParser(); token.TokenType == MINUS || token.TokenType == BANG {
		parser.advanceParser()
		return &Unary{Expression: parser.unary(), Operation: token.TokenType}
	}
	return parser.primary()
}

func (parser *Parser) primary() Node {

	currToken := parser.peekParser()
	parser.advanceParser()

	switch currToken.TokenType {
	case STRING:
		return &Primary{ValType: "string", Value: currToken.Value}
	case NUMBER:
		return &Primary{ValType: "number", Value: currToken.Value}
	case TRUE:
		return &Primary{ValType: "boolean", Value: true}
	case FALSE:
		return &Primary{ValType: "boolean", Value: true}
	case NIL:
		return &Primary{ValType: "nil", Value: nil}
	case PARANLEFT:
		expr := parser.expression()
		parser.consume(PARANRIGHT, fmt.Sprintf("Expected ')' at line %d character %d", currToken.Line, currToken.Character))
		return &Group{Expression: expr}
	}
	return &Primary{}
}

func (parser *Parser) advanceParser() {
	parser.currentToken++

}

func (parser *Parser) match(tokenType int) bool {
	if tokenType == parser.peekParser().TokenType {
		parser.advanceParser()
		return true
	}
	return false
}

func (parser *Parser) consume(tokenType int, message string) {
	if parser.match(tokenType) {
		return
	}
	parser.addError(message)
}

func (parser *Parser) addError(message string) {
	parser.errorList = append(parser.errorList, CompileError{Message: message, Line: parser.peekParser().Line, Char: parser.peekParser().Character})
}

func (parser *Parser) peekPrevious() Token {
	return parser.tokenList[parser.currentToken-1]
}

func (parser *Parser) peekParser() Token {
	return parser.tokenList[parser.currentToken]
}
