package qp

import (
	io "io"
	"log"
	"strconv"
	"strings"
)

func precedence(tokenType Type) int {
	switch tokenType {
	case mulOpType, divOpType:
		return 10
	case addType, subType:
		return 9
	case lessTokenType, lessEqualType, greaterType, greaterEqualType, NoEqualTokenType, EqualType:
		return 8
	case AndType:
		return 7
	case orType:
		return 6
	default:
		return 0
	}
}

//greater or equal
func precedenceGE(first, second Type) bool {
	return precedence(first)-precedence(second) >= 0
}

func isOperatorToken(token Token) bool {
	return token.typ >= addType && token.typ <= NoEqualTokenType
}


const (
	GlobalStatus   = 0
	IfStatus       = 1
	ElseStatus     = 2
	ForStatus      = 3
	FunctionStatus = 4
)

type parser struct {
	history      []Token
	tokens       []Token
	lexer        *lexer
	vmCtx        *VMContext
	closureCheck []*closureCheck
	status       []int
}

func newParser(reader io.Reader) *parser {
	return &parser{
		lexer:  newLexer(reader),
		vmCtx:  newVMContext(),
		status: []int{GlobalStatus},
	}
}

func (p *parser) checkInStatus(status int) bool {
	for i := len(p.status) - 1; i >= 0; i-- {
		if p.status[i] == status {
			return true
		} else if p.status[i] == FunctionStatus {
			break
		}
	}
	return false
}

