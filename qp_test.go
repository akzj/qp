package qp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"testing"
)

type TokenType int

func (t TokenType) String() string {
	switch t {
	case addOperatorTokenType:
		return "+"
	case divOperatorTokenType:
		return "/"
	case subOperatorTokenType:
		return "-"
	case mulOperatorTokenType:
		return "*"
	case intTokenType:
		return "int"
	case leftParenthesisTokenType:
		return "("
	case rightParenthesisTokenType:
		return ")"
	case ifTokenType:
		return "if"
	case elseTokenType:
		return "else"
	case forTokenType:
		return "for"
	case breakTokenType:
		return "break"
	case returnTokenType:
		return "return"
	case leftBraceTokenType:
		return "{"
	case rightBraceTokenType:
		return "}"
	case lessTokenType:
		return "<"
	case lessEqualTokenType:
		return "<="
	case greaterTokenType:
		return ">"
	case greaterEqualTokenType:
		return ">="
	}
	return "unknown token type"
}

const emptyTokenType TokenType = 0
const unknownTokenType TokenType = 1

const addOperatorTokenType TokenType = 101
const subOperatorTokenType TokenType = 102
const mulOperatorTokenType TokenType = 103
const divOperatorTokenType TokenType = 119

const leftParenthesisTokenType TokenType = 120  // (
const rightParenthesisTokenType TokenType = 121 //)
const leftBraceTokenType TokenType = 122        //{
const rightBraceTokenType TokenType = 123       //}
const lessTokenType TokenType = 133             // <
const greaterTokenType TokenType = 134          // >
const lessEqualTokenType TokenType = 135        // <=
const greaterEqualTokenType TokenType = 136     // >=

const ifTokenType TokenType = 230     //if
const elseTokenType TokenType = 331   //else
const funcTokenType TokenType = 332   //else
const returnTokenType TokenType = 333 //else
const breakTokenType TokenType = 334  //else
const forTokenType TokenType = 335    //else

const intTokenType TokenType = 700
const labelTokenType TokenType = 5000

type Token struct {
	typ TokenType
	val string
}

func (t Token) String() string {
	if t.val == "" {
		return fmt.Sprintf("type`%s`", t.typ.String())
	}
	return fmt.Sprintf("type`%s` data `%s`", t.typ.String(), t.val)
}

var (
	emptyToken            = Token{typ: emptyTokenType}
	unknownToken          = Token{typ: unknownTokenType}
	addOperatorToken      = Token{typ: addOperatorTokenType}
	mulOperatorToken      = Token{typ: mulOperatorTokenType}
	leftParenthesisToken  = Token{typ: leftParenthesisTokenType}
	rightParenthesisToken = Token{typ: rightParenthesisTokenType}
	leftBraceToken        = Token{typ: leftBraceTokenType}
	rightBraceToken       = Token{typ: rightBraceTokenType}

)

var Keywords = []string{
	"if", "else", "func", "return", "break", "for",
}

var keywordTokenType = map[string]TokenType{
	"if":     ifTokenType,
	"else":   elseTokenType,
	"func":   funcTokenType,
	"return": returnTokenType,
	"break":  breakTokenType,
	"for":    forTokenType,
}

type Expression interface {
	invoke() (interface{}, error)
}

type Expressions []Expression

func (e *Expressions) invoke() (interface{}, error) {
	var val interface{}
	var err error
	for _, expression := range *e {
		if val, err = expression.invoke(); err != nil {
			return val, err
		}
	}
	return val, err
}

type _leftParenthesisExpression struct {
}
type rightParenthesisExpression struct {
}

func (r rightParenthesisExpression) invoke() (interface{}, error) {
	return r, nil
}

func (l _leftParenthesisExpression) invoke() (interface{}, error) {
	return l, nil
}

type IntExpression struct {
	val int64
}

func (i IntExpression) invoke() (interface{}, error) {
	return i.val, nil
}

type AddExpression struct {
	Left  Expression
	right Expression
}
type MulExpression struct {
	Left  Expression
	right Expression
}

func (expression *MulExpression) invoke() (interface{}, error) {
	fmt.Println("MulExpression invoke")
	l, err := expression.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	r, err := expression.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	switch lVal := l.(type) {
	case int64:
		switch rVal := r.(type) {
		case int64:
			return lVal * rVal, nil
		}
	}
	return nil, fmt.Errorf("unknown operand type")
}

func (a *AddExpression) invoke() (interface{}, error) {
	fmt.Println("AddExpression invoke")
	l, err := a.Left.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	r, err := a.right.invoke()
	if err != nil {
		fmt.Println("invoke left failed", err.Error())
	}
	switch lVal := l.(type) {
	case int64:
		switch rVal := r.(type) {
		case int64:
			return lVal + rVal, nil
		}
	}
	return nil, fmt.Errorf("unknown operand type")
}

type lexer struct {
	scanner *bufio.Scanner
	buffer  bytes.Buffer
	line    int
	token   Token
	err     error
}

func (l *lexer) finish() bool {
	return l.err != nil && l.token.typ == emptyTokenType
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}

