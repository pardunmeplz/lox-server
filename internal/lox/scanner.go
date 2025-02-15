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

func scanNumber(char rune) (bool, error) {
	if !unicode.IsDigit(char) {
		return false, nil
	}
	start := current
	for unicode.IsDigit(peek()) {
		advance()
	}

	if !match('.') {
		value, err := strconv.Atoi((*source)[start-1 : current])
		if err != nil {
			return true, err
		}
		tokens = append(tokens, token{tokenType: NUMBER, line: line, value: value})
		return true, nil
	}

	for unicode.IsDigit(peek()) {
		advance()
	}
	value, err := strconv.ParseFloat((*source)[start-1:current], 64)
	if err != nil {
		return true, err
	}
	tokens = append(tokens, token{tokenType: NUMBER, line: line, value: value})
	return true, nil

}

func scanToken() error {
	char := peek()
	advance()

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
		if match('/') {
			for peekNext() != '\n' && len(*source) > current {
				advance()
			}
			return nil
		}
		tokens = append(tokens, token{tokenType: SLASH, line: line})
	case '=':
		if match('=') {
			tokens = append(tokens, token{tokenType: EQUALEQUAL, line: line})
			return nil
		}
		tokens = append(tokens, token{tokenType: EQUAL, line: line})
	case '!':
		if match('=') {
			tokens = append(tokens, token{tokenType: BANGEQUAL, line: line})
			return nil
		}
		tokens = append(tokens, token{tokenType: BANG, line: line})
	case '<':
		if match('=') {
			tokens = append(tokens, token{tokenType: LESSEQUAL, line: line})
			return nil
		}
		tokens = append(tokens, token{tokenType: LESS, line: line})
	case '>':
		if match('=') {
			tokens = append(tokens, token{tokenType: GREATEREQUAL, line: line})
			return nil
		}
		tokens = append(tokens, token{tokenType: GREATER, line: line})
	case '"':
		start := current
		for len(*source)-1 > current && peek() != '"' {
			if peek() == '\n' {
				line++
			}
			advance()
		}
		err := consume('"', fmt.Sprintf("Missing End of string at line %d", line))
		tokens = append(tokens, token{tokenType: STRING, line: line, value: (*source)[start : current-1]})
		if err != nil {
			return err
		}

	default:
	}
	return nil

}

func advance() {
	current += 1
}

func peek() rune {
	return rune((*source)[current])
}

func peekNext() rune {
	return rune((*source)[current+1])
}

func match(char rune) bool {
	if (*source)[current] == byte(char) {
		advance()
		return true
	}
	return false
}

func consume(char rune, err string) error {
	if match(char) {
		return nil
	}
	return fmt.Errorf("%s", err)
}
