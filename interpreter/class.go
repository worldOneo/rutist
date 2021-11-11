package interpreter

type Class struct {
	natives NativeMap
	members MemberDict
}

type Constructor struct {
	template *Class
}

type Instance struct {
	of      *Class
	members MemberDict
}

var constructorNatives = NativeMap{}
var instanceNatives = NativeMap{}

func (Constructor) Type() String {
	return "builtin+constructor"
}

func init() {
	constructorNatives[NativeRun] = Function(func(r *Runtime, v []Value) (Value, *Error) {
		c := v[0].(Constructor)
		init := c.template.natives[NativeInit]
		inst := &Instance{c.template, MemberDict{}}
		if init == nil {
			return inst, nil
		}
		_, err := r.invokeValue(init, append([]Value{inst}, v[1:]...))
		if err != nil {
			return nil, err
		}
		return inst, nil
	})
}

func (Constructor) Natives() NativeMap {
	return constructorNatives
}

func (I *Instance) Natives() NativeMap {
	return I.of.natives
}

func (Instance) Type() String {
	return "class"
}

func builtinClass(r *Runtime, v []Value) (Value, *Error) {
	class := &Class{
		natives: NativeMap{},
		members: MemberDict{},
	}
	def := func(_ *Runtime, v []Value) (Value, *Error) {
		if len(v) != 2 {
			return builtinThrow(r, []Value{String("def: requires exactly 2 parameter")})
		}
		str, ok := v[0].(String)
		if !ok {
			return builtinThrow(r, []Value{String("def: Parameter 1 must be string")})
		}
		native, ok := magicFuncNativeMap[str]
		if !ok {
			class.members[v[0]] = v[1]
			return nil, nil
		}
		funcDef, ok := v[1].(*FuncDef)
		if !ok {
			class.natives[native] = v[1]
			return nil, nil
		}
		class.natives[native] = Function(funcDef.run)
		return nil, nil
	}
	_, err := r.invokeValue(v[0], []Value{Function(def)})
	if err != nil {
		return nil, err
	}
	if class.natives[NativeGetMember] == nil {
		class.natives[NativeGetMember] = Function(func(r *Runtime, v []Value) (Value, *Error) {
			member, ok := v[0].(*Instance).of.members[v[1]]
			if ok {
				return member, nil
			}
			return v[0].(*Instance).members[v[1]], nil

		})
	}

	if class.natives[NativeSetMember] == nil {
		class.natives[NativeSetMember] = Function(func(r *Runtime, v []Value) (Value, *Error) {
			v[0].(*Instance).members[v[1]] = v[2]
			return nil, nil

		})
	}
	return Constructor{class}, nil
}
