package qp

import (
	"fmt"
	"strconv"
)

type Type int

func (t Type) String() string {
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
	case EOFTokenType:
		return "EOF"
	case labelType:
		return "label"
	case statementType:
		return "statement"
	case statementsType:
		return "statements"
	case expressionType:
		return "expression"
	case varTokenType:
		return "var"
	case assignTokenType:
		return "="
	case varAssignTokenType:
		return "var ="
	case IntObjectType:
		return "IntObject"
	case commaTokenType:
		return ","
	case incOperatorTokenType:
		return "++"
	case callFunctionType:
		return "call"
	default:
		panic("unknown token type " + strconv.Itoa(int(t)))
	}
}

const EOFTokenType Type = 0
const unknownTokenType Type = 1

const incOperatorTokenType Type = 100      // ++
const addOperatorTokenType Type = 101      // +
const subOperatorTokenType Type = 102      // -
const mulOperatorTokenType Type = 103      // *
const divOperatorTokenType Type = 104      // /
const lessTokenType Type = 105             // <
const greaterTokenType Type = 106          // >
const lessEqualTokenType Type = 116        // <=
const greaterEqualTokenType Type = 117     // >=
const leftParenthesisTokenType Type = 120  // (
const rightParenthesisTokenType Type = 121 // )
const leftBraceTokenType Type = 122        // {
const rightBraceTokenType Type = 123       // }
const commaTokenType Type = 124            // ,
const ifTokenType Type = 230               //if
const elseTokenType Type = 331             //else
const funcTokenType Type = 332             //func
const returnTokenType Type = 333           //return
const breakTokenType Type = 334            //break
const forTokenType Type = 335              //for
const elseifTokenType Type = 336           //else if
const varTokenType Type = 400              // var
const assignTokenType Type = 401           // =
const varAssignTokenType Type = 402        // var x =
const intTokenType Type = 700              // int
const labelType Type = 5000                // label
const statementType Type = 6001            // statement
const statementsType Type = 6003           // statementType
const expressionType Type = 6002           // expressionType
const callFunctionType Type = 6004         // call function

type Token struct {
	typ Type
	val string
}

func (t Token) String() string {
	if t.val == "" {
		return fmt.Sprintf("type`%s`", t.typ.String())
	}
	return fmt.Sprintf("type`%s` data `%s`", t.typ.String(), t.val)
}

var (
	emptyToken            = Token{typ: EOFTokenType}
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
	assignToken           = Token{typ: assignTokenType}
	commaToken            = Token{typ: commaTokenType}
	incOperatorToken      = Token{typ: incOperatorTokenType}
)

var Keywords = []string{
	"if", "else", "func", "return", "break", "for", "var",
}

var keywordTokenType = map[string]Type{
	"if":     ifTokenType,
	"else":   elseTokenType,
	"func":   funcTokenType,
	"return": returnTokenType,
	"break":  breakTokenType,
	"for":    forTokenType,
	"var":    varTokenType,
}
