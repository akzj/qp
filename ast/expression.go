package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
	"reflect"
	"time"
)

type Expressions []runtime.Invokable

func (expressions Expressions) String() string {
	panic("implement me")
}

type ParenthesisExpression struct {
	Exp runtime.Invokable
}

func (p ParenthesisExpression) Invoke() runtime.Invokable {
	return p.Exp.Invoke()
}

func (p ParenthesisExpression) GetType() lexer.Type {
	return lexer.LeftParenthesisType
}

func (p ParenthesisExpression) String() string {
	return "(" + p.Exp.String() + ")"
}

func (Expressions) GetType() lexer.Type {
	return lexer.ExpressionType
}

func unwrapObject(expression runtime.Invokable) runtime.Invokable {
	for {
		if obj, ok := expression.(*runtime.Object); ok {
			expression = obj.Pointer
		} else {
			return expression
		}
	}
}

type BinaryOpExpression struct {
	OP    lexer.Type
	Left  runtime.Invokable
	Right runtime.Invokable
}

func (b BinaryOpExpression) String() string {
	return b.Left.String() + b.OP.String() + b.Right.String()
}

func (b BinaryOpExpression) Invoke() runtime.Invokable {
	var left = unwrapObject(b.Left.Invoke())
	var right = unwrapObject(b.Right.Invoke())
	switch lVal := left.(type) {
	case Int:
		switch rVal := right.(type) {
		case Int:
			switch b.OP {
			case lexer.AddType:
				return lVal + rVal
			case lexer.SubType:
				return lVal - rVal
			case lexer.MulOpType:
				return lVal * rVal
			case lexer.DivOpType:
				return lVal / rVal
			case lexer.LessType:
				return Bool(lVal < rVal)
			case lexer.LessEqualType:
				return Bool(lVal < rVal)
			case lexer.GreaterType:
				return Bool(lVal > rVal)
			case lexer.GreaterEqualType:
				return Bool(lVal >= rVal)
			case lexer.EqualType:
				return Bool(lVal == rVal)
			case lexer.NoEqualType:
				return Bool(lVal != rVal)
			}
		case NilObject:
			switch b.OP {
			case lexer.EqualType:
				return FalseObject
			case lexer.NoEqualType:
				return TrueObject
			}
		default:
			panic("no support type " + reflect.TypeOf(lVal).String() +
				"\n" + reflect.TypeOf(rVal).String() + " op type" + b.OP.String())
		}
	case Bool:
		switch rVal := right.(type) {
		case Bool:
			switch b.OP {
			case lexer.EqualType:
				return Bool(lVal == rVal)
			case lexer.NoEqualType:
				return !Bool(lVal == rVal)
			case lexer.OrType:
				return lVal || rVal
			case lexer.AndType:
				return lVal && rVal
			}
		}
	case TimeObject:
		switch rVal := right.(type) {
		case TimeObject:
			switch b.OP {
			case lexer.SubType:
				return DurationObject(time.Time(lVal).Sub(time.Time(rVal)))
			}
		}
	case *FuncStatement:
		switch right.(type) {
		case NilObject:
			switch b.OP {
			case lexer.EqualType:
				return FalseObject
			case lexer.NoEqualType:
				return TrueObject
			}
		}
	case *TypeObject:
		switch right.(type) {
		case NilObject:
			switch b.OP {
			case lexer.NoEqualType:
				return TrueObject
			case lexer.EqualType:
				return FalseObject
			}
		}
	case NilObject:
		switch right.(type) {
		case NilObject:
			switch b.OP {
			case lexer.EqualType:
				return TrueObject
			case lexer.NoEqualType:
				return FalseObject
			}
		}
	default:
		panic("no support type " + reflect.TypeOf(lVal).String() +
			"\n" + reflect.TypeOf(b.Right).String() + b.OP.String())
	}
	panic("no support type\n" + reflect.TypeOf(left).String() +
		"\n" + reflect.TypeOf(right).String() + "\n" + b.OP.String())
}

func (b BinaryOpExpression) GetType() lexer.Type {
	return b.OP
}

type NoStatement struct {
	Exp runtime.Invokable
}

func (n NoStatement) String() string {
	return "nop"
}

func (n NoStatement) Invoke() runtime.Invokable {
	return !n.Exp.Invoke().(Bool)
}

func (n NoStatement) GetType() lexer.Type {
	return lexer.NoType
}

func (expressions Expressions) Invoke() runtime.Invokable {
	var val runtime.Invokable
	for _, expression := range expressions {
		val = expression.Invoke()
		if _, ok := val.(ReturnStatement); ok {
			return val
		}
	}
	return val
}
