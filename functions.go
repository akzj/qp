package qp

import (
	"fmt"
	"reflect"
)

type Function func(arguments ...Expression) (Expression, error)

var builtInFunctionMap = map[string]Function{
	"println": _println,
}

func _println(arguments ...Expression) (Expression, error) {
	fmt.Println("println", len(arguments))
	for _, argument := range arguments {
		object, err := argument.invoke()
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		if object == nil {
			panic(object)
		}
	Loop:
		for {
			switch expression := object.(type) {
			case *Object:
				object, err = expression.invoke()
				if err != nil {
					return nil, err
				}
				if object == nil{
					fmt.Println("expression",reflect.TypeOf(expression).String())
					panic(object)
				}
				continue
			case *IntObject:
				fmt.Print("->", expression.val)
				break Loop
			default:
				panic("unknown type" + reflect.TypeOf(object).String())
			}
		}
	}
	fmt.Println()
	return nil, nil
}
