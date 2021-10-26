package interpreter

import (
	"fmt"

	"github.com/worldOneo/rutist/ast"
)

var builtins = make(map[string]Function)

func init() {
	builtins["print"] = Print
}

func Print(args []Value) Value {
	values := GoTypes(args)
	fmt.Printf(values[0].(string), values[1:]...)
	return nil
}

func GoTypes(args []Value) []interface{} {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		switch value := arg.(type) {
		case ast.ValueString:
			values[i] = value.Value
		case ast.ValueFloat:
			values[i] = value.Value
		case ast.ValueInt:
			values[i] = value.Value
		}
	}
	return values
}
