package builtin

import (
	"fmt"
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/runtime"
	"log"
	"reflect"
	"time"
)

func init() {
	registerArrayFunction()
	registerStringFunction()
	registerGlobalFunction()
}

func registerGlobalFunction() {
	register(runtime.Functions, "println", func(arguments ...runtime.Invokable) runtime.Invokable {
		for index, argument := range arguments {
			if argument == nil{
				panic("argument")
			}
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

	register(runtime.Functions, "now", func(arguments ...runtime.Invokable) runtime.Invokable {
		return ast.TimeObject(time.Now())
	})
}
