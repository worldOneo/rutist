package interpreter

import "strconv"

var floatMap = NativeMap{}

func (Float) Type() String {
	return "builtin+float"
}

func init() {
	floatMap[NativeBool] = Function(func(r *Runtime, v []Value) (Value, *Error) { return Bool(v[0].(Float) != 0.0), nil })
	floatMap[NativeStr] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		return String(strconv.FormatFloat(float64(v[0].(Float)), 'f', 7, 64)), nil
	})
}

func (Float) Natives() NativeMap {
	return floatMap
}
