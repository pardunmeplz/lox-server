package lox

import (
	"fmt"
	"strconv"
	"unicode"
)

type Token struct {
	TokenType int
	Line      int
	Value     any
	Character int
}

type Scanner struct {
	tokens        []Token
	lexicalErrors []CompileError
	line          int
	currChar      int
	current       int
	source        *string
}

func (scannerState *Scanner) initializeScanner(code *string) {
	scannerState.tokens = make([]Token, 0)
	scannerState.lexicalErrors = make([]CompileError, 0)
	scannerState.line = 0
	scannerState.currChar = 0
	scannerState.current = 0
	scannerState.source = code

}

func (scannerState *Scanner) Scan(code string) ([]Token, []CompileError, error) {
	scannerState.initializeScanner(&code)

	for len(*scannerState.source) > scannerState.current {
		err := scannerState.scanToken()
		if err != nil {
			return scannerState.tokens, scannerState.lexicalErrors, err
		}
	}
	scannerState.tokens = append(scannerState.tokens, Token{TokenType: EOF, Line: scannerState.line, Character: scannerState.currChar})

	return scannerState.tokens, scannerState.lexicalErrors, nil
}

var keywords map[string]int = map[string]int{
	"if":     IF,
	"true":   TRUE,
	"false":  FALSE,
	"nil":    NIL,
	"else":   ELSE,
	"for":    FOR,
	"while":  WHILE,
	"fun":    FUN,
	"class":  CLASS,
	"var":    VAR,
	"and":    AND,
	"or":     OR,
	"print":  PRINT,
	"this":   THIS,
	"super":  SUPER,
	"return": RETURN,
}

func (scannerState *Scanner) scanNumber(char rune) (bool, error) {
	if !unicode.IsDigit(char) {
		return false, nil
	}
	start := scannerState.current
	for (len(*scannerState.source) > scannerState.current) && unicode.IsDigit(scannerState.peekScanner()) {
		scannerState.advanceScanner()
	}

	if !scannerState.matchScanner('.') {
		value, err := strconv.Atoi((*scannerState.source)[start-1 : scannerState.current])
		if err != nil {
			return true, err
		}
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: NUMBER, Line: scannerState.line, Character: scannerState.currChar, Value: value})
		return true, nil
	}

	for (len(*scannerState.source) > scannerState.current) && unicode.IsDigit(scannerState.peekScanner()) {
		scannerState.advanceScanner()
	}
	value, err := strconv.ParseFloat((*scannerState.source)[start-1:scannerState.current], 64)
	if err != nil {
		return true, err
	}
	scannerState.tokens = append(scannerState.tokens, Token{TokenType: NUMBER, Line: scannerState.line, Character: scannerState.currChar, Value: value})
	return true, nil

}

func (scannerState *Scanner) scanKeywords(char rune) (bool, error) {
	if !unicode.IsLetter(char) {
		return false, nil
	}

	start := scannerState.current
	startChar := scannerState.currChar - 1
	for len(*scannerState.source) > scannerState.current && (unicode.IsDigit(scannerState.peekScanner()) || unicode.IsLetter(scannerState.peekScanner()) || scannerState.peekScanner() == '_') {
		scannerState.advanceScanner()
	}
	value := (*scannerState.source)[start-1 : scannerState.current]

	tokenType, isKeyword := keywords[value]
	if isKeyword {
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: tokenType, Line: scannerState.line, Character: startChar})
		return true, nil
	}

	scannerState.tokens = append(scannerState.tokens, Token{TokenType: IDENTIFIER, Line: scannerState.line, Character: startChar, Value: value})

	return true, nil
}

