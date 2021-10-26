package ast

import (
	"reflect"
	"testing"

	"github.com/worldOneo/rutist/tokens"
)

func TestParse(t *testing.T) {
	type args struct {
		lexed []tokens.Token
	}
	tests := []struct {
		name    string
		args    args
		want    Node
		wantErr bool
	}{
		{
			"t1",
			args{lexed: tokens.Lexerp(`
			route("/test/", yeet("me", "out"))
			`)},
			Block{
				Body: []Node{
					Expression{"route", []Node{
						ValueString{"/test/"},
						Expression{"yeet", []Node{
							ValueString{"me"},
							ValueString{"out"}},
						},
					},
					},
				},
			},
			false,
		},
		{
			"assignment",
			args{lexed: tokens.Lexerp(`
			test = 100
			`)},
			Block{
				Body: []Node{
					Assignment{"test", ValueInt{100}},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.lexed)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
