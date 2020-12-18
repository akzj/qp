package builtin

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)



type funcWrap struct {
	name     string
	callFunc CallFunc
}

type CallFunc func(arguments ...runtime.Invokable) runtime.Invokable
type FuncObjectMap map[string]*runtime.Object
type RegisterBuiltInFuncHelper func(name string, callFunc CallFunc) RegisterBuiltInFuncHelper

func register(funcObjectMap FuncObjectMap, name string, callFunc CallFunc) RegisterBuiltInFuncHelper {
	var helper RegisterBuiltInFuncHelper
	helper = func(name string, callFunc CallFunc) RegisterBuiltInFuncHelper {
		funcObjectMap[name] = &runtime.Object{
			Label: name,
			Pointer: &funcWrap{
				name:     name,
				callFunc: callFunc,
			},
		}
		return helper
	}
	return helper(name, callFunc)
}

func (b *funcWrap) Call(arguments ...runtime.Invokable) runtime.Invokable {
	return b.callFunc(arguments...)
}

func (b *funcWrap) String() string {
	return b.name
}

func (b *funcWrap) Invoke() runtime.Invokable {
	return b
}

func (b *funcWrap) GetType() lexer.Type {
	return lexer.BuiltInFunctionType
}
