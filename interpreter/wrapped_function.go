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

func (W WrappedFunction) Type() String {
	return "builtin+wrapped_function"
}
