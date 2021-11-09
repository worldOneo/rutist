package tokens

import (
	"reflect"
	"testing"
)

func TestCodeLexer_Lexer(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		want    []Token
		wantErr bool
	}{
		{
			"t1",
			`scope {
				print("test")
			}`,
			[]Token{
				identifierToken("scope", 0), scopeOpenToken(0),
				identifierToken("print", 1), {ParenOpen, "(", 0, 0, 1}, stringToken("test", 1), {ParenClosed, ")", 0, 0, 1},
				scopeClosedToken(2),
			},
			false,
		},
		{
			"single arg int",
			`@{
				print(1)
			}`,
			[]Token{
				{Scoper, "@", 0, 0, 0}, scopeOpenToken(0),
				identifierToken("print", 1), {ParenOpen, "(", 0, 0, 1}, intToken("1", 1, 1), {ParenClosed, ")", 0, 0, 1},
				scopeClosedToken(2),
			},
			false,
		},
		{
			"special char string",
			`
			print("Test\\\n")
			`,
			[]Token{
				identifierToken("print", 1), {ParenOpen, "(", 0, 0, 1}, stringToken("Test\\\n", 1), {ParenClosed, ")", 0, 0, 1},
			},
			false,
		},
		{
			"dot member separator",
			`var.member("call")`,
			[]Token{
				identifierToken("var", 0), {Dot, ".", 0, 0, 0}, identifierToken("member", 0), {ParenOpen, "(", 0, 0, 0}, stringToken("call", 0), {ParenClosed, ")", 0, 0, 0},
			},
			false,
		},
		{
			"operators",
			`a<b c+d c>=d ~a a||b`,
			[]Token{
				identifierToken("a", 0), {OperatorType, "<", OperatorLt, 0, 0}, identifierToken("b", 0),
				identifierToken("c", 0), {OperatorType, "+", OperatorAdd, 0, 0}, identifierToken("d", 0),
				identifierToken("c", 0), {OperatorType, ">=", OperatorGe, 0, 0}, identifierToken("d", 0),
				{OperatorType, "~", OperatorNot, 0, 0}, identifierToken("a", 0),
				identifierToken("a", 0), {OperatorType, "||", OperatorLor, 0, 0}, identifierToken("b", 0),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Lexer(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("WordParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WordParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
