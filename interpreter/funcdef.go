package interpreter

import "github.com/worldOneo/rutist/ast"

var funcdefNatives = NativeMap{}

type FuncDef struct {
	args     []ast.Identifier
	node     ast.Node
	captured map[string]Value
}

func (FuncDef) Type() String {
	return "builtin+funcdef"
}

func init() {
	funcdefNatives[NativeRun] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		return v[0].(*FuncDef).run(r, v[1:])
	})
}

func (FuncDef) Natives() NativeMap {
	return funcdefNatives
}

func (F *FuncDef) run(r *Runtime, v []Value) (Value, *Error) {
	for k, v := range F.captured {
		if r.GetVar(k) != nil {
			continue
		}
		r.CurrentScope().variables[k] = v
	}
	for i := 0; i < len(F.args); i++ {
		if i >= len(v) {
			break
		}
		r.CurrentScope().variables[F.args[i].Name] = v[i]
	}
	return r.Run(F.node)
}
