package qp

import (
	"fmt"
	"log"
	"reflect"
)

type Function interface {
	invoke(arguments ...Expression) (Expression, error)
}

var builtInFunctionMap = map[string]Function{
	"println": &println{},
}

type println struct {
}

func (println) invoke(arguments ...Expression) (Expression, error) {
	for _, argument := range arguments {
		object, err := argument.invoke()
		if err != nil {
			log.Panic(err.Error())
		}
	Loop:
		for {
			switch expression := object.(type) {
			case *Object:
				object, err = expression.invoke()
				if err != nil {
					return nil, err
				}
				if object == nil {
					log.Panic("expression", reflect.TypeOf(expression).String())
				}
				continue
			case *IntObject:
				fmt.Println("------>", expression.val)
				break Loop
			case *StringObject:
				fmt.Println("------>", expression.data)
				break Loop
			default:
				panic("unknown type" + reflect.TypeOf(object).String())
			}
		}
	}
	fmt.Println()
	return nil, nil
}
