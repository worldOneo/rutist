package interpreter

var scoopeNative = NativeMap{}

func (Scoope) Type() String {
	return "builtin+scope"
}

func init() {
	scoopeNative[NativeRun] = Function(func(r *Runtime, v []Value) (Value, *Error) { return r.Run(v[0].(*Scoope).node) })
}

func (Scoope) Natives() NativeMap {
	return scoopeNative
}