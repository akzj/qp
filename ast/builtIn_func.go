package ast

import (
	"fmt"
	"gitlab.com/akzj/qp/vm"
	"log"
	"reflect"
	"time"
)

var (
	Functions       = map[string]*vm.Object{}
	ArrayFunctions  = map[string]*vm.Object{}
	StringFunctions = map[string]*vm.Object{}
)

func init() {
	RegisterArrayFunction()
	RegisterStringFunction()
	RegisterGlobalFunction()
}

func RegisterGlobalFunction() {
	RegisterBuiltInFunc(Functions, "println", func(arguments ...Expression) Expression {
		for index, argument := range arguments {
			if stringer, ok := argument.(fmt.Stringer); ok {
				fmt.Print(stringer)
			} else {
				log.Panicf("unknown type `%s`", reflect.TypeOf(argument).String())
			}
			if index != len(arguments)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
		return nil
	})

	RegisterBuiltInFunc(Functions, "now", func(arguments ...Expression) Expression {
		return TimeObject(time.Now())
	})
}
