package interpreter

import "github.com/worldOneo/rutist/ast"

type String string
type Int int
type Error struct {
	Err error
}

type Float float64
type Function func(*Runtime, []Value) (Value, *Error)
type Bool bool
type Scoope struct {
	node ast.Scope
}

const (
	TypeRun = "__run__"
)

func (S String) Members() map[string]Value {
	return map[string]Value{
		"len": Function(S.len),
	}
}

func (S String) len(_ *Runtime, _ []Value) (Value, *Error) {
	return Int(len(S)), nil
}

func (F Function) Members() map[string]Value {
	return map[string]Value{
		TypeRun: F,
	}
}

func (I Int) Members() map[string]Value {
	return map[string]Value{}
}

func (F Float) Members() map[string]Value {
	return map[string]Value{}
}

func (B Bool) Members() map[string]Value {
	return map[string]Value{}
}

func (S *Scoope) Members() map[string]Value {
	return map[string]Value{}
}

func (E *Error) Members() map[string]Value {
	return map[string]Value{}
}
