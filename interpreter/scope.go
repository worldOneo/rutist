package interpreter

type MemberDict = map[string]Value

type Value interface {
	Members() MemberDict
}

type Scope struct {
	variables map[string]Value
}
