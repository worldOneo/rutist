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
		return R.RunExpression(nil, node)
	case ast.Assignment:
		val, err := R.Run(node.Value)
		if err != nil {
			return nil, err
		}
		v := node.Identifier.(ast.Variable)
		R.CurrentScope().variables[v.Name] = val
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
		v, ok = builtins[name]
		if !ok {
			return null
		}
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

func (R *Runtime) GetMember(v Value, property ast.Node) (Value, *Error) {
	if v == nil {
		return nil, &Error{fmt.Errorf("Member: variable is nil")}
	}
	switch prop := property.(type) {
	case ast.Variable:
		member, ok := v.Members()[prop.Name]
		if !ok {
			return nil, &Error{fmt.Errorf("Member: member %s doesnt exist", property)}
		}
		return member, nil
	case ast.MemberSelector:
		member, ok := v.Members()[prop.Identifier]
		if !ok {
			return nil, &Error{fmt.Errorf("Member: member %s doesnt exist", property)}
		}
		return R.GetMember(member, prop.Property)
	}
	return nil, &Error{fmt.Errorf("Invalid property")}
}

func (R *Runtime) RunExpression(v Value, node ast.Expression) (Value, *Error) {
	if v == nil {
		switch callee := node.Callee.(type) {
		case ast.Variable:
			return R.RunExpression(R.GetVar(callee.Name), node)
		case ast.MemberSelector:
			member, err := R.GetMember(R.GetVar(callee.Identifier), callee.Property)
			if err != nil {
				return nil, err
			}
			return R.RunExpression(member, node)
		}
	}

	run, ok := v.Members()[TypeRun]
	if !ok {
		return nil, &Error{fmt.Errorf("Invalid invocation")}
	}
	fun, ok := run.(Function)
	if !ok {
		return nil, &Error{fmt.Errorf("Invalid invocation")}
	}

	args := make([]Value, 0)
	for _, arg := range node.ArgList {
		v, err := R.Run(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}
	return R.CallFunction(fun, args)
}
