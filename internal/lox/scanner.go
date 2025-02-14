package lox

type token struct {
	tokenType int
	line      int
	value     any
}

var tokens []token
var line int = 0
var current int = 0
var source *string

func scan(code string) []token {
	source = &code

	for len(*source) > current {
		scanToken()
	}
	tokens = append(tokens, token{tokenType: EOF, line: line})

	return tokens
}

func scanToken() {
	char := (*source)[current]
	advance()
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
		line += 1
	case '/':
		if match('/') {
			for peek() != '\n' && len(*source) > current {
				advance()
			}
			return
		}
		tokens = append(tokens, token{tokenType: SLASH, line: line})
	case '=':
		if match('=') {
			tokens = append(tokens, token{tokenType: EQUALEQUAL, line: line})
			return
		}
		tokens = append(tokens, token{tokenType: EQUAL, line: line})
	case '!':
		if match('=') {
			tokens = append(tokens, token{tokenType: BANGEQUAL, line: line})
			return
		}
		tokens = append(tokens, token{tokenType: BANG, line: line})
	case '<':
		if match('=') {
			tokens = append(tokens, token{tokenType: LESSEQUAL, line: line})
			return
		}
		tokens = append(tokens, token{tokenType: LESS, line: line})
	case '>':
		if match('=') {
			tokens = append(tokens, token{tokenType: GREATEREQUAL, line: line})
			return
		}
		tokens = append(tokens, token{tokenType: GREATER, line: line})

	}

}

func advance() {
	current += 1
}

func peek() rune {
	return rune((*source)[current+1])
}

func match(char rune) bool {
	if (*source)[current] == byte(char) {
		advance()
		return true
	}
	return false
}
