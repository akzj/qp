package builtin

import (
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/lexer"
)



type funcWrap struct {
	name     string
	callFunc CallFunc
}

type CallFunc func(arguments ...ast.Expression) ast.Expression
type FuncObjectMap map[string]*ast.Object
type RegisterBuiltInFuncHelper func(name string, callFunc CallFunc) RegisterBuiltInFuncHelper

func register(funcObjectMap FuncObjectMap, name string, callFunc CallFunc) RegisterBuiltInFuncHelper {
	var helper RegisterBuiltInFuncHelper
	helper = func(name string, callFunc CallFunc) RegisterBuiltInFuncHelper {
		funcObjectMap[name] = &ast.Object{
			Label: name,
			Inner: &funcWrap{
				name:     name,
				callFunc: callFunc,
			},
		}
		return helper
	}
	return helper(name, callFunc)
}

func (b *funcWrap) Call(arguments ...ast.Expression) ast.Expression {
	return b.callFunc(arguments...)
}

func (b *funcWrap) String() string {
	return b.name
}

func (b *funcWrap) Invoke() ast.Expression {
	return b
}

func (b *funcWrap) GetType() lexer.Type {
	return lexer.BuiltInFunctionType
}
