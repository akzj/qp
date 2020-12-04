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
	}
	return "unknown token type"
}

const emptyTokenType TokenType = 0
const unknownTokenType TokenType = 1

const operatorTokenType TokenType = 100
const addOperatorTokenType TokenType = 101
const subOperatorTokenType TokenType = 102
const mulOperatorTokenType TokenType = 103
const divOperatorTokenType TokenType = 104
const operatorEndTokenType TokenType = 110

const intTokenType TokenType = 200

type Token struct {
	typ TokenType
	val []byte
}

func (t Token) String() string {
	return fmt.Sprintf("type(%s) data(%s)", t.typ.String(), string(t.val))
}

var (
	emptyToken       = Token{typ: emptyTokenType, val: []byte(`empty`)}
	unknownToken     = Token{typ: unknownTokenType, val: []byte(`unknown`)}
	addOperatorToken = Token{typ: addOperatorTokenType, val: []byte{'+'}}
	mulOperatorToken = Token{typ: mulOperatorTokenType, val: []byte{'*'}}
)

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

func (l *lexer) peek() Token {
	if l.token.typ != emptyTokenType {
		return l.token
	}
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
	case c == '+':
		token = addOperatorToken
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

func (l *lexer) parseNumToken(c byte) Token {
	var result bytes.Buffer
	result.WriteByte(c)
	for {
		c, err := l.buffer.ReadByte()
		if err != nil {
			l.err = err
			break
		}
		if '0' <= c && c <= '9' {
			result.WriteByte(c)
		} else {
			if err := l.buffer.UnreadByte(); err != nil {
				panic(err)
			}
			break
		}
	}
	return Token{
		typ: intTokenType,
		val: result.Bytes(),
	}
}

func (l *lexer) next() {
	l.token = emptyToken
}

func newLexer(reader io.Reader) *lexer {
	return &lexer{
		token:   emptyToken,
		scanner: bufio.NewScanner(reader),
	}
}

func operatorPriority(tokenType TokenType) int {
	switch tokenType {
	case addOperatorTokenType, subOperatorTokenType:
		return 1
	case mulOperatorTokenType, divOperatorTokenType:
		return 2
	default:
		return 0
	}
}

func opGreaterOrEqual(first, second TokenType) bool {
	return operatorPriority(first)-operatorPriority(second) >= 0
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
		fmt.Println(string(token.val))
		switch {
		case token.typ == intTokenType:
			val, err := strconv.ParseInt(string(token.val), 10, 64)
			if err != nil {
				fmt.Println("parse int failed", string(token.val))
				return nil
			}
			expressions = append(expressions, IntExpression{val: val})
		case token.typ > operatorTokenType && token.typ < operatorEndTokenType:
			for len(opStack) != 0 && opGreaterOrEqual(opStack[len(opStack)-1].typ, token.typ) {
				express := makeExpression(opStack[len(opStack)-1], &expressions)
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
	var reader = bytes.NewReader([]byte("1+1"))
	for {
		c, err := reader.ReadByte()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(string(c))
	}
}

func TestLexer(t *testing.T) {
	lexer := newLexer(bytes.NewReader([]byte("1+1")))
	if lexer == nil {
		t.Fatal("lexer nil")
	}
	for lexer.finish() == false {
		fmt.Println(lexer.peek())
		lexer.next()
	}
}

func TestParse(t *testing.T) {
	expression := Parse("1*2+5*3")
	if expression == nil {
		t.Fatal("Parse failed")
	}
	if val, err := expression.invoke(); err != nil {
		t.Errorf(err.Error())
	} else {
		fmt.Println(val)
	}
}
