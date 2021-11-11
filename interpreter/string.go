package interpreter

type String string

var stringNatives = NativeMap{}

func init() {
	stringNatives[NativeStr] = this
	stringNatives[NativeLen] = Function(func(_ *Runtime, v []Value) (Value, *Error) { return Int(len(v[0].(String))), nil })
	stringNatives[NativeBool] = Function(func(_ *Runtime, v []Value) (Value, *Error) { return Bool(v[0].(String) != ""), nil })
	stringNatives[NativeGetMember] = Function(func(_ *Runtime, v []Value) (Value, *Error) {
		if str, ok := v[1].(String); ok && str == "len" {
			return stringNatives[NativeLen], nil
		}
		return nil, nil
	})
}

func (String) Type() String {
	return "builtin+string"
}

func init() {
	stringNatives[NativeLen] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		return Int(len(v[0].(String))), nil
	})
}

func (String) Natives() NativeMap {
	return stringNatives
}
