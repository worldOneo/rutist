package interpreter

type MemberDict = Map

type Value interface {
	Type() String
	Natives() NativeMap
}

type Locals = map[string]Value

type Scope struct {
	variables Locals
}
