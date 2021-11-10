package interpreter

var mapNatives = NativeMap{}

func (Map) Type() String {
	return "builtin+map"
}

func init() {
	mapNatives[NativeGetMember] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		str, ok := v[1].(String)
		if !ok {
			return nil, nil
		}
		if str == "get" {
			return Function(mapGet), nil
		}
		if str == "set" {
			return Function(mapSet), nil
		}
		if str == "has" {
			return Function(mapHas), nil
		}
		return nil, nil
	})
}

func (Map) Natives() NativeMap {
	return mapNatives
}

func mapGet(r *Runtime, v []Value) (Value, *Error) {
	if len(v) != 2 {
		return builtinThrow(r, []Value{String("Map: Get Requires exactly 1 parameter")})
	}
	return v[0].(Map)[v[1]], nil
}

func mapSet(r *Runtime, v []Value) (Value, *Error) {
	if len(v) != 3 {
		return builtinThrow(r, []Value{String("Map: Set Requires exactly 2 parameters")})
	}
	v[0].(Map)[v[1]] = v[2]
	return nil, nil
}

func mapHas(r *Runtime, v []Value) (Value, *Error) {
	if len(v) != 2 {
		return builtinThrow(r, []Value{String("Map: Has Requires exactly 1 parameter")})
	}
	_, ok := v[0].(Map)[v[1]]
	return Bool(ok), nil
}
