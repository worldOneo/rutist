package interpreter

import (
	"github.com/worldOneo/rutist/tokens"
)

type Map map[Value]Value
type Dict Map
type Float float64

const (
	NativeRun = iota
	NativeInit
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
	NativeLor
	NativeLand
	NativeLsh
	NativeRsh
)

type NativeMap [NativeRsh + 1]Value

const (
	TypeRun       = String("__run__")
	TypeInit      = String("__init__")
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

var magicFuncNativeMap = map[String]int{}

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
	operatorMagicType[tokens.OperatorLor] = NativeLor
	operatorMagicType[tokens.OperatorLand] = NativeLand
	operatorMagicType[tokens.OperatorLsh] = NativeLsh
	operatorMagicType[tokens.OperatorRsh] = NativeRsh

	magicFuncNativeMap[TypeRun] = NativeRun
	magicFuncNativeMap[TypeInit] = NativeInit
	magicFuncNativeMap[TypeStr] = NativeStr
	magicFuncNativeMap[TypeLen] = NativeLen
	magicFuncNativeMap[TypeBool] = NativeBool
	magicFuncNativeMap[TypeSetMember] = NativeSetMember
	magicFuncNativeMap[TypeGetMember] = NativeGetMember
	magicFuncNativeMap[TypeAdd] = NativeAdd
	magicFuncNativeMap[TypeSub] = NativeSub
	magicFuncNativeMap[TypeMul] = NativeMul
	magicFuncNativeMap[TypeDiv] = NativeDiv
	magicFuncNativeMap[TypeMod] = NativeMod
	magicFuncNativeMap[TypeOr] = NativeOr
	magicFuncNativeMap[TypeAnd] = NativeAnd
	magicFuncNativeMap[TypeXor] = NativeXor
	magicFuncNativeMap[TypeNot] = NativeNot
	magicFuncNativeMap[TypeEq] = NativeEq
	magicFuncNativeMap[TypeLt] = NativeLt
	magicFuncNativeMap[TypeLe] = NativeLe
	magicFuncNativeMap[TypeGt] = NativeGt
	magicFuncNativeMap[TypeGe] = NativeGe
	magicFuncNativeMap[TypeLsh] = NativeLsh
	magicFuncNativeMap[TypeRsh] = NativeRsh
}

var floatNatives = NativeMap{}

var this = Function(func(_ *Runtime, v []Value) (Value, *Error) {
	return v[0], nil
})
