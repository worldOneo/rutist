package interpreter

type Bool bool

var boolNatives = NativeMap{}

func (Bool) Natives() NativeMap {
	return boolNatives
}

func (Bool) Type() String {
	return "builtin+bool"
}

func boolWrapUnary(f func(a bool) Value) Function {
	return func(r *Runtime, v []Value) (Value, *Error) {
		return f(bool(v[0].(Bool))), nil
	}
}

func init() {
	boolNatives[NativeBool] = this
	boolNatives[NativeStr] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		if v[0].(Bool) {
			return String("true"), nil
		}
		return String("false"), nil
	})

	boolNatives[NativeNot] = boolWrapUnary(func(a bool) Value { return Bool(!a) })
}
