package qp

import (
	"fmt"
	"reflect"
)

type Statements []Statement

type Statement struct {
	expression Expression
}

type IfStatement struct {
	check            Expression
	statement        Statements
	elseIfStatements []IfStatement
	elseStatement    Statements
}

type ReturnStatement struct {
	express Expression
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

type VarAssignStatement struct {
	ctx        *VMContext
	label      string
	expression Expression
}

type IncFieldStatement struct {
	ctx   *VMContext
	label string
}

func (statement *IncFieldStatement) invoke() (Expression, error) {
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
	default:
		panic("unknown type " + reflect.TypeOf(innerObject).String())
	}
	return nil, nil
}

func (statement *IncFieldStatement) getType() Type {
	panic("implement me")
}

func (Statements) getType() Type {
	return statementsType
}

func (s *Statement) getType() Type {
	return statementType
}

func (f *FuncCallStatement) invoke() (Expression, error) {
	fmt.Println("FuncCallStatement invoke")
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
	fmt.Println("fieldStatement invoke")
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
	fmt.Println("VarStatement invoke")
	object := v.ctx.allocObject()
	object.label = v.label
	return nil, nil
}

func (v VarStatement) getType() Type {
	return varTokenType
}

func (v *VarAssignStatement) invoke() (Expression, error) {
	fmt.Println("VarAssignStatement invoke", v.label)
	obj, err := v.expression.invoke()
	if err != nil {
		return nil, err
	}
	object := v.ctx.allocObject()
	object.label = v.label
	object.inner = obj
	object.initType()
	return nil, nil
}

func (v *VarAssignStatement) getType() Type {
	return varAssignTokenType
}

func (r ReturnStatement) invoke() (Expression, error) {
	fmt.Println("ReturnStatement invoke")
	val, err := r.express.invoke()
	if err != nil {
		fmt.Println("invoke return statement failed")
	}
	if val == nil {
		fmt.Println("return nil error")
		return nil, fmt.Errorf("return expression nil")
	}
	fmt.Println("return val", val)
	return val, nil
}

func (ReturnStatement) getType() Type {
	return returnTokenType
}

func (IfStatement) getType() Type {
	return ifTokenType
}

func (s *Statement) invoke() (Expression, error) {
	fmt.Println("Statement invoke")
	if val, err := s.expression.invoke(); err != nil {
		fmt.Println("invoke expression failed")
		return nil, err
	} else {
		return val, nil
	}
}

func (s *Statements) invoke() (Expression, error) {
	var val Expression
	var err error
	for _, it := range *s {
		if val, err = it.invoke(); err != nil {
			return nil, err
		} else if val != nil {
			return val, nil
		}
	}
	return nil, nil
}

func (ifStm *IfStatement) invoke() (Expression, error) {
	fmt.Println("IfStatement invoke")
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
