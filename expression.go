package qp

import (
	"fmt"
	"reflect"
)

type Expression interface {
	invoke() (Expression, error)
	getType() Type
}

type Expressions []Expression

func (Expressions) getType() Type {
	return expressionTokenType
}

func (expressions *Expressions) invoke() (Expression, error) {
	var val Expression
	var err error
	for _, expression := range *expressions {
		if val, err = expression.invoke(); err != nil {
			return val, err
		} else if val != nil {
			return val, nil
		}
	}
	return val, err
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

func (expression *GreaterExpression) invoke() (Expression, error) {
	fmt.Println("GreaterExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	l, err = l.invoke()
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	r, err = r.invoke()
	var val bool
	switch lVal := l.(type) {
	case *IntObject:
		switch rVal := r.(type) {
		case *IntObject:
			val = lVal.val > rVal.val
		default:
			panic(reflect.TypeOf(l).String() + "\n" +
				reflect.TypeOf(r).String())
		}
	default:
		panic(reflect.TypeOf(l).String() + "\n" +
			reflect.TypeOf(r).String())
	}
	return &BoolObject{val: val}, nil
}

func (expression *GreaterEqualExpression) invoke() (Expression, error) {
	fmt.Println("GreaterEqualExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	l, err = l.invoke()
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	r, err = r.invoke()
	var val bool
	switch lVal := l.(type) {
	case *IntObject:
		switch rVal := r.(type) {
		case *IntObject:
			val = lVal.val >= rVal.val
		default:
			panic(reflect.TypeOf(l).String() + "\n" +
				reflect.TypeOf(r).String())
		}
	default:
		panic(reflect.TypeOf(l).String() + "\n" +
			reflect.TypeOf(r).String())
	}
	return &BoolObject{val: val}, nil
}

func (expression *LessEqualExpression) invoke() (Expression, error) {
	fmt.Println("LessEqualExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	l, err = l.invoke()
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	r, err = r.invoke()
	var val bool
	switch lVal := l.(type) {
	case *IntObject:
		switch rVal := r.(type) {
		case *IntObject:
			val = lVal.val <= rVal.val
		default:
			panic(reflect.TypeOf(l).String() + "\n" +
				reflect.TypeOf(r).String())
		}
	default:
		panic(reflect.TypeOf(l).String() + "\n" +
			reflect.TypeOf(r).String())
	}
	return &BoolObject{val: val}, nil
}

func (expression *LessExpression) invoke() (Expression, error) {
	fmt.Println("LessExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	l, err = l.invoke()
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	r, err = r.invoke()
	var val bool
	switch lVal := l.(type) {
	case *IntObject:
		switch rVal := r.(type) {
		case *IntObject:
			val = lVal.val < rVal.val
		default:
			panic(reflect.TypeOf(l).String() + "\n" +
				reflect.TypeOf(r).String())
		}
	default:
		panic(reflect.TypeOf(l).String() + "\n" +
			reflect.TypeOf(r).String())
	}
	return &BoolObject{val: val}, nil
}

func (expression *MulExpression) invoke() (Expression, error) {
	fmt.Println("MulExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	l, err = l.invoke()
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	r, err = r.invoke()
	var val int64
	switch lVal := l.(type) {
	case *IntObject:
		switch rVal := r.(type) {
		case *IntObject:
			val = lVal.val * rVal.val
		default:
			panic(reflect.TypeOf(l).String() + "\n" +
				reflect.TypeOf(r).String())
		}
	default:
		panic(reflect.TypeOf(l).String() + "\n" +
			reflect.TypeOf(r).String())
	}
	return &IntObject{val: val}, nil
}

func (expression *AddExpression) invoke() (Expression, error) {
	fmt.Println("AddExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	l, err = l.invoke()
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
		return nil, err
	}
	r, err = r.invoke()
	var val int64
	switch lVal := l.(type) {
	case *IntObject:
		switch rVal := r.(type) {
		case *IntObject:
			val = lVal.val + rVal.val
		default:
			panic(reflect.TypeOf(l).String() + "\n" +
				reflect.TypeOf(r).String())
		}
	default:
		panic(reflect.TypeOf(l).String() + "\n" +
			reflect.TypeOf(r).String())
	}
	return &IntObject{val: val}, nil
}
