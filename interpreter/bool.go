package interpreter

type Bool bool

var boolNatives = NativeMap{}

func (Bool) Natives() NativeMap {
	return boolNatives
}

func (Bool) Type() String {
	return "builtin+bool"
}

func init() {
	boolNatives[NativeBool] = this
}
