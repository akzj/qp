package qp

import (
	"bytes"
	"log"
	"strconv"
)

/*
   Expression:Expression+Expression
		|Expression-Expression
   		|Expression


   FunCall:id()
   		|lambda()
		|Selector()

   lambda:func(Expressions)BlockStatement

   Expressions:Expression,Expression
   		|Expressions,Expression

   BlockStatement:{Statements}

   Statements:Statement \n Statement
   		|Statements,Statement

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
	status       []int
}

func NewParse2(buffer string) *Parser2 {
	return &Parser2{
		lexer: newLexer(bytes.NewReader([]byte(buffer))),
		vm:    newVMContext(),
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
		if token.typ == ifType {
			p.lexer.next()
			next := p.lexer.peek()
			if next.typ == elseType {
				token.typ = elseifType
			} else {
				p.tokens = append(p.tokens, token)
				continue
			}
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
			return statements
		}
		statements = append(statements, statement)
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
			return p.parseIfStatement()
		case EOFType:
			//log.Println(token)
			return statement
		case NewLineType:
			continue
		case semicolonType:
			continue
		case IDType:
			p.putToken(token)
			return p.parseFactor(0)
		case returnType:
			return p.parseReturn()
		case rightBraceType:
			p.putToken(token)
			return statement
		case forType:
			return p.parseForStatement()
		case breakType:
			return breakObject
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
		for {
			p.expectType(p.nextToken(), varType)
			statement := p.parseVarStatement()
			object.initStatement = append(object.initStatement, statement)
			if ahead := p.ahead(0); ahead.typ == commaType || ahead.typ == semicolonType {
				p.nextToken()
			} else {
				break
			}
		}
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
		expression := p.parseExpression()
		return &VarAssignStatement{
			object:     nil,
			ctx:        p.vm,
			label:      token.val,
			expression: expression,
		}
	}
	p.putToken(token)
	return &VarStatement{
		ctx:   p.vm,
		label: token.val,
	}
}

func (p *Parser2) parseExpression() Expression {
	factor := p.parseFactor(0)
	op := p.ahead(0)
	if op.typ == addType || op.typ == subType {
		p.nextToken()
		return BinaryOpExpression{opType: op.typ,
			Left:  factor,
			right: p.parseFactor(0)}
	} else {
		//		log.Println(op)
		return factor
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
				exp = p.parseCallStatement(exp)
			}
		case rightParenthesisType: //end of parenthesis ()
			if exp == nil {
				log.Panic("parse bool expression failed")
			}
			p.putToken(token)
			return exp
		case periodType:
			token := p.nextToken()
			p.expectType(token, IDType)
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
		case incOperatorTokenType:
			p.assertNoNil(exp)
			exp = &IncFieldStatement{
				ctx: p.vm,
				exp: exp,
			}
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

 */

func (p *Parser2) ahead(index int) Token {
	if len(p.tokens) <= index {
		return Token{typ: EOFType}
	}
	return p.tokens[index]
}

func (p *Parser2) parseCallStatement(left Expression) Expression {
	var call FuncCallStatement
	call.expression = left
	for {
		if p.ahead(0).typ == rightParenthesisType {
			p.nextToken()
			return &call
		}
		call.arguments = append(call.arguments, p.parseExpression())
		if p.ahead(0).typ == commaType { // ,
			p.nextToken()
		}
	}
}

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
	funcS.parameters = p.parseFuncParameters()
	p.expectType(p.nextToken(), rightParenthesisType)
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
	log.Printf("parse expression bool failed,%s", exp.String())
	return exp
}

func (p *Parser2) parseIfStatement() *IfStatement {
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
	p.status = append(p.status, ForStatus)
	defer func() {
		p.status = p.status[:len(p.status)-1]
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
	//check expression
	if token.typ == semicolonType {
		forStatement.checkStatement = &trueObject
	} else {
		p.putToken(token)
		expression := p.parseExpression()
		if p.nextToken().typ != semicolonType {
			log.Panic("expect ; in `for` check expression")
		}
		forStatement.checkStatement = expression
	}

	token = p.nextToken()
	//post expression
	if token.typ == leftBraceType {
		forStatement.postStatement = nopStatement
	} else {
		p.putToken(token)
		expression := p.parseExpression()
		p.expectType(p.nextToken(), leftBraceType)
		forStatement.postStatement = expression
	}
	// statements
	statements := p.ParseStatements()
	forStatement.statements = statements
	return &forStatement
}

func (p *Parser2) ParseStatements() Statements {
	var statements Statements
	if p.ahead(0).typ == rightBraceType {
		return append(statements, NopStatement{})
	}
	for {
		statements = append(statements, p.ParseStatement())
		if p.ahead(0).typ == rightBraceType {
			return statements
		}
	}
}
