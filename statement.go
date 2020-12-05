package qp

import "fmt"

type Statements []Statement

func (Statements) getType() Type {
	return statementsTokenType
}

type Statement struct {
	expression Expression
}

func (s *Statement) getType() Type {
	return statementTokenType
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

type fieldStatement struct {
	ctx   *VMContext
	label string
}

type FuncCallStatement struct {
	vm        *VMContext
	label     string
	arguments Expressions
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
	panic("implement me")
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
	return labelTokenType
}

//just new Object
type VarStatement struct {
	ctx   *VMContext
	label string
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

type VarAssignStatement struct {
	ctx        *VMContext
	label      string
	expression Expression
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
