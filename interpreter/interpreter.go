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
		return R.invokeValue(nil, node)
	case ast.Assignment:
		val, err := R.Run(node.Value)
		if err != nil {
			return nil, err
		}
		v := node.Identifier.(ast.Identifier)
		R.CurrentScope().variables[v.Name] = val
	case ast.Identifier:
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
	case ast.FunctionDefinition:
		return &FuncDef{node.ArgList, node.Scope}, nil
	case ast.MemberSelector:
		v, err := R.Run(node.Object)
		if err != nil {
			return nil, err
		}
		switch prop := node.Property.(type) {
		case ast.Identifier, ast.MemberSelector:
			return R.GetMember(v, prop)
		case ast.Expression:
			return R.invokeExpression(v, prop)
		}
		return nil, nil
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
	case ast.Identifier:
		member, ok := v.Members()[prop.Name]
		if !ok {
			return nil, &Error{fmt.Errorf("Member: member %s doesnt exist", property)}
		}
		return member, nil
	case ast.MemberSelector:
		switch obj := prop.Object.(type) {
		case ast.Identifier:
			member, ok := v.Members()[obj.Name]
			if !ok {
				return nil, &Error{fmt.Errorf("Member: member %s doesnt exist", property)}
			}
			return R.GetMember(member, prop.Property)
		case ast.Expression:
			v, err := R.invokeValue(v, obj)
			if err != nil {
				return nil, err
			}
			return R.GetMember(v, prop.Property)
		}
	case ast.Expression:
		return R.invokeValue(v, prop)
	}
	return nil, &Error{fmt.Errorf("Invalid property")}
}

func (R *Runtime) invokeValue(v Value, node ast.Expression) (Value, *Error) {
	if v == nil {
		switch callee := node.Callee.(type) {
		case ast.Identifier:
			return R.invokeValue(R.GetVar(callee.Name), node)
		case ast.MemberSelector:
			obj, err := R.Run(callee.Object)
			if err != nil {
				return nil, err
			}
			member, err := R.GetMember(obj, callee.Property)
			if err != nil {
				return nil, err
			}
			return R.invokeValue(member, node)
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

func (R *Runtime) invokeExpression(v Value, node ast.Expression) (Value, *Error) {
	fun, ok := v.Members()[node.Callee.(ast.Identifier).Name]
	if !ok {
		return nil, &Error{fmt.Errorf("Invalid invocation field")}
	}
	return R.invokeValue(fun, node)
}
