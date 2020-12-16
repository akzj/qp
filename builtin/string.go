package builtin

import (
	"gitlab.com/akzj/qp/ast"
	"log"
	"reflect"
	"strings"
)

func registerStringFunction() {
	register(ast.StringFunctions, "to_lower", func(arguments ...ast.Expression) ast.Expression {
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
	register(ast.StringFunctions,"clone", func(arguments ...ast.Expression) ast.Expression {
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

