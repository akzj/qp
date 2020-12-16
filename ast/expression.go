package ast

import (
	"fmt"
	"gitlab.com/akzj/qp/lexer"
	"reflect"
	"time"
)

type Expression interface {
	Invoke() Expression
	GetType() lexer.Type
	fmt.Stringer
}

type Expressions []Expression

func (expressions Expressions) String() string {
	panic("implement me")
}

type ParenthesisExpression struct {
	Exp Expression
}

func (p ParenthesisExpression) Invoke() Expression {
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

type AddExpression struct {
	Left  Expression
	right Expression
}

func (expression AddExpression) String() string {
	panic("implement me")
}

type SubExpression struct {
	Left  Expression
	right Expression
}

func (s SubExpression) String() string {
	panic("implement me")
}

type BinaryBoolExpression struct {
	OP    lexer.Type
	Left  Expression
	Right Expression
}

func (b BinaryBoolExpression) String() string {
	return "(" + b.Left.String() +
		b.OP.String() +
		b.Right.String() + ")"
}

func unwrapObject(expression Expression) Expression {
	for {
		if obj, ok := expression.(*Object); ok {
			expression = obj.Inner
		} else {
			return expression
		}
	}
}

func (b BinaryBoolExpression) Invoke() Expression {
	return BinaryOpExpression{
		Left:  b.Left,
		Right: b.Right,
		OP:    b.OP,
	}.Invoke()
}

func (b BinaryBoolExpression) GetType() lexer.Type {
	panic("implement me")
}

type BinaryOpExpression struct {
	OP    lexer.Type
	Left  Expression
	Right Expression
}

func (b BinaryOpExpression) String() string {
	return b.Left.String() + b.OP.String() + b.Right.String()
}

func (b BinaryOpExpression) Invoke() Expression {
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
	panic("implement me")
}

type MulExpression struct {
	Left  Expression
	right Expression
}

func (expression MulExpression) String() string {
	panic("implement me")
}

func (MulExpression) GetType() lexer.Type {
	return lexer.MulOpType
}

type LessExpression struct {
	Left  Expression
	right Expression
}

func (expression LessExpression) String() string {
	panic("implement me")
}

type LessEqualExpression struct {
	Left  Expression
	right Expression
}

func (expression LessEqualExpression) String() string {
	panic("implement me")
}

type GreaterExpression struct {
	Left  Expression
	right Expression
}

func (expression GreaterExpression) String() string {
	panic("implement me")
}

type GreaterEqualExpression struct {
	Left  Expression
	right Expression
}

func (expression GreaterEqualExpression) String() string {
	panic("implement me")
}

type EqualExpression struct {
	Left  Expression
	right Expression
}

func (expression EqualExpression) String() string {
	panic("implement me")
}

type NoEqualExpression struct {
	Left  Expression
	right Expression
}

func (n NoEqualExpression) String() string {
	panic("implement me")
}

type SelectorStatement struct {
	IDs []string
	vm  *VMContext
}

type NoStatement struct {
	Exp Expression
}

func (n NoStatement) String() string {
	panic("implement me")
}

func (n NoStatement) Invoke() Expression {
	return !n.Exp.Invoke().(Bool)
}

func (n NoStatement) GetType() lexer.Type {
	return lexer.NoType
}

func (s SelectorStatement) Invoke() Expression {
	panic("implement me")
}

func (s SelectorStatement) getType() lexer.Type {
	panic("implement me")
}

func (expressions Expressions) Invoke() Expression {
	var val Expression
	for _, expression := range expressions {
		val = expression.Invoke()
		if _, ok := val.(ReturnStatement); ok {
			return val
		}
	}
	return val
}

func (n NoEqualExpression) Invoke() Expression {
	equal := EqualExpression{
		Left:  n.Left,
		right: n.right,
	}
	val := !equal.Invoke().(Bool)
	return val
}

func (n NoEqualExpression) GetType() lexer.Type {
	return lexer.NoEqualType
}

func (expression EqualExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val bool
	switch lVal := left.(type) {
	case Int:
		switch rVal := right.(type) {
		case Int:
			val = lVal == rVal
		}
	case String:
		switch e := right.(type) {
		case String:
			val = lVal == e
		}
	case NilObject:
		switch right.(type) {
		case NilObject:
			val = true
		}
	case Bool:
		switch rVal := right.(type) {
		case Bool:
			val = lVal == rVal
		}
	case *TypeObject:
		switch right.(type) {
		case NilObject:
			val = false
		}
	case *FuncStatement:
		switch right.(type) {
		case NilObject:
			val = false
		}
	default:
		panic(reflect.TypeOf(left).String() + "\n" +
			reflect.TypeOf(right).String())
	}
	return Bool(val)
}

func (EqualExpression) GetType() lexer.Type {
	return lexer.EqualType
}

func (LessExpression) GetType() lexer.Type {
	return lexer.LessType
}

func (AddExpression) GetType() lexer.Type {
	return lexer.AddType
}

func (expression LessEqualExpression) GetType() lexer.Type {
	return lexer.LessEqualType
}

func (GreaterExpression) GetType() lexer.Type {
	return lexer.GreaterType
}

func (GreaterEqualExpression) GetType() lexer.Type {
	return lexer.GreaterEqualType
}
func (expression GreaterExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val = false
	switch lVal := left.(type) {
	case Int:
		switch rVal := right.(type) {
		case Int:
			val = lVal > rVal
			return (Bool)(val)
		}
	case String:
		switch rVal := right.(type) {
		case String:
			val = lVal > rVal
			return Bool(val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression GreaterEqualExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val = false
	switch lVal := left.(type) {
	case Int:
		switch rVal := right.(type) {
		case Int:
			val = lVal >= rVal
			return (Bool)(val)
		}
	case String:
		switch rVal := right.(type) {
		case String:
			val = lVal >= rVal
			return (Bool)(val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression LessEqualExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val = false
	switch lVal := left.(type) {
	case Int:
		switch rVal := right.(type) {
		case Int:
			val = lVal <= rVal
			return (Bool)(val)
		}
	case String:
		switch rVal := right.(type) {
		case String:
			val = lVal <= rVal
			return (Bool)(val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression LessExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val = false
	switch lVal := left.(type) {
	case Int:
		switch rVal := right.(type) {
		case Int:
			val = lVal < rVal
			return (Bool)(val)
		}
	case String:
		switch rVal := right.(type) {
		case String:
			val = lVal < rVal
			return (Bool)(val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression MulExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	switch lVal := left.(type) {
	case Int:
		switch rVal := right.(type) {
		case Int:
			var val = (lVal) * (rVal)
			return (Int)(val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (s SubExpression) Invoke() Expression {
	left := s.Left.Invoke()
	right := s.right.Invoke()
	switch lVal := left.(type) {
	case Int:
		switch rVal := right.(type) {
		case Int:
			val := lVal - rVal
			return (Int)(val)
		}
	case TimeObject:
		switch rVal := right.(type) {
		case TimeObject:
			val := time.Time(lVal).Sub(time.Time(rVal))
			return DurationObject(val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (s SubExpression) GetType() lexer.Type {
	return lexer.SubType
}

func (expression AddExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	if left.GetType() == lexer.IntType && right.GetType() == lexer.IntType {
		return left.(Int) + right.(Int)
	} else if left.GetType() == lexer.StringType && right.GetType() == lexer.StringType {
		return left.(String) + right.(String)
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}
