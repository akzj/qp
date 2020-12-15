package qp

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

var (
	builtInFunctions       = map[string]*Object{}
	arrayBuiltInFunctions  = map[string]*Object{}
	stringBuiltInFunctions = map[string]*Object{}
)

func init() {
	registerArrayFunction()
	registerStringFunction()
	registerGlobalFunction()
}

func registerGlobalFunction() {
	registerBuiltInFunc(builtInFunctions, "println", func(arguments ...Expression) Expression {
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

	registerBuiltInFunc(builtInFunctions, "now", func(arguments ...Expression) Expression {
		return TimeObject(time.Now())
	})
}
