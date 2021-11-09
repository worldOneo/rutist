package interpreter

import (
	"fmt"
	"strconv"

	"github.com/worldOneo/rutist/ast"
)

type Runtime struct {
	File          string
	Scopes        []*Scope
	ScopeIndex    int
	SpecialFields Map
}

const (
	SpecialfFieldExport     = String("export")
	SpecialFieldTypeMembers = String("typemembers")
)

func New(file string) *Runtime {
	return &Runtime{
		file,
		[]*Scope{NewScope(map[string]Value{})},
		0,
		map[Value]Value{
			SpecialfFieldExport: Dict{},
			SpecialFieldTypeMembers: Map{
				Map{}.Type(): Map{
					String("get"): Function(mapGet),
					String("set"): Function(mapSet),
					String("has"): Function(mapHas),
				},
				Dict{}.Type(): Map{
					TypeSetMember: Function(func(r *Runtime, v []Value) (Value, *Error) {
						v[0].(Dict)[v[1]] = v[2]
						return nil, nil
					}),
					TypeGetMember: Function(func(r *Runtime, v []Value) (Value, *Error) {
						return v[0].(Dict)[v[1]], nil
					}),
				},
				String("").Type(): Map{
					String("len"): Function(func(_ *Runtime, v []Value) (Value, *Error) { return Int(len(v[0].(String))), nil }),
					TypeStr:       Function(func(_ *Runtime, v []Value) (Value, *Error) { return v[0], nil }),
					TypeLen:       Function(func(_ *Runtime, v []Value) (Value, *Error) { return Int(len(v[0].(String))), nil }),
					TypeBool:      Function(func(_ *Runtime, v []Value) (Value, *Error) { return Bool(v[0].(String) != ""), nil }),
				},
				Int(1).Type(): Map{
					TypeBool: Function(func(r *Runtime, v []Value) (Value, *Error) { return Bool(v[0].(Int) != 0), nil }),
					TypeStr:  Function(func(r *Runtime, v []Value) (Value, *Error) { return String(strconv.Itoa(int(v[0].(Int)))), nil }),
				},
				Float(0).Type(): Map{
					TypeBool: Function(func(r *Runtime, v []Value) (Value, *Error) { return Bool(v[0].(Float) != 0.0), nil }),
					TypeStr: Function(func(r *Runtime, v []Value) (Value, *Error) {
						return String(strconv.FormatFloat(float64(v[0].(Float)), 'f', 7, 64)), nil
					}),
				},
				Bool(true).Type(): Map{},
				Function(builtinImport).Type(): Map{
					TypeRun: Function(func(r *Runtime, v []Value) (Value, *Error) { return v[0].(Function)(r, v[1:]) }),
				},
				FuncDef{}.Type(): Map{
					TypeRun: Function(func(r *Runtime, v []Value) (Value, *Error) {
						return v[0].(*FuncDef).run(r, v[1:])
					}),
				},
				Scoope{}.Type(): Map{
					TypeRun: Function(func(r *Runtime, v []Value) (Value, *Error) { return r.Run(v[0].(*Scoope).node) }),
				},
				WrappedFunction{}.Type(): Map{
					TypeRun: Function(func(r *Runtime, v []Value) (Value, *Error) {
						w := v[0].(WrappedFunction)
						args := append([]Value{w.this}, v[1:]...)
						return w.f(r, args)
					}),
				},
			},
		},
	}
}

