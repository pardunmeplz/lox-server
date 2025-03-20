package lox

import (
	"fmt"
	"slices"
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

var NativeFunctions []string = []string{"clock"}

type SymbolMap struct {
	currentTable    map[string]Token
	previous        *SymbolMap
	currScope       int
	definitions     []Token
	scopeRange      ScopeRange
	functionContext int // global / function / method
	classContext    int // global / class
}

type ScopeRange struct {
	StartLine       int
	StartChar       int
	EndLine         int
	EndChar         int
	ScopeContext    int // if scope is a function / class / loop ... etc
	FunctionContext int // global / function / method
	ClassContext    int // global / class
}

type Parser struct {
	tokenList       []Token
	errorList       []CompileError
	currentToken    int
	symbolMap       *SymbolMap
	identifierNodes []Node
	references      map[Token][]Token
	scopeTable      map[ScopeRange][]Token
	scopeRanges     []ScopeRange
	panicMode       bool
}

func (parser *Parser) initialize(input []Token) {
	parser.tokenList = input
	parser.currentToken = 0
	parser.errorList = make([]CompileError, 0)
	parser.symbolMap = &SymbolMap{
		currentTable:    make(map[string]Token),
		previous:        nil,
		currScope:       0,
		functionContext: GLOBAL_CONTEXT,
	}
	parser.references = make(map[Token][]Token)
	parser.scopeTable = map[ScopeRange][]Token{}
}

func (parser *Parser) isGlobal() bool { return parser.symbolMap.currScope == 0 }

func (parser *Parser) raiseScope(startLine int, startChar int, scopeContext int) {
	functionContext := parser.symbolMap.functionContext
	if scopeContext == GLOBAL_CONTEXT || scopeContext == METHOD_CONTEXT || scopeContext == FUNCTION_CONTEXT {
		functionContext = scopeContext
	}
	classContext := parser.symbolMap.classContext
	if scopeContext == GLOBAL_CONTEXT || scopeContext == CLASS_CONTEXT {
		classContext = scopeContext
	}
	parser.symbolMap = &SymbolMap{
		currentTable:    make(map[string]Token),
		definitions:     make([]Token, 0),
		previous:        parser.symbolMap,
		currScope:       parser.symbolMap.currScope + 1,
		functionContext: functionContext,
		classContext:    classContext,
		scopeRange: ScopeRange{
			StartChar:       startChar,
			StartLine:       startLine,
			ScopeContext:    scopeContext,
			FunctionContext: functionContext,
			ClassContext:    classContext,
		},
	}
}

func (parser *Parser) closeScope(endLine int, endChar int) {
	// append all scoped definitions to scopeTable
	parser.symbolMap.scopeRange.EndChar = endChar
	parser.symbolMap.scopeRange.EndLine = endLine
	parser.scopeTable[parser.symbolMap.scopeRange] = parser.symbolMap.definitions
	// rollback to previous scope symbolMap
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
			parser.addError(fmt.Sprintf("%s is already declared in this scope at line %d", name, definition.Line+1), ERROR_RESOLVER)
		}
		parser.symbolMap.currentTable[name] = token
		parser.references[token] = []Token{}
	}
	// add definitions to scope
	parser.symbolMap.definitions = append(parser.symbolMap.definitions, token)
}

func (parser *Parser) Parse(input []Token) ([]Node, []Node, map[Token][]Token, map[ScopeRange][]Token, []CompileError) {
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
	// add global definitions to scope table
	parser.scopeTable[ScopeRange{FunctionContext: GLOBAL_CONTEXT, ClassContext: GLOBAL_CONTEXT, ScopeContext: GLOBAL_CONTEXT}] = parser.symbolMap.definitions
	return program, parser.identifierNodes, parser.references, parser.scopeTable, parser.errorList
}

