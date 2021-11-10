package interpreter

var errorNatives = NativeMap{}

func (Error) Type() String {
	return "buitin+error"
}

func (Error) Natives() NativeMap {
	return errorNatives
}