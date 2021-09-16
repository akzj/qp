package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type VarAssignStatement struct {
	Ctx  *runtime.VMRuntime //global or stack var
	Name string             //var Name : var a,`a` is the Name
	Exp  runtime.Invokable  // Init Exp : var a = 1+1
}

func (statement VarAssignStatement) String() string {
	return "var " + statement.Name + " = " + statement.Exp.String()
}

func (expression VarAssignStatement) Invoke() runtime.Invokable {
	obj := expression.Exp.Invoke()
	var object = expression.Ctx.AllocObject(expression.Name)
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

func (expression VarAssignStatement) GetType() lexer.Type {
	return lexer.VarAssignType
}
