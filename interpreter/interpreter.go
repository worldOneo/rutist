package interpreter

import (
	"fmt"

	"github.com/worldOneo/rutist/ast"
)

type Runtime struct {
	Scopes     []*Scope
	ScopeIndex int
}

func Run(ast ast.Node) (Value, error) {
	runtime := Runtime{
		[]*Scope{NewScope(map[string]Value{})},
		0,
	}
	val, err := runtime.Run(ast)
	if err == nil {
		return val, nil
	}
	return val, err.Err
}

func NewScope(capture map[string]Value) *Scope {
	scope := &Scope{
		variables: map[string]Value{},
	}
	for k, v := range capture {
		scope.variables[k] = v
	}
	return scope
}

func (R *Runtime) Run(program ast.Node) (Value, *Error) {
	switch node := program.(type) {
	case ast.Block:
		for i := 0; i < len(node.Body); i++ {
			_, err := R.Run(node.Body[i])
			if err != nil {
				return nil, err
			}
		}
	case ast.Expression:
		function, ok := R.CurrentScope().functions[node.Identifier]
		if !ok {
			function, ok = builtins[node.Identifier]
			if !ok {
				return nil, &Error{fmt.Errorf("Function %s doesnt exist ", node.Identifier)}
			}
		}
		args := make([]Value, 0)
		for _, arg := range node.ArgList {
			v, err := R.Run(arg)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		return R.CallFunction(function, args)
	case ast.Assignment:
		val, err := R.Run(node.Value)
		if err != nil {
			return nil, err
		}
		R.CurrentScope().variables[node.Identifier] = val
	case ast.Variable:
		return R.CurrentScope().variables[node.Name], nil
	case ast.Float:
		return Float(node.Value), nil
	case ast.Int:
		return Int(node.Value), nil
	case ast.Bool:
		return Bool(node.Value), nil
	case ast.String:
		return String(node.Value), nil
	case ast.Scope:
		return &Scoope{node}, nil
	}
	return nil, nil
}

const null = Int(0)

func (R *Runtime) GetVar(name string) Value {
	v, ok := R.CurrentScope().variables[name]
	if !ok {
		return null
	}
	return v
}

func (R *Runtime) CurrentScope() *Scope {
	return R.Scopes[R.ScopeIndex]
}

func (R *Runtime) CallFunction(function Function, args []Value) (Value, *Error) {
	R.ScopeIndex++
	if R.ScopeIndex >= len(R.Scopes) {
		R.Scopes = append(R.Scopes, nil)
	}
	R.Scopes[R.ScopeIndex] = NewScope(R.Scopes[R.ScopeIndex-1].variables)
	val, err := function(R, args)
	R.ScopeIndex--
	return val, err
}
