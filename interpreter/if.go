package interpreter

type ifState struct {
	fullfilled bool
	val        Value
}

func (ifState) Type() String {
	return "builtin+ifState"
}

var ifStateNatives = NativeMap{}

func init() {
	ifStateNatives[NativeGetMember] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		state := v[0].(*ifState)
		if member, str := v[1].(String); str {
			if member == "elseif" {
				if state.fullfilled {
					return createTrashCan(v[0]), nil
				}
				return Function(builtinElse), nil
			} else if member == "else" {
				if state.fullfilled {
					return createTrashCan(v[0]), nil
				}
				return Function(func(r *Runtime, v []Value) (Value, *Error) {
					val, err := builtinRun(r, v[1:])
					return &ifState{true, val}, err
				}), nil
			} else if member == "value" {
				return state.val, nil
			}
		}
		return nil, nil
	})
}

func (ifState) Natives() NativeMap {
	return ifStateNatives
}

func builtinElse(r *Runtime, args []Value) (Value, *Error) {
	return builtinIf(r, args[1:])
}

func builtinIf(r *Runtime, args []Value) (Value, *Error) {
	if len(args) != 2 {
		return builtinThrow(r, []Value{String("If: requires exactly 2 args")})
	}
	statement := r.getNativeField(args[0], NativeRun)
	if statement == nil {
		return builtinThrow(r, []Value{String("If: arg 1 must be runnable")})
	}
	if fn, ok := statement.(Function); ok {
		val, err := fn(r, []Value{args[0]})
		if err != nil {
			return nil, err
		}
		boolable := r.getNativeField(val, NativeBool)
		if boolable == nil {
			return builtinThrow(r, []Value{String("If: arg 1 must be boolable")})
		}
		boolFn, ok := boolable.(Function)
		if !ok {
			return builtinThrow(r, []Value{String("If: arg 1 must be boolable")})
		}
		boolVal, err := boolFn(r, []Value{val})
		if err != nil {
			return nil, err
		}
		if boolVal, ok := boolVal.(Bool); ok {
			if boolVal {
				val, err := builtinRun(r, []Value{args[1]})
				if err != nil {
					return nil, err
				}
				return &ifState{true, val}, nil
			}
		}
	}
	return &ifState{false, nil}, nil
}
