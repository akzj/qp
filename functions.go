package qp

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

type Function interface {
	Expression
	call(arguments ...Expression) Expression
}

type printlnFunc struct{}

func (p printlnFunc) Invoke() Expression {
	return p
}

func (printlnFunc) getType() Type {
	return builtInFunctionType
}

func (printlnFunc) call(arguments ...Expression) Expression {
	for index, argument := range arguments {
		if stringer, ok := argument.(fmt.Stringer); ok {
			fmt.Print(stringer)
		} else {
			log.Panicf("unknown type `%s`", reflect.TypeOf(argument).String())
		}
		if index != len(arguments) -1{
			fmt.Print(" ")
		}
	}
	fmt.Println()
	return nil
}

type NowFunc struct {
}

func (n NowFunc) Invoke() Expression {
	return n
}

func (n NowFunc) getType() Type {
	return builtInFunctionType
}

func (n NowFunc) call(arguments ...Expression) Expression {
	return TimeObject(time.Now())
}
