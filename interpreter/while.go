package interpreter

func builtinWhile(r *Runtime, args []Value) (Value, *Error) {
	if len(args) != 2 {
		return builtinThrow(r, []Value{String("While: Requires exactly 2 parameters")})
	}
	statement := r.getNativeField(args[0], NativeRun)
	r.lowerScope()
	defer r.raiseScope()
	if statement == nil {
		return builtinThrow(r, []Value{String("While: arg 1 must be runnable")})
	}
	if fn, ok := statement.(Function); ok {
		var res Value
		for true {
			val, err := fn(r, []Value{args[0]})
			if err != nil {
				return nil, err
			}
			boolable := r.getNativeField(val, NativeBool)
			if boolable == nil {
				return builtinThrow(r, []Value{String("While: arg 1 must be boolable")})
			}
			boolFn, ok := boolable.(Function)
			if !ok {
				return builtinThrow(r, []Value{String("While: arg 1 must be boolable")})
			}
			boolVal, err := boolFn(r, []Value{val})
			if err != nil {
				return nil, err
			}
			if boolVal, ok := boolVal.(Bool); ok {
				if !boolVal {
					break
				}
				res, err = builtinRun(r, []Value{args[1]})
				if err != nil {
					return nil, err
				}
			}
		}
		return res, nil
	}
	return builtinThrow(r, []Value{String("While:")})
}