func (parser *Parser) declaration() Node {
	parser.panicMode = false
	switch {
	case parser.match(VAR):
		return parser.varDeclaration()
	case parser.match(FUN):
		return parser.funcDeclaration(FUNCTION_CONTEXT)
	case parser.match(CLASS):
		return parser.classDeclaration()
	case parser.match(NEWLINE):
		return &NewLine{Token: parser.peekPrevious()}
	case parser.match(COMMENT):
		return &Comment{Comment: parser.peekPrevious()}
	default:
		return parser.statement(GLOBAL_CONTEXT)
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
	brace := parser.peekPrevious()
	parser.raiseScope(brace.Line, brace.Character, CLASS_CONTEXT)
	methods := make([]Node, 0)
	for token := parser.peekParser().TokenType; token != BRACERIGHT && token != EOF; token = parser.peekParser().TokenType {
		if token == NEWLINE {
			parser.match(NEWLINE)
			methods = append(methods, &NewLine{Token: parser.peekPrevious()})
			continue
		}
		if token == COMMENT {
			parser.match(COMMENT)
			methods = append(methods, &Comment{Comment: parser.peekPrevious(), Inline: false})
			continue
		}
		method := parser.funcDeclaration(METHOD_CONTEXT)
		if method == nil {
			continue
		}
		methods = append(methods, method)
	}

	parser.consume(BRACERIGHT, "Expect '}' at end of class declaration")
	brace = parser.peekPrevious()
	parser.closeScope(brace.Line, brace.Character)

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

func (parser *Parser) funcDeclaration(functionContext int) Node {
	identifier := parser.peekParser()
	parser.addDefinition(identifier)
	parser.consume(IDENTIFIER, "Expected identifier for function name")

	parser.consume(PARANLEFT, "Expected ( after function name")
	parameters := make([]Node, 0)
	if !parser.match(PARANRIGHT) {
		parameters = parser.parameters()
	}

	parser.consume(BRACELEFT, "Expected { at start of function body")

	body := parser.block(functionContext)

	return &FuncDecl{Name: identifier, Body: body, Parameters: parameters, FunctionType: functionContext}

}

func (parser *Parser) parameters() []Node {
	parameters := make([]Node, 0)
	consumed := parser.consume(IDENTIFIER, "Expected Parameter Name")
	parameters = append(parameters, &Variable{Identifier: parser.peekPrevious()})
	if consumed {
		parser.addDefinition(parser.peekPrevious())
	}

	for parser.match(COMMA) {
		consumed = parser.consume(IDENTIFIER, "Expected Parameter Name")
		parameters = append(parameters, &Variable{Identifier: parser.peekPrevious()})
		if consumed {
			parser.addDefinition(parser.peekPrevious())
		}
	}

	parser.consume(PARANRIGHT, "Expected ')' before function body")

	return parameters

}

func (parser *Parser) statement(scopeContext int) Node {
	fmt.Printf(">>%s", parser.peekParser().Value)
	switch {
	case parser.match(PRINT):
		expr := parser.expression()
		parser.consume(SEMICOLON, "Expected ; at end of statement")
		return &PrintStmt{Expr: expr}

	case parser.match(RETURN):
		if parser.symbolMap.functionContext == GLOBAL_CONTEXT {
			parser.addError("Unexpected Return statement outside of functions or methods", ERROR_RESOLVER)
		}
		if parser.match(SEMICOLON) {
			return &ReturnStmt{Expr: &Primary{ValType: "nil", Value: nil}, ReturnsValue: false}
		}
		expr := parser.expression()
		parser.consume(SEMICOLON, "Expected ; at end of statement")
		return &ReturnStmt{Expr: expr, ReturnsValue: true}

	case parser.match(BRACELEFT):
		if scopeContext == GLOBAL_CONTEXT {
			return parser.block(BLOCK_CONTEXT)
		}
		return parser.block(scopeContext)

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

func (parser *Parser) block(scopeContext int) Node {
	brace := parser.peekPrevious()
	parser.raiseScope(brace.Line, brace.Character, scopeContext)
	body := make([]Node, 0)
	for token := parser.peekParser(); token.TokenType != EOF && token.TokenType != BRACERIGHT; token = parser.peekParser() {
		body = append(body, parser.declaration())
	}
	parser.consume(BRACERIGHT, "Expected '}' at end of block")
	brace = parser.peekPrevious()
	parser.closeScope(brace.Line, brace.Character)
	return &BlockStmt{Body: body, BlockContext: scopeContext}
}

func (parser *Parser) ifStmt() Node {
	parser.consume(PARANLEFT, "Expected '(' after if")
	condition := parser.expression()
	parser.consume(PARANRIGHT, "Expected ')' after condition")

	thenBranch := parser.statement(IF_CONTEXT)
	var elseBranch Node = nil

	if parser.match(ELSE) {
		elseBranch = parser.statement(IF_CONTEXT)
	}

	return &IfStmt{Condition: condition, Then: thenBranch, Else: elseBranch}
}

func (parser *Parser) whileStmt() Node {
	parser.consume(PARANLEFT, "Expected '(' after while")
	condition := parser.expression()
	parser.consume(PARANRIGHT, "Expected ')' after condition")

	body := parser.statement(WHILE_CONTEXT)

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
		assign = expr
	}

	parser.consume(PARANRIGHT, "Expected ')' before body")
	body := parser.statement(FOR_CONTEXT)

	return &ForStmt{Initializer: initializer, Condition: condition, Assignment: assign, Body: body}
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
			parser.addErrorAt("Invalid assignment target", token.Line, token.Character, ERROR_PARSER)
			return expr
		}
		expr = &Assignment{Identifier: variable, Value: value}
	}
	return expr
}

func (parser *Parser) logicalOr() Node {
	expr := parser.logicalAnd()

	for token := parser.peekParser(); token.TokenType == OR; token = parser.peekParser() {
		parser.advanceParser(true)
		right := parser.logicalAnd()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) logicalAnd() Node {
	expr := parser.equality()

	for token := parser.peekParser(); token.TokenType == AND; token = parser.peekParser() {
		parser.advanceParser(true)
		right := parser.equality()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) equality() Node {
	expr := parser.comparison()

	for token := parser.peekParser(); token.TokenType == EQUALEQUAL || token.TokenType == BANGEQUAL; token = parser.peekParser() {
		parser.advanceParser(true)
		right := parser.comparison()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) comparison() Node {
	expr := parser.term()

	for token := parser.peekParser(); token.TokenType == GREATER || token.TokenType == GREATEREQUAL ||
		token.TokenType == LESS || token.TokenType == LESSEQUAL; token = parser.peekParser() {
		parser.advanceParser(true)
		right := parser.term()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) term() Node {
	expr := parser.factor()

	for token := parser.peekParser(); token.TokenType == PLUS || token.TokenType == MINUS; token = parser.peekParser() {
		parser.advanceParser(true)
		right := parser.factor()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) factor() Node {
	expr := parser.unary()

	for token := parser.peekParser(); token.TokenType == STAR || token.TokenType == SLASH; token = parser.peekParser() {
		parser.advanceParser(true)
		right := parser.unary()
		expr = &Binary{Left: expr, Right: right, Operation: token.TokenType}
	}

	return expr
}

func (parser *Parser) unary() Node {
	if token := parser.peekParser(); token.TokenType == MINUS || token.TokenType == BANG {
		parser.advanceParser(true)
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
		parser.addError("Can't have more than 255 arguments", ERROR_RESOLVER)
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
		if parser.symbolMap.classContext != CLASS_CONTEXT {
			parser.addError("Invalid use of 'this' keyword outside of class context ", ERROR_RESOLVER)
		}
		return &This{Identifier: currToken}
	case parser.match(SUPER):
		if parser.symbolMap.classContext != CLASS_CONTEXT {
			parser.addError("Invalid use of 'super' keyword outside of class context ", ERROR_RESOLVER)
		}
		parser.consume(DOT, "Expected '.' after super")
		parser.consume(IDENTIFIER, "Expected method name for super-class")
		return &Super{Identifier: currToken, Property: parser.peekPrevious()}
	case parser.match(IDENTIFIER):
		name, ok := currToken.Value.(string)
		var definition Token
		if ok {
			if slices.Contains(NativeFunctions, name) {
				return &Variable{Identifier: currToken}
			}
			definition, ok = parser.getDefinition(name)
			result := Variable{Identifier: currToken, Definition: definition}
			if !ok {
				parser.addError(fmt.Sprintf("%s is not defined in current scope", name), ERROR_RESOLVER)
			} else {
				parser.addIdentifier(&result, definition, currToken)
			}
			return &result
		} else {
			return &Variable{Identifier: currToken}
		}
	case parser.match(PARANLEFT):
		expr := parser.expression()
		parser.consume(PARANRIGHT, fmt.Sprintf("Expected ')' at line %d character %d", currToken.Line+1, currToken.Character+1))
		return &Group{Expression: expr}

	case parser.peekParser().TokenType == (EOF):
		parser.addError("Unexpected end of file", ERROR_PARSER)

	default:
		parser.addError(fmt.Sprintf("Unexpected token %d at line %d character %d", currToken.TokenType, currToken.Line+1, currToken.Character+1), ERROR_PARSER)
		parser.advanceParser(true)
	}
	return &Primary{}
}

func (parser *Parser) advanceParser(ignoreNewline bool) {
	if ignoreNewline && parser.peekParser().TokenType == NEWLINE {
		for parser.peekParser().TokenType == NEWLINE {
			parser.currentToken++
		}
		return
	}
	parser.currentToken++
}

func (parser *Parser) match(tokenType int) bool {

	if tokenType == parser.peekParser().TokenType {
		parser.advanceParser(false)
		return true
	}

	// if tokenType == parser.peekParserIgnoreNewline().TokenType {
	// 	parser.advanceParser(true)
	// 	return true
	// }
	return false
}

func (parser *Parser) consume(tokenType int, message string) bool {
	if parser.match(tokenType) {
		return true
	}
	parser.addError(message, ERROR_PARSER)
	return false
}

func (parser *Parser) addError(message string, source int) {
	if parser.panicMode {
		return
	}
	parser.errorList = append(parser.errorList, CompileError{Message: message, Line: parser.peekParser().Line, Char: parser.peekParser().Character, Severity: 1, Source: source})
	parser.panicMode = true
}

func (parser *Parser) addWarning(message string) {
	if parser.panicMode {
		return
	}
	parser.errorList = append(parser.errorList, CompileError{Message: message, Line: parser.peekParser().Line, Char: parser.peekParser().Character, Severity: 2, Source: ERROR_WARNING})
}

func (parser *Parser) addWarningAt(message string, line int, char int) {
	if parser.panicMode {
		return
	}
	parser.errorList = append(parser.errorList, CompileError{Message: message, Line: line, Char: char, Severity: 2, Source: ERROR_WARNING})
}

func (parser *Parser) addErrorAt(message string, line int, char int, source int) {
	if parser.panicMode {
		return
	}
	parser.errorList = append(parser.errorList, CompileError{Message: message, Line: line, Char: char, Severity: 1, Source: source})
	parser.panicMode = true
}

func (parser *Parser) peekPrevious() Token {
	return parser.tokenList[parser.currentToken-1]
}

func (parser *Parser) peekParser() Token {
	return parser.tokenList[parser.currentToken]
}

func (parser *Parser) peekParserIgnoreNewline() Token {
	var peek int
	for peek = parser.currentToken; parser.tokenList[peek].TokenType == NEWLINE; peek++ {
	}
	return parser.tokenList[peek]
}

func (parser *Parser) peekNext() (Token, bool) {
	if len(parser.tokenList) >= parser.currentToken+1 {
		return Token{}, false
	}
	return parser.tokenList[parser.currentToken+1], true
}
