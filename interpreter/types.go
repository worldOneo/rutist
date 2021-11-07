package interpreter

import "github.com/worldOneo/rutist/ast"

type String string
type Int int
type Error struct {
	Err error
}

type Float float64
type Function func(*Runtime, []Value) (Value, *Error)
type Bool bool
type Scoope struct {
	node ast.Node
}
type FuncDef struct {
	args []ast.Identifier
	node ast.Node
}

const (
	TypeRun  = "__run__"
	TypeStr  = "__str__"
	TypeLen  = "__len__"
	TypeBool = "__bool__"
)

func (S String) Members() MemberDict {
	return MemberDict{
		"len":    Function(S.len),
		TypeStr:  Function(func(_ *Runtime, _ []Value) (Value, *Error) { return S, nil }),
		TypeLen:  Function(S.len),
		TypeBool: Function(func(_ *Runtime, _ []Value) (Value, *Error) { return Bool(S != ""), nil }),
	}
}

func (F FuncDef) Members() MemberDict {
	return MemberDict{
		TypeRun: Function(F.run),
	}
}

func (F *FuncDef) run(r *Runtime, v []Value) (Value, *Error) {
	for i := 0; i < len(F.args); i++ {
		if i > len(v) {
			break
		}
		r.CurrentScope().variables[F.args[i].Name] = v[i]
	}
	return r.Run(F.node)
}

func (S String) len(_ *Runtime, _ []Value) (Value, *Error) {
	return Int(len(S)), nil
}

func (F Function) Members() MemberDict {
	return MemberDict{
		TypeRun: F,
	}
}

func (I Int) Members() MemberDict {
	return MemberDict{
		TypeBool: Function(func(r *Runtime, v []Value) (Value, *Error) { return Bool(I != 0), nil }),
	}
}

func (F Float) Members() MemberDict {
	return MemberDict{
		TypeBool: Function(func(r *Runtime, v []Value) (Value, *Error) { return Bool(F != 0.0), nil }),
	}
}

func (B Bool) Members() MemberDict {
	return MemberDict{
		TypeBool: Function(func(r *Runtime, v []Value) (Value, *Error) { return B, nil }),
	}
}

func (S *Scoope) Members() MemberDict {
	return MemberDict{
		TypeRun: Function(func(r *Runtime, v []Value) (Value, *Error) { return r.Run(S.node) }),
	}
}

func (E *Error) Members() MemberDict {
	return MemberDict{}
}
