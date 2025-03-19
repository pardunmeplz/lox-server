package lox

type Literal struct {
	value   any
	valType string
}

type CompileError struct {
	Message  string
	Line     int
	Char     int
	Severity int
	Source   int
}

const (
	ERROR_SCANNER = iota
	ERROR_PARSER
	ERROR_RESOLVER
	ERROR_WARNING
	ERROR_NONE
)
