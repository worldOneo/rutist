package interpreter

type Value interface{}

type Function func(*Runtime, []Value) (Value, error)

type Scope struct {
	variables map[string]Value
	functions map[string]Function
}
