package ast

import (
	"log"
	"reflect"

	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type CallStatement struct {
	ParentExp runtime.Invokable
	Function  runtime.Invokable
	Arguments Expressions
}

func (f *CallStatement) GetType() lexer.Type {
	return lexer.CallType
}

func (f *CallStatement) String() string {
	var str = f.Function.String() + "("
	for index, statement := range f.Arguments {
		if index != 0 {
			str += ","
		}
		str += statement.String()
	}
	return str + ") "
}

func (f *CallStatement) Invoke() runtime.Invokable {
	exp := f.Function.Invoke()
	switch obj := exp.(type) {
	case *runtime.Object:
		exp = obj.Invoke()
	case ReturnStatement:
		exp = obj.Val
	}
	if exp == nil {
		log.Panic("Function nil")
	}
	var arguments []runtime.Invokable
	if Func, ok := exp.(*FuncExpression); f.ParentExp != nil && (ok == false || Func.Closure == false) {
		switch argument := f.ParentExp.Invoke().(type) {
		case *runtime.Object:
			if argument.Pointer == nil {
				panic(argument.Label)
			}
			arguments = append(arguments, argument.Pointer)
		default:
			if argument == nil {
				panic("argument nil")
			}
			arguments = append(arguments, argument)
		}
	}

	if function, ok := exp.(Function); ok {
		for _, argument := range f.Arguments {
			switch job := argument.Invoke().(type) {
			case *runtime.Object:
				if job.Pointer == nil {
					panic(job.Label + " " + f.Function.String())
				}
				arguments = append(arguments, job.Pointer)
			default:
				if job == nil {
					panic("argument nil")
				}
				arguments = append(arguments, job)
			}
		}
		return function.Call(arguments...)
	}
	log.Panicf("Exp`%s` `%s` is no callable", exp.String(), reflect.TypeOf(exp).String())
	return nil
}
