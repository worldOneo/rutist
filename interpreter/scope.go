package interpreter

type MemberDict = Map

type Value interface {
	Type() String
}

type Scope struct {
	variables map[string]Value
}
