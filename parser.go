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
	history []Token
	tokens  []Token
	lexer   *lexer
	vmCtx   *VMContext
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

func (p *parser) nextToken() Token {
	if len(p.tokens) == 0 {
		for i := 0; i < 100; i++ {
			p.tokens = append(p.tokens, p.lexer.peek())
			p.lexer.next()
			if p.lexer.finish() {
				break
			}
		}
	}
	if len(p.tokens) == 0 {
		return emptyToken
	}
	token := p.tokens[0]
	p.tokens = p.tokens[1:]
	p.history = append(p.history, token)
	return token
}

func (p *parser) historyToken(index int) Token {
	index = len(p.history) - 1 - index
	if index >= 0 && index < len(p.history) {
		return p.history[index]
	}
	return emptyToken
}

func (p *parser) putToken(token Token) {
	p.tokens = append([]Token{token}, p.tokens...)
}

func (p *parser) parseExpression() Expression {
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
	for {
		token := p.nextToken()
		if token == emptyToken {
			break
		}
		fmt.Println("parseExpression", token)
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
				fmt.Println("break")
				p.putToken(token)
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
			//function call
		case token.typ == labelType:
			his := p.historyToken(1)
			fmt.Println(his, "historyToken")
			if isOperatorToken(his) ||
				his.typ == leftParenthesisTokenType ||
				his.typ == returnTokenType {
				label := token.val
				next := p.nextToken()
				//function call
				if next.typ == leftParenthesisTokenType {
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
					p.putToken(next)
					expressions = append(expressions, &fieldStatement{
						ctx:   p.vmCtx,
						label: label,
					})
				}
			} else {
				p.putToken(token)
				break Loop
			}
		default:
			p.putToken(token)
			break Loop
		}
	}
	doMakeExpression()
	return (*Expressions)(&expressions)
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
	for {
		token := p.nextToken()
		if token == emptyToken {
			break
		}
		fmt.Println(token)
		switch {
		case token.typ == ifTokenType:
			p.nextToken()
			expression := p.parseExpression()
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
			next := p.nextToken()
			if next.typ == ifTokenType {
				next.typ = elseifTokenType
				checkExp := p.parseExpression()
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
				p.putToken(next)
				statement := p.parseStatement()
				if statement == nil {
					fmt.Println("parseStatement failed")
					return nil
				}
				ifStatement.elseStatement = statement
			}
		case token.typ == forTokenType:
			if len(expressions) != 0 {
				doMakeExpression()
			}
			fmt.Println("forTokenType begin ----------")
			expression := p.parseForStatement()
			if expression == nil {
				fmt.Println("parse for statement failed")
				return nil
			}
			expressions = append(expressions, expression)
		case token.typ == returnTokenType:
			fmt.Println(token, "returnTokenType")
			doMakeExpression()
			statement := ReturnStatement{}
			expression := p.parseExpression()
			if expression == nil {
				fmt.Println("parse return expression failed")
				return nil
			}
			statement.express = expression
			expressions = append(expressions, statement)
			continue
		case token.typ == varTokenType:
			token = p.nextToken()
			if token.typ != labelType {
				fmt.Println("error ,expect label", token)
				return nil
			}
			var label = token.val
			token = p.nextToken()
			if token.typ == assignTokenType {
				fmt.Println(assignTokenType, "begin")
				expression := p.parseExpression()
				if expression == nil {
					fmt.Println("parse assign expression failed")
					return nil
				}
				fmt.Println("assignTokenType end")
				expressions = append(expressions, &VarAssignStatement{
					ctx:        p.vmCtx,
					label:      label,
					expression: expression,
				})
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
		case token.typ == semicolonTokenType: // ; end of expression
			break Loop
		case token.typ == commaTokenType: // end of expression
			p.nextToken()
			break Loop
		case token.typ == labelType:
			label := token.val
			next := p.nextToken()
			//function call
			if next.typ == leftParenthesisTokenType {
				fmt.Println("parseFunCallArguments")
				expression := p.parseFunCallArguments()
				if expression == nil {
					fmt.Println("get argument failed")
					return nil
				}
				expressions = append(expressions, &FuncCallStatement{
					vm:        p.vmCtx,
					label:     label,
					arguments: expression})
				// ++ increase
			} else if next.typ == incOperatorTokenType {
				p.nextToken()
				fmt.Println("incOperatorTokenType", label)
				expressions = append(expressions, &IncFieldStatement{
					ctx:   p.vmCtx,
					label: label,
				})
			} else if next.typ == assignTokenType {
				fmt.Println("assignTokenType")
				p.nextToken()
				expression := p.parseExpression()
				fmt.Println("assignTokenType end")
				if expression == nil {
					fmt.Println("get assign expression failed")
					return nil
				}
				expressions = append(expressions, &AssignStatement{
					ctx:        p.vmCtx,
					label:      label,
					expression: expression,
				})
			} else {
				p.putToken(next)
				expressions = append(expressions, &fieldStatement{
					ctx:   p.vmCtx,
					label: label,
				})
			}
			continue
		}
	}
	doMakeExpression()
	return (*Expressions)(&expressions)
}

