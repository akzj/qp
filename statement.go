package qp

import "fmt"

type Statements []Statement

func (Statements) getType() TokenType {
	return statementsTokenType
}

type Statement struct {
	expression Expression
}

func (s *Statement) getType() TokenType {
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

func (r ReturnStatement) invoke() (interface{}, error) {
	fmt.Println("ReturnStatement invoke")
	val, err := r.express.invoke()
	if err != nil {
		fmt.Println("invoke return statement failed")
	}
	if val == nil {
		fmt.Println("return val nil")
	}
	fmt.Println("return val", val)
	return val, nil
}

func (ReturnStatement) getType() TokenType {
	return returnTokenType
}

func (IfStatement) getType() TokenType {
	return ifTokenType
}

func (s *Statement) invoke() (interface{}, error) {
	fmt.Println("Statement invoke")
	if val, err := s.expression.invoke(); err != nil {
		fmt.Println("invoke expression failed")
		return nil, err
	} else {
		return val, nil
	}
}

func (s *Statements) invoke() (interface{}, error) {
	var val interface{}
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

func (ifStm *IfStatement) invoke() (interface{}, error) {
	fmt.Println("IfStatement invoke")
	check, err := ifStm.check.invoke()
	if err != nil {
		fmt.Println("IfStatement check error", err.Error())
		return nil, err
	}
	if check.(bool) {
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
			check, err := stm.check.invoke()
			if err != nil {
				fmt.Println("IfStatement check error", err.Error())
				return nil, err
			}
			if check.(bool) {
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
