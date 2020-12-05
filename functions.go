package qp

import (
	"fmt"
	"reflect"
)

type BuiltInFunction func(arguments ...Expression) (Expression, error)

var builtInFunctionMap = map[string]BuiltInFunction{
	"println": _println,
}

func _println(arguments ...Expression) (Expression, error) {
	fmt.Println("arguments size", len(arguments))
	for _, argument := range arguments {
		fmt.Println("argument type", argument.getType())
		object, err := argument.invoke()
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		if object == nil {
			panic(object)
		}
		switch expression := object.(type) {
		case *IntObject:
			fmt.Print("->", expression.val)
		default:
			panic("unknown type" + reflect.TypeOf(object).String())
		}
	}
	fmt.Println()
	return nil, nil
}
