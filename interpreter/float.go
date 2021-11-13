package interpreter

import "strconv"

var floatMap = NativeMap{}

func (Float) Type() String {
	return "builtin+float"
}

func floatWrapUnary(f func(a float64) Value) Function {
	return func(r *Runtime, v []Value) (Value, *Error) {
		return f(float64(v[0].(Float))), nil
	}
}

func floatWrapBinary(f func(a, b float64) Value) Function {
	return func(r *Runtime, v []Value) (Value, *Error) {
		a := v[0].(Float)
		if b, ok := v[1].(Float); ok {
			return f(float64(a), float64(b)), nil
		}
		return builtinThrow(r, []Value{String("Invalid right hand type")})
	}
}

func init() {
	floatMap[NativeBool] = Function(func(r *Runtime, v []Value) (Value, *Error) { return Bool(v[0].(Float) != 0.0), nil })
	floatMap[NativeStr] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		return String(strconv.FormatFloat(float64(v[0].(Float)), 'f', 7, 64)), nil
	})

	intNatives[NativeEq] = floatWrapBinary(func(a, b float64) Value { return Bool(a == b) })
	intNatives[NativeLt] = floatWrapBinary(func(a, b float64) Value { return Bool(a < b) })
	intNatives[NativeLe] = floatWrapBinary(func(a, b float64) Value { return Bool(a <= b) })
	intNatives[NativeGt] = floatWrapBinary(func(a, b float64) Value { return Bool(a > b) })
	intNatives[NativeGe] = floatWrapBinary(func(a, b float64) Value { return Bool(a >= b) })

	floatNatives[NativeAdd] = floatWrapBinary(floatAdd)
	floatNatives[NativeSub] = floatWrapBinary(floatSub)
	floatNatives[NativeMul] = floatWrapBinary(floatMul)
	floatNatives[NativeDiv] = floatWrapBinary(floatDiv)
}

func floatAdd(a, b float64) Value { return Float(a + b) }
func floatSub(a, b float64) Value { return Float(a - b) }
func floatMul(a, b float64) Value { return Float(a * b) }
func floatDiv(a, b float64) Value { return Float(a / b) }

func (Float) Natives() NativeMap {
	return floatMap
}
