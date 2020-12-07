package qp

import (
	"bytes"
	"fmt"
	io "io"
	"strconv"
	"strings"
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
	p.history = p.history[:len(p.history)-1]
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
		if token.typ == commentTokenType {
			continue
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
				his.typ == semicolonTokenType ||
				his.typ == commaTokenType ||
				his.typ == ifTokenType ||
				his.typ == elseifTokenType ||
				his.typ == leftParenthesisTokenType ||
				his.typ == assignTokenType ||
				his.typ == returnTokenType {
				label := token.val
				next := p.nextToken()
				//function call
				if next.typ == leftParenthesisTokenType {
					expression := p.parseFunctionCall(label)
					if expression == nil {
						fmt.Println("parseFunctionCall failed")
						return nil
					}
					expressions = append(expressions, expression)
				} else if next.typ == incOperatorTokenType {
					fmt.Println("incOperatorTokenType", label)
					expressions = append(expressions, &IncFieldStatement{
						ctx:   p.vmCtx,
						label: label,
					})
				} else if next.typ == periodTokenType { // label.
					expression := p.parsePeriodStatement(token.val)
					if expression == nil {
						fmt.Println("parsePeriodStatement failed")
						return nil
					}
					expressions = append(expressions, expression)
				} else if next.typ == leftBraceTokenType { //eg: User {
					fmt.Println(token.val, "{")
					statement := p.parseObjectStructInit(token.val)
					if statement == nil {
						fmt.Println("parseObjectStructInit failed")
						return nil
					}
					return statement
				} else {
					p.putToken(next)
					fmt.Println("getVarStatement")
					expressions = append(expressions, &getVarStatement{
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

func (p *parser) parse() Statements {
	var statements Statements
	statements = append(statements, &NopStatement{})
Loop:
	for {
		token := p.nextToken()
		if token == emptyToken {
			break
		}
		fmt.Println(token)
		if token.typ == commentTokenType {
			continue
		}
		switch {
		case token.typ == ifTokenType:
			expression := p.parseIfStatement()
			if expression == nil {
				fmt.Println("parseIfStatement failed")
				return nil
			}
			statements = append(statements, expression)
		case token.typ == forTokenType:
			fmt.Println("forTokenType begin ----------")
			expression := p.parseForStatement()
			if expression == nil {
				fmt.Println("parse for statement failed")
				return nil
			}
			statements = append(statements, expression)
		case token.typ == returnTokenType:
			fmt.Println(token, "returnTokenType")
			statement := ReturnStatement{}
			expression := p.parseExpression()
			if expression == nil {
				fmt.Println("parse return expression failed")
				return nil
			}
			statement.express = expression
			statements = append(statements, &statement)
			continue
		case token.typ == funcTokenType:
			if functionStatement := p.parseFunctionStatement(); functionStatement == nil {
				fmt.Println("parseFunctionStatement failed")
				return nil
			} else {
				p.vmCtx.addUserFunction(functionStatement)
			}
		case token.typ == typeTokenType:
			if structObject := p.parseStructObject(); structObject == nil {
				fmt.Println("parseStructObject failed")
				return nil
			} else {
				if err := p.vmCtx.addStructObject(structObject); err != nil {
					fmt.Println("addStructObject failed", err)
					return nil
				}
			}

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
				statements = append(statements, &VarAssignStatement{
					ctx:        p.vmCtx,
					label:      label,
					expression: expression,
				})
			} else {
				p.putToken(token)
				statements = append(statements, &VarStatement{
					ctx:   p.vmCtx,
					label: label,
				})
			}
			continue
		case token.typ == leftBraceTokenType: //{ end of expression
			p.putToken(token)
			break Loop
		case token.typ == rightBraceTokenType: // } end of statement
			p.putToken(token)
			break Loop
		case token.typ == semicolonTokenType: // ; end of expression
			p.putToken(token)
			break Loop
		case token.typ == commaTokenType: // end of expression
			p.putToken(token)
			break Loop
		case token.typ == labelType:
			label := token.val
			next := p.nextToken()
			//function call
			if next.typ == leftParenthesisTokenType {
				fmt.Println("parseFunctionCall")
				statement := p.parseFunctionCall(label)
				if statement == nil {
					fmt.Println("parseFunctionCall failed")
					return nil
				}
				statements = append(statements, statement)
			} else if next.typ == incOperatorTokenType {
				fmt.Println("incOperatorTokenType", label)
				statements = append(statements, &IncFieldStatement{
					ctx:   p.vmCtx,
					label: label,
				})
			} else if next.typ == periodTokenType { // label.
				statement := p.parsePeriodStatement(token.val)
				if statement == nil {
					fmt.Println("parsePeriodStatement failed")
					return nil
				}
				statements = append(statements, statement)
			} else if next.typ == assignTokenType {
				fmt.Println("assignTokenType")
				expression := p.parseExpression()
				fmt.Println("assignTokenType end")
				if expression == nil {
					fmt.Println("get assign expression failed")
					return nil
				}
				statements = append(statements, &AssignStatement{
					ctx:        p.vmCtx,
					label:      label,
					expression: expression,
				})
			} else {
				p.putToken(next)
				statements = append(statements, &getVarStatement{
					ctx:   p.vmCtx,
					label: label,
				})
			}
			continue
		}
	}
	return statements
}

func (p *parser) parseStatement() Statements {
	fmt.Println("parseStatement")
	var statements Statements
	var leftBrace []Token

	token := p.nextToken()
	if token.typ != leftBraceTokenType {
		fmt.Println("error ,expect { ")
		return nil
	}
	leftBrace = append(leftBrace, token)
	//check empty statement
	if token = p.nextToken(); token.typ == rightParenthesisTokenType {
		statements = append(statements, &NopStatement{})
		return statements
	}
	p.putToken(token)

	for {
		statement := p.parse()
		if statement == nil {
			fmt.Println("parse statement failed")
			return nil
		}
		statements = append(statements, statement...)
		token := p.nextToken()
		fmt.Println("nextToken", token)
		if token.typ == rightBraceTokenType {
			if leftBrace = leftBrace[:len(leftBrace)-1]; len(leftBrace) == 0 {
				fmt.Println("parseStatement break")
				break
			}
		}
	}
	return statements
}

func (p *parser) parseFunctionCall(label string) *FuncCallStatement {
	var statement FuncCallStatement
	statement.label = label
	statement.vm = p.vmCtx
	for {
		expression := p.parseExpression()
		if expression == nil {
			fmt.Println("parse expression failed")
			return nil
		}
		statement.arguments = append(statement.arguments, expression)
		token := p.nextToken()
		if token.typ == rightParenthesisTokenType {
			break
		} else if token.typ == commaTokenType {
			// next arguments
			continue
		} else {
			p.putToken(token)
		}
	}
	return &statement
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
		expression := p.parseExpression()
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
		p.putToken(token)
		expression := p.parseExpression()
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

func (p *parser) parseIfStatement() *IfStatement {
	ifStem := IfStatement{
		elseIfStatements: nil,
		elseStatement:    nil,
	}
	if ifStem.check = p.parseExpression(); ifStem.check == nil {
		fmt.Println("parse checkExpression failed")
		return nil
	}
	if ifStem.statement = p.parseStatement(); ifStem.statement == nil {
		fmt.Println("parseStatement failed")
		return nil
	}
	fmt.Println("parseIfStatement check else if")
	for {
		token := p.nextToken()
		fmt.Println(token)
		if token.typ == elseTokenType {
			next := p.nextToken()
			if next.typ == ifTokenType {
				token.typ = elseifTokenType
			} else {
				p.putToken(next)
			}
		}
		//check else or else if
		if token.typ == elseifTokenType {
			elseIfStem := IfStatement{
				elseIfStatements: nil,
				elseStatement:    nil,
			}
			if elseIfStem.check = p.parseExpression(); elseIfStem.check == nil {
				fmt.Println("parse checkExpression failed")
				return nil
			}
			if elseIfStem.statement = p.parseStatement(); elseIfStem.statement == nil {
				fmt.Println("parseStatement failed")
				return nil
			}
			elseIfStem.elseIfStatements = append(elseIfStem.elseIfStatements, &elseIfStem)
		} else if token.typ == elseTokenType {
			if ifStem.elseStatement = p.parseStatement(); ifStem.elseStatement == nil {
				fmt.Println("parse else statement failed")
				return nil
			}
			return &ifStem
		} else {
			p.putToken(token)
			return &ifStem
		}
	}
}

func (p *parser) parseFunctionStatement() *FuncStatement {
	fmt.Println("--parseFunctionStatement--")
	var functionStatement FuncStatement

	//function name
	token := p.nextToken()
	if token.typ != labelType {
		fmt.Println("function declare expect label here")
		return nil
	}
	functionStatement.vm = p.vmCtx
	functionStatement.label = token.val

	fmt.Println("function label", token.val)

	if p.nextToken().typ != leftParenthesisTokenType {
		fmt.Println("function declare require `(`,error")
		return nil
	}
	for {
		token := p.nextToken()
		if token.typ == rightParenthesisTokenType {
			// end of argument list
			break
		} else if token.typ == commaTokenType {
			//next argument
			continue
		} else if token.typ == labelType {
			functionStatement.arguments = append(functionStatement.arguments, token.val)
			fmt.Println("find argument", token.val)
		} else {
			fmt.Println("unknown argument token", token)
			return nil
		}
	}
	statement := p.parseStatement()
	if statement == nil {
		fmt.Println("parseForStatement for function failed")
		return nil
	}
	fmt.Println("--parseFunctionStatement end --- ")
	functionStatement.statements = statement
	return &functionStatement
}

func (p *parser) parseStructObject() *StructObject {
	var object = &StructObject{}
	token := p.nextToken()
	if token.typ != labelType {
		fmt.Println("expect label follow type")
		return nil
	}
	object.label = token.val
	object.vm = p.vmCtx
	if token = p.nextToken(); token.typ != structTokenType {
		fmt.Println("expect `struct` keyword ")
		return nil
	}
	statements := p.parseStatement()
	if statements == nil {
		fmt.Println("object struct parseStatement failed")
		return nil
	}
	return object
}

func (p *parser) parseObjectStructInit(label string) *StructObjectInitStatement {
	fmt.Println("parseObjectStructInit")
	var statement StructObjectInitStatement
	var leftBrace []int

	statement.label = label
	statement.vm = p.vmCtx
	leftBrace = append(leftBrace, 1)
	//check empty statement
	if token := p.nextToken(); token.typ == rightParenthesisTokenType {
		statement.initStatements = append(statement.initStatements, &NopStatement{})
		return &statement
	} else {
		p.putToken(token)
	}

	for {
		token := p.nextToken()
		if token.typ != labelType {
			fmt.Println("expect label,error", token)
			return nil
		}
		if next := p.nextToken(); next.typ != colonTokenType {
			fmt.Println("expect colon `:` ,error", token)
			return nil
		}
		expression := p.parseExpression()
		if expression == nil {
			fmt.Println("parse expression failed")
			return nil
		}
		statement.initStatements = append(statement.initStatements, &VarAssignStatement{
			ctx:        p.vmCtx,
			label:      token.val,
			expression: expression,
		})
		//check end
		token = p.nextToken()
		if token.typ == rightBraceTokenType {
			fmt.Println("struct object init end", label)
			break
		}
		p.putToken(token)
	}
	return &statement
}

// a.b.c = 1 // assign
// a.b.c()   // function call
// a.b.c + 1 // get val statement
func (p *parser) parsePeriodStatement(label string) Statement {
	fmt.Println("parsePeriodStatement", label)
	var labels = []string{label}
	token := p.nextToken()
	if token.typ != labelType {
		fmt.Println("expect label ", token)
		return nil
	}
	labels = append(labels, token.val)
	for {
		next := p.nextToken()
		if next.typ == periodTokenType {
			token = p.nextToken()
			continue
			// a.b.c(1) //function call
		} else if next.typ == leftParenthesisTokenType {
			statement := p.parseFunctionCall(strings.Join(labels, "."))
			if statement == nil {
				fmt.Println("parseFunctionCall failed")
				return nil
			}
			statement.getObject = &getStructObjectStatement{
				vmContext: p.vmCtx,
				labels:    labels,
			}
			return statement
		} else if next.typ == assignTokenType { // a.b = 1
			expression := p.parseExpression()
			if expression == nil {
				fmt.Println("parseExpression for assign statement failed", token)
				return nil
			}
			return &AssignStatement{
				ctx:   p.vmCtx,
				label: strings.Join(labels, "."),
				getObject: &getStructObjectStatement{
					vmContext: p.vmCtx,
					labels:    labels,
				},
				expression: expression,
			}
		} else {
			// a.b.c +  // expression
			// var c = a.b.c //end of statement
			fmt.Println("getStructObjectStatement", labels)
			p.putToken(next)
			return &getStructObjectStatement{
				vmContext: p.vmCtx,
				labels:    labels,
			}
		}
	}
}

func Parse(data string) Statements {
	parser := newParser(bytes.NewReader([]byte(data)))
	return parser.parse()
}