func makeExpression(opToken Token, expressions *Expressions) Expression {
	var expression Expression
	switch opToken.typ {
	case subType:
		expression = SubExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case addType:
		expression = AddExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case mulOpType:
		expression = MulExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case lessTokenType:
		expression = LessExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case lessEqualType:
		expression = LessEqualExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case greaterType:
		expression = GreaterExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case greaterEqualType:
		expression = GreaterEqualExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case EqualType:
		expression = EqualExpression{
			Left:  (*expressions)[len(*expressions)-2],
			right: (*expressions)[len(*expressions)-1],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case NoEqualTokenType:
		expression = NoEqualExpression{
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
	var expressions Expressions
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
		if token.typ == commentTokenType {
			continue
		}
		switch {
		case token.typ == intType:
			val, err := strconv.ParseInt(string(token.val), 10, 64)
			if err != nil {
				log.Panic("parse int failed", string(token.val))
			}
			expressions = append(expressions, (Int)(val))
		case token.typ == stringType:
			expressions = append(expressions, String(token.val))
		case token.typ == nilType:
			expressions = append(expressions, nilObject)
		case token.typ == leftParenthesisType:
			opStack = append(opStack, token)
		case token.typ == TrueType:
			expressions = append(expressions, &trueObject)
		case token.typ == falseType:
			expressions = append(expressions, &falseObject)
		case token.typ == rightParenthesisType:
			var find = false
			for _, opCode := range opStack {
				if opCode.typ == leftParenthesisType {
					find = true
				}
			}
			//end of expression
			if find == false {
				p.putToken(token)
				break Loop
			}
			if find {
				for len(opStack) != 0 && opStack[len(opStack)-1].typ != leftParenthesisType {
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
					log.Panic("make expression failed", opStack[len(opStack)-1])
				}
				expressions = append(expressions, express)
				opStack = opStack[:len(opStack)-1]
			}
			opStack = append(opStack, token)
			//lambda func
		case token.typ == funcType:
			last := p.historyToken(1)
			if isOperatorToken(last) || //1 + func(){}()
				last.typ == assignType || //var a = func(){}
				last.typ == leftBracketTokenType || //var  = [func(){}]
				last.typ == commaType || //var  = [1,func(){}]
				last.typ == returnType { //return func(){}
				expression := p.parseFunctionStatement()
				expressions = append(expressions, expression)
			} else {
				p.putToken(token)
				break Loop
			}
		//function call
		case token.typ == periodType:
			//translation to `this.`
			expression := p.parsePeriodStatement("this")
			expressions = append(expressions, expression)
			// []
		case token.typ == leftBracketTokenType:
			expressions = append(expressions, p.parseArrayInit())
		case token.typ == IDType:
			his := p.historyToken(1)
			if isOperatorToken(his) ||
				his.typ == semicolonType ||
				his.typ == commaType ||
				his.typ == ifType ||
				his.typ == elseifType ||
				his.typ == leftParenthesisType ||
				his.typ == leftBraceType ||
				his.typ == assignType ||
				his.typ == colonTokenType ||
				his.typ == returnType {
				label := token.val
				next := p.nextToken(true)
				//function call
				if next.typ == leftParenthesisType {
					expressions = append(expressions, p.parseFunctionCall([]string{label}))
				} else if next.typ == incOperatorTokenType {
					p.closureCheckVisit(token.val)
					expressions = append(expressions, &IncFieldStatement{
						ctx:   p.vmCtx,
						//label: label,
					})
				} else if next.typ == periodType { // label.
					p.closureCheckVisit(token.val)
					expressions = append(expressions, p.parsePeriodStatement(token.val))
				} else if next.typ == leftBraceType { //eg: User {
					return p.parseObjectStructInit(token.val)
				} else {
					p.putToken(next)
					p.closureCheckVisit(token.val)
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
	return expressions
}

func (p *parser) parse() Statements {
	var statements Statements
	statements = append(statements, nopStatement)
Loop:
	for {
		token := p.nextToken(true)
		if token == emptyToken {
			break
		}
		if token.typ == commentTokenType {
			continue
		}
		if token.typ == EOFType {
			break
		}
		switch {
		case token.typ == ifType:
			expression := p.parseIfStatement()
			if expression == nil {
				log.Panic("parseIfStatement failed")
			}
			statements = append(statements, expression)
		case token.typ == forType:
			statements = append(statements, p.parseForStatement())
		case token.typ == returnType:
			statement := ReturnStatement{}
			statement.express = p.parseExpression()
			statements = append(statements, &statement)
			continue
		case token.typ == funcType:
			statement := p.parseFunctionStatement()
			if function, ok := statement.(*FuncStatement); ok {
				p.vmCtx.addUserFunction(function)
			} else {
				statements = append(statements, statement)
			}
		case token.typ == typeType:
			p.vmCtx.addStructObject(p.parseTypeObject())
		case token.typ == varType:
			token = p.nextToken(true)
			if token.typ != IDType {
				log.Panic("error ,expect label", token)
				return nil
			}
			//closure check
			p.closureCheckAddVar(token.val)
			var label = token.val
			token = p.nextToken(true)
			if token.typ == assignType {
				expression := p.parseExpression()
				if expression == nil {
					log.Panic("parse assign expression failed")
					return nil
				}
				statements = append(statements, VarAssignStatement{
					ctx:        p.vmCtx,
					label:      label,
					expression: expression,
				})
			} else {
				p.putToken(token)
				statements = append(statements, VarStatement{
					ctx:   p.vmCtx,
					label: label,
				})
			}
			continue
		case token.typ == leftBraceType: //{ end of expression
			p.putToken(token)
			break Loop
		case token.typ == rightBraceType: // } end of statement
			p.putToken(token)
			break Loop
		case token.typ == semicolonType: // ; end of expression
			p.putToken(token)
			break Loop
		case token.typ == commaType: // end of expression
			p.putToken(token)
			break Loop
		case token.typ == periodType:
			statements = append(statements, p.parsePeriodStatement("this"))
		case token.typ == breakType:
			if p.checkInStatus(ForStatus) == false {
				log.Panicf("break no in `for` block error line:%d", token.line)
			}
			statements = append(statements, breakObject)
		case token.typ == TrueType:
			statements = append(statements, &trueObject)
		case token.typ == falseType:
			statements = append(statements, &falseObject)
		case token.typ == IDType:
			label := token.val
			next := p.nextToken(true)
			//function call
			if next.typ == leftParenthesisType {
				statement := p.parseFunctionCall([]string{label})
				if statement == nil {
					log.Panic("parseFunctionCall failed")
					return nil
				}
				statements = append(statements, statement)
			} else if next.typ == incOperatorTokenType {
				p.closureCheckVisit(token.val)
				statements = append(statements, &IncFieldStatement{
					ctx:   p.vmCtx,
					//label: label,
				})
			} else if next.typ == leftBraceType { //eg: User {
				statements = append(statements, p.parseObjectStructInit(token.val))
			} else if next.typ == periodType { // label.
				p.closureCheckVisit(token.val)
				statement := p.parsePeriodStatement(token.val)
				if statement == nil {
					log.Panic("parsePeriodStatement failed")
					return nil
				}
				statements = append(statements, statement)
			} else if next.typ == assignType {
				p.closureCheckVisit(token.val)
				expression := p.parseExpression()
				statements = append(statements, &AssignStatement{
					ctx:        p.vmCtx,
					label:      label,
					expression: expression,
				})
			} else {
				panic(token)
			}
		default:
			panic("unknown token" + token.String())
		}
	}
	return statements
}

func (p *parser) parseStatement() Statements {
	var statements Statements
	var leftBrace []Token

	token := p.nextToken(true)
	if token.typ != leftBraceType {
		log.Panic("error ,expect { ", token)
	}
	leftBrace = append(leftBrace, token)
	//check empty statement
	if token = p.nextToken(true); token.typ == rightParenthesisType {
		statements = append(statements, nopStatement)
		return statements
	}
	p.putToken(token)

	for {
		statement := p.parse()
		if statement == nil {
			log.Panic("parse statement failed")
		}
		statements = append(statements, statement...)
		token := p.nextToken(true)
		if token.typ == rightBraceType {
			if leftBrace = leftBrace[:len(leftBrace)-1]; len(leftBrace) == 0 {
				break
			}
		}
	}
	return statements
}

func (p parser) getStatus() int {
	return p.status[len(p.status)-1]
}

func (p *parser) parseForStatement() *ForStatement {
	p.status = append(p.status, ForStatus)
	defer func() {
		p.status = p.status[:len(p.status)-1]
	}()
	var forStatement = ForStatement{
		vm: p.vmCtx,
	}
	token := p.nextToken(true)
	if token.typ == semicolonType {
		forStatement.preStatement = nopStatement
	} else if token.typ == leftBraceType {
		p.putToken(token)
		forStatement.preStatement = nopStatement
		forStatement.postStatement = nopStatement
		forStatement.checkStatement = &trueObject
		statements := p.parseStatement()
		forStatement.statements = statements
		return &forStatement
	} else {
		//support var= ;
		p.putToken(token)
		expression := p.parse()
		if p.nextToken(true).typ != semicolonType {
			log.Panic("expect ; in `for` statement")
			return nil
		}
		forStatement.preStatement = expression
	}
	token = p.nextToken(true)
	//check expression
	if token.typ == semicolonType {
		forStatement.checkStatement = &trueObject
	} else {
		p.putToken(token)
		expression := p.parseExpression()
		if p.nextToken(true).typ != semicolonType {
			log.Panic("expect ; in `for` check expression")
		}
		forStatement.checkStatement = expression
	}

	token = p.nextToken(true)
	//post expression
	if token.typ == leftBraceType {
		forStatement.postStatement = nopStatement
		p.putToken(token)
	} else {
		p.putToken(token)
		expression := p.parse()
		if next := p.nextToken(true); next.typ != leftBraceType {
			log.Panicf("expect { in `for` post expression get `%s`", next)
		} else {
			p.putToken(next)
		}
		forStatement.postStatement = expression
	}
	// statements
	statements := p.parseStatement()
	forStatement.statements = statements
	return &forStatement
}

func (p *parser) parseBoolExpression() Expression {
	var expressions = p.parseExpression().(Expressions)
	if len(expressions) != 1 {
		log.Panic("parse bool expression failed")
	}
	if lessTokenType <= expressions[0].getType() &&
		expressions[0].getType() <= NoEqualTokenType {
		return expressions
	}
	log.Panic("parseBoolExpression failed")
	return nil
}

func (p *parser) parseIfStatement() *IfStatement {
	ifStem := IfStatement{
		vm: p.vmCtx,
	}
	if ifStem.check = p.parseBoolExpression(); ifStem.check == nil {
		log.Panic("parse checkExpression failed")
		return nil
	}
	if ifStem.statement = p.parseStatement(); ifStem.statement == nil {
		log.Panic("parseStatement failed")
		return nil
	}
	for {
		token := p.nextToken(true)
		if token.typ == elseType {
			next := p.nextToken(true)
			if next.typ == ifType {
				token.typ = elseifType
			} else {
				p.putToken(next)
			}
		}
		//check else or else if
		if token.typ == elseifType {
			elseIfStem := IfStatement{vm: p.vmCtx}
			if elseIfStem.check = p.parseExpression(); elseIfStem.check == nil {
				log.Panic("parse checkExpression failed")
				return nil
			}
			if elseIfStem.statement = p.parseStatement(); elseIfStem.statement == nil {
				log.Panic("parseStatement failed")
				return nil
			}
			elseIfStem.elseIfStatements = append(elseIfStem.elseIfStatements, &elseIfStem)
		} else if token.typ == elseType {
			if ifStem.elseStatement = p.parseStatement(); ifStem.elseStatement == nil {
				log.Panic("parse else statement failed")
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
	/*var statement FuncCallStatement
	if len(labels) == 1 {
		statement.label = labels[0]
	} else {
		// bind self
		bindSelf := &getObjectPropStatement{
			this: true,
			getObject: &getObjectObjectStatement{
				vmContext: p.vmCtx,
				labels:    labels[:len(labels)-1],
			},
		}
		statement.arguments = append(statement.arguments, bindSelf)
	}
	statement.vm = p.vmCtx

	//empty arguments
	if token := p.nextToken(true); token.typ == rightParenthesisType {
		return &statement
	} else {
		p.putToken(token)
	}
	for {
		expression := p.parseExpression()
		if expression.getType() == nopStatementType && len(statement.arguments) == 1 {
		} else {
			statement.arguments = append(statement.arguments, expression)
		}
		token := p.nextToken(true)
		if token.typ == rightParenthesisType {
			break
		} else if token.typ == commaType {
			// next parameters
			continue
		} else {
			p.putToken(token)
		}
	}
	return &statement*/
	return nil
}

func (p *parser) parseFunctionStatement() Statement {
	var functionStatement FuncStatement
	functionStatement.vm = p.vmCtx

	//function name
	if token := p.nextToken(true); token.typ == IDType {
		functionStatement.label = token.val
		for {
			next := p.nextToken(true)
			if next.typ == periodType { //type objects function eg:user.get(){}
				functionStatement.labels = append(functionStatement.labels, token.val)
				token = p.nextToken(true)
				continue
			} else if next.typ == leftParenthesisType {
				if functionStatement.labels != nil {
					functionStatement.labels = append(functionStatement.labels, token.val)
				}
				break
			} else {
				log.Panic("expect label or . ,error")
			}
		}
	} else {
		functionStatement.closure = true
		p.pushClosureCheck()
	}
	//bind struct objects to `this` argument
	if functionStatement.labels != nil {
		functionStatement.parameters = append(functionStatement.parameters, "this")
	}
	//parse argument list
	for {
		token := p.nextToken(true)
		if token.typ == rightParenthesisType {
			// end of argument list
			break
		} else if token.typ == commaType {
			//next argument
			continue
		} else if token.typ == IDType {
			p.closureCheckAddVar(token.val)
			functionStatement.parameters = append(functionStatement.parameters, token.val)
		} else {
			log.Panic("unknown argument token", token)
		}
	}
	functionStatement.statements = p.parseStatement()
	if functionStatement.closure {
		functionStatement.closureLabel = p.popClosureLabels()
	}
	/*
		func(){
			return func(){
				return func(){
					return func(){
						printlnFunc(1)
					}
				}
			}
		}()()()()
	*/

	token := p.nextToken(true)
	if token.typ == leftParenthesisType {
		var funcCallQueueStatement FuncCallQueueStatement
		for {
			expression := p.parseFunctionCall([]string{"lambda"})
			//expression.function = &functionStatement
			funcCallQueueStatement.statement = append(funcCallQueueStatement.statement, expression)
			token = p.nextToken(true)
			if token.typ != leftParenthesisType {
				p.putToken(token)
				break
			}
		}
		return &funcCallQueueStatement
	} else {
		p.putToken(token)
	}
	return &functionStatement
}

func (p *parser) parseTypeObject() *TypeObject {
	var object = &TypeObject{}
	token := p.nextToken(true)
	if token.typ != IDType {
		log.Panic("expect label follow type")
		return nil
	}
	object.label = token.val
	object.vm = p.vmCtx
	statements := p.parseStatement()
	if statements == nil {
		log.Panic("objects struct parseStatement failed")
		return nil
	}
	object.initStatement = statements
	return object
}

func (p *parser) parseObjectStructInit(label string) *StructObjectInitStatement {
	var statement StructObjectInitStatement
	var leftBrace []int

	statement.label = label
	statement.vm = p.vmCtx
	leftBrace = append(leftBrace, 1)
	//check empty statement
	if token := p.nextToken(true); token.typ == rightBraceType {
		statement.initStatements = append(statement.initStatements, nopStatement)
		return &statement
	} else {
		p.putToken(token)
	}

	for {
		token := p.nextToken(true)
		if token.typ != IDType {
			log.Panic("expect label,error", token)
		}
		if token = p.nextToken(true); token.typ != colonTokenType {
			log.Panic("expect colon `:` ,error ", token)
		}
		statement.initStatements = append(statement.initStatements, VarAssignStatement{
			ctx:        p.vmCtx,
			label:      token.val,
			expression: p.parseExpression(),
		})
		//check end
		token = p.nextToken(true)
		if token.typ == rightBraceType {
			break
		} else if token.typ == commaType {
			continue
		}
		p.putToken(token)
	}
	return &statement
}

// a.b.c = 1 // assign
// a.b.c()   // function call
// a.b.c + 1 // get val statement
func (p *parser) parsePeriodStatement(label string) Statement {
	var labels = []string{label}
	token := p.nextToken(true)
	if token.typ != IDType {
		log.Panic("expect label ", token)
	}
	labels = append(labels, token.val)
	for {
		next := p.nextToken(true)
		if next.typ == periodType {
			token = p.nextToken(true)
			if token.typ != IDType {
				log.Panic("expect label ", token)
			}
			labels = append(labels, token.val)
			continue
			// a.b.c(1) //function call
		} else if next.typ == leftParenthesisType {
			statement := p.parseFunctionCall(labels)
			if statement == nil {
				log.Panic("parseFunctionCall failed")
				return nil
			}
			/*statement.getObject = &getObjectPropStatement{
				getObject: &getObjectObjectStatement{
					vmContext: p.vmCtx,
					labels:    labels,
				},
			}*/
			return statement
		} else if next.typ == assignType { // a.b = 1
			expression := p.parseExpression()
			if expression == nil {
				log.Panic("parseExpression for assign statement failed", token)
				return nil
			}
			return &AssignStatement{
				ctx:   p.vmCtx,
				label: strings.Join(labels, "."),
				getObject: &getObjectObjectStatement{
					vmContext: p.vmCtx,
					labels:    labels,
				},
				expression: expression,
			}
		} else {
			// a.b.c +  // expression
			// var c = a.b.c //end of statement
			p.putToken(next)
			return &getObjectPropStatement{
				getObject: &getObjectObjectStatement{
					vmContext: p.vmCtx,
					labels:    labels},
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

// var a = []
func (p *parser) parseArrayInit() *makeArrayStatement {
	var statement = &makeArrayStatement{
		vm:             p.vmCtx,
		initStatements: nil,
	}
	for {
		token := p.nextToken(true)
		if token.typ == rightBracketTokenType {
			return statement
		} else if token.typ == commaType {
			statement.initStatements = append(statement.initStatements,
				p.parseExpression())
		} else {
			p.putToken(token)
			statement.initStatements = append(statement.initStatements,
				p.parseExpression())
		}
	}
}

/*func Parse(data string) Statements {
	parser := newParser(bytes.NewReader([]byte(data)))
	return parser.parse()
}*/
