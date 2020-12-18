package builtin

import (
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/runtime"
	"log"
	"reflect"
	"strings"
)

func registerStringFunction() {
	register(runtime.StringFunctions, "to_lower", func(arguments ...runtime.Invokable) runtime.Invokable {
		if len(arguments) > 1 {
			log.Panicln("only one Arguments")
		}
		for {
			switch inner := arguments[0].(type) {
			case ast.String:
				inner = ast.String(strings.ToLower(string(inner)))
				return inner
			default:
				log.Panicln("type error", reflect.TypeOf(arguments[0]).String())
			}
		}
	})
	register(runtime.StringFunctions,"clone", func(arguments ...runtime.Invokable) runtime.Invokable {
		if len(arguments) > 1 {
			log.Panicln("only one Arguments")
		}
		for {
			switch inner := arguments[0].(type) {
			case ast.String:
				return inner.Clone()
			default:
				log.Panicln("type error", reflect.TypeOf(arguments[0]).String())
			}
		}
	})
}

