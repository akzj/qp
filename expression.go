package qp

import "fmt"

type Expression interface {
	invoke() (interface{}, error)
	getType() TokenType
}

type Expressions []Expression

func (Expressions) getType() TokenType {
	return expressionTokenType
}

func (e *Expressions) invoke() (interface{}, error) {
	var val interface{}
	var err error
	for _, expression := range *e {
		if val, err = expression.invoke(); err != nil {
			return val, err
		}
	}
	return val, err
}

type IntExpression struct {
	val int64
}

type AddExpression struct {
	Left  Expression
	right Expression
}

type MulExpression struct {
	Left  Expression
	right Expression
}

func (MulExpression) getType() TokenType {
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

func (LessExpression) getType() TokenType {
	return lessTokenType
}

func (AddExpression) getType() TokenType {
	return addOperatorTokenType
}

func (expression LessEqualExpression) getType() TokenType {
	return lessEqualTokenType
}

func (IntExpression) getType() TokenType {
	return intTokenType
}

func (GreaterExpression) getType() TokenType {
	return greaterTokenType
}
func (GreaterEqualExpression) getType() TokenType {
	return greaterEqualTokenType
}

func (expression *GreaterExpression) invoke() (interface{}, error) {
	fmt.Println("GreaterExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	fmt.Println("l,r", l, r)
	switch lVal := l.(type) {
	case int64:
		switch rVal := r.(type) {
		case int64:
			fmt.Println(lVal, rVal)
			return lVal > rVal, nil
		}
	}
	return nil, fmt.Errorf("unknown operand type")
}

func (expression *GreaterEqualExpression) invoke() (interface{}, error) {
	fmt.Println("GreaterEqualExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	switch lVal := l.(type) {
	case int64:
		switch rVal := r.(type) {
		case int64:
			return lVal >= rVal, nil
		}
	}
	return nil, fmt.Errorf("unknown operand type")
}

func (expression *LessEqualExpression) invoke() (interface{}, error) {
	fmt.Println("LessEqualExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	switch lVal := l.(type) {
	case int64:
		switch rVal := r.(type) {
		case int64:
			return lVal <= rVal, nil
		}
	}
	return nil, fmt.Errorf("unknown operand type")
}

func (expression *LessExpression) invoke() (interface{}, error) {
	fmt.Println("LessExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	switch lVal := l.(type) {
	case int64:
		switch rVal := r.(type) {
		case int64:
			fmt.Println(lVal, rVal)
			return lVal < rVal, nil
		}
	}
	return nil, fmt.Errorf("unknown operand type")
}

func (i IntExpression) invoke() (interface{}, error) {
	return i.val, nil
}

func (expression *MulExpression) invoke() (interface{}, error) {
	fmt.Println("MulExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	switch lVal := l.(type) {
	case int64:
		switch rVal := r.(type) {
		case int64:
			return lVal * rVal, nil
		}
	}
	return nil, fmt.Errorf("unknown operand type")
}

func (a *AddExpression) invoke() (interface{}, error) {
	fmt.Println("AddExpression invoke")
	l, err := a.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	r, err := a.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	switch lVal := l.(type) {
	case int64:
		switch rVal := r.(type) {
		case int64:
			return lVal + rVal, nil
		}
	}
	return nil, fmt.Errorf("unknown operand type")
}
