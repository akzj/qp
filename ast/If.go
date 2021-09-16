package ast

import (
	"log"
	"reflect"

	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type IfExpression struct {
	VM         *runtime.VMRuntime
	Check      runtime.Invokable
	Statements Expressions
	ElseIf     []IfExpression
	Else       Expressions
}

func (exp IfExpression) String() string {
	codes := "if " + exp.Check.String() +
		" {\n\t" + exp.Statements.String() + "\n}"

	for _, elseif := range exp.ElseIf {
		codes += " else " + elseif.String()
	}
	if exp.Else != nil {
		codes += " else {\n\t" + exp.Else.String() + "\n}"
	}
	return codes
}
func (IfExpression) GetType() lexer.Type {
	return lexer.IfType
}

func (exp IfExpression) Invoke() runtime.Invokable {
	check := exp.Check.Invoke()
	if _, ok := check.(Bool); !ok {
		log.Panic("if Statements Check require boolObject", reflect.TypeOf(check).String())
	}
	if check.(Bool) {
		exp.VM.PushStackFrame(false) //make  if brock stack
		val := exp.Statements.Invoke()
		exp.VM.PopStackFrame() //release  if brock stack
		return val
	} else {
		for _, stm := range exp.ElseIf {
			elseIf := stm.Check.Invoke()
			if _, ok := elseIf.(Bool); !ok {
				log.Panicln("else if require bool result")
			}
			if elseIf.(Bool) {
				exp.VM.PushStackFrame(false) //make  if brock stack
				val := stm.Statements.Invoke()
				exp.VM.PopStackFrame() //release  if brock stack
				return val
			}
		}
		if exp.Else != nil {
			exp.VM.PushStackFrame(false) //make  brock stack
			val := exp.Else.Invoke()
			exp.VM.PopStackFrame() //release  if brock stack
			return val
		}
	}
	return nil
}
