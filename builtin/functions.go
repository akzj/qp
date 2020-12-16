package builtin

import (
	"fmt"
	"gitlab.com/akzj/qp/ast"
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
	register(ast.Functions, "println", func(arguments ...ast.Expression) ast.Expression {
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

	register(ast.Functions, "now", func(arguments ...ast.Expression) ast.Expression {
		return ast.TimeObject(time.Now())
	})
}
