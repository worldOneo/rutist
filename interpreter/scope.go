package interpreter

type Value interface {
	Members() map[string]Value
}

type Scope struct {
	variables map[string]Value
}
