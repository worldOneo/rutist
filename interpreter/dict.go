package interpreter

func (Dict) Type() String {
	return "builtin+dict"
}

var dictNatives = NativeMap{}

func init() {
	dictNatives[NativeSetMember] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		v[0].(Dict)[v[1]] = v[2]
		return nil, nil
	})
	dictNatives[NativeGetMember] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		return v[0].(Dict)[v[1]], nil
	})
}

func (Dict) Natives() NativeMap {
	return dictNatives
}
