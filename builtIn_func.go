package qp

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

var (
	BuiltInFunctions       = map[string]*Object{}
	ArrayBuiltInFunctions  = map[string]*Object{}
	StringBuiltInFunctions = map[string]*Object{}
)

func init() {
	RegisterArrayFunction()
	RegisterStringFunction()
	RegisterGlobalFunction()
}

func RegisterGlobalFunction() {
	registerBuiltInFunc(BuiltInFunctions, "println", func(arguments ...Expression) Expression {
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

	registerBuiltInFunc(BuiltInFunctions, "now", func(arguments ...Expression) Expression {
		return TimeObject(time.Now())
	})
}
