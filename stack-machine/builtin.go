package stackmachine

import (
	"fmt"
	"log"
)

type Function struct {
	Name string
	Call func(object ...Object)
}

func __println(object ...Object) {
	for _, obj := range object {
		fmt.Print(obj," ")
	}
	fmt.Println()
}

func __print(object ...Object) {
	fmt.Print(object)
}

func __panic(object ...Object) {
	log.Panicln(object)
}

var BuiltInFunctions = []Function{
	{
		Name: "println",
		Call: __println,
	},
	{
		Name: "print",
		Call: __print,
	},
	{
		Name: "panic",
		Call: __panic,
	},
}
