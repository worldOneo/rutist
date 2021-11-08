package interpreter

type MemberDict = Map

type Value interface {
	Members() MemberDict
}

type Scope struct {
	variables map[string]Value
}
