package qp

type Function interface {
	Expression
	call(arguments ...Expression) Expression
}

type BuiltInFunctionHandler struct {
	name     string
	callFunc CallFunc
}

type CallFunc func(arguments ...Expression) Expression
type funcObjectMap map[string]*Object
type registerBuiltInFuncHelper func(name string, callFunc CallFunc) registerBuiltInFuncHelper

func registerBuiltInFunc(funcObjectMap funcObjectMap, name string, callFunc CallFunc) registerBuiltInFuncHelper {
	var helper registerBuiltInFuncHelper
	helper = func(name string, callFunc CallFunc) registerBuiltInFuncHelper {
		funcObjectMap[name] = &Object{
			label: name,
			inner: &BuiltInFunctionHandler{
				name:     name,
				callFunc: callFunc,
			},
		}
		return helper
	}
	return helper(name, callFunc)
}

func (b *BuiltInFunctionHandler) call(arguments ...Expression) Expression {
	return b.callFunc(arguments...)
}

func (b *BuiltInFunctionHandler) String() string {
	return b.name
}

func (b *BuiltInFunctionHandler) Invoke() Expression {
	return b
}

func (b *BuiltInFunctionHandler) GetType() Type {
	return BuiltInFunctionType
}
