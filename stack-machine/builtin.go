package stackmachine

import (
	"fmt"
	"log"
	"strings"
)

type Function struct {
	Name string
	Call func(object ...Object) []Object
}

func __println(object ...Object) []Object {
	for index, obj := range object {
		if index != 0 {
			fmt.Print(" ")
		}
		fmt.Print(obj)
	}
	fmt.Println()
	return nil
}

func __tLower_(object ...Object) []Object {
	o := strings.ToLower(object[0].Str)
	return []Object{
		{
			Str:   o,
			VType: String,
		},
	}
}

func __print(object ...Object) []Object {
	fmt.Print(object)
	return nil
}

func __panic(object ...Object) []Object {
	log.Panicln(object)
	return nil
}

var BuiltInFunctionsIndex = map[string]int64{}

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
	{
		Name: "string.to_lower",
		Call: __tLower_,
	},
}

func init() {
	for index, function := range BuiltInFunctions {
		BuiltInFunctionsIndex[function.Name] = int64(index)
	}
}
