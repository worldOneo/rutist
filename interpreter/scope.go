package interpreter

type Value interface{}

type Function func([]Value) Value

type Scope struct {
	variables map[string]Value
	functions map[string]Function
}
