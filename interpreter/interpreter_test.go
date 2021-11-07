package interpreter

import (
	"fmt"
	"reflect"
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
		},
		{
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
		{
			"member invoke",
			args{
				ast.Parsep(tokens.Lexerp(`
				varString = "test"
				l = varString.len()
				`)),
			},
			func(r *Runtime) bool {
				v, ok := r.CurrentScope().variables["l"]
				return ok && v == Int(4)
			},
			false,
		},
		{
			"Result member access invoke",
			args{
				ast.Parsep(tokens.Lexerp(`
				varString = "test"
				l = str(varString).len()
				`)),
			},
			func(r *Runtime) bool {
				v, ok := r.CurrentScope().variables["l"]
				return ok && v == Int(4)
			},
			false,
		},
		{
			"function definition",
			args{ast.Parsep(tokens.Lexerp(`
				handle=@(err){
					print("Err: %s", err)
				}
			`))},
			func(r *Runtime) bool {
				handle, ok := r.CurrentScope().variables["handle"]
				if !ok {
					return false
				}
				h, o := handle.(*FuncDef)
				if !o {
					return false
				}
				return reflect.DeepEqual(h.args, []ast.Identifier{{"err"}})
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