func (p *parser) parseStatement() []Statement {
	var statement []Statement
	var leftBrace []Token

	token := p.nextToken()
	if token.typ != leftBraceTokenType {
		fmt.Println("error ,expect { ")
		return nil
	}
	leftBrace = append(leftBrace, token)

	for {
		expression := p.parse()
		if expression == nil {
			fmt.Println("parse expression failed")
			return nil
		}
		statement = append(statement, Statement{
			expression: expression,
		})
		token := p.nextToken()
		if token.typ == rightBraceTokenType {
			p.nextToken() //drop }
			if leftBrace = leftBrace[:len(leftBrace)-1]; len(leftBrace) == 0 {
				break
			}
		}
	}
	return statement
}

func (p *parser) parseFunCallArguments() Expressions {
	var expressions Expressions
	for {
		expression := p.parseExpression()
		if expression == nil {
			fmt.Println("parse expression failed")
			return nil
		}
		expressions = append(expressions, expression)
		token := p.nextToken()
		if token.typ == rightParenthesisTokenType {
			break
		}
	}
	return expressions
}

func (p *parser) parseForStatement() *ForStatement {
	var forStatement ForStatement
	token := p.nextToken()
	if token.typ == semicolonTokenType {
		forStatement.preStatement = &NopStatement{}
		fmt.Println("semicolonTokenType")
	} else if token.typ == leftBraceTokenType {
		forStatement.preStatement = &NopStatement{}
		forStatement.postStatement = &NopStatement{}
		forStatement.checkStatement = &BoolObject{val: true}
		statements := p.parseStatement()
		if statements == nil {
			fmt.Println("parse `for` statements failed")
			return nil
		}
		forStatement.statements = statements
		return &forStatement
	} else {
		expression := p.parse()
		if expression == nil {
			fmt.Println("parse `for` preStatement failed")
			return nil
		}
		if p.nextToken().typ != semicolonTokenType {
			fmt.Println("expect ; in `for` statement")
			return nil
		}
		p.nextToken()
		forStatement.preStatement = expression
	}
	token = p.nextToken()
	//check expression
	if token.typ == semicolonTokenType {
		forStatement.checkStatement = &BoolObject{val: true}
		fmt.Println("semicolonTokenType ccccccccccccccccccccc")
	} else {
		p.putToken(token)
		expression := p.parseExpression()
		if expression == nil {
			fmt.Println("parse `for` checkStatement failed")
			return nil
		}
		if p.nextToken().typ != semicolonTokenType {
			fmt.Println("expect ; in `for` check expression")
			return nil
		}
		forStatement.checkStatement = expression
	}

	token = p.nextToken()
	//post expression
	if token.typ == leftBraceTokenType {
		forStatement.postStatement = &NopStatement{}
		fmt.Println("leftBraceTokenType xxxxxxxxxxxxxx")
		p.putToken(token)
	} else {
		expression := p.parse()
		if expression == nil {
			fmt.Println("parse `for` checkStatement failed")
			return nil
		}

		if next := p.nextToken(); next.typ != leftBraceTokenType {
			fmt.Println("expect { in `for` post expression")
			return nil
		} else {
			p.putToken(next)
		}
		forStatement.postStatement = expression
	}
	// statements
	statements := p.parseStatement()
	if statements == nil {
		fmt.Println("parse `for` statements failed")
		return nil
	}
	forStatement.statements = statements
	return &forStatement
}

func Parse(data string) Expression {
	parser := newParser(bytes.NewReader([]byte(data)))
	return parser.parse()
}
