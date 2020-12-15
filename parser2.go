package qp

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
)

/*
   Expression:Expression+Expression
		|Expression-Expression
   		|Expression
		|Factor


   FunCall:id()
   		|lambda()
		|Selector()

   lambda:func(Expressions)BlockStatement

   Expressions:Expression,Expression
   		|Expressions,Expression

   BlockStatement:{Statements}

   Statements:Statement \n Statement
   		|Statement;Statement

Statement:IfStatement
		|forStatement
		|assignStatement
		|breakStatement
		|varStatement

varStatement:var ID = Expression

assignStatement:ID = Expression
		|Selector = Expression



	Selector:ID.ID
		|Selector.ID
		|Factor().ID

Factor:|Factor*Factor
   	|!Factor
   	|^Factor
   	|funCall
   	|(Expressions)
	|FunCall

BoolOperator:==
	| !=
	| >=
	| <=
	| <
	| >

BoolExpression:Factor BoolOperator Factor
	|FallCall
	|!BoolExpression
	|(BoolExpression)
	|BoolExpression == BoolExpression
	|BoolExpression != BoolExpression

*/

type Parser2 struct {
	vm           *VMContext
	lexer        *lexer
	tokens       []Token
	hTokens      []Token
	pStack       []int //parenthesis stack
	closureCheck []*closureCheck
	status       []PStatus
}

type PStatus int

