package stackmachine

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Function struct {
	Name string
	Call func(object ...Element) []Element
}

func __println(object ...Element) []Element {
	for index, obj := range object {
		if index != 0 {
			fmt.Print(" ")
		}
		fmt.Print(obj)
	}
	fmt.Println()
	return nil
}

func __tLower_(object ...Element) []Element {
	o := strings.ToLower(object[0].Obj.(string))
	return []Element{
		{
			Obj:  o,
			Type: String,
		},
	}
}

func __print(object ...Element) []Element {
	fmt.Print(object)
	return nil
}

func __panic(object ...Element) []Element {
	log.Panicln(object)
	return nil
}

func __now__(object ...Element) []Element {
	now := time.Now()
	return []Element{
		{
			Type: Time,
			Obj:  now,
		},
	}
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
		Name: "now",
		Call: __now__,
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
