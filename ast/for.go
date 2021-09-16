package ast

import (
	"log"

	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type ForStatement struct {
	VM         *runtime.VMContext
	Pre        runtime.Invokable
	Check      runtime.Invokable
	Post       runtime.Invokable
	Statements Statements
}

func (f ForStatement) String() string {
	var codes = "for "
	codes += f.Pre.String() + ";" + f.Check.String() + ";" + f.Post.String() + " {\n\t"
	codes += f.Statements.String() + "\n}"
	return codes
}

func (f ForStatement) Invoke() runtime.Invokable {
	f.VM.PushStackFrame(false) //make stack frame

	//make for brock stack
	f.Pre.Invoke()

	for {
		val, ok := f.Check.Invoke().(Bool)
		if !ok {
			log.Panic("for Check expect Bool")
		}
		if !val {
			f.VM.PopStackFrame() //end of for
			return nil
		}
		f.VM.PushStackFrame(false) //make stack frame for `{` brock
		for _, statement := range f.Statements {
			val := statement.Invoke()
			if val == BreakObj {
				return nil
			}
			if _, ok := val.(ReturnStatement); ok {
				return val
			}
		}
		f.VM.PopStackFrame()
		f.Post.Invoke()
	}
}

func (f ForStatement) GetType() lexer.Type {
	return lexer.ForType
}
