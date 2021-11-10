package ast

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/worldOneo/rutist/tokens"
)

var meta = &Meta{tokens.Token{}}

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
					Expression{Identifier{"route", meta}, []Node{
						String{"/test/", meta},
						Expression{Identifier{"yeet", meta}, []Node{
							String{"me", meta},
							String{"out", meta},
						},
							meta,
						},
					},
						meta,
					},
				},
				Meta: meta,
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
					Assignment{Identifier{"test", meta}, Int{100, meta}, meta},
					Assignment{Identifier{"test", meta}, Float{100.5, meta}, meta},
				},
				Meta: meta,
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
					Assignment{Identifier{"test", meta}, Bool{true, meta}, meta},
					Assignment{Identifier{"test", meta}, Bool{false, meta}, meta},
				},
				Meta: meta,
			},
			false,
		},
		{
			"scoper",
			args{lexed: tokens.Lexerp(`
			try({
				print(test)
			}, {

			})
			`)},
			Block{
				Body: []Node{
					Expression{Identifier{"try", meta}, []Node{
						Scope{Block{[]Node{
							Expression{Identifier{"print", meta}, []Node{Identifier{"test", meta}}, meta},
						}, meta}, meta},
						Scope{Block{[]Node{}, meta}, meta},
					}, meta},
				},
				Meta: meta,
			},
			false,
		},
		{
			"return assign",
			args{tokens.Lexerp(`
			err = try({
				print(1)
			})`)},
			Block{
				[]Node{
					Assignment{
						Identifier{"err", meta},
						Expression{
							Identifier{"try", meta},
							[]Node{
								Scope{
									Block{
										[]Node{
											Expression{
												Identifier{"print", meta},
												[]Node{Int{1, meta}},
												meta,
											},
										},
										meta},
									meta},
							},
							meta},
						meta,
					},
				},
				meta,
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
					Expression{MemberSelector{Identifier{"var", meta}, Identifier{"member", meta}, meta}, []Node{}, meta},
				},
				meta,
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
						Identifier{"l", meta},
						Expression{
							MemberSelector{
								Expression{Identifier{"str", meta}, []Node{Identifier{"varString", meta}}, meta},
								Identifier{"len", meta},
								meta,
							},
							[]Node{},
							meta,
						},
						meta,
					},
				},
				meta,
			},
			false,
		},
		{
			"function definition",
			args{tokens.Lexerp(`
				handle= (err){
					print("Err: %s", err)
				}
			`)},
			Block{
				[]Node{
					Assignment{
						Identifier{"handle", meta},
						FunctionDefinition{
							Block{
								[]Node{
									Expression{
										Identifier{"print", meta},
										[]Node{String{"Err: %s", meta}, Identifier{"err", meta}},
										meta,
									},
								},
								meta,
							}, []Identifier{{"err", meta}},
							meta},
						meta},
				},
				meta,
			},
			false,
		},
		{
			"inline func call",
			args{tokens.Lexerp(`
				var = {"test"}()
			`)},
			Block{
				[]Node{
					Assignment{
						Identifier{"var", meta},
						Expression{
							Scope{
								Block{[]Node{String{"test", meta}}, meta},
								meta,
							},
							[]Node{},
							meta,
						},
						meta,
					},
				},
				meta,
			},
			false,
		},
		{
			"Access and Expression",
			args{tokens.Lexerp(`
			v = a.value
			print("Magik: %v", v)
			`)},
			Block{
				[]Node{
					Assignment{
						Identifier{"v", meta},
						MemberSelector{Identifier{"a", meta}, Identifier{"value", meta}, meta},
						meta,
					},
					Expression{
						Identifier{"print", meta},
						[]Node{
							String{"Magik: %v", meta},
							Identifier{"v", meta},
						},
						meta,
					},
				},
				meta,
			},
			false,
		},
		{
			"Deep nest",
			args{tokens.Lexerp(`
			a.c().value = 1
			`)},
			Block{
				[]Node{
					Assignment{
						MemberSelector{
							Expression{
								MemberSelector{
									Identifier{"a", meta},
									Identifier{"c", meta},
									meta,
								},
								[]Node{},
								meta},
							Identifier{"value", meta},
							meta,
						},
						Int{1, meta},
						meta,
					},
				},
				meta,
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
			walkTree(got, func(node Node) {
				node.SetToken(meta.At) // Nulling meta for testing, would be to anoying
			})
			if !reflect.DeepEqual(got, tt.want) {
				json.NewEncoder(os.Stdout).Encode(got)
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
