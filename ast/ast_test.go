package ast

import (
	"encoding/json"
	"os"
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
					Expression{Identifier{"route"}, []Node{
						String{"/test/"},
						Expression{Identifier{"yeet"}, []Node{
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
					Assignment{Identifier{"test"}, Int{100}},
					Assignment{Identifier{"test"}, Float{100.5}},
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
					Assignment{Identifier{"test"}, Bool{true}},
					Assignment{Identifier{"test"}, Bool{false}},
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
					Expression{Identifier{"try"}, []Node{
						Scope{Block{[]Node{
							Expression{Identifier{"print"}, []Node{Identifier{"test"}}},
						}}},
						Scope{Block{[]Node{}}},
					}},
				},
			},
			false,
		},
		{
			"return assign",
			args{tokens.Lexerp(`
			err = try(@{
				print(1)
			})`)},
			Block{
				[]Node{
					Assignment{
						Identifier{"err"},
						Expression{Identifier{"try"}, []Node{Scope{Block{[]Node{
							Expression{Identifier{"print"}, []Node{Int{1}}},
						}}}}},
					},
				},
			},
			false,
		},
		{
			"member call",
			args{tokens.Lexerp(`
			var.member()
			`)},
			Block{
				[]Node{
					Expression{MemberSelector{Identifier{"var"}, Identifier{"member"}}, []Node{}},
				},
			},
			false,
		},
		{
			"result call",
			args{tokens.Lexerp(`
			l = str(varString).len()
			`)},
			Block{
				[]Node{
					Assignment{
						Identifier{"l"},
						MemberSelector{
							Expression{Identifier{"str"}, []Node{Identifier{"varString"}}},
							Expression{Identifier{"len"}, []Node{}},
						},
					},
				},
			},
			false,
		},
		{
			"function definition",
			args{tokens.Lexerp(`
				handle=@(err){
					print("Err: %s", err)
				}
			`)},
			Block{
				[]Node{
					Assignment{
						Identifier{"handle"},
						FunctionDefinition{Block{
							[]Node{
								Expression{Identifier{"print"}, []Node{String{"Err: %s"}, Identifier{"err"}}},
							},
						}, []Identifier{{"err"}}},
					},
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
				json.NewEncoder(os.Stdout).Encode(got)
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
