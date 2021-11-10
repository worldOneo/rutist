package interpreter

var funcdefNatives = NativeMap{}

func (FuncDef) Type() String {
	return "builtin+funcdef"
}

func init() {
	funcdefNatives[NativeRun] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		return v[0].(*FuncDef).run(r, v[1:])
	})
}

func (FuncDef) Natives() NativeMap {
	return funcdefNatives
}

func (F *FuncDef) run(r *Runtime, v []Value) (Value, *Error) {
	for i := 0; i < len(F.args); i++ {
		if i > len(v) {
			break
		}
		r.CurrentScope().variables[F.args[i].Name] = v[i]
	}
	return r.Run(F.node)
}
