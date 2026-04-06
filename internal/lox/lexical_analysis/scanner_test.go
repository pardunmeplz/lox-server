package lexical_analysis

import (
	"reflect"
	"strings"
	"testing"
)

func TestScanner_Blanket_PositiveProgram(t *testing.T) {
	var s Scanner
	code := "var a = 1+2;\nprint a;"

	tokens, diags, err := s.Scan(code)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %d: %+v", len(diags), diags)
	}

	got := make([]int, 0, len(tokens))
	for _, tok := range tokens {
		got = append(got, tok.TokenType)
	}

	want := []int{
		VAR, IDENTIFIER, EQUAL, NUMBER, PLUS, NUMBER, SEMICOLON,
		NEWLINE,
		PRINT, IDENTIFIER, SEMICOLON,
		EOF,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("token types mismatch\n got=%v\nwant=%v", got, want)
	}

	// A couple value sanity checks (not exhaustive).
	if tokens[1].Value != "a" {
		t.Fatalf("expected identifier value %q, got %#v", "a", tokens[1].Value)
	}
	if tokens[3].Value != 1 {
		t.Fatalf("expected number value %v, got %#v", 1, tokens[3].Value)
	}
}

func TestScanner_Blanket_DiagnosticsProgram(t *testing.T) {
	var s Scanner
	code := "print @;\n\"unterminated"

	tokens, diags, err := s.Scan(code)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(diags) != 2 {
		t.Fatalf("expected 2 diagnostics, got %d: %+v", len(diags), diags)
	}
	if !strings.Contains(diags[0].Message, "Unexpected character") && !strings.Contains(diags[1].Message, "Unexpected character") {
		t.Fatalf("expected an Unexpected character diagnostic, got: %+v", diags)
	}
	if !strings.Contains(diags[0].Message, "Expected \"") && !strings.Contains(diags[1].Message, "Expected \"") {
		t.Fatalf("expected an unterminated string diagnostic, got: %+v", diags)
	}

	got := make([]int, 0, len(tokens))
	for _, tok := range tokens {
		got = append(got, tok.TokenType)
	}

	// We mainly care that tokenization continues, produces a STRING, and ends in EOF.
	want := []int{
		PRINT, // print
		SEMICOLON,
		NEWLINE,
		STRING,
		EOF,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("token types mismatch\n got=%v\nwant=%v", got, want)
	}
}

func TestScanner_Blanket_EdgeCasesChunk(t *testing.T) {
	var s Scanner
	code := "class Foo { fun init(a,b){ this.a=a; super.init(); return nil; } } " +
		"var x=true and false or !false; " +
		"if(x!=false){print 1;}else{print 2;} " +
		"for(var i=0;i<3;i=i+1){x=x-1*2/3;} " +
		"while(x>=0){x=x-1;} " +
		"// cmt\n" +
		"x<=10==10;"

	tokens, diags, err := s.Scan(code)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %d: %+v", len(diags), diags)
	}

	got := make([]int, 0, len(tokens))
	for _, tok := range tokens {
		got = append(got, tok.TokenType)
	}

	want := []int{
		CLASS, IDENTIFIER, BRACELEFT,
		FUN, IDENTIFIER, PARANLEFT, IDENTIFIER, COMMA, IDENTIFIER, PARANRIGHT, BRACELEFT,
		THIS, DOT, IDENTIFIER, EQUAL, IDENTIFIER, SEMICOLON,
		SUPER, DOT, IDENTIFIER, PARANLEFT, PARANRIGHT, SEMICOLON,
		RETURN, NIL, SEMICOLON,
		BRACERIGHT, BRACERIGHT,

		VAR, IDENTIFIER, EQUAL, TRUE, AND, FALSE, OR, BANG, FALSE, SEMICOLON,

		IF, PARANLEFT, IDENTIFIER, BANGEQUAL, FALSE, PARANRIGHT, BRACELEFT,
		PRINT, NUMBER, SEMICOLON,
		BRACERIGHT,
		ELSE, BRACELEFT,
		PRINT, NUMBER, SEMICOLON,
		BRACERIGHT,

		FOR, PARANLEFT,
		VAR, IDENTIFIER, EQUAL, NUMBER, SEMICOLON,
		IDENTIFIER, LESS, NUMBER, SEMICOLON,
		IDENTIFIER, EQUAL, IDENTIFIER, PLUS, NUMBER,
		PARANRIGHT, BRACELEFT,
		IDENTIFIER, EQUAL, IDENTIFIER, MINUS, NUMBER, STAR, NUMBER, SLASH, NUMBER, SEMICOLON,
		BRACERIGHT,

		WHILE, PARANLEFT, IDENTIFIER, GREATEREQUAL, NUMBER, PARANRIGHT, BRACELEFT,
		IDENTIFIER, EQUAL, IDENTIFIER, MINUS, NUMBER, SEMICOLON,
		BRACERIGHT,

		COMMENT, NEWLINE,
		IDENTIFIER, LESSEQUAL, NUMBER, EQUALEQUAL, NUMBER, SEMICOLON,
		EOF,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("token types mismatch\n got=%v\nwant=%v", got, want)
	}
}
