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



type println struct {
}

func (p *println) Invoke() Expression {
	return p
}

func (p *println) getType() Type {
	return TypeObjectType
}

func (println) call(arguments ...Expression) Expression {
	for _, argument := range arguments {
	Loop:
		for {
			switch expression := argument.(type) {
			case *IntObject,
				*StringObject,
				*NilObject:
				fmt.Print(expression)
				break Loop
			default:
				log.Panic("unknown type" + reflect.TypeOf(argument).String())
			}
		}
	}
	fmt.Println()
	return nil
}
