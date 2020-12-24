package stackmachine

import "fmt"

type Function struct {
	Name string
	Call func(object ...Object)
}

func __println(object ...Object) {
	fmt.Println(object)
}
func __print(object ...Object) {
	fmt.Print(object)
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
}
