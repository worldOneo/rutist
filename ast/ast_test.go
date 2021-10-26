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
						String{"/test/"},
						Expression{"yeet", []Node{
							String{"me"},
							String{"out"}},
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
			test = 100.5
			`)},
			Block{
				Body: []Node{
					Assignment{"test", Int{100}},
					Assignment{"test", Float{100.5}},
				},
			},
			false,
		},
		{
			"assignment bool",
			args{lexed: tokens.Lexerp(`
			test = true
			test = false
			`)},
			Block{
				Body: []Node{
					Assignment{"test", Bool{true}},
					Assignment{"test", Bool{false}},
				},
			},
			false,
		},
		{
			"scoper",
			args{lexed: tokens.Lexerp(`
			try(@{
				print(test)
			}, @{

			})
			`)},
			Block{
				Body: []Node{
					Expression{"try", []Node{
						Scope{Block{[]Node{
							Expression{"print", []Node{Variable{"test"}}},
						}}},
						Scope{Block{[]Node{}}},
					}},
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
