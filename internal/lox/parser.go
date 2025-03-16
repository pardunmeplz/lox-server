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
   unary          → ( "!" | "-" ) unary | call ;
   call           → primary ( "(" arguments? ")" )* | getExpression;
   getExpression  → primary ( "." IDENTIFIER )*;
   arguments      → expression ( "," expression )*;
   primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | "super" "." IDENTIFIER ;
*/

type SymbolMap struct {
	currentTable map[string]Token
	previous     *SymbolMap
	currScope    int
}

type Parser struct {
	tokenList       []Token
	errorList       []CompileError
	currentToken    int
	symbolMap       *SymbolMap
	identifierNodes []Node
	references      map[Token][]Token
}

func (parser *Parser) initialize(input []Token) {
	parser.tokenList = input
	parser.currentToken = 0
	parser.errorList = make([]CompileError, 0)
	parser.symbolMap = &SymbolMap{
		currentTable: make(map[string]Token),
		previous:     nil,
		currScope:    0,
	}
	parser.references = make(map[Token][]Token)
}

func (parser *Parser) isGlobal() bool { return parser.symbolMap.currScope == 0 }

func (parser *Parser) raiseScope() {
	parser.symbolMap = &SymbolMap{
		currentTable: make(map[string]Token),
		previous:     parser.symbolMap,
		currScope:    parser.symbolMap.currScope + 1,
	}
}

func (parser *Parser) closeScope() {
	if parser.symbolMap != nil && parser.symbolMap.previous != nil {
		parser.symbolMap = parser.symbolMap.previous
	}
}

func (parser *Parser) getDefinition(name string) (Token, bool) {
	current := parser.symbolMap
	for current != nil {
		if token, isPresent := current.currentTable[name]; isPresent {
			return token, true
		}
		current = current.previous
	}
	return Token{}, false
}

func (parser *Parser) getDefinitionInScope(name string) (Token, bool) {
	token, isPresent := parser.symbolMap.currentTable[name]
	return token, isPresent
}

func (parser *Parser) addIdentifier(node Node, definition Token, reference Token) {
	parser.identifierNodes = append(parser.identifierNodes, node)
	parser.references[definition] = append(parser.references[definition], reference)
}

func (parser *Parser) addDefinition(token Token) {

	name, ok := token.Value.(string)
	if ok {
		definition, isPresent := parser.getDefinitionInScope(name)
		if isPresent && parser.isGlobal() {
			parser.addWarning(fmt.Sprintf("%s is already declared in this scope at line %d", name, definition.Line+1))
		} else if isPresent {
			parser.addError(fmt.Sprintf("%s is already declared in this scope at line %d", name, definition.Line+1))
		}
		parser.symbolMap.currentTable[name] = token
		parser.references[token] = []Token{}
	}
}

func (parser *Parser) Parse(input []Token) ([]Node, []Node, map[Token][]Token, []CompileError) {
	parser.initialize(input)
	program := make([]Node, 0)
	for token := parser.peekParser(); token.TokenType != EOF; token = parser.peekParser() {
		program = append(program, parser.declaration())
	}
	for name := range parser.references {
		if len(parser.references[name]) == 0 {
			parser.addWarningAt("No usages after definition", name.Line, name.Character)
		}
	}
	return program, parser.identifierNodes, parser.references, parser.errorList
}

func (parser *Parser) declaration() Node {
	switch {
	case parser.match(VAR):
		return parser.varDeclaration()
	case parser.match(FUN):
		return parser.funcDeclaration()
	case parser.match(CLASS):
		return parser.classDeclaration()
	case parser.match(NEWLINE):
		return &NewLine{Token: parser.peekPrevious()}
	case parser.match(COMMENT):
		return &Comment{Comment: parser.peekPrevious()}
	default:
		return parser.statement()
	}
}