func Run(file string, ast ast.Node) (Value, error) {
	runtime := New(file)
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
		return R.invokeExpression(node)
	case ast.Assignment:
		val, err := R.Run(node.Value)
		if err != nil {
			return nil, err
		}
		return R.assignValue(val, node.Identifier)
	case ast.BinaryExpression:
		left, err := R.Run(node.Right)
		if err != nil {
			return nil, err
		}
		right, err := R.Run(node.Right)
		if err != nil {
			return nil, err
		}
		operator, ok := R.getNativeField(left.Type(), operatorMagicType[node.Operation])
		if !ok {
			return nil, R.error("Invalid operator", node)
		}
		fn, ok := operator.(Function)
		if !ok {
			return nil, R.error("Invalid operator", node)
		}
		return R.CallFunction(fn, []Value{right, left})
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
		return R.resolveMemberSelector(node)
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

func (R *Runtime) getNativeFields(t String) Map {
	typeFields := R.SpecialFields[SpecialFieldTypeMembers].(Map)[t].(Map)
	return typeFields
}

func (R *Runtime) getNativeField(t String, field String) (Value, bool) {
	n, ok := R.getNativeFields(t)[field]
	return n, ok
}

func (R *Runtime) getDynamicMember(v Value, property Value) (Value, *Error) {
	getMember, ok := R.getNativeField(v.Type(), TypeGetMember)
	if !ok {
		return nil, nil
	}
	f, ok := getMember.(Function)
	return f(R, []Value{v, property})
}

func (R *Runtime) getMember(v Value, field String) (Value, *Error) {
	member, ok := R.getNativeField(v.Type(), field)
	var err *Error
	if !ok {
		member, err = R.getDynamicMember(v, field)
		if err != nil {
			return nil, err
		}
	}
	fn, ok := member.(Function)
	if !ok {
		return member, nil
	}
	return wrappMemberFunction(v, fn), nil
}

func (R *Runtime) getMemberProperty(v Value, property ast.Node) (Value, *Error) {
	if v == nil {
		return nil, R.error("Member: value is nil exist", property)
	}
	switch prop := property.(type) {
	case ast.Identifier:
		field := String(prop.Name)
		return R.getMember(v, field)
	case ast.Expression:
		member, err := R.getMemberProperty(v, prop.Callee)
		if err != nil {
			return nil, err
		}
		return R.invokeValue(member)
	case ast.MemberSelector:
		member, err := R.getMemberProperty(v, prop.Object)
		if err != nil {
			return nil, err
		}
		return R.getMemberProperty(member, prop.Property)
	}
	return nil, R.error("Invalid property", property)
}

func (R *Runtime) buildArgs(this Value, args []ast.Node) ([]Value, *Error) {
	valArgs := make([]Value, len(args)+1)
	valArgs[0] = this
	for i := 0; i < len(args); i++ {
		val, err := R.Run(args[i])
		if err != nil {
			return nil, err
		}
		valArgs[i+1] = val
	}
	return valArgs, nil
}

func (R *Runtime) invokeFunction(function Function, args []Value) (Value, *Error) {
	return R.CallFunction(function, args)
}

func (R *Runtime) invokeExpression(node ast.Expression) (Value, *Error) {
	value, err := R.Run(node.Callee)
	if err != nil {
		return nil, err
	}
	runnable, ok := R.getNativeField(value.Type(), TypeRun)
	if !ok {
		return nil, R.error("Invalid invocation", node)
	}
	function, ok := runnable.(Function)
	if !ok {
		return nil, R.error("Invalid invocation", node)
	}
	args, err := R.buildArgs(value, node.ArgList)
	if err != nil {
		return nil, err
	}
	return R.invokeFunction(function, args)
}

func (R *Runtime) invokeValue(value Value) (Value, *Error) {
	runnable, ok := R.getNativeField(value.Type(), TypeRun)
	if !ok {
		return nil, &Error{fmt.Errorf("Invalid invocation")}
	}
	function, ok := runnable.(Function)
	if !ok {
		return nil, &Error{fmt.Errorf("Invalid invocation")}
	}
	return R.invokeFunction(function, []Value{value})
}

func (R *Runtime) resolveMemberSelector(node ast.MemberSelector) (Value, *Error) {
	v, err := R.Run(node.Object)
	if err != nil {
		return nil, err
	}
	return R.getMemberProperty(v, node.Property)
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
		assign, err := R.getMember(obj, TypeSetMember)
		if err != nil {
			return nil, err
		}
		assignFn, ok := R.getNativeField(assign.Type(), TypeRun)
		if !ok {
			return nil, R.error("Invalid assignment", node)
		}
		fn, ok := assignFn.(Function)
		if !ok {
			return nil, R.error("Invalid assignment", node)
		}
		return R.CallFunction(fn, []Value{assign, String(prop.Name), val})
	}
	return nil, R.error("Invalid assignment", node)
}

func (R *Runtime) error(msg string, node ast.Node) *Error {
	return &Error{fmt.Errorf("%s at %s:%d", msg, R.File, node.Token().Line)}
}
