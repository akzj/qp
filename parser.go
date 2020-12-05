package qp

import (
	"bytes"
	"fmt"
	io "io"
	"strconv"
)

func precedence(tokenType Type) int {
	switch tokenType {
	case addOperatorTokenType, subOperatorTokenType:
		return 100
	case mulOperatorTokenType, divOperatorTokenType:
		return 200
	case lessTokenType, lessEqualTokenType, greaterTokenType, greaterEqualTokenType:
		return 1
	default:
		return 0
	}
}

//greater or equal
func precedenceGE(first, second Type) bool {
	return precedence(first)-precedence(second) >= 0
}

func isOperatorToken(token Token) bool {
	return token.typ >= addOperatorTokenType && token.typ <= greaterEqualTokenType
}

type parser struct {
	lexer *lexer
	vmCtx *VMContext
}

func newParser(reader io.Reader) *parser {
	return &parser{
		lexer: newLexer(reader),
		vmCtx: newVMContext(),
	}
}

func makeExpression(opToken Token, expressions *[]Expression) Expression {
	var expression Expression
	switch opToken.typ {
	case addOperatorTokenType:
		expression = &AddExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case mulOperatorTokenType:
		expression = &MulExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case lessTokenType:
		expression = &LessExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case lessEqualTokenType:
		expression = &LessEqualExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case greaterTokenType:
		expression = &GreaterExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case greaterEqualTokenType:
		expression = &GreaterEqualExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	default:
		panic(opToken)
	}
	return expression
}

