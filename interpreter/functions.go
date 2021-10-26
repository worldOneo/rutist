package interpreter

import (
	"fmt"

	"github.com/worldOneo/rutist/ast"
)

var builtins = make(map[string]Function)

func init() {
	builtins["print"] = Print
}

func Print(_ *Runtime, args []Value) (Value, error) {
	values := GoTypes(args)
	str, ok := values[0].(string)
	if !ok {
		return nil, fmt.Errorf("Print: Parameter1 must be type string")
	}
	fmt.Printf(str, values[1:]...)
	return nil, nil
}

func GoTypes(args []Value) []interface{} {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		switch value := arg.(type) {
		case ast.String:
			values[i] = value.Value
		case ast.Float:
			values[i] = value.Value
		case ast.Int:
			values[i] = value.Value
		case ast.Bool:
			values[i] = value.Value
		}
	}
	return values
}
