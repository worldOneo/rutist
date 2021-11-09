package interpreter

import (
	"github.com/worldOneo/rutist/ast"
	"github.com/worldOneo/rutist/tokens"
)

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
	TypeGetMember = String("__getmember__")
	TypeAdd       = String("__add__")
	TypeSub       = String("__sub__")
	TypeMul       = String("__mul__")
	TypeDiv       = String("__div__")
	TypeMod       = String("__mod__")
	TypeOr        = String("__or__")
	TypeAnd       = String("__and__")
	TypeXor       = String("__xor__")
	TypeNot       = String("__not__")
	TypeEq        = String("__eq__")
	TypeLt        = String("__lt__")
	TypeLe        = String("__le__")
	TypeGt        = String("__gt__")
	TypeGe        = String("__ge__")
	TypeLsh       = String("__lsh__")
	TypeRsh       = String("__rsh__")
)

var operatorMagicType = map[tokens.Operator]String{}

func init() {
	operatorMagicType[tokens.OperatorAdd] = TypeAdd
	operatorMagicType[tokens.OperatorSub] = TypeSub
	operatorMagicType[tokens.OperatorMul] = TypeMul
	operatorMagicType[tokens.OperatorDiv] = TypeDiv
	operatorMagicType[tokens.OperatorMod] = TypeMod
	operatorMagicType[tokens.OperatorOr] = TypeOr
	operatorMagicType[tokens.OperatorAnd] = TypeAnd
	operatorMagicType[tokens.OperatorXor] = TypeXor
	operatorMagicType[tokens.OperatorNot] = TypeNot
	operatorMagicType[tokens.OperatorEq] = TypeEq
	operatorMagicType[tokens.OperatorLt] = TypeLt
	operatorMagicType[tokens.OperatorLe] = TypeLe
	operatorMagicType[tokens.OperatorGt] = TypeGt
	operatorMagicType[tokens.OperatorGe] = TypeGe
	operatorMagicType[tokens.OperatorLsh] = TypeLsh
	operatorMagicType[tokens.OperatorRsh] = TypeRsh
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

func (FuncDef) Type() String {
	return "builtin+funcdef"
}

func (String) Type() String {
	return "builtin+string"
}

func (Function) Type() String {
	return "builtin+function"
}

func (Int) Type() String {
	return "builtin+integer"
}

func (Float) Type() String {
	return "builtin+float"
}

func (Bool) Type() String {
	return "builtin+bool"
}

func (Scoope) Type() String {
	return "builtin+scope"
}

func (Error) Type() String {
	return "buitin+error"
}

func mapGet(r *Runtime, v []Value) (Value, *Error) {
	if len(v) != 2 {
		return builtinThrow(r, []Value{String("Map: Get Requires exactly 1 parameter")})
	}
	return v[0].(Map)[v[1]], nil
}

func mapSet(r *Runtime, v []Value) (Value, *Error) {
	if len(v) != 3 {
		return builtinThrow(r, []Value{String("Map: Set Requires exactly 2 parameters")})
	}
	v[0].(Map)[v[1]] = v[2]
	return nil, nil
}

func mapHas(r *Runtime, v []Value) (Value, *Error) {
	if len(v) != 2 {
		return builtinThrow(r, []Value{String("Map: Has Requires exactly 1 parameter")})
	}
	_, ok := v[0].(Map)[v[1]]
	return Bool(ok), nil
}

func (Map) Type() String {
	return "builtin+map"
}

func (Dict) Type() String {
	return "builtin+dict"
}
