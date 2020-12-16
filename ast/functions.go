package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/vm"
)

type Function interface {
	Expression
	Call(arguments ...Expression) Expression
}

type BuiltInFunctionHandler struct {
	name     string
	callFunc CallFunc
}

type CallFunc func(arguments ...Expression) Expression
type FuncObjectMap map[string]*vm.Object
type RegisterBuiltInFuncHelper func(name string, callFunc CallFunc) RegisterBuiltInFuncHelper

func RegisterBuiltInFunc(funcObjectMap FuncObjectMap, name string, callFunc CallFunc) RegisterBuiltInFuncHelper {
	var helper RegisterBuiltInFuncHelper
	helper = func(name string, callFunc CallFunc) RegisterBuiltInFuncHelper {
		funcObjectMap[name] = &vm.Object{
			Label: name,
			Inner: &BuiltInFunctionHandler{
				name:     name,
				callFunc: callFunc,
			},
		}
		return helper
	}
	return helper(name, callFunc)
}

func (b *BuiltInFunctionHandler) Call(arguments ...Expression) Expression {
	return b.callFunc(arguments...)
}

func (b *BuiltInFunctionHandler) String() string {
	return b.name
}

func (b *BuiltInFunctionHandler) Invoke() Expression {
	return b
}

func (b *BuiltInFunctionHandler) GetType() lexer.Type {
	return lexer.BuiltInFunctionType
}
