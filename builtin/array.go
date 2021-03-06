package builtin

import (
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/runtime"
	"log"
)

func registerArrayFunction() {
	register(runtime.ArrayFunctions, "append", func(arguments ...runtime.Invokable) runtime.Invokable {
		array := arguments[0].Invoke().(*ast.Array)
		for _, exp := range arguments[1:] {
			array.Data = append(array.Data, exp.Invoke())
		}
		return array
	})("size", func(arguments ...runtime.Invokable) runtime.Invokable {
		return ast.Int(len(arguments[0].Invoke().(*ast.Array).Data))
	})("Get", func(arguments ...runtime.Invokable) runtime.Invokable {
		if len(arguments) != 2 {
			log.Panic("array Get() Arguments error")
		}
		array, ok := arguments[0].Invoke().(*ast.Array)
		if ok == false {
			log.Panic("Exp not array type")
		}
		i, ok := arguments[1].(ast.Int)
		if ok == false {
			log.Panic("is not array Arguments error")
		}
		if len(array.Data) <= int(i) {
			log.Panic("index out of range")
		}
		return array.Data[i]
	})
}

