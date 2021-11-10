package interpreter

type WrappedFunction struct {
	this Value
	f    Function
}

func wrappMemberFunction(v Value, f Function) WrappedFunction {
	return WrappedFunction{
		v,
		f,
	}
}

func (WrappedFunction) Type() String {
	return "builtin+wrapped_function"
}

var wrappedFunctionNative = NativeMap{}

func init() {
	wrappedFunctionNative[NativeRun] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		w := v[0].(WrappedFunction)
		args := append([]Value{w.this}, v[1:]...)
		return w.f(r, args)
	})
}

func (WrappedFunction) Natives() NativeMap {
	return wrappedFunctionNative
}