func (p *parser) parse() Expression {
	var opStack []Token
	var expressions []Expression

	doMakeExpression := func() {
		for len(opStack) != 0 {
			express := makeExpression(opStack[len(opStack)-1], &expressions)
			expressions = append(expressions, express)
			opStack = opStack[:len(opStack)-1]
		}
	}
Loop:
	for p.lexer.finish() == false {
		token := p.lexer.peek()
		if token.typ == EOFTokenType {
			break
		}
		switch {
		case token.typ == intTokenType:
			val, err := strconv.ParseInt(string(token.val), 10, 64)
			if err != nil {
				fmt.Println("parse int failed", string(token.val))
				return nil
			}
			expressions = append(expressions, &IntObject{val: val})
		case token.typ == leftParenthesisTokenType:
			opStack = append(opStack, token)
		case token.typ == rightParenthesisTokenType:
			var find = false
			for _, opCode := range opStack {
				if opCode.typ == leftParenthesisTokenType {
					find = true
				}
			}
			//end of expression
			if find == false {
				fmt.Println("break Loop")
				break Loop
			}
			if find {
				for len(opStack) != 0 && opStack[len(opStack)-1].typ != leftParenthesisTokenType {
					express := makeExpression(opStack[len(opStack)-1], &expressions)
					expressions = append(expressions, express)
					opStack = opStack[:len(opStack)-1]
				}
				opStack = opStack[:len(opStack)-1]
			}
		case isOperatorToken(token):
			for len(opStack) != 0 &&
				isOperatorToken(opStack[len(opStack)-1]) &&
				precedenceGE(opStack[len(opStack)-1].typ, token.typ) {
				express := makeExpression(opStack[len(opStack)-1], &expressions)
				if express == nil {
					fmt.Println("make expression failed", opStack[len(opStack)-1])
					return nil
				}
				expressions = append(expressions, express)
				opStack = opStack[:len(opStack)-1]
			}
			opStack = append(opStack, token)

		case token.typ == ifTokenType:
			if len(expressions) != 0 {
				doMakeExpression()
			}
			fmt.Println("ifTokenType begin")
			p.lexer.next()
			expression := p.parse()
			if expression == nil {
				fmt.Println("parse expression failed")
				return nil
			}
			statements := p.parseStatement()
			if statements == nil {
				fmt.Println("parseStatement failed")
				return nil
			}
			ifStem := IfStatement{
				check:            expression,
				statement:        statements,
				elseIfStatements: nil,
				elseStatement:    nil,
			}
			expressions = append(expressions, &ifStem)
		case token.typ == elseTokenType:
			if len(expressions) == 0 {
				fmt.Println("no statement")
				return nil
			}
			ifStatement, ok := expressions[len(expressions)-1].(*IfStatement)
			if ok == false {
				fmt.Println("no find if statement for else")
				return nil
			}
			p.lexer.next()
			token = p.lexer.peek()
			if token.typ == ifTokenType {
				token.typ = elseifTokenType
				p.lexer.next()
				checkExp := p.parse()
				if checkExp == nil {
					fmt.Println("parse expression failed")
					return nil
				}
				statement := p.parseStatement()
				if statement == nil {
					fmt.Println("parseStatement failed")
					return nil
				}
				ifStatement.elseIfStatements = append(ifStatement.elseIfStatements,
					IfStatement{
						check:     checkExp,
						statement: statement,
					})
			} else {
				statement := p.parseStatement()
				if statement == nil {
					fmt.Println("parseStatement failed")
					return nil
				}
				ifStatement.elseStatement = statement
			}
		case token.typ == returnTokenType:
			doMakeExpression()
			statement := ReturnStatement{}
			p.lexer.next()
			expression := p.parse()
			if expression == nil {
				fmt.Println("parse return expression failed")
				return nil
			}
			statement.express = expression
			expressions = append(expressions, statement)
			continue
		case token.typ == varTokenType:
			doMakeExpression()
			//todo check in the same line
			p.lexer.next()
			token = p.lexer.peek()
			if token.typ != labelTokenType {
				fmt.Println("error ,expect label")
				return nil
			}
			var label = token.val
			p.lexer.next()
			token = p.lexer.peek()
			if token.typ == assignTokenType {
				p.lexer.next()
				token = p.lexer.peek()
				expression := p.parse()
				if expression == nil {
					fmt.Println("parse assign expression failed")
					return nil
				}
				expressionList := *expression.(*Expressions)
				expressions = append(expressions, &VarAssignStatement{
					ctx:        p.vmCtx,
					label:      label,
					expression: expressionList[0],
				})
				expressions = append(expressions, expressionList[1:]...)
			} else {
				expressions = append(expressions, &VarStatement{
					ctx:   p.vmCtx,
					label: label,
				})
			}
			continue
		case token.typ == leftBraceTokenType: //{ end of expression
			break Loop
		case token.typ == rightBraceTokenType: // } end of statement
			break Loop
		case token.typ == labelTokenType:
			label := token.val
			p.lexer.next()
			token = p.lexer.peek()
			//function call
			if token.typ == leftParenthesisTokenType {
				fmt.Println("find function call")
				expression := p.parseFunCallArguments()
				if expression == nil {
					fmt.Println("get argument failed")
					return nil
				}
				expressions = append(expressions, &FuncCallStatement{
					vm:        p.vmCtx,
					label:     label,
					arguments: expression})
			} else {
				expressions = append(expressions, &fieldStatement{
					ctx:   p.vmCtx,
					label: label,
				})
			}
			continue
		}
		p.lexer.next()
	}
	doMakeExpression()
	return (*Expressions)(&expressions)
}

func (p *parser) parseStatement() []Statement {
	var statement []Statement
	var leftBrace []Token
	if p.lexer.finish() == false {
		token := p.lexer.peek()
		if token.typ != leftBraceTokenType {
			fmt.Println("error ,expect { ")
		}
		leftBrace = append(leftBrace, token)
		p.lexer.next()
	}
	for p.lexer.finish() == false {
		expression := p.parse()
		if expression == nil {
			fmt.Println("parse expression failed")
			return nil
		}
		statement = append(statement, Statement{
			expression: expression,
		})
		token := p.lexer.peek()
		if token.typ == rightBraceTokenType {
			leftBrace = leftBrace[:len(leftBrace)-1]
		}
		if len(leftBrace) == 0 {
			break
		}
	}
	return statement
}

func (p *parser) parseFunCallArguments() Expressions {
	var expressions Expressions
	p.lexer.next()
	for p.lexer.finish() == false {
		expression := p.parse()
		if expression == nil {
			fmt.Println("parse expression failed")
			return nil
		}
		expressions = append(expressions, expression)
		token := p.lexer.peek()
		p.lexer.next()
		if token.typ == rightParenthesisTokenType {
			fmt.Println("end of arguments", len(expressions), expressions.getType())
			break
		}
	}
	return expressions
}

func Parse(data string) Expression {
	parser := newParser(bytes.NewReader([]byte(data)))
	return parser.parse()
}
