package parser

import (
	"bytes"
	"log"
	"strconv"
	"strings"

	"gitlab.com/akzj/qp/ast"
	_ "gitlab.com/akzj/qp/builtin"
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

/*
   Invokable:Invokable+Invokable
		|Invokable-Invokable
   		|Invokable
		|Factor


   FunCall:id()
   		|lambda()
		|Selector()

   lambda:func(Expressions)BlockStatement

   Expressions:Invokable,Invokable
   		|Expressions,Invokable

   BlockStatement:{Statements}

   Statements:Statements \n Statements
   		|Statements;Statements

Statements:IfStatement
		|forStatement
		|assignStatement
		|breakStatement
		|varStatement

varStatement:var ID = Invokable

assignStatement:ID = Invokable
		|Selector = Invokable



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
	vm           *runtime.VMRuntime
	lexer        *lexer.Lexer
	tokens       []lexer.Token
	hTokens      []lexer.Token
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

func precedence(tokenType lexer.Type) int {
	switch tokenType {
	case lexer.MulOpType, lexer.DivOpType:
		return 10
	case lexer.AddType, lexer.SubType:
		return 9
	case lexer.LessType,
		lexer.LessEqualType,
		lexer.GreaterType,
		lexer.GreaterEqualType,
		lexer.NoEqualType,
		lexer.EqualType:
		return 8
	case lexer.AndType:
		return 7
	case lexer.OrType:
		return 6
	default:
		return 0
	}
}

func New(buffer string) *Parser {
	return &Parser{
		status: []PStatus{GlobalStatus},
		lexer:  lexer.New(bytes.NewReader([]byte(buffer))),
		vm:     runtime.New(),
	}
}

func Parse(data string) ast.Expressions {
	return New(data).Parse()
}
func (p *Parser) GetVMContext() *runtime.VMRuntime {
	return p.vm
}
func (p *Parser) putToken(token lexer.Token) {
	p.hTokens = p.hTokens[:len(p.hTokens)-1]
	p.tokens = append([]lexer.Token{token}, p.tokens...)
}

func (p *Parser) initTokens() *Parser {
	for {
		token := p.lexer.Peek()
		if token.Typ == lexer.EOFType {
			p.tokens = append(p.tokens, token)
			break
		}
		if token.Typ == lexer.ElseType {
			p.lexer.Next()
			next := p.lexer.Peek()
			p.lexer.Next()
			if next.Typ == lexer.IfType {
				token.Typ = lexer.ElseifType
				p.tokens = append(p.tokens, token)
			} else {
				p.tokens = append(p.tokens, token, next)
			}
			continue
		}
		p.tokens = append(p.tokens, token)
		p.lexer.Next()
	}
	return p
}

func (p *Parser) Parse() ast.Expressions {
	p.initTokens()
	var statements ast.Expressions
	for {
		if statement := p.ParseStatement(); statement != nil {
			statements = append(statements, statement)
		}
		if p.ahead(0).Typ == lexer.EOFType {
			return statements
		}
	}
}

/*
	AssignStatement:
		|ID = Invokable
		|selector = Invokable

	FunctionCallStatement:
		|ID()
		|selector()
		|FunctionCallStatement()
		|lambda
*/
func (p *Parser) ParseIDPrefixStatement(token lexer.Token) ast.Expression {
	var exp runtime.Invokable = ast.GetVarStatement{
		VM:    p.vm,
		Label: token.Val,
	}
	var parentExp runtime.Invokable
	for {
		token := p.nextToken()
		switch token.Typ {
		case lexer.AssignType:
			return p.parseAssignStatement(exp)
		case lexer.IncType:
			return ast.IncFieldStatement{
				Exp: exp,
			}
		case lexer.ColonType:
			log.Println("ColonType")
			return ast.AssignStatement{Left: exp, Exp: p.parseFactor(0)}
		case lexer.PeriodType:
			token := p.nextToken()
			p.expectType(token, lexer.IDType)
			parentExp = exp
			exp = ast.PeriodStatement{
				Val: token.Val,
				Exp: exp,
			}
		case lexer.LeftParenthesisType:
			exp = p.parseCallStatement(parentExp, exp)
			if p.ahead(0).Typ != lexer.LeftParenthesisType {
				return exp
			}
		default:
			log.Panic("unexpect token ", token)
		}
	}
}

func (p *Parser) ParseStatement() ast.Expression {
	var statement ast.Expression
	for {
		token := p.nextToken()
		//		log.Println(token)
		switch token.Typ {
		case lexer.TypeType:
			typeObject := p.parseTypeStatement()
			p.vm.AddStructObject(&runtime.Object{
				Pointer: typeObject,
				Label:   typeObject.Label,
			})
		case lexer.FuncType:
			//function
			if p.ahead(0).Typ == lexer.IDType {
				p.addUserFunction(p.parseFuncStatement())
			} else if p.ahead(0).Typ == lexer.LeftParenthesisType { //func(){} lambda
				funcStatement := p.parseLambdaStatement()
				//function Call
				for p.ahead(0).Typ == lexer.LeftParenthesisType {
					p.nextToken()
					funcStatement = p.parseCallStatement(nil, funcStatement)
				}
				return funcStatement
			}
		case lexer.VarType:
			return p.parseVarStatement()
		case lexer.IfType:
			return p.parseIfStatement(false)
		case lexer.EOFType:
			if statement == nil {
				return ast.NopStatement{}
			}
			return statement
		case lexer.NewLineType:
			continue
		case lexer.SemicolonType:
			continue
		case lexer.IDType:
			return p.ParseIDPrefixStatement(token)
		case lexer.ReturnType:
			return p.parseReturn()
		case lexer.RightBraceType:
			p.putToken(token)
			if statement == nil {
				return ast.NopStatement{}
			}
			return statement
		case lexer.ForType:
			return p.parseForStatement()
		case lexer.BreakType:
			return p.parseBreakStatement()
		default:
			log.Panic(token)
		}
	}
}

func (p *Parser) nextToken() lexer.Token {
	if len(p.tokens) == 0 {
		return lexer.Token{Typ: lexer.EOFType}
	}
	token := p.tokens[0]
	p.tokens = p.tokens[1:]
	p.hTokens = append(p.hTokens, token)
	if token.Typ == lexer.CommentType {
		return p.nextToken()
	}
	return token
}

/*
// user define object
type User{
}
*/
func (p *Parser) parseTypeObjectInit() []ast.TypeObjectPropTemplate {
	var objectPropTemplates []ast.TypeObjectPropTemplate
	for {
		token := p.nextToken()
		p.expectType(p.nextToken(), lexer.ColonType)
		exp := p.parseFactor(0) //delay bind
		objectPropTemplates = append(objectPropTemplates, ast.TypeObjectPropTemplate{
			Name: token.Val,
			Exp:  exp,
		})
		if ahead := p.ahead(0); ahead.Typ == lexer.CommaType || ahead.Typ == lexer.SemicolonType {
			p.nextToken()
		} else {
			//log.Println(p.ahead(0))
			if p.ahead(0).Typ == lexer.IDType {
				if p.historyToken(1).Line != p.ahead(0).Line {
					continue
				}
				log.Panic("require new Line")
			} else {
				//log.Println(p.ahead(0))
				break
			}
		}
	}
	return objectPropTemplates
}

func (p *Parser) parseTypeStatement() *ast.TypeObject {
	var object ast.TypeObject
	token := p.nextToken()
	//	log.Println(token)
	p.expectType(token, lexer.IDType)
	p.expectType(p.nextToken(), lexer.LeftBraceType) //{
	object.Label = token.Val
	if ahead := p.ahead(0); ahead.Typ == lexer.RightBraceType {
		p.nextToken()
		return &object
	} else {
		object.TypeObjectPropTemplates = p.parseTypeObjectInit()
	}
	p.expectType(p.nextToken(), lexer.RightBraceType) //}
	return &object
}

func (p *Parser) expectType(token lexer.Token, expect lexer.Type) {
	if token.Typ != expect {
		log.Panicf("unexpect %s token", token)
	}
}

func (p *Parser) parseVarStatement() ast.Expression {
	token := p.nextToken()
	//log.Println(token)
	p.expectType(token, lexer.IDType)
	next := p.nextToken()
	//var id = ....
	p.closureCheckAddVar(token.Val)
	if next.Typ == lexer.AssignType {
		expression := p.parseFactor(0)
		return ast.VarAssignStatement{
			Ctx:  p.vm,
			Name: token.Val,
			Exp:  expression,
		}
	}
	p.putToken(next)
	return ast.VarStatement{
		VM:    p.vm,
		Label: token.Val,
	}
}

/*
func(){}()

*/
func (p *Parser) parseLambdaStatement() ast.Expression {
	p.pushStatus(FunctionStatus)
	defer func() {
		p.assertTrue(p.popStatus() == FunctionStatus)
	}()
	var funcS ast.FuncExpression
	funcS.Closure = true
	funcS.VM = p.vm
	token := p.nextToken()

	p.pushClosureCheck()

	p.expectType(token, lexer.LeftParenthesisType)
	funcS.Parameters = p.parseFuncParameters()
	p.expectType(p.nextToken(), lexer.LeftBraceType)

	if p.ahead(0).Typ == lexer.RightBraceType { //empty body
		p.nextToken()
		return &funcS
	}
	for {
		funcS.Statements = append(funcS.Statements, p.ParseStatement())
		if p.ahead(0).Typ == lexer.RightBraceType {
			p.nextToken()
			break
		}
		//log.Println("lambda Next token", p.ahead(0))
	}
	funcS.ClosureLabel = p.popClosureLabels()
	return &funcS
}

func (p *Parser) parseFactor(pre int) runtime.Invokable {
	var exp runtime.Invokable
	var parentExp runtime.Invokable
	for {
		token := p.nextToken()
		switch token.Typ {
		case lexer.FalseType:
			p.assertNil(exp)
			exp = ast.Bool(false)
		case lexer.TrueType:
			p.assertNil(exp)
			exp = ast.Bool(true)
		case lexer.LeftParenthesisType:
			if exp == nil {
				p.pStack = append(p.pStack, token.Line)
				exp = ast.ParenthesisExpression{Exp: p.parseFactor(pre)}
				p.expectType(p.nextToken(), lexer.RightParenthesisType)
			} else {
				//log.Println("parseCallStatement")
				exp = p.parseCallStatement(parentExp, exp)
			}
		case lexer.RightParenthesisType: //end of parenthesis ()
			if exp == nil {
				log.Panic("parse bool exp failed")
			}
			p.putToken(token)
			return exp
		case lexer.PeriodType:
			token := p.nextToken()
			p.expectType(token, lexer.IDType)
			parentExp = exp
			exp = ast.PeriodStatement{
				Val: token.Val,
				Exp: exp,
			}
		case lexer.MulOpType, lexer.DivOpType:
			if exp == nil {
				log.Panicf("unexpect token %s", token)
			}
			p.nextToken()
			return ast.BinaryOpExpression{
				OP:    token.Typ,
				Left:  exp,
				Right: p.parseFactor(precedence(token.Typ)),
			}
		case lexer.NoType:
			p.assertNoNil(exp)
			exp = ast.NoStatement{Exp: p.parseFactor(pre)}
		case lexer.IDType:
			if exp != nil {
				p.putToken(token)
				if p.historyToken(1).Line != token.Line {
					return exp
				}
			}
			p.assertNil(exp)
			exp = ast.GetVarStatement{
				VM:    p.vm,
				Label: token.Val,
			}
			p.closureCheckVisit(token.Val)
		case lexer.StringType:
			if exp != nil {
				p.putToken(token)
				if p.historyToken(1).Line != token.Line {
					return exp
				}
			}
			p.assertNil(exp)
			exp = ast.String(token.Val)
		case lexer.IntType:
			p.assertNil(exp)
			val, _ := strconv.ParseInt(token.Val, 10, 64)
			exp = ast.Int(val)
		case lexer.FuncType: // func(){}()
			if exp != nil {
				p.putToken(token)
				if p.historyToken(1).Line != token.Line {
					return exp
				}
			}
			exp = p.parseLambdaStatement()
		case lexer.EqualType, // ==
			lexer.NoEqualType,      // !=
			lexer.GreaterEqualType, // >=
			lexer.GreaterType,      // >
			lexer.LessEqualType,    // <=
			lexer.SubType,          // -
			lexer.AddType,          // +
			lexer.LessType,         // <
			lexer.AndType,          // &&
			lexer.OrType:           // ||
			if exp == nil {
				log.Panic("exp nil")
			}
			if pre >= precedence(token.Typ) {
				//				log.Println("return", pre, precedence(token.Typ))
				p.putToken(token)
				return exp
			}
			exp = ast.BinaryOpExpression{
				OP:    token.Typ,
				Left:  exp,
				Right: p.parseFactor(precedence(token.Typ)),
			}
		case lexer.IncType:
			p.assertNoNil(exp)
			exp = ast.IncFieldStatement{
				Exp: exp,
			}
		case lexer.NilType:
			p.assertNil(exp)
			exp = ast.NilObject{}
		case lexer.LeftBraceType:
			if status := p.getStatus(); status == IfStatus || status == ForStatus {
				p.putToken(token)
				if exp == nil {
					log.Panicf("unexpect token %s", token)
				}
				//log.Println("return {")
				return exp
			}
			return p.ParseObjInitStatement(exp)
		case lexer.LeftBracketType:
			exp = p.parseBracketStatement(exp)
		case lexer.EOFType:
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

func (p *Parser) parseBracketStatement(exp runtime.Invokable) runtime.Invokable {
	//init array object
	if exp == nil {
		var arrayStatement ast.MakeArrayStatement
		for {
			if p.ahead(0).Typ == lexer.RightBracketType {
				p.nextToken()
				return &arrayStatement
			}
			arrayStatement.Inits = append(arrayStatement.Inits, p.parseFactor(0))
			if p.ahead(0).Typ == lexer.CommaType { // ,
				p.nextToken()
			}
		}
	} else { //Get array field
		index := p.parseFactor(0)
		p.expectType(p.nextToken(), lexer.RightBracketType)
		return ast.ArrayGetElement{
			Exp:   exp,
			Index: index,
		}
	}
}

func (p *Parser) ahead(index int) lexer.Token {
	if len(p.tokens) <= index {
		return lexer.Token{Typ: lexer.EOFType}
	}
	return p.tokens[index]
}

func (p *Parser) parseCallStatement(parentExp, function runtime.Invokable) runtime.Invokable {
	var call ast.CallStatement
	call.Function = function
	call.ParentExp = parentExp
	for {
		if p.ahead(0).Typ == lexer.RightParenthesisType {
			p.nextToken()
			return &call
		}
		call.Arguments = append(call.Arguments, p.parseFactor(0))
		if p.ahead(0).Typ == lexer.CommaType { // ,
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
func (p *Parser) parseFuncStatement() *ast.FuncExpression {
	var funcS ast.FuncExpression
	token := p.nextToken()
	p.expectType(token, lexer.IDType)

	if ahead := p.ahead(0); ahead.Typ == lexer.PeriodType {
		p.nextToken()
		next := p.nextToken()
		p.expectType(next, lexer.IDType)
		funcS.Labels = []string{token.Val, next.Val}
	} else {
		funcS.Label = token.Val
	}
	p.expectType(p.nextToken(), lexer.LeftParenthesisType)
	if len(funcS.Labels) != 0 {
		funcS.Parameters = append(funcS.Parameters, "this")
	}
	funcS.Parameters = append(funcS.Parameters, p.parseFuncParameters()...)
	p.expectType(p.nextToken(), lexer.LeftBraceType)
	for {
		if p.ahead(0).Typ == lexer.RightBraceType {
			p.nextToken()
			break
		}
		funcS.Statements = append(funcS.Statements, p.ParseStatement())
	}
	funcS.VM = p.vm
	return &funcS
}

func (p *Parser) parseFuncParameters() []string {
	var parameters []string
	if p.ahead(0).Typ == lexer.RightParenthesisType {
		p.nextToken()
		return nil
	}
	for {
		token := p.nextToken()
		p.expectType(token, lexer.IDType)
		parameters = append(parameters, token.Val)
		p.closureCheckAddVar(token.Val)
		if ahead := p.ahead(0); ahead.Typ == lexer.CommaType {
			p.nextToken()
			continue
		} else if ahead.Typ == lexer.RightParenthesisType {
			p.nextToken()
			break
		}
	}
	return parameters
}

func (p *Parser) parseBoolExpression(pre int) runtime.Invokable {
	var exp runtime.Invokable
	exp = p.parseFactor(0)
	if _, ok := exp.(ast.BinaryOpExpression); ok {
		return exp
	}
	log.Printf("parse exp bool failed,%s", exp.String())
	return exp
}

func (p *Parser) parseIfStatement(elseif bool) ast.Expression {
	p.pushStatus(IfStatus)
	//log.Println("enter parseIfStatement")
	defer func() {
		p.assertTrue(p.popStatus() == IfStatus)
		//log.Println("out parseIfStatement")
	}()
	var ifS = ast.IfExpression{
		VM: p.vm,
	}
	if p.ahead(0).Typ != lexer.LeftBraceType {
		ifS.Check = p.parseBoolExpression(0)
	}
	p.expectType(p.nextToken(), lexer.LeftBraceType)

	if p.ahead(0).Typ == lexer.RightBraceType {
		p.nextToken()
		return ifS
	}
	ifS.VM = p.vm
	ifS.Statements = p.ParseStatements()
	p.expectType(p.nextToken(), lexer.RightBraceType)

	for !elseif {
		next := p.ahead(0)
		if next.Typ == lexer.ElseifType {
			p.nextToken()
			statement := p.parseIfStatement(true)
			ifS.ElseIf = append(ifS.ElseIf, statement.(ast.IfExpression))
		} else if next.Typ == lexer.ElseType {
			p.expectType(p.nextToken(), lexer.ElseType)
			p.expectType(p.nextToken(), lexer.LeftBraceType)
			ifS.Else = p.ParseStatements()
			p.expectType(p.nextToken(), lexer.RightBraceType)
		} else {
			break
		}
	}
	return ifS

}
func (p *Parser) assertNil(exp runtime.Invokable) {
	if exp != nil {
		log.Panicf("expect nil")
	}
}

func (p *Parser) assertNoNil(exp runtime.Invokable) {
	if exp == nil {
		log.Panicf("expect nil")
	}
}

func (p *Parser) isTerminateToken(next lexer.Token) bool {
	if next.Typ == lexer.IfType ||
		next.Typ == lexer.ForType ||
		next.Typ == lexer.LeftBraceType ||
		next.Typ == lexer.RightBraceType ||
		next.Typ == lexer.SemicolonType ||
		next.Typ == lexer.RightBracketType ||
		next.Typ == lexer.CommaType ||
		next.Typ == lexer.VarType ||
		next.Typ == lexer.BreakType ||
		next.Typ == lexer.ReturnType ||
		next.Typ == lexer.TypeType ||
		next.Typ == lexer.EOFType {
		return true
	}
	return false
}

func (p *Parser) historyToken(i int) lexer.Token {
	index := len(p.hTokens) - i - 1
	if index < len(p.hTokens) && index >= 0 {
		return p.hTokens[index]
	}
	log.Panic("out of history tokens range")
	return lexer.EmptyToken
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

func (p *Parser) parseReturn() ast.Expression {
	if p.ahead(0).Typ == lexer.RightBraceType {
		return ast.ReturnStatement{
			Exp: ast.NilObject{},
			Val: ast.NilObject{},
		}
	}
	return ast.ReturnStatement{
		Exp: p.parseFactor(0),
	}
}

func (p *Parser) parseForStatement() ast.Expression {
	p.pushStatus(ForStatus)
	defer func() {
		p.assertTrue(p.popStatus() == ForStatus)
	}()
	var forStatement = ast.ForExpression{
		VM: p.vm,
	}
	token := p.nextToken()
	//log.Println(token)
	if token.Typ == lexer.SemicolonType {
		forStatement.Pre = ast.NopStatement{}
	} else if token.Typ == lexer.LeftBraceType {
		forStatement.Pre = ast.NopStatement{}
		forStatement.Post = ast.NopStatement{}
		forStatement.Check = &ast.TrueObject
		statements := p.ParseStatements()
		p.expectType(p.nextToken(), lexer.RightBraceType)
		forStatement.Statements = statements
		return &forStatement
	} else {
		//support var= ;
		expression := p.parseVarStatement()
		if p.nextToken().Typ != lexer.SemicolonType {
			log.Panic("expect ; in `for` statement")
			return nil
		}
		forStatement.Pre = expression
	}
	token = p.nextToken()
	//check exp
	if token.Typ == lexer.SemicolonType {
		forStatement.Check = &ast.TrueObject
	} else {
		p.putToken(token)
		expression := p.parseFactor(0)
		if p.nextToken().Typ != lexer.SemicolonType {
			log.Panic("expect ; in `for` check exp")
		}
		forStatement.Check = expression
	}

	token = p.nextToken()
	//post exp
	if token.Typ == lexer.LeftBraceType {
		forStatement.Post = ast.NopStatement{}
	} else {
		p.putToken(token)
		expression := p.ParseStatement()
		p.expectType(p.nextToken(), lexer.LeftBraceType)
		forStatement.Post = expression
	}
	// statements
	statements := p.ParseStatements()
	forStatement.Statements = statements
	p.expectType(p.nextToken(), lexer.RightBraceType)
	return forStatement
}

func (p *Parser) ParseStatements() ast.Expressions {
	var statements ast.Expressions
	if p.ahead(0).Typ == lexer.RightBraceType {
		return append(statements, ast.NopStatement{})
	}
	for {
		statement := p.ParseStatement()
		statements = append(statements, statement)
		if p.ahead(0).Typ == lexer.RightBraceType {
			if p.getStatus() == GlobalStatus {
				p.nextToken()
				continue
			}
			return statements
		}
	}
}

func (p *Parser) ParseObjInitStatement(exp runtime.Invokable) runtime.Invokable {
	var statement ast.ObjectInitStatement
	statement.Exp = exp
	statement.VM = p.vm
	if ahead := p.ahead(0); ahead.Typ == lexer.RightBraceType {
		p.nextToken()
		return statement
	} else {
		statement.PropTemplates = p.parseTypeObjectInit()
	}
	p.expectType(p.nextToken(), lexer.RightBraceType)
	return statement
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

func (p *Parser) parseBreakStatement() ast.Expression {
	if p.checkInStatus(ForStatus) == false {
		log.Panic("current no in `for` brock ")
	}
	return ast.BreakObj
}

func (p *Parser) parseAssignStatement(exp runtime.Invokable) ast.AssignStatement {
	return ast.AssignStatement{
		Exp:  p.parseFactor(0),
		Left: exp,
	}
}

func (p *Parser) addUserFunction(function *ast.FuncExpression) {
	if function.Labels != nil {
		structObject := p.vm.GetTypeObject(function.Labels[0])
		if structObject == nil { //todo fixme
			log.Panic("no find structObject", function.Labels[0])
		}
		structObject.Pointer.(*ast.TypeObject).AddObject(function.Labels[1], &runtime.Object{
			Pointer: function,
			Label:   strings.Join(function.Labels, "."),
		})
		return
	}
	if _, ok := runtime.Functions[function.Label]; ok {
		log.Panic("function name conflict with built in function", function.Label)
	}
	p.vm.AddGlobalFunction(&runtime.Object{
		Pointer: function,
		Label:   function.Label,
	})
}
