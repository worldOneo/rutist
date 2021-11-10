package interpreter

type MemberDict = Map

type Value interface {
	Type() String
	Natives() NativeMap
}

type Scope struct {
	variables map[string]Value
}
