package interpreter

var errorNatives = NativeMap{}

type Error struct {
	Err error
}

func (Error) Type() String {
	return "buitin+error"
}

func (Error) Natives() NativeMap {
	return errorNatives
}