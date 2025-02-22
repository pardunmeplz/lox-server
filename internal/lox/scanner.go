package lox

import (
	"fmt"
	"strconv"
	"unicode"
)

type token struct {
	tokenType int
	line      int
	value     any
	character int
}

var scannerState struct {
	tokens        []token
	lexicalErrors []CompileError
	line          int
	currChar      int
	current       int
	source        *string
}

func initializeScanner(code *string) {
	scannerState.tokens = make([]token, 0)
	scannerState.lexicalErrors = make([]CompileError, 0)
	scannerState.line = 1
	scannerState.currChar = 1
	scannerState.current = 0
	scannerState.source = code

}

func scan(code string) ([]token, []CompileError, error) {
	initializeScanner(&code)

	for len(*scannerState.source) > scannerState.current {
		err := scanToken()
		if err != nil {
			return scannerState.tokens, scannerState.lexicalErrors, err
		}
	}
	scannerState.tokens = append(scannerState.tokens, token{tokenType: EOF, line: scannerState.line, character: scannerState.currChar})

	return scannerState.tokens, scannerState.lexicalErrors, nil
}

var keywords map[string]int = map[string]int{
	"if":    IF,
	"true":  TRUE,
	"false": FALSE,
	"nil":   NIL,
	"else":  ELSE,
	"for":   FOR,
	"while": WHILE,
	"fun":   FUN,
	"class": CLASS,
	"var":   VAR,
	"and":   AND,
	"or":    OR,
	"print": PRINT,
}

func scanNumber(char rune) (bool, error) {
	if !unicode.IsDigit(char) {
		return false, nil
	}
	start := scannerState.current
	for (len(*scannerState.source) > scannerState.current) && unicode.IsDigit(peekScanner()) {
		advanceScanner()
	}

	if !matchScanner('.') {
		value, err := strconv.Atoi((*scannerState.source)[start-1 : scannerState.current])
		if err != nil {
			return true, err
		}
		scannerState.tokens = append(scannerState.tokens, token{tokenType: NUMBER, line: scannerState.line, character: scannerState.currChar, value: value})
		return true, nil
	}

	for (len(*scannerState.source) > scannerState.current) && unicode.IsDigit(peekScanner()) {
		advanceScanner()
	}
	value, err := strconv.ParseFloat((*scannerState.source)[start-1:scannerState.current], 64)
	if err != nil {
		return true, err
	}
	scannerState.tokens = append(scannerState.tokens, token{tokenType: NUMBER, line: scannerState.line, character: scannerState.currChar, value: value})
	return true, nil

}

func scanKeywords(char rune) (bool, error) {
	if !unicode.IsLetter(char) {
		return false, nil
	}

	start := scannerState.current
	for len(*scannerState.source) > scannerState.current && (unicode.IsDigit(peekScanner()) || unicode.IsLetter(peekScanner()) || peekScanner() == '_') {
		advanceScanner()
	}
	value := (*scannerState.source)[start-1 : scannerState.current]

	tokenType, isKeyword := keywords[value]
	if isKeyword {
		scannerState.tokens = append(scannerState.tokens, token{tokenType: tokenType, line: scannerState.line, character: scannerState.currChar})
		return true, nil
	}

	scannerState.tokens = append(scannerState.tokens, token{tokenType: IDENTIFIER, line: scannerState.line, character: scannerState.currChar, value: value})

	return true, nil
}

