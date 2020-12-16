package qp

import (
	"fmt"
	"reflect"
	"time"
)

type Expression interface {
	Invoke() Expression
	GetType() Type
	fmt.Stringer
}

type Expressions []Expression

func (expressions Expressions) String() string {
	panic("implement me")
}

type ParenthesisExpression struct {
	exp Expression
}

func (p ParenthesisExpression) Invoke() Expression {
	return p.exp.Invoke()
}

func (p ParenthesisExpression) GetType() Type {
	return LeftParenthesisType
}

func (p ParenthesisExpression) String() string {
	return "(" + p.exp.String() + ")"
}

func (Expressions) GetType() Type {
	return ExpressionType
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
	opType Type
	Left   Expression
	right  Expression
}

func (b BinaryBoolExpression) String() string {
	return "(" + b.Left.String() +
		b.opType.String() +
		b.right.String() + ")"
}

func unwrapObject(expression Expression) Expression {
	for {
		if obj, ok := expression.(*Object); ok {
			expression = obj.inner
		} else {
			return expression
		}
	}
}

func (b BinaryBoolExpression) Invoke() Expression {
	return BinaryOpExpression{
		Left:   b.Left,
		right:  b.right,
		opType: b.opType,
	}.Invoke()
}

func (b BinaryBoolExpression) GetType() Type {
	panic("implement me")
}

type BinaryOpExpression struct {
	opType Type
	Left   Expression
	right  Expression
}

func (b BinaryOpExpression) String() string {
	return b.Left.String() + b.opType.String() + b.right.String()
}

func (b BinaryOpExpression) Invoke() Expression {
	var left = unwrapObject(b.Left.Invoke())
	var right = unwrapObject(b.right.Invoke())
	switch lVal := left.(type) {
	case Int:
		switch rVal := right.(type) {
		case Int:
			switch b.opType {
			case AddType:
				return lVal + rVal
			case SubType:
				return lVal - rVal
			case MulOpType:
				return lVal * rVal
			case DivOpType:
				return lVal / rVal
			case LessType:
				return Bool(lVal < rVal)
			case LessEqualType:
				return Bool(lVal < rVal)
			case GreaterType:
				return Bool(lVal > rVal)
			case GreaterEqualType:
				return Bool(lVal >= rVal)
			case EqualType:
				return Bool(lVal == rVal)
			case NoEqualType:
				return Bool(lVal != rVal)
			}
		default:
			panic("no support type " + reflect.TypeOf(lVal).String() +
				"\n" + reflect.TypeOf(rVal).String() + " op type" + b.opType.String())
		}
	case Bool:
		switch rVal := right.(type) {
		case Bool:
			switch b.opType {
			case EqualType:
				return Bool(lVal == rVal)
			case NoEqualType:
				return !Bool(lVal == rVal)
			}
		}
	case TimeObject:
		switch rVal := right.(type) {
		case TimeObject:
			switch b.opType {
			case SubType:
				return DurationObject(time.Time(lVal).Sub(time.Time(rVal)))
			}
		}
	case *FuncStatement:
		switch right.(type) {
		case NilObject:
			switch b.opType {
			case EqualType:
				return falseObject
			case NoEqualType:
				return trueObject
			}
		}
	case *TypeObject:
		switch right.(type) {
		case NilObject:
			switch b.opType {
			case NoEqualType:
				return trueObject
			case EqualType:
				return falseObject
			}
		}
	case NilObject:
		switch right.(type) {
		case NilObject:
			switch b.opType {
			case EqualType:
				return trueObject
			case NoEqualType:
				return falseObject
			}
		}
	default:
		panic("no support type " + reflect.TypeOf(lVal).String() +
			"\n" + reflect.TypeOf(b.right).String() + b.opType.String())
	}
	panic("no support type\n" + reflect.TypeOf(left).String() +
		"\n" + reflect.TypeOf(right).String() + "\n" + b.opType.String())
}

func (b BinaryOpExpression) GetType() Type {
	panic("implement me")
}

type MulExpression struct {
	Left  Expression
	right Expression
}

func (expression MulExpression) String() string {
	panic("implement me")
}

func (MulExpression) GetType() Type {
	return MulOpType
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
	exp Expression
}

func (n NoStatement) String() string {
	panic("implement me")
}

func (n NoStatement) Invoke() Expression {
	return !n.exp.Invoke().(Bool)
}

func (n NoStatement) GetType() Type {
	return NoType
}

func (s SelectorStatement) Invoke() Expression {
	panic("implement me")
}

func (s SelectorStatement) getType() Type {
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

func (n NoEqualExpression) GetType() Type {
	return NoEqualType
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

func (EqualExpression) GetType() Type {
	return EqualType
}

func (LessExpression) GetType() Type {
	return LessType
}

func (AddExpression) GetType() Type {
	return AddType
}

func (expression LessEqualExpression) GetType() Type {
	return LessEqualType
}

func (GreaterExpression) GetType() Type {
	return GreaterType
}

func (GreaterEqualExpression) GetType() Type {
	return GreaterEqualType
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

func (s SubExpression) GetType() Type {
	return SubType
}

func (expression AddExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	if left.GetType() == IntType && right.GetType() == IntType {
		return left.(Int) + right.(Int)
	} else if left.GetType() == StringType && right.GetType() == StringType {
		return left.(String) + right.(String)
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}
