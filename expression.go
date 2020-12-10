package qp

import (
	"reflect"
	"time"
)

type Expression interface {
	Invoke() Expression
	getType() Type
}

type Expressions []Expression

func (Expressions) getType() Type {
	return expressionType
}

type AddExpression struct {
	Left  Expression
	right Expression
}

type SubExpression struct {
	Left  Expression
	right Expression
}
type MulExpression struct {
	Left  Expression
	right Expression
}

func (MulExpression) getType() Type {
	return mulOperatorTokenType
}

type LessExpression struct {
	Left  Expression
	right Expression
}

type LessEqualExpression struct {
	Left  Expression
	right Expression
}

type GreaterExpression struct {
	Left  Expression
	right Expression
}

type GreaterEqualExpression struct {
	Left  Expression
	right Expression
}

type EqualExpression struct {
	Left  Expression
	right Expression
}

type NoEqualExpression struct {
	Left  Expression
	right Expression
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

func (n NoEqualExpression) getType() Type {
	return NoEqualTokenType
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

func (EqualExpression) getType() Type {
	return EqualTokenType
}

func (LessExpression) getType() Type {
	return lessTokenType
}

func (AddExpression) getType() Type {
	return addOperatorTokenType
}

func (expression LessEqualExpression) getType() Type {
	return lessEqualTokenType
}

func (GreaterExpression) getType() Type {
	return greaterTokenType
}

func (GreaterEqualExpression) getType() Type {
	return greaterEqualTokenType
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

func (s SubExpression) getType() Type {
	return subOperatorTokenType
}

func (expression AddExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	if left.getType() == IntType && right.getType() == IntType {
		return left.(Int) + right.(Int)
	} else if left.getType() == stringTokenType && right.getType() == stringTokenType {
		return left.(String) + right.(String)
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}
