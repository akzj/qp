package qp

import (
	"fmt"
	"reflect"
)

type Statements []Statement

type Statement interface {
	Expression
}

type IfStatement struct {
	check            Expression
	statement        Statements
	elseIfStatements []*IfStatement
	elseStatement    Statements
}

type ReturnStatement struct {
	express   Expression
	returnVal Expression
}

//just new Object
type VarStatement struct {
	ctx   *VMContext
	label string
}

type fieldStatement struct {
	ctx   *VMContext
	label string
}

type FuncCallStatement struct {
	vm        *VMContext
	label     string
	arguments Expressions
}

type AssignStatement struct {
	ctx        *VMContext
	label      string
	expression Expression
}

type VarAssignStatement struct {
	ctx        *VMContext
	label      string
	expression Expression
}

type IncFieldStatement struct {
	ctx   *VMContext
	label string
}

type BreakStatement struct {
}

type NopStatement struct {
}

func (expression *AssignStatement) invoke() (Expression, error) {
	fmt.Println("AssignStatement")
	val, err := expression.expression.invoke()
	if err != nil {
		fmt.Println("AssignStatement .expression.invoke() failed", err.Error())
		return nil, err
	}
	fmt.Println(val.getType())
	fmt.Println(val.(*IntObject).val)
	object := expression.ctx.getObject(expression.label)
	if object == nil {
		fmt.Println("AssignStatement .expression.getObject failed", object.label)
		return nil, err
	}
	object.inner = val
	return nil, nil
}

func (expression *AssignStatement) getType() Type {
	return AssignStatementType
}

func (n *NopStatement) invoke() (Expression, error) {
	return nil, nil
}

func (n *NopStatement) getType() Type {
	return nopStatementType
}

type ForStatement struct {
	preStatement   Expression
	checkStatement Expression
	postStatement  Expression
	statements     Statements
}

func (f *ForStatement) invoke() (Expression, error) {
	val, err := f.preStatement.invoke()
	if err != nil {
		fmt.Println("for preStatement.invoke() error", err)
		return nil, err
	}
	if val != nil {
		fmt.Println("for preStatement.invoke() must nil")
		return nil, fmt.Errorf("for preStatement.invoke() must nil")
	}

	for ; ; {
		val, err := f.checkStatement.invoke()
		if err != nil {
			fmt.Println("for checkStatement.invoke() error", err)
			return nil, err
		}
		bObj, ok := val.(*BoolObject)
		if ok == false {
			fmt.Errorf("for checkStatement expect BoolObject")
			return nil, fmt.Errorf("for checkStatement expect BoolObject")
		}
		if bObj.val == false {
			return nil, nil
		}
		for _, statement := range f.statements {
			val, err := statement.invoke()
			if err != nil {
				fmt.Println("for checkStatement.invoke() error", err)
				return nil, err
			}
			if val == breakObject {
				fmt.Println("break from for")
				return nil, nil
			}
			if _, ok := val.(*ReturnStatement); ok {
				return val, nil
			}
		}
		if _, err = f.postStatement.invoke(); err != nil {
			fmt.Println("for postStatement.invoke() error", err)
			return nil, err
		}
	}
}

func (f *ForStatement) getType() Type {
	return forTokenType
}

func (statement *IncFieldStatement) invoke() (Expression, error) {
	fmt.Println("IncFieldStatement")
	object := statement.ctx.getObject(statement.label)
	if object == nil {
		return nil, fmt.Errorf("no find Object with label `%s`", statement.label)
	}
	innerObject, err := object.invoke()
	if err != nil {
		panic(err)
	}
	switch obj := innerObject.(type) {
	case *IntObject:
		obj.val++
		return object, nil
	default:
		panic("unknown type " + reflect.TypeOf(innerObject).String())
	}
	return nil, nil
}

func (statement *IncFieldStatement) getType() Type {
	return incOperatorTokenType
}

func (Statements) getType() Type {
	return statementsType
}

func (f *FuncCallStatement) invoke() (Expression, error) {
	//fmt.Println("FuncCallStatement invoke")
	function, ok := builtInFunctionMap[f.label]
	if ok {
		return function(f.arguments...)
	}
	return nil, fmt.Errorf("no find function")
}

func (f *FuncCallStatement) getType() Type {
	return funcTokenType
}

func (f *fieldStatement) invoke() (Expression, error) {
	//fmt.Println("fieldStatement invoke")
	object := f.ctx.getObject(f.label)
	if object == nil {
		return nil, fmt.Errorf("no find Object with label `%s`", f.label)
	}
	return object.invoke()
}

func (f *fieldStatement) getType() Type {
	return labelType
}

func (v *VarStatement) invoke() (Expression, error) {
	//fmt.Println("VarStatement invoke")
	v.ctx.allocObject(v.label)
	return nil, nil
}

func (v VarStatement) getType() Type {
	return varTokenType
}

func (expression *VarAssignStatement) invoke() (Expression, error) {
	fmt.Println("VarAssignStatement invoke", expression.label)
	obj, err := expression.expression.invoke()
	if err != nil {
		return nil, err
	}
	fmt.Println(obj.getType())
	object := expression.ctx.allocObject(expression.label)
	object.inner = obj
	object.initType()
	return nil, nil
}

func (expression *VarAssignStatement) getType() Type {
	return varAssignTokenType
}

func (r *ReturnStatement) invoke() (Expression, error) {
	//fmt.Println("ReturnStatement invoke")
	if r.returnVal != nil {
		return r, nil
	}
	val, err := r.express.invoke()
	if err != nil {
		fmt.Println("invoke return statement failed")
		return nil, err
	}
	if val == nil {
		fmt.Println("return nil error")
		return nil, fmt.Errorf("return expression nil")
	}
	fmt.Println("return val", val)
	r.returnVal = val
	return r, nil
}

func (ReturnStatement) getType() Type {
	return returnTokenType
}

func (IfStatement) getType() Type {
	return ifTokenType
}

func (statements Statements) invoke() (Expression, error) {
	var val Expression
	var err error
	for _, statement := range statements {
		fmt.Println("statement type", statement.getType())
		if val, err = statement.invoke(); err != nil {
			return nil, err
		} else if val != nil {
			if _, ok := val.(*ReturnStatement); ok {
				return val, nil
			} else {
				fmt.Println("---------", val.getType())
			}
		}
	}
	return nil, nil
}

func (ifStm *IfStatement) invoke() (Expression, error) {
	//fmt.Println("IfStatement invoke")
	check, err := ifStm.check.invoke()
	if err != nil {
		fmt.Println("IfStatement check error", err.Error())
		return nil, err
	}
	if check.(*BoolObject).val {
		fmt.Println("true")
		val, err := ifStm.statement.invoke()
		if err != nil {
			fmt.Println("IfStatement statement error", err.Error())
			return nil, err
		}
		fmt.Println(ifStm.statement.getType())
		return val, nil
	} else {
		fmt.Println(ifStm.check)
		for _, stm := range ifStm.elseIfStatements {
			check2, err := stm.check.invoke()
			if err != nil {
				fmt.Println("IfStatement check error", err.Error())
				return nil, err
			}
			if check2.(*BoolObject).val {
				fmt.Println("else if true")
				val, err := stm.statement.invoke()
				if err != nil {
					fmt.Println("IfStatement statement error", err.Error())
					return nil, err
				}
				return val, nil
			} else {
				fmt.Println("false")
			}
		}
		if ifStm.elseStatement != nil {
			fmt.Println("else")
			val, err := ifStm.elseStatement.invoke()
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			return val, nil
		}
	}
	return nil, err
}