func (parser *Parser) classDeclaration() Node {
	identifier := parser.peekParser()
	parser.addDefinition(identifier)
	parser.consume(IDENTIFIER, "Expected identifier for class name")

	var parent *Token
	if parser.match(LESS) {
		token := parser.peekParser()
		parser.consume(IDENTIFIER, "Expected identifier for class name")
		parent = &token
	}

	parser.consume(BRACELEFT, "Expected '{' before class body")
	parser.raiseScope()
	methods := make([]Node, 0)
	for token := parser.peekParser().TokenType; token != BRACERIGHT && token != EOF; token = parser.peekParser().TokenType {
		method := parser.funcDeclaration()
		if method == nil {
			continue
		}
		methods = append(methods, method)
	}

	parser.consume(BRACERIGHT, "Expect '}' at end of class declaration")
	parser.closeScope()

	return &ClassDecl{Body: methods, Name: identifier, Parent: parent}
}

func (parser *Parser) varDeclaration() Node {
	identifier := parser.peekParser()
	parser.addDefinition(identifier)
	parser.consume(IDENTIFIER, "Expected identifier after var declaration")

	var value Node = &Primary{ValType: "nil", Value: nil}
	initialzied := false
	if parser.match(EQUAL) {
		value = parser.expression()
		initialzied = true
	}
	parser.consume(SEMICOLON, "Expected ; at end of statement")

	return &VarDecl{Identifier: identifier, Value: value, Initialized: initialzied}

}

func (parser *Parser) funcDeclaration() Node {
	identifier := parser.peekParser()
	parser.addDefinition(identifier)
	parser.consume(IDENTIFIER, "Expected identifier for function name")

	parser.consume(PARANLEFT, "Expected ( after function name")
	parameters := make([]Node, 0)
	if !parser.match(PARANRIGHT) {
		parameters = parser.parameters()
	}

	parser.consume(BRACELEFT, "Expected { at start of function body")
	body := parser.block()

	return &FuncDecl{Name: identifier, Body: body, Parameters: parameters}

}

func (parser *Parser) parameters() []Node {
	parameters := make([]Node, 0)

	parser.consume(IDENTIFIER, "Expected Parameter Name")
	parameters = append(parameters, &Variable{Identifier: parser.peekPrevious()})

	for parser.match(COMMA) {
		parser.consume(IDENTIFIER, "Expected Parameter Name")
		parameters = append(parameters, &Variable{Identifier: parser.peekPrevious()})
	}

	parser.consume(PARANRIGHT, "Expected ')' before function body")

	return parameters

}

func (parser *Parser) statement() Node {
	fmt.Printf(">>%s", parser.peekParser().Value)
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
	parser.raiseScope()
	body := make([]Node, 0)
	for token := parser.peekParser(); token.TokenType != EOF && token.TokenType != BRACERIGHT; token = parser.peekParser() {
		body = append(body, parser.declaration())
	}
	parser.consume(BRACERIGHT, "Expected '}' at end of block")
	parser.closeScope()
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

	var assign Node = nil
	if parser.peekParser().TokenType != PARANRIGHT {
		expr := parser.expression()
		assign = &ExpressionStmt{Expr: expr}
	}

	parser.consume(PARANRIGHT, "Expected ')' before body")

	loop := &WhileStmt{Condition: condition, Then: parser.statement()}

	if assign != nil {
		loop.Then = &BlockStmt{Body: []Node{loop.Then, assign}}
	}

	if initializer == nil {
		return loop
	}

	return &BlockStmt{Body: []Node{initializer, loop}}
}

func (parser *Parser) expression() Node {
	return parser.assignment()
}

func (parser *Parser) assignment() Node {
	expr := parser.logicalOr()

	if parser.match(EQUAL) {

		token := parser.peekPrevious()
		value := parser.assignment()

		variable, ok := expr.(*Variable)
		if !ok {
			parser.addErrorAt("Invalid assignment target", token.Line, token.Character)
			return expr
		}
		expr = &Assignment{Identifier: variable, Value: value}
	}
	return expr
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
	return parser.call()
}

