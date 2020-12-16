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

type Parser struct {
	vm           *VMContext
	lexer        *Lexer
	tokens       []Token
	hTokens      []Token
	pStack       []int //parenthesis stack
	closureCheck []*ClosureCheck
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
	case MulOpType, DivOpType:
		return 10
	case AddType, SubType:
		return 9
	case LessType, LessEqualType, GreaterType, GreaterEqualType, NoEqualType, EqualType:
		return 8
	case AndType:
		return 7
	case OrType:
		return 6
	default:
		return 0
	}
}

func NewParse2(buffer string) *Parser {
	return &Parser{
		status: []PStatus{GlobalStatus},
		lexer:  newLexer(bytes.NewReader([]byte(buffer))),
		vm:     newVMContext(),
	}
}

func Parse(data string) Statements {
	return NewParse2(data).Parse()
}

func (p *Parser) putToken(token Token) {
	p.hTokens = p.hTokens[:len(p.hTokens)-1]
	p.tokens = append([]Token{token}, p.tokens...)
}

func (p *Parser) initTokens() *Parser {
	for {
		token := p.lexer.Peek()
		if token.typ == EOFType {
			p.tokens = append(p.tokens, token)
			break
		}
		if token.typ == ElseType {
			p.lexer.next()
			next := p.lexer.Peek()
			p.lexer.next()
			if next.typ == IfType {
				token.typ = ElseifType
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

func (p *Parser) Parse() Statements {
	p.initTokens()
	var statements Statements
	for {

		if statement := p.ParseStatement(); statement != nil {
			statements = append(statements, statement)
		}
		if p.ahead(0).typ == EOFType {
			return statements
		}
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
func (p *Parser) ParseIDPrefixStatement(token Token) Statement {
	var exp Expression = getVarStatement{
		ctx:   p.vm,
		label: token.val,
	}
	var parentExp Expression
	for {
		token := p.nextToken()
		switch token.typ {
		case AssignType:
			return p.parseAssignStatement(exp)
		case IncType:
			return IncFieldStatement{
				exp: exp,
			}
		case ColonType:
			log.Println("ColonType")
			return AssignStatement{left: exp, exp: p.parseFactor(0)}
		case PeriodType:
			token := p.nextToken()
			p.expectType(token, IDType)
			parentExp = exp
			exp = periodStatement{
				val: token.val,
				exp: exp,
			}
		case LeftParenthesisType:
			exp = p.parseCallStatement(parentExp, exp)
			if p.ahead(0).typ != LeftParenthesisType {
				return exp
			}
		default:
			log.Panic("unexpect token ", token)
		}
	}
}

func (p *Parser) ParseStatement() Statement {
	var statement Statement
	for {
		token := p.nextToken()
		//		log.Println(token)
		switch token.typ {
		case TypeType:
			p.vm.addStructObject(p.parseTypeStatement())
		case FuncType:
			//function
			if p.ahead(0).typ == IDType {
				p.vm.addUserFunction(p.parseFuncStatement())
			} else if p.ahead(0).typ == LeftParenthesisType { //func(){} lambda
				funcStatement := p.parseLambdaStatement()
				//function call
				for p.ahead(0).typ == LeftParenthesisType {
					p.nextToken()
					funcStatement = p.parseCallStatement(nil, funcStatement)
				}
				return funcStatement
			}
		case VarType:
			return p.parseVarStatement()
		case IfType:
			return p.parseIfStatement(false)
		case EOFType:
			if statement == nil {
				return nopStatement
			}
			return statement
		case NewLineType:
			continue
		case SemicolonType:
			continue
		case IDType:
			return p.ParseIDPrefixStatement(token)
		case ReturnType:
			return p.parseReturn()
		case RightBraceType:
			p.putToken(token)
			if statement == nil {
				return nopStatement
			}
			return statement
		case ForType:
			return p.parseForStatement()
		case BreakType:
			return p.parseBreakStatement()
		default:
			log.Panic(token)
		}
	}
}

func (p *Parser) nextToken() Token {
	if len(p.tokens) == 0 {
		return Token{typ: EOFType}
	}
	token := p.tokens[0]
	p.tokens = p.tokens[1:]
	p.hTokens = append(p.hTokens, token)
	if token.typ == CommentType {
		return p.nextToken()
	}
	return token
}

/*
// user define object
type User{
}
*/
func (p *Parser) parseTypeObjectInit() []TypeObjectPropTemplate {
	var objectPropTemplates []TypeObjectPropTemplate
	for {
		token := p.nextToken()
		p.expectType(p.nextToken(), ColonType)
		exp := p.parseFactor(0) //delay bind
		objectPropTemplates = append(objectPropTemplates, TypeObjectPropTemplate{
			name: token.val,
			exp:  exp,
		})
		if ahead := p.ahead(0); ahead.typ == CommaType || ahead.typ == SemicolonType {
			p.nextToken()
		} else {
			//log.Println(p.ahead(0))
			if p.ahead(0).typ == IDType {
				if p.historyToken(1).line != p.ahead(0).line {
					continue
				}
				log.Panic("require new line")
			} else {
				//log.Println(p.ahead(0))
				break
			}
		}
	}
	return objectPropTemplates
}

func (p *Parser) parseTypeStatement() *TypeObject {
	var object TypeObject
	token := p.nextToken()
	//	log.Println(token)
	p.expectType(token, IDType)
	p.expectType(p.nextToken(), LeftBraceType) //{
	object.label = token.val
	if ahead := p.ahead(0); ahead.typ == RightBraceType {
		p.nextToken()
		return &object
	} else {
		object.typeObjectPropTemplates = p.parseTypeObjectInit()
	}
	p.expectType(p.nextToken(), RightBraceType) //}
	return &object
}

func (p *Parser) expectType(token Token, expect Type) {
	if token.typ != expect {
		log.Panicf("unexpect %s token", token)
	}
}

func (p *Parser) parseVarStatement() Statement {
	token := p.nextToken()
	//log.Println(token)
	p.expectType(token, IDType)
	next := p.nextToken()
	//var id = ....
	p.closureCheckAddVar(token.val)
	if next.typ == AssignType {
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
func (p *Parser) parseLambdaStatement() Statement {
	p.pushStatus(FunctionStatus)
	defer func() {
		p.assertTrue(p.popStatus() == FunctionStatus)
	}()
	var funcS FuncStatement
	funcS.closure = true
	funcS.vm = p.vm
	token := p.nextToken()

	p.pushClosureCheck()

	p.expectType(token, LeftParenthesisType)
	funcS.parameters = p.parseFuncParameters()
	p.expectType(p.nextToken(), LeftBraceType)

	if p.ahead(0).typ == RightBraceType { //empty body
		p.nextToken()
		return &funcS
	}
	for {
		funcS.statements = append(funcS.statements, p.ParseStatement())
		if p.ahead(0).typ == RightBraceType {
			p.nextToken()
			break
		}
		//log.Println("lambda next token", p.ahead(0))
	}
	funcS.closureLabel = p.popClosureLabels()
	return &funcS
}

func (p *Parser) parseFactor(pre int) Expression {
	var exp Expression
	var parentExp Expression
	for {
		token := p.nextToken()
		switch token.typ {
		case FalseType:
			p.assertNil(exp)
			exp = Bool(false)
		case TrueType:
			p.assertNil(exp)
			exp = Bool(true)
		case LeftParenthesisType:
			if exp == nil {
				p.pStack = append(p.pStack, token.line)
				exp = ParenthesisExpression{exp: p.parseFactor(pre)}
				p.expectType(p.nextToken(), RightParenthesisType)
			} else {
				//log.Println("parseCallStatement")
				exp = p.parseCallStatement(parentExp, exp)
			}
		case RightParenthesisType: //end of parenthesis ()
			if exp == nil {
				log.Panic("parse bool exp failed")
			}
			p.putToken(token)
			return exp
		case PeriodType:
			token := p.nextToken()
			p.expectType(token, IDType)
			parentExp = exp
			exp = periodStatement{
				val: token.val,
				exp: exp,
			}
		case MulOpType, DivOpType:
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
		case StringType:
			if exp != nil {
				p.putToken(token)
				if p.historyToken(1).line != token.line {
					return exp
				}
			}
			p.assertNil(exp)
			exp = String(token.val)
		case IntType:
			p.assertNil(exp)
			val, _ := strconv.ParseInt(token.val, 10, 64)
			exp = Int(val)
		case FuncType: // func(){}()
			if exp != nil {
				p.putToken(token)
				if p.historyToken(1).line != token.line {
					return exp
				}
			}
			exp = p.parseLambdaStatement()
		case EqualType, // ==
			NoEqualType,      // !=
			GreaterEqualType, // >=
			GreaterType,      // >
			LessEqualType,    // <=
			SubType,          // -
			AddType,          // +
			LessType,         // <
			AndType,          // ||
			OrType:           // ||
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
		case IncType:
			p.assertNoNil(exp)
			exp = &IncFieldStatement{
				exp: exp,
			}
		case NilType:
			p.assertNil(exp)
			exp = nilObject
		case LeftBraceType:
			if status := p.getStatus(); status == IfStatus || status == ForStatus {
				p.putToken(token)
				if exp == nil {
					log.Panicf("unexpect token %s", token)
				}
				//log.Println("return {")
				return exp
			}
			return p.ParseObjInitStatement(exp)
		case LeftBracketType:
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

func (p *Parser) parseBracketStatement(exp Expression) Expression {
	//init array object
	if exp == nil {
		var arrayStatement makeArrayStatement
		for {
			if p.ahead(0).typ == RightBracketType {
				p.nextToken()
				return &arrayStatement
			}
			arrayStatement.initStatements = append(arrayStatement.initStatements, p.parseFactor(0))
			if p.ahead(0).typ == CommaType { // ,
				p.nextToken()
			}
		}
	} else { //Get array field
		index := p.parseFactor(0)
		p.expectType(p.nextToken(), RightBracketType)
		return getArrayElement{
			arrayExp: exp,
			indexExp: index,
		}
	}
}

func (p *Parser) ahead(index int) Token {
	if len(p.tokens) <= index {
		return Token{typ: EOFType}
	}
	return p.tokens[index]
}

func (p *Parser) parseCallStatement(parentExp, function Expression) Expression {
	var call FuncCallStatement
	call.function = function
	call.parentExp = parentExp
	for {
		if p.ahead(0).typ == RightParenthesisType {
			p.nextToken()
			return &call
		}
		call.arguments = append(call.arguments, p.parseFactor(0))
		if p.ahead(0).typ == CommaType { // ,
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
func (p *Parser) parseFuncStatement() *FuncStatement {
	var funcS FuncStatement
	token := p.nextToken()
	p.expectType(token, IDType)

	if ahead := p.ahead(0); ahead.typ == PeriodType {
		p.nextToken()
		next := p.nextToken()
		p.expectType(next, IDType)
		funcS.labels = []string{token.val, next.val}
	} else {
		funcS.label = token.val
	}
	p.expectType(p.nextToken(), LeftParenthesisType)
	if len(funcS.labels) != 0 {
		funcS.parameters = append(funcS.parameters, "this")
	}
	funcS.parameters = append(funcS.parameters, p.parseFuncParameters()...)
	p.expectType(p.nextToken(), LeftBraceType)
	for {
		if p.ahead(0).typ == RightBraceType {
			p.nextToken()
			break
		}
		funcS.statements = append(funcS.statements, p.ParseStatement())
	}
	funcS.vm = p.vm
	return &funcS
}

func (p *Parser) parseFuncParameters() []string {
	var parameters []string
	if p.ahead(0).typ == RightParenthesisType {
		p.nextToken()
		return nil
	}
	for {
		token := p.nextToken()
		p.expectType(token, IDType)
		parameters = append(parameters, token.val)
		p.closureCheckAddVar(token.val)
		if ahead := p.ahead(0); ahead.typ == CommaType {
			p.nextToken()
			continue
		} else if ahead.typ == RightParenthesisType {
			p.nextToken()
			break
		}
	}
	return parameters
}

func (p *Parser) parseBoolExpression(pre int) Expression {
	var exp Expression
	exp = p.parseFactor(0)
	if _, ok := exp.(BinaryBoolExpression); ok {
		return exp
	}
	log.Printf("parse exp bool failed,%s", exp.String())
	return exp
}

func (p *Parser) parseIfStatement(elseif bool) *IfStatement {
	p.pushStatus(IfStatus)
	//log.Println("enter parseIfStatement")
	defer func() {
		p.assertTrue(p.popStatus() == IfStatus)
		//log.Println("out parseIfStatement")
	}()
	var ifS IfStatement
	if p.ahead(0).typ != LeftBraceType {
		ifS.check = p.parseBoolExpression(0)
	}
	p.expectType(p.nextToken(), LeftBraceType)

	if p.ahead(0).typ == RightBraceType {
		p.nextToken()
		return &ifS
	}
	ifS.vm = p.vm
	ifS.statement = p.ParseStatements()
	p.expectType(p.nextToken(), RightBraceType)

	for elseif == false {
		next := p.ahead(0)
		if next.typ == ElseifType {
			p.nextToken()
			statement := p.parseIfStatement(true)
			ifS.elseIfStatements = append(ifS.elseIfStatements, statement)
		} else if next.typ == ElseType {
			p.expectType(p.nextToken(), ElseType)
			p.expectType(p.nextToken(), LeftBraceType)
			ifS.elseStatement = p.ParseStatements()
			p.expectType(p.nextToken(), RightBraceType)
		} else {
			break
		}
	}
	return &ifS

}
func (p *Parser) assertNil(exp Expression) {
	if exp != nil {
		log.Panicf("expect nil")
	}
}

func (p *Parser) assertNoNil(exp Expression) {
	if exp == nil {
		log.Panicf("expect nil")
	}
}

func (p *Parser) isTerminateToken(next Token) bool {
	if next.typ == IfType ||
		next.typ == ForType ||
		next.typ == LeftBraceType ||
		next.typ == RightBraceType ||
		next.typ == SemicolonType ||
		next.typ == RightBracketType ||
		next.typ == CommaType ||
		next.typ == VarType ||
		next.typ == BreakType ||
		next.typ == ReturnType ||
		next.typ == EOFType {
		return true
	}
	return false
}

func (p *Parser) historyToken(i int) Token {
	index := len(p.hTokens) - i - 1
	if index < len(p.hTokens) && index >= 0 {
		return p.hTokens[index]
	}
	log.Panic("out of history tokens range")
	return EmptyToken
}

func (p *Parser) closureCheckAddVar(data string) {
	if len(p.closureCheck) != 0 {
		p.closureCheck[len(p.closureCheck)-1].AddVar(data)
	}
}

func (p *Parser) closureCheckVisit(data string) {
	for i := len(p.closureCheck) - 1; i >= 0; i-- {
		closure := p.closureCheck[i]
		if closure.Visit(data) == false {
			break
		}
	}
}

func (p *Parser) pushClosureCheck() {
	p.closureCheck = append(p.closureCheck, NewClosureCheck())
}

func (p *Parser) popClosureLabels() []string {
	if len(p.closureCheck) != 0 {
		closureLabel := p.closureCheck[len(p.closureCheck)-1].closures
		p.closureCheck = p.closureCheck[:len(p.closureCheck)-1]
		return closureLabel
	}
	return nil
}

func (p *Parser) parseReturn() Statement {
	if p.ahead(0).typ == RightBraceType {
		return ReturnStatement{
			express:   nilObject,
			returnVal: nilObject,
		}
	}
	return ReturnStatement{
		express: p.parseFactor(0),
	}
}

func (p *Parser) parseForStatement() Statement {
	p.pushStatus(ForStatus)
	defer func() {
		p.assertTrue(p.popStatus() == ForStatus)
	}()
	var forStatement = ForStatement{
		vm: p.vm,
	}
	token := p.nextToken()
	//log.Println(token)
	if token.typ == SemicolonType {
		forStatement.preStatement = nopStatement
	} else if token.typ == LeftBraceType {
		forStatement.preStatement = nopStatement
		forStatement.postStatement = nopStatement
		forStatement.checkStatement = &trueObject
		statements := p.ParseStatements()
		p.expectType(p.nextToken(), RightBraceType)
		forStatement.statements = statements
		return &forStatement
	} else {
		//support var= ;
		expression := p.parseVarStatement()
		if p.nextToken().typ != SemicolonType {
			log.Panic("expect ; in `for` statement")
			return nil
		}
		forStatement.preStatement = expression
	}
	token = p.nextToken()
	//check exp
	if token.typ == SemicolonType {
		forStatement.checkStatement = &trueObject
	} else {
		p.putToken(token)
		expression := p.parseFactor(0)
		if p.nextToken().typ != SemicolonType {
			log.Panic("expect ; in `for` check exp")
		}
		forStatement.checkStatement = expression
	}

	token = p.nextToken()
	//post exp
	if token.typ == LeftBraceType {
		forStatement.postStatement = nopStatement
	} else {
		p.putToken(token)
		expression := p.ParseStatement()
		p.expectType(p.nextToken(), LeftBraceType)
		forStatement.postStatement = expression
	}
	// statements
	statements := p.ParseStatements()
	forStatement.statements = statements
	p.expectType(p.nextToken(), RightBraceType)
	return &forStatement
}

func (p *Parser) ParseStatements() Statements {
	var statements Statements
	if p.ahead(0).typ == RightBraceType {
		return append(statements, NopStatement{})
	}
	for {
		statement := p.ParseStatement()
		statements = append(statements, statement)
		if p.ahead(0).typ == RightBraceType {
			if p.getStatus() == GlobalStatus {
				p.nextToken()
				continue
			}
			return statements
		}
	}
}

func (p *Parser) ParseObjInitStatement(exp Expression) Expression {
	var statement objectInitStatement
	statement.exp = exp
	statement.vm = p.vm
	if ahead := p.ahead(0); ahead.typ == RightBraceType {
		p.nextToken()
		return &statement
	} else {
		statement.propTemplates = p.parseTypeObjectInit()
	}
	p.expectType(p.nextToken(), RightBraceType)
	return &statement
}

func (p *Parser) pushStatus(status PStatus) {
	p.status = append(p.status, status)
}
func (p *Parser) popStatus() PStatus {
	if len(p.status) == 0 {
		log.Panic("status stack empty")
	}
	status := p.status[len(p.status)-1]
	p.status = p.status[:len(p.status)-1]
	return status
}

func (p *Parser) assertTrue(val bool) {
	if val == false {
		panic("assert failed")
	}
}

func (p *Parser) checkInStatus(status PStatus) bool {
	for i := len(p.status) - 1; i >= 0; i-- {
		if p.status[i] == status {
			return true
		} else if p.status[i] == FunctionStatus {
			break
		}
	}
	return false
}

func (p Parser) getStatus() PStatus {
	if len(p.status) != 0 {
		return p.status[len(p.status)-1]
	}
	log.Panic("status stack empty")
	return 0
}

func (p *Parser) parseBreakStatement() Statement {
	if p.checkInStatus(ForStatus) == false {
		log.Panic("current no in `for` brock ")
	}
	return breakObject
}

func (p *Parser) parseAssignStatement(exp Expression) AssignStatement {
	return AssignStatement{
		exp:  p.parseFactor(0),
		left: exp,
	}
}
