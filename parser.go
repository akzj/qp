package qp

import (
	"bytes"
	"fmt"
	io "io"
	"log"
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

type closureCheck struct {
	vars     map[string]bool
	closures []string
}

func newClosureCheck() *closureCheck {
	return &closureCheck{
		vars:     map[string]bool{},
		closures: nil,
	}
}

func (c *closureCheck) addVar(label string) {
	c.vars[label] = true
}
func (c *closureCheck) visit(label string) bool {
	var closure bool
	if c.vars[label] == false {
		c.closures = append(c.closures, label)
		closure = true
	}
	c.addVar(label)
	return closure
}

type parser struct {
	history      []Token
	tokens       []Token
	lexer        *lexer
	vmCtx        *VMContext
	closureCheck []*closureCheck
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

func (p *parser) nextToken(skipComment bool) Token {
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
	if token.typ == commentTokenType && skipComment {
		return p.nextToken(skipComment)
	}
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
		token := p.nextToken(false)
		if token == emptyToken {
			break
		}
		fmt.Println("parseExpression", token)
		if token.typ == commentTokenType {
			continue
		}
		switch {
		case token.typ == intTokenType:
			val, err := strconv.ParseInt(string(token.data), 10, 64)
			if err != nil {
				fmt.Println("parse int failed", string(token.data))
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
			//lambda func
		case token.typ == funcTokenType:
			last := p.historyToken(1)
			if isOperatorToken(last) ||
				last.typ == assignTokenType {
				expression := p.parseFunctionStatement()
				if expression == nil {
					fmt.Println("p.parseFunctionStatement() failed")
					return nil
				}
				expressions = append(expressions, &Object{
					inner: expression,
					label: "lambda",
					typ:   FuncStatementType,
				})
			} else {
				p.putToken(token)
				break Loop
			}
		//function call
		case token.typ == periodTokenType:
			//translation to `this.`
			expression := p.parsePeriodStatement("this")
			if expression == nil {
				fmt.Println("parsePeriodStatement failed")
				return nil
			}
			expressions = append(expressions, expression)
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
				label := token.data
				next := p.nextToken(true)
				//function call
				if next.typ == leftParenthesisTokenType {
					expression := p.parseFunctionCall([]string{label})
					if expression == nil {
						fmt.Println("parseFunctionCall failed")
						return nil
					}
					expressions = append(expressions, expression)
				} else if next.typ == incOperatorTokenType {
					p.closureCheckVisit(token.data)
					fmt.Println("incOperatorTokenType", label)
					expressions = append(expressions, &IncFieldStatement{
						ctx:   p.vmCtx,
						label: label,
					})
				} else if next.typ == periodTokenType { // label.
					p.closureCheckVisit(token.data)
					expression := p.parsePeriodStatement(token.data)
					if expression == nil {
						fmt.Println("parsePeriodStatement failed")
						return nil
					}
					expressions = append(expressions, expression)
				} else if next.typ == leftBraceTokenType { //eg: User {
					fmt.Println(token.data, "{")
					statement := p.parseObjectStructInit(token.data)
					if statement == nil {
						fmt.Println("parseObjectStructInit failed")
						return nil
					}
					return statement
				} else {
					p.putToken(next)
					p.closureCheckVisit(token.data)
					expressions = append(expressions, &getVarStatement{
						ctx:   p.vmCtx,
						label: label,
					})
					fmt.Println("getVarStatement", token.data)
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
		token := p.nextToken(true)
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
			log.Println("parseFunctionStatement------------------")
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
			token = p.nextToken(true)
			if token.typ != labelType {
				fmt.Println("error ,expect label", token)
				return nil
			}
			//closure check
			p.closureCheckAddVar(token.data)
			var label = token.data
			token = p.nextToken(true)
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
			label := token.data
			next := p.nextToken(true)
			//function call
			if next.typ == leftParenthesisTokenType {
				fmt.Println("parseFunctionCall")
				statement := p.parseFunctionCall([]string{label})
				if statement == nil {
					fmt.Println("parseFunctionCall failed")
					return nil
				}
				statements = append(statements, statement)
			} else if next.typ == incOperatorTokenType {
				p.closureCheckVisit(token.data)
				fmt.Println("incOperatorTokenType", label)
				statements = append(statements, &IncFieldStatement{
					ctx:   p.vmCtx,
					label: label,
				})
			} else if next.typ == periodTokenType { // label.
				p.closureCheckVisit(token.data)
				statement := p.parsePeriodStatement(token.data)
				if statement == nil {
					fmt.Println("parsePeriodStatement failed")
					return nil
				}
				statements = append(statements, statement)
			} else if next.typ == assignTokenType {
				p.closureCheckVisit(token.data)
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
				panic(token)
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

	token := p.nextToken(true)
	if token.typ != leftBraceTokenType {
		fmt.Println("error ,expect { ")
		return nil
	}
	leftBrace = append(leftBrace, token)
	//check empty statement
	if token = p.nextToken(true); token.typ == rightParenthesisTokenType {
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
		token := p.nextToken(true)
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

func (p *parser) parseForStatement() *ForStatement {
	var forStatement = ForStatement{
		vm: p.vmCtx,
	}
	token := p.nextToken(true)
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
		//support var= ;
		p.putToken(token)
		expression := p.parse()
		if expression == nil {
			fmt.Println("parse `for` preStatement failed")
			return nil
		}
		if p.nextToken(true).typ != semicolonTokenType {
			fmt.Println("expect ; in `for` statement")
			return nil
		}
		forStatement.preStatement = expression
	}
	token = p.nextToken(true)
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
		if p.nextToken(true).typ != semicolonTokenType {
			fmt.Println("expect ; in `for` check expression")
			return nil
		}
		forStatement.checkStatement = expression
	}

	token = p.nextToken(true)
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

		if next := p.nextToken(true); next.typ != leftBraceTokenType {
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
		vm: p.vmCtx,
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
		token := p.nextToken(true)
		fmt.Println(token)
		if token.typ == elseTokenType {
			next := p.nextToken(true)
			if next.typ == ifTokenType {
				token.typ = elseifTokenType
			} else {
				p.putToken(next)
			}
		}
		//check else or else if
		if token.typ == elseifTokenType {
			elseIfStem := IfStatement{vm: p.vmCtx}
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

func (p *parser) parseFunctionCall(labels []string) *FuncCallStatement {
	fmt.Println("parseFunctionCall", labels)
	var statement FuncCallStatement
	if len(labels) == 1 {
		statement.label = labels[0]
	} else {
		// bind self
		bindSelf := &getStructObjectStatement{
			this:      true,
			vmContext: p.vmCtx,
			labels:    labels[:len(labels)-1],
		}
		statement.arguments = append(statement.arguments, bindSelf)
	}
	statement.vm = p.vmCtx

	//empty arguments
	if token := p.nextToken(true); token.typ == rightParenthesisTokenType {
		fmt.Println("empty arguments")
		return &statement
	} else {
		p.putToken(token)
	}
	for {
		expression := p.parseExpression()
		if expression == nil {
			fmt.Println("parse expression failed")
			return nil
		}
		if expression.getType() == nopStatementType && len(statement.arguments) == 1 {
		} else {
			statement.arguments = append(statement.arguments, expression)
		}
		token := p.nextToken(true)
		if token.typ == rightParenthesisTokenType {
			break
		} else if token.typ == commaTokenType {
			// next parameters
			continue
		} else {
			p.putToken(token)
		}
	}
	fmt.Println("-------arguments-------", len(statement.arguments))
	return &statement
}

func (p *parser) parseFunctionStatement() *FuncStatement {
	fmt.Println("--parseFunctionStatement--")
	var functionStatement FuncStatement

	functionStatement.vm = p.vmCtx

	//function name
	if token := p.nextToken(true); token.typ == labelType {
		functionStatement.label = token.data
		for {
			next := p.nextToken(true)
			fmt.Println("next", next)
			if next.typ == periodTokenType { //type object function eg:user.get(){}
				functionStatement.labels = append(functionStatement.labels, token.data)
				token = p.nextToken(true)
				continue
			} else if next.typ == leftParenthesisTokenType {
				if functionStatement.labels != nil {
					functionStatement.labels = append(functionStatement.labels, token.data)
				}
				break
			} else {
				fmt.Println("expect label or . ,error")
				return nil
			}
		}
	} else {
		functionStatement.closure = true
		p.pushClosureCheck()
	}
	if functionStatement.labels != nil {
		fmt.Println("get function label",
			strings.Join(functionStatement.labels, "."))
	} else {
		fmt.Println("get function label", functionStatement.label)
	}
	//bind struct object to `this` argument
	if functionStatement.labels != nil {
		functionStatement.parameters = append(functionStatement.parameters, "this")
	}
	//parse argument list
	for {
		token := p.nextToken(true)
		if token.typ == rightParenthesisTokenType {
			// end of argument list
			break
		} else if token.typ == commaTokenType {
			//next argument
			continue
		} else if token.typ == labelType {
			p.closureCheckAddVar(token.data)
			functionStatement.parameters = append(functionStatement.parameters, token.data)
			fmt.Println("find argument", token.data)
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
	if functionStatement.closure {
		functionStatement.closureLabel = p.popClosureLabels()
	}

	return &functionStatement
}

func (p *parser) parseStructObject() *TypeObject {
	var object = &TypeObject{}
	token := p.nextToken(true)
	if token.typ != labelType {
		fmt.Println("expect label follow type")
		return nil
	}
	object.label = token.data
	object.vm = p.vmCtx
	statements := p.parseStatement()
	if statements == nil {
		fmt.Println("object struct parseStatement failed")
		return nil
	}
	object.initStatement = statements
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
	if token := p.nextToken(true); token.typ == rightBraceTokenType {
		statement.initStatements = append(statement.initStatements, &NopStatement{})
		return &statement
	} else {
		p.putToken(token)
	}

	for {
		token := p.nextToken(true)
		if token.typ != labelType {
			fmt.Println("expect label,error", token)
			return nil
		}
		if next := p.nextToken(true); next.typ != colonTokenType {
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
			label:      token.data,
			expression: expression,
		})
		//check end
		token = p.nextToken(true)
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
// a.b.c + 1 // get data statement
func (p *parser) parsePeriodStatement(label string) Statement {
	fmt.Println("parsePeriodStatement", label)
	var labels = []string{label}
	token := p.nextToken(true)
	if token.typ != labelType {
		fmt.Println("expect label ", token)
		return nil
	}
	labels = append(labels, token.data)
	for {
		next := p.nextToken(true)
		if next.typ == periodTokenType {
			token = p.nextToken(true)
			continue
			// a.b.c(1) //function call
		} else if next.typ == leftParenthesisTokenType {
			statement := p.parseFunctionCall(labels)
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

func (p *parser) closureCheckAddVar(data string) {
	if len(p.closureCheck) != 0 {
		p.closureCheck[len(p.closureCheck)-1].addVar(data)
	}
}

func (p *parser) closureCheckVisit(data string) {
	for i := len(p.closureCheck) - 1; i >= 0; i-- {
		closure := p.closureCheck[i]
		if closure.visit(data) == false {
			break
		}else{
			log.Println("-------------closureCheckVisit",data,"-------------")
		}
	}
}

func (p *parser) pushClosureCheck() {
	p.closureCheck = append(p.closureCheck, newClosureCheck())
}

func (p *parser) popClosureLabels() []string {
	if len(p.closureCheck) != 0 {
		closureLabel := p.closureCheck[len(p.closureCheck)-1].closures
		p.closureCheck = p.closureCheck[:len(p.closureCheck)-1]
		return closureLabel
	}
	return nil
}

func Parse(data string) Statements {
	parser := newParser(bytes.NewReader([]byte(data)))
	return parser.parse()
}
