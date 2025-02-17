package lsp

type Position struct {
	Line      uint `json:"line"`
	Character uint `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}
