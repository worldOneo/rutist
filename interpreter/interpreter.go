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
		[]*Scope{NewScope()},
		0,
	}
	return runtime.Run(ast)
}

func NewScope() *Scope {
	return &Scope{
		variables: map[string]Value{},
	}
}

func (R *Runtime) Run(program ast.Node) (Value, error) {
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
				return nil, fmt.Errorf("Function %s doesnt exist ", node.Identifier)
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
		return R.CallFunction(function, args), nil
	case ast.Assignment:
		val, err := R.Run(node.Value)
		if err != nil {
			return nil, err
		}
		R.CurrentScope().variables[node.Identifier] = val
	case ast.Variable:
		return R.CurrentScope().variables[node.Name], nil
	case ast.ValueFloat, ast.ValueInt, ast.ValueString:
		return node, nil
	}
	return nil, nil
}

func (R *Runtime) CurrentScope() *Scope {
	return R.Scopes[R.ScopeIndex]
}

func (R *Runtime) CallFunction(function Function, args []Value) Value {
	R.ScopeIndex++
	if R.ScopeIndex >= len(R.Scopes) {
		R.Scopes = append(R.Scopes, NewScope())
	}
	R.Scopes[R.ScopeIndex] = NewScope()
	val := function(args)
	R.ScopeIndex--
	return val
}
