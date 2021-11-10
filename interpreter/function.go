package interpreter

var functionNatives = NativeMap{}

type Function func(*Runtime, []Value) (Value, *Error)

func (Function) Type() String {
	return "builtin+function"
}

func (Function) Natives() NativeMap {
	return functionNatives
}

func init() {
	functionNatives[NativeRun] = Function(func(r *Runtime, v []Value) (Value, *Error) { return v[0].(Function)(r, v[1:]) })
}
