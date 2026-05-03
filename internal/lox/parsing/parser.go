package parsing

import (
	lox_dg "lox-server/internal/lox/diagnostics"
	lox_sc "lox-server/internal/lox/lexical_analysis"
)

/*
   program        → declaration* EOF ;

   declaration    → varDecl | statement | funcDecl | classDecl ;

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

type ParserState struct {
	// token list to parse
	Tokens []lox_sc.Token
	// current token
	Index int
	// error list
	Errors []lox_dg.CompileError
}

func (parserState *ParserState) Parse() []Node {
	var ast = make([]Node, 0)
	for !peekMatch(parserState, lox_sc.EOF) {
		ast = append(ast, declaration(parserState))
	}
	return ast
}

func declaration(parserState *ParserState) Node {

}

func expression(parserState *ParserState) Node {

}

func primary(parserState *ParserState) Node {
	currToken := advance(parserState)

	switch currToken.TokenType {
	case lox_sc.NUMBER, lox_sc.STRING:
		return &Primary{Value: currToken.Value, Token: currToken}
	case lox_sc.TRUE:
		return &Primary{Value: true, Token: currToken}
	case lox_sc.FALSE:
		return &Primary{Value: false, Token: currToken}
	case lox_sc.NIL:
		return &Primary{Value: nil, Token: currToken}
	case lox_sc.PARANLEFT:
		expr := expression(parserState)
		closeToken := consume(parserState, lox_sc.PARANRIGHT, "Missing closing Paranthseis ')'")
		return &Group{OpenParan: currToken, CloseParan: closeToken, Expression: expr}
	case lox_sc.SUPER:
		superExpr := &Super{Super: currToken}
		if peekMatch(parserState, lox_sc.DOT) {
			advance(parserState)
			superExpr.Identifier = consume(parserState, lox_sc.IDENTIFIER, "Expected Identifier after 'Super.'")
		}
		addError(parserState, "Unexpected Token", 1, currToken.Line, currToken.Character)
		return superExpr
	}

	return nil
}

// helper functions

func peekMatch(parserState *ParserState, tokenTp int) bool {
	return peek(parserState).TokenType == tokenTp
}

func peek(parserState *ParserState) *lox_sc.Token {
	return &parserState.Tokens[parserState.Index]
}

func consume(parserState *ParserState, tokenTp int, message string) *lox_sc.Token {
	if !peekMatch(parserState, tokenTp) {
		token := peek(parserState)
		addError(parserState, message, 1, token.Line, token.Character)
		return token
	} else {
		advance(parserState)
		return nil
	}
}

func advance(parserState *ParserState) *lox_sc.Token {
	parserState.Index += 1
	return &parserState.Tokens[parserState.Index-1]
}

func addError(parserState *ParserState, message string, severity int, line int, char int) {
	compileErr := lox_dg.CompileError{Message: message, Line: line, Char: char, Severity: severity, Source: lox_dg.ERROR_PARSER}
	parserState.Errors = append(parserState.Errors, compileErr)
}
