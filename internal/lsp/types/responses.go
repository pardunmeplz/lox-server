package lsp

type Position struct {
	line      uint
	character uint
}

type Range struct {
	start Position
	end   Position
}

type Diagnostic struct {
	/*
	   1 = Error
	   2 = Warning
	   3 = Info
	   4 = Hint
	*/
	severity int
	errRange Range `json:"range"`
	message  string
}

type PublishDiagnosticParams struct {
	uri         string
	version     int
	diagnostics Diagnostic
}