func (scannerState *Scanner) scanToken() error {
	char := scannerState.peekScanner()
	scannerState.advanceScanner()

	isNum, err := scannerState.scanNumber(char)
	if err != nil {
		return err
	}
	if isNum {
		return nil
	}

	switch char {
	case '+':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: PLUS, Line: scannerState.line, Character: scannerState.currChar})
	case '-':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: MINUS, Line: scannerState.line, Character: scannerState.currChar})
	case '*':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: STAR, Line: scannerState.line, Character: scannerState.currChar})
	case ';':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: SEMICOLON, Line: scannerState.line, Character: scannerState.currChar})
	case '}':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: BRACERIGHT, Line: scannerState.line, Character: scannerState.currChar})
	case '{':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: BRACELEFT, Line: scannerState.line, Character: scannerState.currChar})
	case '(':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: PARANLEFT, Line: scannerState.line, Character: scannerState.currChar})
	case ')':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: PARANRIGHT, Line: scannerState.line, Character: scannerState.currChar})
	case '.':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: DOT, Line: scannerState.line, Character: scannerState.currChar})
	case ',':
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: COMMA, Line: scannerState.line, Character: scannerState.currChar})
	case ' ':
	case '\t':
	case '\r':
	case '\n':
		scannerState.line++
		scannerState.currChar = 0
	case '/':
		if scannerState.matchScanner('/') {
			for scannerState.peekScannerNext() != '\n' && len(*scannerState.source) > scannerState.current {
				scannerState.advanceScanner()
			}
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: SLASH, Line: scannerState.line, Character: scannerState.currChar})
	case '=':
		if scannerState.matchScanner('=') {
			scannerState.tokens = append(scannerState.tokens, Token{TokenType: EQUALEQUAL, Line: scannerState.line, Character: scannerState.currChar})
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: EQUAL, Line: scannerState.line, Character: scannerState.currChar})
	case '!':
		if scannerState.matchScanner('=') {
			scannerState.tokens = append(scannerState.tokens, Token{TokenType: BANGEQUAL, Line: scannerState.line, Character: scannerState.currChar})
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: BANG, Line: scannerState.line, Character: scannerState.currChar})
	case '<':
		if scannerState.matchScanner('=') {
			scannerState.tokens = append(scannerState.tokens, Token{TokenType: LESSEQUAL, Line: scannerState.line, Character: scannerState.currChar})
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: LESS, Line: scannerState.line, Character: scannerState.currChar})
	case '>':
		if scannerState.matchScanner('=') {
			scannerState.tokens = append(scannerState.tokens, Token{TokenType: GREATEREQUAL, Line: scannerState.line, Character: scannerState.currChar})
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: GREATER, Line: scannerState.line, Character: scannerState.currChar})
	case '"':
		start := scannerState.current
		for len(*scannerState.source)-1 > scannerState.current && scannerState.peekScanner() != '"' {
			if scannerState.peekScanner() == '\n' {
				scannerState.line++
				scannerState.currChar = 0
			}
			scannerState.advanceScanner()
		}
		scannerState.consumeScanner('"', fmt.Sprintf("Expected \" at end of string at line %d column %d", scannerState.line, scannerState.currChar))
		scannerState.tokens = append(scannerState.tokens, Token{TokenType: STRING, Line: scannerState.line, Character: scannerState.currChar, Value: (*scannerState.source)[start : scannerState.current-1]})

	default:
		isKeyword, err := scannerState.scanKeywords(char)
		if isKeyword {
			if err != nil {
				return err
			}
			return nil
		}
		scannerState.lexicalErrors = append(scannerState.lexicalErrors, CompileError{Line: scannerState.line, Char: scannerState.currChar, Message: fmt.Sprintf("Unexpected token %c at line %d column %d", char, scannerState.line+1, scannerState.currChar+1), Severity: 1})
		return nil
	}
	return nil

}

func (scannerState *Scanner) advanceScanner() {
	scannerState.currChar++
	scannerState.current++
}

func (scannerState *Scanner) peekScanner() rune {
	return rune((*scannerState.source)[scannerState.current])
}

func (scannerState *Scanner) peekScannerNext() rune {
	return rune((*scannerState.source)[scannerState.current+1])
}

func (scannerState *Scanner) matchScanner(char rune) bool {
	if (*scannerState.source)[scannerState.current] == byte(char) {
		scannerState.advanceScanner()
		return true
	}
	return false
}

func (scannerState *Scanner) consumeScanner(char rune, err string) {
	if scannerState.matchScanner(char) {
		return
	}
	scannerState.lexicalErrors = append(scannerState.lexicalErrors, CompileError{Line: scannerState.line + 1, Char: scannerState.currChar + 1, Message: fmt.Sprintf("%s", err), Severity: 1})
}
