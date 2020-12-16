package qp

import (
	"log"
	"reflect"
	"strings"
)

type String string

func (s String) getObject(label string) *Object {
	return StringBuiltInFunctions[label]
}

func (s String) allocObject(label string) *Object {
	return StringBuiltInFunctions[label]
}

func (s String) addObject(k string, v *Object) {
	panic("implement me")
}

func (s String) String() string {
	return string(s)
}

func (s String) Invoke() Expression {
	return s
}

func (s String) clone() BaseObject {
	return String(string(s))
}

func (s String) GetType() Type {
	return StringType
}

func RegisterStringFunction() {
	registerBuiltInFunc(StringBuiltInFunctions, "to_lower", func(arguments ...Expression) Expression {
		if len(arguments) > 1 {
			log.Panicln("only one arguments")
		}
		for {
			switch inner := arguments[0].(type) {
			case String:
				inner = String(strings.ToLower(string(inner)))
				return inner
			default:
				log.Panicln("type error", reflect.TypeOf(arguments[0]).String())
			}
		}
	})

	registerBuiltInFunc(StringBuiltInFunctions, "clone", func(arguments ...Expression) Expression {
		if len(arguments) > 1 {
			log.Panicln("only one arguments")
		}
		for {
			switch inner := arguments[0].(type) {
			case String:
				return inner.clone()
			default:
				log.Panicln("type error", reflect.TypeOf(arguments[0]).String())
			}
		}
	})
}
