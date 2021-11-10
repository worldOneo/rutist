package interpreter

import (
	"github.com/worldOneo/rutist/ast"
	"github.com/worldOneo/rutist/tokens"
)

type String string
type Error struct {
	Err error
}

type Map map[Value]Value
type Dict Map
type Float float64
type Scoope struct {
	node ast.Node
}
type FuncDef struct {
	args []ast.Identifier
	node ast.Node
}

const (
	NativeRun = iota
	NativeStr
	NativeLen
	NativeBool
	NativeSetMember
	NativeGetMember
	NativeAdd
	NativeSub
	NativeMul
	NativeDiv
	NativeMod
	NativeOr
	NativeAnd
	NativeXor
	NativeNot
	NativeEq
	NativeLt
	NativeLe
	NativeGt
	NativeGe
	NativeLsh
	NativeRsh
)

type NativeMap [NativeRsh]Value

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

var operatorMagicType = map[tokens.Operator]int{}

func init() {
	operatorMagicType[tokens.OperatorAdd] = NativeAdd
	operatorMagicType[tokens.OperatorSub] = NativeSub
	operatorMagicType[tokens.OperatorMul] = NativeMul
	operatorMagicType[tokens.OperatorDiv] = NativeDiv
	operatorMagicType[tokens.OperatorMod] = NativeMod
	operatorMagicType[tokens.OperatorOr] = NativeOr
	operatorMagicType[tokens.OperatorAnd] = NativeAnd
	operatorMagicType[tokens.OperatorXor] = NativeXor
	operatorMagicType[tokens.OperatorNot] = NativeNot
	operatorMagicType[tokens.OperatorEq] = NativeEq
	operatorMagicType[tokens.OperatorLt] = NativeLt
	operatorMagicType[tokens.OperatorLe] = NativeLe
	operatorMagicType[tokens.OperatorGt] = NativeGt
	operatorMagicType[tokens.OperatorGe] = NativeGe
	operatorMagicType[tokens.OperatorLsh] = NativeLsh
	operatorMagicType[tokens.OperatorRsh] = NativeRsh
}

var floatNatives = NativeMap{}

var this = Function(func (_ *Runtime, v []Value) (Value, *Error) {
	return v[0], nil
})
