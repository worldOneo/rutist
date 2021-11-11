package interpreter

type Bool bool

var boolNatives = NativeMap{}

func (Bool) Natives() NativeMap {
	return boolNatives
}

func (Bool) Type() String {
	return "builtin+bool"
}

func init() {
	boolNatives[NativeBool] = this
	boolNatives[NativeStr] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		if v[0].(Bool) {
			return String("true"), nil
		}
		return String("false"), nil
	})
}
