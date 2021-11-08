package interpreter

import "github.com/worldOneo/rutist/ast"

type String string
type Int int
type Error struct {
	Err error
}

type Map map[Value]Value
type Dict Map
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
	TypeRun       = String("__run__")
	TypeStr       = String("__str__")
	TypeLen       = String("__len__")
	TypeBool      = String("__bool__")
	TypeSetMember = String("__setmember__")
)

func (S String) Members() MemberDict {
	return MemberDict{
		String("len"): Function(S.len),
		TypeStr:       Function(func(_ *Runtime, _ []Value) (Value, *Error) { return S, nil }),
		TypeLen:       Function(S.len),
		TypeBool:      Function(func(_ *Runtime, _ []Value) (Value, *Error) { return Bool(S != ""), nil }),
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

func (M Map) Members() MemberDict {
	return MemberDict{
		String("get"): Function(M.Get),
		String("set"): Function(M.Set),
		String("has"): Function(M.Has),
	}
}

func (M Map) Get(r *Runtime, v []Value) (Value, *Error) {
	if len(v) != 1 {
		return builtinThrow(r, []Value{String("Map: Get Requires exactly 1 parameter")})
	}
	return M[v[0]], nil
}

func (M Map) Set(r *Runtime, v []Value) (Value, *Error) {
	if len(v) != 2 {
		return builtinThrow(r, []Value{String("Map: Set Requires exactly 2 parameters")})
	}
	M[v[0]] = v[1]
	return nil, nil
}

func (M Map) Has(r *Runtime, v []Value) (Value, *Error) {
	if len(v) != 1 {
		return builtinThrow(r, []Value{String("Map: Has Requires exactly 1 parameter")})
	}
	_, ok := M[v[0]]
	return Bool(ok), nil
}

func (D Dict) Members() MemberDict {
	D[TypeSetMember] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		D[v[0]] = v[1]
		return nil, nil
	})
	return MemberDict(D)
}
