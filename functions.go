package qp

import (
	"fmt"
	"log"
	"reflect"
)

type Function interface {
	Expression
	call(arguments ...Expression) Expression
}

var builtInFunctionMap = map[string]Function{
	"println": &println{},
}

type println struct {
}

func (p *println) invoke() Expression {
	return p
}

func (p *println) getType() Type {
	return TypeObjectType
}

func (println) call(arguments ...Expression) Expression {
	for _, argument := range arguments {
		object := argument.invoke()
	Loop:
		for {
			switch expression := object.(type) {
			case *IntObject,
				*StringObject,
				*NilObject:
				fmt.Print(expression)
				break Loop
			default:
				log.Panic("unknown type" + reflect.TypeOf(object).String())
			}
		}
	}
	fmt.Println()
	return nil
}
