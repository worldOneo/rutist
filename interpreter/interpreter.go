package interpreter

import (
	"fmt"

	"github.com/worldOneo/rutist/ast"
)

type Runtime struct {
	File          string
	Scopes        []*Scope
	ScopeIndex    int
	SpecialFields Map
}

const (
	SpecialfFieldExport = String("export")
)

func New(file string) *Runtime {
	return &Runtime{
		file,
		[]*Scope{NewScope(map[string]Value{})},
		0,
		map[Value]Value{
			SpecialfFieldExport: Dict{},
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
		left, err := R.Run(node.Left)
		if err != nil {
			return nil, err
		}
		right, err := R.Run(node.Right)
		if err != nil {
			return nil, err
		}
		operator := R.getNativeField(left, operatorMagicType[node.Operation])
		if operator == nil {
			return nil, R.error("Invalid operator", node)
		}
		fn, ok := operator.(Function)
		if !ok {
			return nil, R.error("Invalid operator", node)
		}
		return R.CallFunction(fn, []Value{left, right})
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
		return &FuncDef{[]ast.Identifier{}, node.Body}, nil
	case ast.FunctionDefinition:
		return &FuncDef{node.ArgList, node.Scope}, nil
	case ast.MemberSelector:
		return R.resolveMemberSelector(node)
	}
	return nil, nil
}

func (R *Runtime) GetVar(name string) Value {
	v, ok := R.CurrentScope().variables[name]
	if !ok {
		v, ok = builtins[name]
		if !ok {
			return nil
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

func (R *Runtime) getNativeField(v Value, field int) Value {
	n := v.Natives()[field]
	return n
}

func (R *Runtime) getDynamicMember(v Value, property Value) (Value, *Error) {
	getMember := R.getNativeField(v, NativeGetMember)
	if getMember == nil {
		return nil, nil
	}
	f, ok := getMember.(Function)
	if !ok {
		return nil, &Error{fmt.Errorf("Invalid member")}
	}
	return f(R, []Value{v, property})
}

func (R *Runtime) getMember(v Value, field String) (Value, *Error) {
	member, err := R.getDynamicMember(v, field)
	if err != nil {
		return nil, err
	}
	fd, ok := member.(*FuncDef)
	if ok {
		return wrappMemberFunction(v, fd.run), nil
	}
	fn, ok := member.(Function)
	if !ok {
		return member, nil
	}
	return wrappMemberFunction(v, fn), nil
}

func (R *Runtime) getMemberProperty(v Value, property ast.Node) (Value, *Error) {
	if v == nil {
		return nil, R.error("Member: value is nil", property)
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
		return R.invokeValue(member, []Value{})
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
	runnable := R.getNativeField(value, NativeRun)
	if runnable == nil {
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

func (R *Runtime) invokeValue(value Value, args []Value) (Value, *Error) {
	runnable := R.getNativeField(value, NativeRun)
	if runnable == nil {
		return nil, &Error{fmt.Errorf("Invalid invocation")}
	}
	function, ok := runnable.(Function)
	if !ok {
		return nil, &Error{fmt.Errorf("Invalid invocation")}
	}
	return R.invokeFunction(function, append([]Value{value}, args...))
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
			return R.assignObjectProperty(obj, val, v.Property)
		}
		return R.assignObject(obj, val, String(prop.Name))
	}
	return nil, R.error("Invalid assignment", node)
}

func (R *Runtime) assignObjectProperty(obj Value, val Value, node ast.Node) (Value, *Error) {
	switch v := node.(type) {
	case ast.Identifier:
		return R.assignObject(obj, val, String(v.Name))
	case ast.MemberSelector:
		key, ok := v.Object.(ast.Identifier)
		if !ok {
			return nil, R.error("Invalid assignment", node)
		}
		member, err := R.getDynamicMember(obj, String(key.Name))
		if err != nil {
			return nil, err
		}
		prop, ok := v.Property.(ast.Identifier)
		if !ok {
			return R.assignObjectProperty(member, val, v.Property)
		}
		return R.assignObject(member, val, String(prop.Name))
	}
	return nil, R.error("Invalid assignment", node)
}

func (R *Runtime) assignObject(obj Value, val Value, prop Value) (Value, *Error) {
	assign := R.getNativeField(obj, NativeSetMember)
	if assign == nil {
		return nil, &Error{fmt.Errorf("Invalid assignment")}
	}
	assignFn := R.getNativeField(assign, NativeRun)
	if assignFn == nil {
		return nil, &Error{fmt.Errorf("Invalid assignment")}
	}
	fn, ok := assignFn.(Function)
	if !ok {
		return nil, &Error{fmt.Errorf("Invalid assignment")}
	}
	return R.CallFunction(fn, []Value{assignFn, assign, obj, prop, val})
}

func (R *Runtime) error(msg string, node ast.Node) *Error {
	return &Error{fmt.Errorf("%s at %s:%d", msg, R.File, node.Token().Line)}
}
