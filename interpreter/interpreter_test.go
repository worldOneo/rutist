package interpreter

import (
	"fmt"
	"testing"

	"github.com/worldOneo/rutist/ast"
	"github.com/worldOneo/rutist/tokens"
)

func TestRun_Status(t *testing.T) {
	type args struct {
		ast ast.Node
	}
	tests := []struct {
		name    string
		args    args
		want    func(*Runtime) bool
		wantErr bool
	}{
		{
			"Invoke",
			args{
				ast: ast.Parsep(tokens.Lexerp(`
					print("test")
				`)),
			},
			func(r *Runtime) bool {
				return true
			},
			false,
		}, {
			"TryCatch",
			args{
				ast: ast.Parsep(tokens.Lexerp(`
				err = try(@{
					throw("This is an error")
				})
				`)),
			},
			func(r *Runtime) bool {
				fmt.Printf("%v", r.CurrentScope())
				v := r.GetVar("err").(*Error)
				_, err := builtinThrow(r, []Value{String("This is an error")})
				return v.Err.Error() == err.Err.Error()
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runtime{
				[]*Scope{NewScope(map[string]Value{})},
				0,
			}
			_, err := r.Run(tt.args.ast)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.want(r) {
				t.Error("Condition failed")
			}
		})
	}
}
