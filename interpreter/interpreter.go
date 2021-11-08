package interpreter

import (
	"fmt"

	"github.com/worldOneo/rutist/ast"
)

type Runtime struct {
	Scopes        []*Scope
	ScopeIndex    int
	SpecialFields Map
}

const (
	SpecialfFieldExport = String("export")
)

func New() *Runtime {
	return &Runtime{
		[]*Scope{NewScope(map[string]Value{})},
		0,
		map[Value]Value{
			SpecialfFieldExport: Map{},
		},
	}
}

func Run(ast ast.Node) (Value, error) {
	runtime := New()
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
		var lastVal Value
		var err *Error
		for i := 0; i < len(node.Body); i++ {
			lastVal, err = R.Run(node.Body[i])
			if err != nil {
				return nil, err
			}
		}
		if lastVal != nil {
			return lastVal, nil
		}
	case ast.Expression:
		return R.invokeValue(nil, node)
	case ast.Assignment:
		val, err := R.Run(node.Value)
		if err != nil {
			return nil, err
		}
		return R.assignValue(val, node.Identifier)
	case ast.Identifier:
		return R.GetVar(node.Name), nil
	case ast.Float:
		return Float(node.Value), nil
	case ast.Int:
		return Int(node.Value), nil
	case ast.Bool:
		return Bool(node.Value), nil
	case ast.String:
		return String(node.Value), nil
	case ast.Scope:
		return &Scoope{node.Body}, nil
	case ast.FunctionDefinition:
		return &FuncDef{node.ArgList, node.Scope}, nil
	case ast.MemberSelector:
		return R.getMember(node)
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
		member, ok := v.Members()[String(prop.Name)]
		if !ok {
			return nil, &Error{fmt.Errorf("Member: member %s doesnt exist", property)}
		}
		return member, nil
	case ast.MemberSelector:
		switch obj := prop.Object.(type) {
		case ast.Identifier:
			member, ok := v.Members()[String(obj.Name)]
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
		default:
			v, err := R.Run(callee)
			if err != nil {
				return nil, err
			}
			return R.invokeValue(v, node)
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
	fun, ok := v.Members()[String(node.Callee.(ast.Identifier).Name)]
	if !ok {
		return nil, &Error{fmt.Errorf("Invalid invocation field")}
	}
	return R.invokeValue(fun, node)
}

func (R *Runtime) getMember(node ast.MemberSelector) (Value, *Error) {
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

func (R *Runtime) assignValue(val Value, node ast.Node) (Value, *Error) {
	switch v := node.(type) {
	case ast.Identifier:
		R.CurrentScope().variables[v.Name] = val
		return nil, nil
	case ast.MemberSelector:
		obj, err := R.Run(v.Object)
		if err != nil {
			return nil, err
		}
		prop, ok := v.Property.(ast.Identifier)
		if !ok {
			return R.assignValue(val, v.Property)
		}
		assign, ok := obj.Members()[TypeSetMember]
		err = R.error("Invalid assignment", node)
		if !ok {
			return nil, err
		}
		assignFn, ok := assign.Members()[TypeRun]
		if !ok {
			return nil, err
		}
		fn, ok := assignFn.(Function)
		if !ok {
			return nil, err
		}
		return R.CallFunction(fn, []Value{String(prop.Name), val})
	}
	return nil, R.error("Invalid assignment", node)
}

func (R *Runtime) error(msg string, node ast.Node) *Error {
	return &Error{fmt.Errorf(msg)}
}