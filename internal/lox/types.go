package lox

type Literal struct {
	value   any
	valType string
}

type CompileError struct {
	Message string
	Line    int
	Char    int
}
