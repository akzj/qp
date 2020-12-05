package qp

import (
	"fmt"
	"strconv"
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
	case elseifTokenType:
		return "else if"
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
	case emptyTokenType:
		return "empty"
	case labelTokenType:
		return "label"
	case statementTokenType:
		return "statement"
	case statementsTokenType:
		return "statements"
	default:
		panic("unknown token type " + strconv.Itoa(int(t)))
	}
}

const emptyTokenType TokenType = 0
const unknownTokenType TokenType = 1

const addOperatorTokenType TokenType = 101  // +
const subOperatorTokenType TokenType = 102  // -
const mulOperatorTokenType TokenType = 103  // *
const divOperatorTokenType TokenType = 104  // /
const lessTokenType TokenType = 105         // <
const greaterTokenType TokenType = 106      // >
const lessEqualTokenType TokenType = 116    // <=
const greaterEqualTokenType TokenType = 117 // >=

const leftParenthesisTokenType TokenType = 120  // (
const rightParenthesisTokenType TokenType = 121 //)
const leftBraceTokenType TokenType = 122        //{
const rightBraceTokenType TokenType = 123       //}

const ifTokenType TokenType = 230     //if
const elseTokenType TokenType = 331   //else
const funcTokenType TokenType = 332   //func
const returnTokenType TokenType = 333 //return
const breakTokenType TokenType = 334  //break
const forTokenType TokenType = 335    //for
const elseifTokenType TokenType = 336 //else if

const intTokenType TokenType = 700
const labelTokenType TokenType = 5000
const statementTokenType TokenType = 6001
const statementsTokenType TokenType = 6003
const expressionTokenType TokenType = 6002

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
	lessToken             = Token{typ: lessTokenType}
	lessEqualToken        = Token{typ: lessEqualTokenType}
	greaterToken          = Token{typ: greaterTokenType}
	greaterEqualToken     = Token{typ: greaterEqualTokenType}
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
