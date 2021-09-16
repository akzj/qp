package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type VarStatement struct {
	VM    *runtime.VMRuntime
	Label string
	Exp   runtime.Invokable
}

func (v VarStatement) String() string {
	return "var " + v.Label + " = " + v.Exp.String()
}

func (v VarStatement) Invoke() runtime.Invokable {
	if v.Exp != nil {
		v.VM.AllocObject(v.Label).Pointer = v.Exp.Invoke()
	} else {
		v.VM.AllocObject(v.Label).Pointer = NilObj
	}
	return nil
}

func (v VarStatement) GetType() lexer.Type {
	return lexer.VarType
}
