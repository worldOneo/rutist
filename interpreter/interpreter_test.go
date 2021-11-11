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
				err = try({
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
				handle = (err){
					print("Err: %s", err)
				}
			`))},
			func(r *Runtime) bool {
				handle := r.GetVar("handle")
				if handle == nil {
					return false
				}
				h, o := handle.(*FuncDef)
				if !o {
					return false
				}
				return reflect.DeepEqual(h.args, []ast.Identifier{{"err", ast.NewMeta(tokens.Token{tokens.Identifier, "err", 0, 0, 1})}})
			},
			false,
		},
		{
			"inline func call",
			args{ast.Parsep(tokens.Lexerp(`
				var = {"test"}()
			`))},
			func(r *Runtime) bool {
				v := r.GetVar("var")
				if v == nil {
					return false
				}
				return v.(String) == String("test")
			},
			false,
		},
		{
			"module export test",
			args{ast.Parsep(tokens.Lexerp(`
			module((export) {
				export("value", 1)
			})
			`))},
			func(r *Runtime) bool {
				v, ok := r.SpecialFields[SpecialfFieldExport].(Dict)[String("value")]
				if !ok || v == nil {
					return false
				}
				return v == Int(1)
			},
			false,
		},
		{
			"set attribute",
			args{ast.Parsep(tokens.Lexerp(`
				a = Dict()
				a.value = 1
				v = a.value
				print("Magik: %v", v)
			`))},
			func(r *Runtime) bool {
				v := r.GetVar("v")
				return v == Int(1)
			},
			false,
		},
		{
			"wrapped function",
			args{ast.Parsep(tokens.Lexerp(`
				a = Dict()
				b = "test"
				a.dict = Dict()
				a.dict.value = b.len
				v = a.dict.value()
				print("Magik: %v", v)
			`))},
			func(r *Runtime) bool {
				v := r.GetVar("v")
				return v == Int(4)
			},
			false,
		},
		{
			"deep nest",
			args{ast.Parsep(tokens.Lexerp(`
				a = Dict()
				b = Dict()
				c = (){ b }
				a.c = c
				a.c().value = "test".len
				v = a.c().value()
				print("Magik: %v", v)
			`))},
			func(r *Runtime) bool {
				v := r.GetVar("v")
				return v == Int(4)
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New("test.go")
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

func TestRun_Abstract(t *testing.T) {
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
			"Class Def",
			args{
				ast: ast.Parsep(tokens.Lexerp(`
					init = class((def) {
						def("__init__", (self, name) {
							self._name = name
						})

						def("name", (self) { self._name })
					})

					inst = init("test")
					test = inst.name()
				`)),
			},
			func(r *Runtime) bool {
				t := r.GetVar("test")
				return t != nil && t.(String) == "test"
			},
			false,
		},
		{
			"Native Overload",
			args{
				ast: ast.Parsep(tokens.Lexerp(`
					init = class((def) {
						def("__getmember__", (self, member) { member })
					})

					inst = init()
					test = inst.name
				`)),
			},
			func(r *Runtime) bool {
				t := r.GetVar("test")
				return t != nil && t.(String) == "name"
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New("test.go")
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