func isLetter(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func (l *lexer) peek() Token {
	if l.token.typ != emptyTokenType {
		return l.token
	}
	for {
		if l.buffer.Len() == 0 {
			if l.scanner.Scan() {
				l.buffer.Write(l.scanner.Bytes())
				l.line++
			} else {
				return emptyToken
			}
		}
		c, err := l.buffer.ReadByte()
		if err != nil {
			l.err = err
			return emptyToken
		}
		var token Token
		switch {
		case isSpace(c):
			continue
		case isLetter(c):
			token = l.parseLabel(c)
		case c == '+':
			token = addOperatorToken
		case c == '(':
			token = leftParenthesisToken
		case c == ')':
			token = rightParenthesisToken
		case c == '{':
			token = leftBraceToken
		case c == '}':
			token = rightBraceToken
		case c == '*':
			token = mulOperatorToken
		case '0' <= c && c <= '9':
			token = l.parseNumToken(c)
		default:
			token = unknownToken
		}
		l.token = token
		return token
	}
}

func (l *lexer) parseNumToken(c byte) Token {
	var buf bytes.Buffer
	buf.WriteByte(c)
	for {
		c, err := l.buffer.ReadByte()
		if err != nil {
			l.err = err
			break
		}
		if isDigit(c) {
			buf.WriteByte(c)
		} else {
			if err := l.buffer.UnreadByte(); err != nil {
				panic(err)
			}
			break
		}
	}
	return Token{
		typ: intTokenType,
		val: buf.String(),
	}
}

func (l *lexer) next() {
	l.token = emptyToken
}

func (l *lexer) parseLabel(c byte) Token {
	fmt.Println(string(c))
	var buf bytes.Buffer
	buf.WriteByte(c)
	for {
		c, err := l.buffer.ReadByte()
		if err != nil {
			l.err = err
			break
		}
		fmt.Println(string(c))
		if isLetter(c) {
			buf.WriteByte(c)
		} else {
			if err := l.buffer.UnreadByte(); err != nil {
				panic(err)
			}
			break
		}
	}
	fmt.Println(buf.String())
	for _, keyword := range Keywords {
		if keyword == buf.String() {
			return Token{
				typ: keywordTokenType[keyword],
			}
		}
	}

	return Token{
		typ: labelTokenType,
		val: buf.String(),
	}
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func newLexer(reader io.Reader) *lexer {
	return &lexer{
		token:   emptyToken,
		scanner: bufio.NewScanner(reader),
	}
}

func precedence(tokenType TokenType) int {
	switch tokenType {
	case addOperatorTokenType, subOperatorTokenType:
		return 1
	case mulOperatorTokenType, divOperatorTokenType:
		return 2
	default:
		return 0
	}
}

//greater or equal
func precedenceGE(first, second TokenType) bool {
	return precedence(first)-precedence(second) >= 0
}

func isOperatorToken(token Token) bool {
	return token.typ >= addOperatorTokenType && token.typ < divOperatorTokenType
}

type parser struct {
	stack []Token
	lexer *lexer
}

func (p *parser) EOf() bool {
	return false
}

func newParser(reader io.Reader) *parser {
	return &parser{
		lexer: newLexer(reader),
	}
}

func makeExpression(opToken Token, expressions *[]Expression) Expression {
	var expression Expression
	switch opToken.typ {
	case addOperatorTokenType:
		expression = &AddExpression{
			Left:  (*expressions)[len(*expressions)-1],
			right: (*expressions)[len(*expressions)-2],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	case mulOperatorTokenType:
		expression = &MulExpression{
			Left:  (*expressions)[len(*expressions)-1],
			right: (*expressions)[len(*expressions)-2],
		}
		*expressions = (*expressions)[:len(*expressions)-2]
	}
	return expression
}

func (p *parser) parse() Expression {
	var opStack []Token
	var expressions []Expression

	for p.lexer.finish() == false {
		token := p.lexer.peek()
		fmt.Println(token)
		switch {
		case token.typ == intTokenType:
			val, err := strconv.ParseInt(string(token.val), 10, 64)
			if err != nil {
				fmt.Println("parse int failed", string(token.val))
				return nil
			}
			expressions = append(expressions, IntExpression{val: val})
		case token.typ == leftParenthesisTokenType:
			//expressions = append(expressions, leftParenthesisExpression)
			opStack = append(opStack, token)
		case token.typ == rightParenthesisTokenType:
			fmt.Println(opStack)
			for len(opStack) != 0 && opStack[len(opStack)-1].typ != leftParenthesisTokenType {
				express := makeExpression(opStack[len(opStack)-1], &expressions)
				expressions = append(expressions, express)
				opStack = opStack[:len(opStack)-1]
			}
			opStack = opStack[:len(opStack)-1]
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
		}
		p.lexer.next()
	}
	fmt.Println("===========")
	for len(opStack) != 0 {
		express := makeExpression(opStack[len(opStack)-1], &expressions)
		expressions = append(expressions, express)
		opStack = opStack[:len(opStack)-1]
	}
	return (*Expressions)(&expressions)
}

func Parse(data string) Expression {
	parser := newParser(bytes.NewReader([]byte(data)))
	return parser.parse()
}

func TestName(t *testing.T) {
	fmt.Println("hello qp")
}

func TestBuffer(t *testing.T) {
	var reader = bytes.NewReader([]byte("1+1 if else"))
	for {
		c, err := reader.ReadByte()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(string(c))
		fmt.Println(isLetter(c))
		fmt.Println(c, 'a')
		fmt.Println(c, 'i')
	}
}

func TestLexer(t *testing.T) {
	lexer := newLexer(bytes.NewReader([]byte(`1+1 if else break return for {}
if 1+1 > 3{
	print()
}else{
	print()
}
`)))
	if lexer == nil {
		t.Fatal("lexer nil")
	}
	var count = 100
	for lexer.finish() == false && count > 0 {
		fmt.Println(lexer.peek())
		lexer.next()
		count--
	}
}

func TestParse(t *testing.T) {
	expression := Parse("1*(5+5+5)*2")
	if expression == nil {
		t.Fatal("Parse failed")
	}
	if val, err := expression.invoke(); err != nil {
		t.Errorf(err.Error())
	} else {
		fmt.Println(val)
	}
}
