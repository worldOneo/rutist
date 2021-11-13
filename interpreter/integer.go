package interpreter

import "strconv"

type Int int

var intNatives = NativeMap{}

func (Int) Type() String {
	return "builtin+integer"
}

func intWrapBinary(f func(a, b int) Value, floatBinary func(a, b float64) Value) Function {
	return func(r *Runtime, v []Value) (Value, *Error) {
		a := v[0].(Int)
		if b, ok := v[1].(Int); ok {
			return f(int(a), int(b)), nil
		}
		if b, ok := v[1].(Float); ok {
			return floatBinary(float64(a), float64(b)), nil
		}
		return builtinThrow(r, []Value{String("Invalid right hand value")})
	}
}

func intWrapUnary(f func(a int) Value) Function {
	return func(r *Runtime, v []Value) (Value, *Error) {
		return f(int(v[0].(Int))), nil
	}
}

func init() {
	intNatives[NativeBool] = intWrapUnary(func(a int) Value { return Bool(a != 0) })
	intNatives[NativeNot] = intWrapUnary(func(a int) Value { return Bool(a == 0) })

	intNatives[NativeStr] = intWrapUnary(func(a int) Value { return String(strconv.Itoa(a)) })

	intNatives[NativeEq] = intWrapBinary(func(a, b int) Value { return Bool(a == b) }, nil)
	intNatives[NativeLt] = intWrapBinary(func(a, b int) Value { return Bool(a < b) }, nil)
	intNatives[NativeLe] = intWrapBinary(func(a, b int) Value { return Bool(a <= b) }, nil)
	intNatives[NativeGt] = intWrapBinary(func(a, b int) Value { return Bool(a > b) }, nil)
	intNatives[NativeGe] = intWrapBinary(func(a, b int) Value { return Bool(a >= b) }, nil)

	intNatives[NativeOr] = intWrapBinary(func(a, b int) Value { return Int(a | b) }, nil)
	intNatives[NativeAnd] = intWrapBinary(func(a, b int) Value { return Int(a & b) }, nil)
	intNatives[NativeXor] = intWrapBinary(func(a, b int) Value { return Int(a ^ b) }, nil)
	intNatives[NativeLsh] = intWrapBinary(func(a, b int) Value { return Int(a << b) }, nil)
	intNatives[NativeRsh] = intWrapBinary(func(a, b int) Value { return Int(a >> b) }, nil)
	intNatives[NativeMod] = intWrapBinary(func(a, b int) Value { return Int(a % b) }, nil)
	intNatives[NativeOr] = intWrapBinary(func(a, b int) Value { return Int(a | b) }, nil)

	intNatives[NativeAdd] = intWrapBinary(func(a, b int) Value { return Int(a + b) }, floatAdd)
	intNatives[NativeSub] = intWrapBinary(func(a, b int) Value { return Int(a - b) }, floatSub)
	intNatives[NativeMul] = intWrapBinary(func(a, b int) Value { return Int(a * b) }, floatMul)
	intNatives[NativeDiv] = intWrapBinary(func(a, b int) Value { return Int(a / b) }, floatDiv)
}

func (Int) Natives() NativeMap {
	return intNatives
}