func (parser *Parser) call() Node {
	expr := parser.primary()

	for parser.match(PARANLEFT) || parser.match(DOT) {
		if parser.peekPrevious().TokenType == PARANLEFT {
			expr = parser.finishCall(expr)
		} else {
			expr = parser.getExpression(expr)
		}
	}
	return expr
}

func (parser *Parser) getExpression(object Node) Node {
	property := parser.peekParser()
	parser.consume(IDENTIFIER, "Expected Property name")

	return &GetExpr{Object: object, Property: property}
}

func (parser *Parser) finishCall(callee Node) Node {
	if parser.match(PARANRIGHT) {
		return &Call{Callee: callee, Argument: make([]Node, 0)}
	}
	arguments := parser.arguments()
	parser.consume(PARANRIGHT, "Expected ')' and end of function call")
	if len(arguments) > 255 {
		parser.addError("Can't have more than 255 arguments")
	}
	return &Call{Callee: callee, Argument: arguments}
}

func (parser *Parser) arguments() []Node {
	response := make([]Node, 0)
	response = append(response, parser.expression())
	for parser.match(COMMA) {
		response = append(response, parser.expression())
	}

	return response
}

func (parser *Parser) primary() Node {

	currToken := parser.peekParser()

	switch {
	case parser.match(STRING):
		return &Primary{ValType: "string", Value: currToken.Value}
	case parser.match(NUMBER):
		return &Primary{ValType: "number", Value: currToken.Value}
	case parser.match(TRUE):
		return &Primary{ValType: "boolean", Value: true}
	case parser.match(FALSE):
		return &Primary{ValType: "boolean", Value: false}
	case parser.match(NIL):
		return &Primary{ValType: "nil", Value: nil}
	case parser.match(THIS):
		return &This{Identifier: currToken}
	case parser.match(SUPER):
		parser.consume(DOT, "Expected '.' after super")
		parser.consume(IDENTIFIER, "Expected method name for super-class")
		return &Super{Identifier: currToken, Property: parser.peekPrevious()}
	case parser.match(IDENTIFIER):
		name, ok := currToken.Value.(string)
		var definition Token
		if ok {
			definition, ok = parser.getDefinition(name)
		}
		result := Variable{Identifier: currToken, Definition: definition}

		if !ok {
			parser.addError(fmt.Sprintf("%s is not defined in current scope", name))
		} else {
			parser.addIdentifier(&result, definition, currToken)
		}

		return &result
	case parser.match(PARANLEFT):
		expr := parser.expression()
		parser.consume(PARANRIGHT, fmt.Sprintf("Expected ')' at line %d character %d", currToken.Line, currToken.Character))
		return &Group{Expression: expr}

	case parser.peekParser().TokenType == (EOF):
		parser.addError("Unexpected end of file")

	default:
		parser.addError(fmt.Sprintf("Unexpedted token at line %d character %d", currToken.Line, currToken.Character))
		parser.advanceParser()
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
	parser.errorList = append(parser.errorList, CompileError{Message: message, Line: parser.peekParser().Line, Char: parser.peekParser().Character, Severity: 1})
}

func (parser *Parser) addWarning(message string) {
	parser.errorList = append(parser.errorList, CompileError{Message: message, Line: parser.peekParser().Line, Char: parser.peekParser().Character, Severity: 2})
}

func (parser *Parser) addWarningAt(message string, line int, char int) {
	parser.errorList = append(parser.errorList, CompileError{Message: message, Line: line, Char: char, Severity: 2})
}

func (parser *Parser) addErrorAt(message string, line int, char int) {
	parser.errorList = append(parser.errorList, CompileError{Message: message, Line: line, Char: char, Severity: 1})
}

func (parser *Parser) peekPrevious() Token {
	return parser.tokenList[parser.currentToken-1]
}

func (parser *Parser) peekParser() Token {
	return parser.tokenList[parser.currentToken]
}

func (parser *Parser) peekNext() (Token, bool) {
	if len(parser.tokenList) >= parser.currentToken+1 {
		return Token{}, false
	}
	return parser.tokenList[parser.currentToken+1], true
}
