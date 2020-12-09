package qp

import (
	"reflect"
)

type Expression interface {
	Invoke() Expression
	getType() Type
}

type Expressions []Expression

func (Expressions) getType() Type {
	return expressionType
}

func (expressions *Expressions) Invoke() Expression {
	var val Expression
	for _, expression := range *expressions {
		val = expression.Invoke()
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

func (expression *EqualExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val bool
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			val = *lVal == *rVal
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
	case *TypeObject:
		switch right.(type) {
		case *NilObject:
			val = false
		default:
			val = false
		}
	default:
		panic(reflect.TypeOf(left).String() + "\n" +
			reflect.TypeOf(right).String())
	}
	return (*BoolObject)(&val)
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

func (expression *GreaterExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val = false
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			val = *lVal > *rVal
			return (*BoolObject)(&val)
		}
	case *StringObject:
		switch rVal := right.(type) {
		case *StringObject:
			val = lVal.data > rVal.data
			return (*BoolObject)(&val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *GreaterEqualExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val = false
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			val = *lVal >= *rVal
			return (*BoolObject)(&val)
		}
	case *StringObject:
		switch rVal := right.(type) {
		case *StringObject:
			val = lVal.data >= rVal.data
			return (*BoolObject)(&val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *LessEqualExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val = false
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			val = *lVal <= *rVal
			return (*BoolObject)(&val)
		}
	case *StringObject:
		switch rVal := right.(type) {
		case *StringObject:
			val = lVal.data <= rVal.data
			return (*BoolObject)(&val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *LessExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	var val = false
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			val = *lVal < *rVal
			return (*BoolObject)(&val)
		}
	case *StringObject:
		switch rVal := right.(type) {
		case *StringObject:
			val = lVal.data < rVal.data
			return (*BoolObject)(&val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *MulExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			var val = (*lVal) * (*rVal)
			return (*IntObject)(&val)
		}
	}
	panic(reflect.TypeOf(left).String() + "\n" +
		reflect.TypeOf(right).String())
}

func (expression *AddExpression) Invoke() Expression {
	left := expression.Left.Invoke()
	right := expression.right.Invoke()
	switch lVal := left.(type) {
	case *IntObject:
		switch rVal := right.(type) {
		case *IntObject:
			val := *lVal + *rVal
			return (*IntObject)(&val)
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
