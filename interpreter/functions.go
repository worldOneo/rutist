package interpreter

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/worldOneo/rutist/ast"
	"github.com/worldOneo/rutist/tokens"
)

var builtins = map[string]Function{}

func init() {
	builtins["print"] = builtinPrint
	builtins["try"] = builtinTrycatch
	builtins["throw"] = builtinThrow
	builtins["run"] = builtinRun
	builtins["str"] = builtinStr
	builtins["module"] = builtinModule
	builtins["import"] = builtinImport
	builtins["Map"] = func(r *Runtime, v []Value) (Value, *Error) { return Map{}, nil }
	builtins["Dict"] = func(r *Runtime, v []Value) (Value, *Error) { return Dict{}, nil }
}

func builtinImport(r *Runtime, args []Value) (Value, *Error) {
	if len(args) != 1 {
		return builtinThrow(r, []Value{String("Import: Requires exactly 1 parameter")})
	}

	str, err := builtinStr(r, args)
	if err != nil {
		return nil, err
	}
	fileVar, ok := str.(String)
	if !ok {
		return builtinThrow(r, []Value{String("Import: Arg1 must be string")})
	}
	file := string(fileVar)
	dir := filepath.Dir(r.File)
	if !filepath.IsAbs(file) {
		file = filepath.Join(dir, file)
	}
	content, e := ioutil.ReadFile(file)
	if e != nil {
		return nil, &Error{e}
	}
	code := string(content)
	tokens, e := tokens.Lexer(code)
	if err != nil {
		return nil, &Error{e}
	}
	parsed, e := ast.Parse(tokens)
	if err != nil {
		return nil, &Error{e}
	}

	runtime := New(file)
	_, err = runtime.Run(parsed)
	if err != nil {
		return nil, err
	}
	return runtime.SpecialFields[String(SpecialfFieldExport)], nil
}

func builtinModule(r *Runtime, args []Value) (Value, *Error) {
	if len(args) != 1 {
		return builtinThrow(r, []Value{String("Module: Requires exactly 1 parameter")})
	}
	f, ok := args[0].Members()[TypeRun]
	if !ok {
		return builtinThrow(r, []Value{String("Module: Parameter 1 must be runnable")})
	}
	fn, ok := f.(Function)
	if !ok {
		return builtinThrow(r, []Value{String("Module: Parameter 1 must be runnable")})
	}

	return r.CallFunction(fn, []Value{
		Function(func(r *Runtime, v []Value) (Value, *Error) {
			if len(v) != 2 {
				return builtinThrow(r, []Value{String("Module: Export requires exactly 2 parameters")})
			}

			r.SpecialFields[SpecialfFieldExport].(Dict)[v[0]] = v[1]
			return nil, nil
		}),
	})
}

func builtinStr(R *Runtime, args []Value) (Value, *Error) {
	if len(args) != 1 {
		return builtinThrow(R, []Value{String("Str: Require exactly 1 parameter")})
	}
	str, ok := args[0].Members()[TypeStr]
	if !ok {
		return String(fmt.Sprintf("%v", args[0])), nil
	}
	strFunc, funcOk := str.(Function)
	if !funcOk {
		return str, nil
	}
	return strFunc(R, args)
}

func builtinRun(R *Runtime, args []Value) (Value, *Error) {
	if len(args) != 1 {
		return builtinThrow(R, []Value{String("Run: Require exactly 1 parameter")})
	}

	scope, ok := args[0].(*Scoope)
	if !ok {
		return builtinThrow(R, []Value{String("Run: Parameter1 must be scope")})
	}
	return R.Run(scope.node)
}

func builtinThrow(_ *Runtime, args []Value) (Value, *Error) {
	if len(args) == 0 {
		return nil, &Error{fmt.Errorf("")}
	}
	arg := goNativeTypes(args)
	return nil, &Error{fmt.Errorf("%v", arg[0])}
}

func builtinTrycatch(r *Runtime, args []Value) (Value, *Error) {
	if len(args) == 0 {
		return nil, nil
	}

	try, ok := args[0].(*Scoope)
	if !ok {
		return nil, &Error{fmt.Errorf("Try: Parameter1 must be type scope")}
	}

	if len(args) == 1 {
		_, err := builtinRun(r, []Value{try})
		if err != nil {
			return err, nil
		}
		return nil, nil
	}
	catch, ok := args[1].(*FuncDef)
	if !ok {
		return nil, &Error{fmt.Errorf("Try-Catch: Parameter2 must be type funcdef")}
	}
	_, err := builtinRun(r, []Value{try})
	if err != nil {
		return catch.run(r, []Value{err})
	}
	return nil, nil
}

func builtinPrint(_ *Runtime, args []Value) (Value, *Error) {
	if len(args) == 0 {
		fmt.Println()
		return nil, nil
	}
	values := goNativeTypes(args)

	str, ok := values[0].(string)
	if !ok {
		fmt.Print(values...)
		return nil, nil
	}
	fmt.Printf(str, values[1:]...)
	return nil, nil
}

func goNativeTypes(args []Value) []interface{} {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		switch value := arg.(type) {
		case String:
			values[i] = string(value)
		case Float:
			values[i] = float64(value)
		case Int:
			values[i] = int(value)
		case Bool:
			values[i] = bool(value)
		case *Error:
			values[i] = value.Err
		default:
			values[i] = value
		}
	}
	return values
}
