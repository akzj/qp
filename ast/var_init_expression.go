package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type VarInitExpression struct {
	Ctx  *runtime.VMRuntime //global or stack var
	Name string             //var Name : var a,`a` is the Name
	Exp  runtime.Invokable  // Init Exp : a := 1+1
}

func (exp VarInitExpression) String() string {
	return exp.Name + " := " + exp.Exp.String()
}

func (exp VarInitExpression) Invoke() runtime.Invokable {
	obj := exp.Exp.Invoke()
	var object = exp.Ctx.AllocObject(exp.Name)
	if obj == nil {
		panic(obj)
	}
	switch obj := obj.(type) {
	case *runtime.Object:
		object.Pointer = obj.Pointer
	default:
		object.Pointer = obj
	}
	return nil
}

func (exp VarInitExpression) GetType() lexer.Type {
	return lexer.VarAssignType
}
