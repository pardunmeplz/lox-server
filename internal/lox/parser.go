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

var tokenList []token
var currentToken int

func parse(input []token) {
	tokenList = input
}

func expression() {

}

func primary() {
	currToken := peekParser()

	switch currToken.tokenType {
	case STRING:
		return
	case NUMBER:

	}

}

func advanceParser() {
	currentToken++
}

func peekParser() token {
	return tokenList[currentToken]
}
