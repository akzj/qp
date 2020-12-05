package qp

import "fmt"

type BuiltInFunction func(arguments ...Expression) (Expression, error)

var builtInFunctionMap = map[string]BuiltInFunction{
	"println": _println,
}

func _println(arguments ...Expression) (Expression, error) {
	fmt.Println("arguments size", len(arguments))
	for _, argument := range arguments {
		fmt.Println("argument type",argument.getType())
		object, err := argument.invoke()
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		switch object := object.(type) {
		case *IntObject:
			fmt.Print("->",object.val)
		}
	}
	fmt.Println()
	return nil, nil
}