func scanToken() error {
	char := peekScanner()
	advanceScanner()

	isNum, err := scanNumber(char)
	if err != nil {
		return err
	}
	if isNum {
		return nil
	}

	switch char {
	case '+':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: PLUS, line: scannerState.line, character: scannerState.currChar})
	case '-':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: MINUS, line: scannerState.line, character: scannerState.currChar})
	case '*':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: STAR, line: scannerState.line, character: scannerState.currChar})
	case ';':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: SEMICOLON, line: scannerState.line, character: scannerState.currChar})
	case '}':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: BRACERIGHT, line: scannerState.line, character: scannerState.currChar})
	case '{':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: BRACELEFT, line: scannerState.line, character: scannerState.currChar})
	case '(':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: PARANLEFT, line: scannerState.line, character: scannerState.currChar})
	case ')':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: PARANRIGHT, line: scannerState.line, character: scannerState.currChar})
	case '.':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: DOT, line: scannerState.line, character: scannerState.currChar})
	case ',':
		scannerState.tokens = append(scannerState.tokens, token{tokenType: COMMA, line: scannerState.line, character: scannerState.currChar})
	case ' ':
	case '\t':
	case '\r':
	case '\n':
		scannerState.line++
		scannerState.currChar = 0
	case '/':
		if matchScanner('/') {
			for peekScannerNext() != '\n' && len(*scannerState.source) > scannerState.current {
				advanceScanner()
			}
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, token{tokenType: SLASH, line: scannerState.line, character: scannerState.currChar})
	case '=':
		if matchScanner('=') {
			scannerState.tokens = append(scannerState.tokens, token{tokenType: EQUALEQUAL, line: scannerState.line, character: scannerState.currChar})
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, token{tokenType: EQUAL, line: scannerState.line, character: scannerState.currChar})
	case '!':
		if matchScanner('=') {
			scannerState.tokens = append(scannerState.tokens, token{tokenType: BANGEQUAL, line: scannerState.line, character: scannerState.currChar})
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, token{tokenType: BANG, line: scannerState.line, character: scannerState.currChar})
	case '<':
		if matchScanner('=') {
			scannerState.tokens = append(scannerState.tokens, token{tokenType: LESSEQUAL, line: scannerState.line, character: scannerState.currChar})
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, token{tokenType: LESS, line: scannerState.line, character: scannerState.currChar})
	case '>':
		if matchScanner('=') {
			scannerState.tokens = append(scannerState.tokens, token{tokenType: GREATEREQUAL, line: scannerState.line, character: scannerState.currChar})
			return nil
		}
		scannerState.tokens = append(scannerState.tokens, token{tokenType: GREATER, line: scannerState.line, character: scannerState.currChar})
	case '"':
		start := scannerState.current
		for len(*scannerState.source)-1 > scannerState.current && peekScanner() != '"' {
			if peekScanner() == '\n' {
				scannerState.line++
				scannerState.currChar = 0
			}
			advanceScanner()
		}
		consumeScanner('"', fmt.Sprintf("Expected \" at end of string at line %d column %d", scannerState.line, scannerState.currChar))
		scannerState.tokens = append(scannerState.tokens, token{tokenType: STRING, line: scannerState.line, character: scannerState.currChar, value: (*scannerState.source)[start : scannerState.current-1]})

	default:
		isKeyword, err := scanKeywords(char)
		if isKeyword {
			if err != nil {
				return err
			}
			return nil
		}
		scannerState.lexicalErrors = append(scannerState.lexicalErrors, CompileError{Line: scannerState.line, Char: scannerState.currChar, Message: fmt.Sprintf("Unexpected token %c at line %d column %d", char, scannerState.line, scannerState.currChar)})
		return nil
	}
	return nil

}

func advanceScanner() {
	scannerState.currChar++
	scannerState.current++
}

func peekScanner() rune {
	return rune((*scannerState.source)[scannerState.current])
}

func peekScannerNext() rune {
	return rune((*scannerState.source)[scannerState.current+1])
}

func matchScanner(char rune) bool {
	if (*scannerState.source)[scannerState.current] == byte(char) {
		advanceScanner()
		return true
	}
	return false
}

func consumeScanner(char rune, err string) {
	if matchScanner(char) {
		return
	}
	scannerState.lexicalErrors = append(scannerState.lexicalErrors, CompileError{Line: scannerState.line, Char: scannerState.currChar, Message: fmt.Sprintf("%s", err)})
}
