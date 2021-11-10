package interpreter

import "strconv"

type Int int

var intNatives = NativeMap{}

func (Int) Type() String {
	return "builtin+integer"
}

func init() {
	intNatives[NativeBool] = Function(func(r *Runtime, v []Value) (Value, *Error) { return Bool(v[0].(Int) != 0), nil })
	intNatives[NativeStr] = Function(func(r *Runtime, v []Value) (Value, *Error) { return String(strconv.Itoa(int(v[0].(Int)))), nil })
}

func (Int) Natives() NativeMap {
	return intNatives
}
