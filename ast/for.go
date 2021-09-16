package ast

import (
	"log"

	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type ForExpression struct {
	VM         *runtime.VMRuntime
	Pre        runtime.Invokable
	Check      runtime.Invokable
	Post       runtime.Invokable
	Statements Expressions
}

func (exp ForExpression) String() string {
	var codes = "for "
	codes += exp.Pre.String() + ";" + exp.Check.String() + ";" + exp.Post.String() + " {\n\t"
	codes += exp.Statements.String() + "\n}"
	return codes
}

func (exp ForExpression) Invoke() runtime.Invokable {

	//make stack frame
	exp.VM.PushStackFrame(false)
	exp.Pre.Invoke()
	for {
		val, ok := exp.Check.Invoke().(Bool)
		if !ok {
			log.Panic("for Check expect Bool")
		}
		if !val {
			exp.VM.PopStackFrame() //end of for
			return nil
		}
		exp.VM.PushStackFrame(false) //make stack frame for `{` brock
		for _, statement := range exp.Statements {
			val := statement.Invoke()
			if val == BreakObj {
				return nil
			}
			if _, ok := val.(ReturnStatement); ok {
				return val
			}
		}
		exp.VM.PopStackFrame()
		exp.Post.Invoke()
	}
}

func (f ForExpression) GetType() lexer.Type {
	return lexer.ForType
}
