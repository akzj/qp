package qp

import (
	"reflect"
)

type Expression interface {
	invoke() Expression
	getType() Type
}

type Expressions []Expression

func (Expressions) getType() Type {
	return expressionType
}

func (expressions *Expressions) invoke() Expression {
	var val Expression
	for _, expression := range *expressions {
		val = expression.invoke()
		if _, ok := val.(*ReturnStatement); ok {
			return val
		}
	}
	return val
}

type AddExpression struct {
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

func (expression *EqualExpression) invoke() Expression {
	left := expression.Left.invoke()
	right := expression.right.invoke()
	var val bool
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			val = lVal.val == rVal.val
		default:
			val = false
		}
	case *StringObject:
		switch e := right.(type) {
		case *StringObject:
			val = lVal.data == e.data
		default:
			val = false
		}
	case *NilObject:
		switch right.(type) {
		case *NilObject:
			val = true
		default:
			val = false
		}
	default:
		panic(reflect.TypeOf(left).String() + "\n" +
			reflect.TypeOf(right).String())
	}
	return &BoolObject{val: val}
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

func (expression *GreaterExpression) invoke() Expression {
	left := expression.Left.invoke()
	right := expression.right.invoke()
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			return &BoolObject{val: lVal.val > rVal.val}
		}
	case *StringObject:
		switch rVal := right.(type) {
		case *StringObject:
			return &BoolObject{val: lVal.data > rVal.data}
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *GreaterEqualExpression) invoke() Expression {
	left := expression.Left.invoke()
	right := expression.right.invoke()
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			return &BoolObject{val: lVal.val >= rVal.val}
		}
	case *StringObject:
		switch rVal := right.(type) {
		case *StringObject:
			return &BoolObject{val: lVal.data >= rVal.data}
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *LessEqualExpression) invoke() Expression {
	left := expression.Left.invoke()
	right := expression.right.invoke()
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			return &BoolObject{val: lVal.val <= rVal.val}
		}
	case *StringObject:
		switch rVal := right.(type) {
		case *StringObject:
			return &BoolObject{val: lVal.data <= rVal.data}
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *LessExpression) invoke() Expression {
	left := expression.Left.invoke()
	right := expression.right.invoke()
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			return &BoolObject{val: lVal.val < rVal.val}
		}
	case *StringObject:
		switch rVal := right.(type) {
		case *StringObject:
			return &BoolObject{val: lVal.data < rVal.data}
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *MulExpression) invoke() Expression {
	left := expression.Left.invoke()
	right := expression.right.invoke()
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			return &IntObject{val: lVal.val * rVal.val}
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *AddExpression) invoke() Expression {
	left := expression.Left.invoke()
	right := expression.right.invoke()
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			return &IntObject{val: lVal.val + rVal.val}
		}
	case *StringObject:
		switch e := right.(type) {
		case *StringObject:
			return &StringObject{data: lVal.data + e.data}
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}