const (
	GlobalStatus   = 0
	IfStatus       = 1
	ElseStatus     = 2
	ForStatus      = 3
	FunctionStatus = 4
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

func NewParse2(buffer string) *Parser2 {
	return &Parser2{
		status: []PStatus{GlobalStatus},
		lexer:  newLexer(bytes.NewReader([]byte(buffer))),
		vm:     newVMContext(),
	}
}

func Parse(data string) Statements {
	return NewParse2(data).Parse()
}

func (p *Parser2) putToken(token Token) {
	p.hTokens = p.hTokens[:len(p.hTokens)-1]
	p.tokens = append([]Token{token}, p.tokens...)
}

func (p *Parser2) initTokens() *Parser2 {
	for {
		token := p.lexer.peek()
		if token.typ == EOFType {
			p.tokens = append(p.tokens, token)
			break
		}
		if token.typ == elseType {
			p.lexer.next()
			next := p.lexer.peek()
			p.lexer.next()
			if next.typ == ifType {
				token.typ = elseifType
				p.tokens = append(p.tokens, token)
			} else {
				p.tokens = append(p.tokens, token, next)
			}
			continue
		}
		p.tokens = append(p.tokens, token)
		p.lexer.next()
	}
	return p
}

func (p *Parser2) Parse() Statements {
	p.initTokens()
	var statements Statements
	for {
		statement := p.ParseStatement()
		if statement == nil {
			if p.ahead(0).typ == EOFType {
				return statements
			}
			continue
		}
		statements = append(statements, statement)
	}
}

/*
	AssignStatement:
		|ID = Expression
		|selector = Expression

	FunctionCallStatement:
		|ID()
		|selector()
		|FunctionCallStatement()
		|lambda
*/
func (p *Parser2) ParseIDPrefixStatement(token Token) Statement {
	var exp Expression = getVarStatement{
		ctx:   p.vm,
		label: token.val,
	}
	var parentExp Expression
	for {
		token := p.nextToken()
		switch token.typ {
		case assignType:
			return p.parseAssignStatement(exp)
		case incType:
			return IncFieldStatement{
				exp: exp,
			}
		case colonTokenType:
			log.Println("colonTokenType")
			return AssignStatement{left: exp, exp: p.parseFactor(0)}
		case periodType:
			token := p.nextToken()
			p.expectType(token, IDType)
			parentExp = exp
			exp = periodStatement{
				val: token.val,
				exp: exp,
			}
		case leftParenthesisType:
			exp = p.parseCallStatement(parentExp, exp)
			if p.ahead(0).typ != leftParenthesisType {
				return exp
			}
		default:
			log.Panic("unexpect token ", token)
		}
	}
}

func (p *Parser2) ParseStatement() Statement {
	var statement Statement
	for {
		token := p.nextToken()
		log.Println(token)
		switch token.typ {
		case typeType:
			p.vm.addStructObject(p.parseTypeStatement())
		case funcType:
			p.vm.addUserFunction(p.parseFuncStatement())
		case varType:
			return p.parseVarStatement()
		case ifType:
			return p.parseIfStatement(false)
		case EOFType:
			return statement
		case NewLineType:
			continue
		case semicolonType:
			continue
		case IDType:
			return p.ParseIDPrefixStatement(token)
		case returnType:
			return p.parseReturn()
		case rightBraceType:
			p.putToken(token)
			return statement
		case forType:
			return p.parseForStatement()
		case breakType:
			return p.parseBreakStatement()
		default:
			log.Panic(token)
		}
	}
}

func (p *Parser2) nextToken() Token {
	if len(p.tokens) == 0 {
		return Token{typ: EOFType}
	}
	token := p.tokens[0]
	p.tokens = p.tokens[1:]
	p.hTokens = append(p.hTokens, token)
	if token.typ == commentTokenType {
		return p.nextToken()
	}
	return token
}

/*
// user define object
type User{
}
*/
func (p *Parser2) parseTypeObjectInit() []TypeObjectPropTemplate {
	var objectPropTemplates []TypeObjectPropTemplate
	for {
		token := p.nextToken()
		p.expectType(p.nextToken(), colonTokenType)
		exp := p.parseFactor(0) //delay bind
		objectPropTemplates = append(objectPropTemplates, TypeObjectPropTemplate{
			name: token.val,
			exp:  exp,
		})
		if ahead := p.ahead(0); ahead.typ == commaType || ahead.typ == semicolonType {
			p.nextToken()
		} else {
			log.Println(p.ahead(0))
			if p.ahead(0).typ == IDType {
				if p.historyToken(1).line != p.ahead(0).line {
					continue
				}
				log.Panic("require new line")
			} else {
				log.Println(p.ahead(0))
				break
			}
		}
	}
	return objectPropTemplates
}

func (p *Parser2) parseTypeStatement() *TypeObject {
	var object TypeObject
	token := p.nextToken()
	//	log.Println(token)
	p.expectType(token, IDType)
	p.expectType(p.nextToken(), leftBraceType) //{
	object.label = token.val
	if ahead := p.ahead(0); ahead.typ == rightBraceType {
		p.nextToken()
		return &object
	} else {
		object.typeObjectPropTemplates = p.parseTypeObjectInit()
	}
	p.expectType(p.nextToken(), rightBraceType) //}
	return &object
}

func (p *Parser2) expectType(token Token, expect Type) {
	if token.typ != expect {
		log.Panicf("unexpect %s token", token)
	}
}

func (p *Parser2) parseVarStatement() Statement {
	token := p.nextToken()
	log.Println(token)
	p.expectType(token, IDType)
	next := p.nextToken()
	//var id = ....
	p.closureCheckAddVar(token.val)
	if next.typ == assignType {
		expression := p.parseFactor(0)
		return VarAssignStatement{
			ctx:        p.vm,
			name:       token.val,
			expression: expression,
		}
	}
	p.putToken(next)
	return VarStatement{
		ctx:   p.vm,
		label: token.val,
	}
}

/*
func(){}()

*/
func (p *Parser2) parseLambdaStatement() Statement {
	var funcS FuncStatement
	funcS.closure = true
	funcS.vm = p.vm
	token := p.nextToken()

	p.pushClosureCheck()

	p.expectType(token, leftParenthesisType)
	funcS.parameters = p.parseFuncParameters()
	p.expectType(p.nextToken(), leftBraceType)

	if p.ahead(0).typ == rightBraceType { //empty body
		p.nextToken()
		return &funcS
	}
	for {
		funcS.statements = append(funcS.statements, p.ParseStatement())
		if p.ahead(0).typ == rightBraceType {
			p.nextToken()
			break
		}
	}
	funcS.closureLabel = p.popClosureLabels()
	return &funcS
}

func (p *Parser2) parseFactor(pre int) Expression {
	var exp Expression
	var parentExp Expression
	for {
		token := p.nextToken()
		log.Println(token)
		switch token.typ {
		case falseType:
			p.assertNil(exp)
			exp = Bool(false)
		case TrueType:
			p.assertNil(exp)
			exp = Bool(true)
		case leftParenthesisType:
			if exp == nil {
				p.pStack = append(p.pStack, token.line)
				exp = ParenthesisExpression{exp: p.parseFactor(pre)}
				p.expectType(p.nextToken(), rightParenthesisType)
			} else {
				log.Println("parseCallStatement")
				exp = p.parseCallStatement(parentExp, exp)
			}
		case rightParenthesisType: //end of parenthesis ()
			if exp == nil {
				log.Panic("parse bool exp failed")
			}
			p.putToken(token)
			return exp
		case periodType:
			token := p.nextToken()
			p.expectType(token, IDType)
			parentExp = exp
			exp = periodStatement{
				val: token.val,
				exp: exp,
			}
		case mulOpType, divOpType:
			if exp == nil {
				log.Panicf("unexpect token %s", token)
			}
			p.nextToken()
			return BinaryOpExpression{
				opType: token.typ,
				Left:   exp,
				right:  p.parseFactor(precedence(token.typ)),
			}
		case NoType:
			p.assertNoNil(exp)
			exp = NoStatement{exp: p.parseFactor(pre)}
		case IDType:
			if exp != nil {
				p.putToken(token)
				if p.historyToken(1).line != token.line {
					return exp
				}
			}
			p.assertNil(exp)
			exp = &getVarStatement{
				ctx:   p.vm,
				label: token.val,
			}
			p.closureCheckVisit(token.val)
		case stringType:
			if exp != nil {
				p.putToken(token)
				if p.historyToken(1).line != token.line {
					return exp
				}
			}
			p.assertNil(exp)
			exp = String(token.val)
		case intType:
			p.assertNil(exp)
			val, _ := strconv.ParseInt(token.val, 10, 64)
			exp = Int(val)
		case funcType: // func(){}()
			if exp != nil {
				p.putToken(token)
				if p.historyToken(1).line != token.line {
					return exp
				}
			}
			exp = p.parseLambdaStatement()
		case EqualType, // ==
			NoEqualTokenType, // !=
			greaterEqualType, // >=
			greaterType,      // >
			lessEqualType,    // <=
			subType,          // -
			addType,          // +
			lessTokenType,    // <
			AndType,          // ||
			orType:           // ||
			if exp == nil {
				log.Panic("exp nil")
			}
			if pre >= precedence(token.typ) {
				//				log.Println("return", pre, precedence(token.typ))
				p.putToken(token)
				return exp
			}
			exp = BinaryBoolExpression{
				opType: token.typ,
				Left:   exp,
				right:  p.parseFactor(precedence(token.typ)),
			}
		case incType:
			p.assertNoNil(exp)
			exp = &IncFieldStatement{
				exp: exp,
			}
		case nilType:
			p.assertNil(exp)
			exp = nilObject
		case leftBraceType:
			if status := p.getStatus(); status == IfStatus || status == ForStatus {
				p.putToken(token)
				if exp == nil {
					log.Panicf("unexpect token %s", token)
				}
				log.Println("return {")
				return exp
			} else {
				log.Println(status)
			}
			return p.ParseObjInitStatement(exp)
		case leftBracketTokenType:
			exp = p.parseBracketStatement(exp)
		case EOFType:
			return exp
		default:
			if p.isTerminateToken(token) == false {
				log.Panicf("unexpect token %s", token)
			}
			p.putToken(token)
			if exp == nil {
				log.Panicf("unexpect token %s", token)
			}
			return exp
		}
	}
}

/*
	FactorList:Factor,Factor
			|FactorList,Factor

	getBracketStatement:ID[]
		|ID[Factor]
		|ID[FactorList]

	initBracketStatement:[Factor]
			|[FactorList]

*/

func (p *Parser2) parseBracketStatement(exp Expression) Expression {
	//init array object
	if exp == nil {
		var arrayStatement makeArrayStatement
		for {
			if p.ahead(0).typ == rightBracketType {
				p.nextToken()
				return &arrayStatement
			}
			arrayStatement.initStatements = append(arrayStatement.initStatements, p.parseFactor(0))
			if p.ahead(0).typ == commaType { // ,
				p.nextToken()
			}
		}
	} else { //get array field
		index := p.parseFactor(0)
		p.expectType(p.nextToken(), rightBracketType)
		return getArrayElement{
			arrayExp: exp,
			indexExp: index,
		}
	}
}

func (p *Parser2) ahead(index int) Token {
	if len(p.tokens) <= index {
		return Token{typ: EOFType}
	}
	return p.tokens[index]
}

func (p *Parser2) parseCallStatement(parentExp, function Expression) Expression {
	var call FuncCallStatement
	call.function = function
	call.parentExp = parentExp
	for {
		if p.ahead(0).typ == rightParenthesisType {
			p.nextToken()
			return &call
		}
		call.arguments = append(call.arguments, p.parseFactor(0))
		if p.ahead(0).typ == commaType { // ,
			p.nextToken()
		}
	}
}

/*
BoolExpression:Factor BoolOperator Factor
	|FallCall
	|!BoolExpression
	|(BoolExpression)
	|BoolExpression BoolOperator BoolExpression


if !(0 != 2) == true {

}

precedence
|| : 1
&& : 2
> >= < <= == != : 3
*/

/*
func hello(){}
func User.add(){

}
*/
func (p *Parser2) parseFuncStatement() *FuncStatement {
	var funcS FuncStatement
	token := p.nextToken()
	p.expectType(token, IDType)

	if ahead := p.ahead(0); ahead.typ == periodType {
		p.nextToken()
		next := p.nextToken()
		p.expectType(next, IDType)
		funcS.labels = []string{token.val, next.val}
	} else {
		funcS.label = token.val
	}
	p.expectType(p.nextToken(), leftParenthesisType)
	if len(funcS.labels) != 0 {
		funcS.parameters = append(funcS.parameters, "this")
	}
	funcS.parameters = append(funcS.parameters, p.parseFuncParameters()...)
	p.expectType(p.nextToken(), leftBraceType)
	for {
		funcS.statements = append(funcS.statements, p.ParseStatement())
		if p.ahead(0).typ == rightBraceType {
			p.nextToken()
			break
		}
	}
	funcS.vm = p.vm
	return &funcS
}

func (p *Parser2) parseFuncParameters() []string {
	var parameters []string
	if p.ahead(0).typ == rightParenthesisType {
		p.nextToken()
		return nil
	}
	for {
		token := p.nextToken()
		p.expectType(token, IDType)
		parameters = append(parameters, token.val)
		p.closureCheckAddVar(token.val)
		if ahead := p.ahead(0); ahead.typ == commaType {
			p.nextToken()
			continue
		} else if ahead.typ == rightParenthesisType {
			p.nextToken()
			break
		}
	}
	return parameters
}

func (p *Parser2) parseBoolExpression(pre int) Expression {
	var exp Expression
	exp = p.parseFactor(0)
	if _, ok := exp.(BinaryBoolExpression); ok {
		return exp
	}
	log.Printf("parse exp bool failed,%s", exp.String())
	return exp
}

func (p *Parser2) parseIfStatement(elseif bool) *IfStatement {
	p.pushStatus(IfStatus)
	log.Println("enter parseIfStatement")
	defer func() {
		p.assertTrue(p.popStatus() == IfStatus)
		log.Println("out parseIfStatement")
	}()
	var ifS IfStatement
	if p.ahead(0).typ != leftBraceType {
		ifS.check = p.parseBoolExpression(0)
	}
	p.expectType(p.nextToken(), leftBraceType)

	if p.ahead(0).typ == rightBraceType {
		p.nextToken()
		return &ifS
	}
	ifS.vm = p.vm
	ifS.statement = p.ParseStatements()
	p.expectType(p.nextToken(), rightBraceType)

	for elseif == false {
		next := p.ahead(0)
		if next.typ == elseifType {
			p.nextToken()
			statement := p.parseIfStatement(true)
			ifS.elseIfStatements = append(ifS.elseIfStatements, statement)
		} else if next.typ == elseType {
			p.expectType(p.nextToken(), elseType)
			p.expectType(p.nextToken(), leftBraceType)
			ifS.elseStatement = p.ParseStatements()
			p.expectType(p.nextToken(), rightBraceType)
		} else {
			break
		}
	}
	return &ifS

}
func (p *Parser2) assertNil(exp Expression) {
	if exp != nil {
		log.Panicf("expect nil")
	}
}

func (p *Parser2) assertNoNil(exp Expression) {
	if exp == nil {
		log.Panicf("expect nil")
	}
}

func (p *Parser2) isTerminateToken(next Token) bool {
	if next.typ == ifType ||
		next.typ == forType ||
		next.typ == leftBraceType ||
		next.typ == rightBraceType ||
		next.typ == semicolonType ||
		next.typ == commaType ||
		next.typ == varType ||
		next.typ == breakType ||
		next.typ == returnType ||
		next.typ == EOFType {
		return true
	}
	return false
}

func (p *Parser2) historyToken(i int) Token {
	index := len(p.hTokens) - i - 1
	if index < len(p.hTokens) && index >= 0 {
		return p.hTokens[index]
	}
	log.Panic("out of history tokens range")
	return emptyToken
}

func (p *Parser2) closureCheckAddVar(data string) {
	if len(p.closureCheck) != 0 {
		p.closureCheck[len(p.closureCheck)-1].addVar(data)
	}
}

func (p *Parser2) closureCheckVisit(data string) {
	for i := len(p.closureCheck) - 1; i >= 0; i-- {
		closure := p.closureCheck[i]
		if closure.visit(data) == false {
			break
		}
	}
}

func (p *Parser2) pushClosureCheck() {
	p.closureCheck = append(p.closureCheck, newClosureCheck())
}

func (p *Parser2) popClosureLabels() []string {
	if len(p.closureCheck) != 0 {
		closureLabel := p.closureCheck[len(p.closureCheck)-1].closures
		p.closureCheck = p.closureCheck[:len(p.closureCheck)-1]
		return closureLabel
	}
	return nil
}

func (p *Parser2) parseReturn() Statement {
	if p.ahead(0).typ == rightBraceType {
		return ReturnStatement{
			express:   nilObject,
			returnVal: nilObject,
		}
	}
	return ReturnStatement{
		express: p.parseFactor(0),
	}
}

func (p *Parser2) parseForStatement() Statement {
	log.Println("parseForStatement")
	p.pushStatus(ForStatus)
	defer func() {
		p.assertTrue(p.popStatus() == ForStatus)
	}()
	var forStatement = ForStatement{
		vm: p.vm,
	}
	token := p.nextToken()
	log.Println(token)
	if token.typ == semicolonType {
		forStatement.preStatement = nopStatement
	} else if token.typ == leftBraceType {
		forStatement.preStatement = nopStatement
		forStatement.postStatement = nopStatement
		forStatement.checkStatement = &trueObject
		statements := p.ParseStatements()
		p.expectType(p.nextToken(), rightBraceType)
		forStatement.statements = statements
		return &forStatement
	} else {
		//support var= ;
		expression := p.parseVarStatement()
		if p.nextToken().typ != semicolonType {
			log.Panic("expect ; in `for` statement")
			return nil
		}
		forStatement.preStatement = expression
	}
	token = p.nextToken()
	//check exp
	if token.typ == semicolonType {
		forStatement.checkStatement = &trueObject
	} else {
		p.putToken(token)
		expression := p.parseFactor(0)
		if p.nextToken().typ != semicolonType {
			log.Panic("expect ; in `for` check exp")
		}
		forStatement.checkStatement = expression
	}

	token = p.nextToken()
	//post exp
	if token.typ == leftBraceType {
		forStatement.postStatement = nopStatement
	} else {
		p.putToken(token)
		expression := p.ParseStatement()
		p.expectType(p.nextToken(), leftBraceType)
		forStatement.postStatement = expression
	}
	// statements
	statements := p.ParseStatements()
	forStatement.statements = statements
	p.expectType(p.nextToken(), rightBraceType)
	return &forStatement
}

func (p *Parser2) ParseStatements() Statements {
	var statements Statements
	if p.ahead(0).typ == rightBraceType {
		return append(statements, NopStatement{})
	}
	for {
		statement := p.ParseStatement()
		fmt.Println(statement.String())
		statements = append(statements, statement)

		if p.ahead(0).typ == rightBraceType {
			log.Println(p.getStatus())
			if p.getStatus() == GlobalStatus {
				p.nextToken()
				continue
			}
			return statements
		}
	}
}

func (p *Parser2) ParseObjInitStatement(exp Expression) Expression {
	var statement objectInitStatement
	statement.exp = exp
	statement.vm = p.vm
	if ahead := p.ahead(0); ahead.typ == rightBraceType {
		p.nextToken()
		return &statement
	} else {
		statement.propTemplates = p.parseTypeObjectInit()
	}
	p.expectType(p.nextToken(), rightBraceType)
	return &statement
}

func (p *Parser2) pushStatus(status PStatus) {
	p.status = append(p.status, status)
}
func (p *Parser2) popStatus() PStatus {
	if len(p.status) == 0 {
		log.Panic("status stack empty")
	}
	status := p.status[len(p.status)-1]
	p.status = p.status[:len(p.status)-1]
	return status
}

func (p *Parser2) assertTrue(val bool) {
	if val == false {
		panic("assert failed")
	}
}

func (p *Parser2) checkInStatus(status PStatus) bool {
	for i := len(p.status) - 1; i >= 0; i-- {
		if p.status[i] == status {
			return true
		} else if p.status[i] == FunctionStatus {
			break
		}
	}
	return false
}

func (p Parser2) getStatus() PStatus {
	if len(p.status) != 0 {
		return p.status[len(p.status)-1]
	}
	log.Panic("status stack empty")
	return 0
}

func (p *Parser2) parseBreakStatement() Statement {
	if p.checkInStatus(ForStatus) == false {
		log.Panic("current no in `for` brock ")
	}
	return breakObject
}

func (p *Parser2) parseAssignStatement(exp Expression) AssignStatement {
	return AssignStatement{
		exp:  p.parseFactor(0),
		left: exp,
	}
}
