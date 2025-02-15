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
}

var tokens []token
var line int = 0
var current int = 0
var source *string

func scan(code string) ([]token, error) {
	source = &code

	for len(*source) > current {
		err := scanToken()
		if err != nil {
			return tokens, err
		}
	}
	tokens = append(tokens, token{tokenType: EOF, line: line})

	return tokens, nil
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
	start := current
	for (len(*source) > current) && unicode.IsDigit(peekScanner()) {
		advanceScanner()
	}

	if !matchScanner('.') {
		value, err := strconv.Atoi((*source)[start-1 : current])
		if err != nil {
			return true, err
		}
		tokens = append(tokens, token{tokenType: NUMBER, line: line, value: value})
		return true, nil
	}

	for (len(*source) > current) && unicode.IsDigit(peekScanner()) {
		advanceScanner()
	}
	value, err := strconv.ParseFloat((*source)[start-1:current], 64)
	if err != nil {
		return true, err
	}
	tokens = append(tokens, token{tokenType: NUMBER, line: line, value: value})
	return true, nil

}

func scanKeywords(char rune) (bool, error) {
	if !unicode.IsLetter(char) {
		return false, nil
	}

	start := current
	for len(*source) > current && (unicode.IsDigit(peekScanner()) || unicode.IsLetter(peekScanner()) || peekScanner() == '_') {
		advanceScanner()
	}
	value := (*source)[start-1 : current]

	tokenType, isKeyword := keywords[value]
	if isKeyword {
		tokens = append(tokens, token{tokenType: tokenType, line: line})
		return true, nil
	}

	tokens = append(tokens, token{tokenType: IDENTIFIER, line: line, value: value})

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
		tokens = append(tokens, token{tokenType: PLUS, line: line})
	case '-':
		tokens = append(tokens, token{tokenType: MINUS, line: line})
	case '*':
		tokens = append(tokens, token{tokenType: STAR, line: line})
	case ';':
		tokens = append(tokens, token{tokenType: SEMICOLON, line: line})
	case '}':
		tokens = append(tokens, token{tokenType: BRACERIGHT, line: line})
	case '{':
		tokens = append(tokens, token{tokenType: BRACELEFT, line: line})
	case '(':
		tokens = append(tokens, token{tokenType: PARANLEFT, line: line})
	case ')':
		tokens = append(tokens, token{tokenType: PARANRIGHT, line: line})
	case '.':
		tokens = append(tokens, token{tokenType: DOT, line: line})
	case ',':
		tokens = append(tokens, token{tokenType: COMMA, line: line})
	case ' ':
	case '\t':
	case '\n':
		line++
	case '/':
		if matchScanner('/') {
			for peekScannerNext() != '\n' && len(*source) > current {
				advanceScanner()
			}
			return nil
		}
		tokens = append(tokens, token{tokenType: SLASH, line: line})
	case '=':
		if matchScanner('=') {
			tokens = append(tokens, token{tokenType: EQUALEQUAL, line: line})
			return nil
		}
		tokens = append(tokens, token{tokenType: EQUAL, line: line})
	case '!':
		if matchScanner('=') {
			tokens = append(tokens, token{tokenType: BANGEQUAL, line: line})
			return nil
		}
		tokens = append(tokens, token{tokenType: BANG, line: line})
	case '<':
		if matchScanner('=') {
			tokens = append(tokens, token{tokenType: LESSEQUAL, line: line})
			return nil
		}
		tokens = append(tokens, token{tokenType: LESS, line: line})
	case '>':
		if matchScanner('=') {
			tokens = append(tokens, token{tokenType: GREATEREQUAL, line: line})
			return nil
		}
		tokens = append(tokens, token{tokenType: GREATER, line: line})
	case '"':
		start := current
		for len(*source)-1 > current && peekScanner() != '"' {
			if peekScanner() == '\n' {
				line++
			}
			advanceScanner()
		}
		err := consumeScanner('"', fmt.Sprintf("Missing End of string at line %d", line))
		tokens = append(tokens, token{tokenType: STRING, line: line, value: (*source)[start : current-1]})
		if err != nil {
			return err
		}

	default:
		isKeyword, err := scanKeywords(char)
		if isKeyword {
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("Unexpected token %c at line %d", char, line)
	}
	return nil

}

func advanceScanner() {
	current += 1
}

func peekScanner() rune {
	return rune((*source)[current])
}

func peekScannerNext() rune {
	return rune((*source)[current+1])
}

func matchScanner(char rune) bool {
	if (*source)[current] == byte(char) {
		advanceScanner()
		return true
	}
	return false
}

func consumeScanner(char rune, err string) error {
	if matchScanner(char) {
		return nil
	}
	return fmt.Errorf("%s", err)
}
