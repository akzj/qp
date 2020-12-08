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
	case semicolonTokenType:
		return ";"
	case assignStatementType:
		return "assignStatement"
	case funcTokenType:
		return "func"
	case ObjectType:
		return "object"
	case mapObjectType:
		return "map"
	case arrayObjectType:
		return "Array"
	case commentTokenType:
		return "comment"
	case typeTokenType:
		return "type"
	case nopStatementType:
		return "nop"
	case typeObjectInitStatementType:
		return "typeObjectInitStatementType"
	case TypeObjectType:
		return "TypeObjectType"
	case periodTokenType:
		return "."
	case getTypeObjectVarStatementType:
		return "getTypeObjectVarStatementType"
	case FuncStatementType:
		return "FuncStatementType"
	default:
		panic("unknown token type " + strconv.Itoa(int(t)))
	}
}

const EOFTokenType Type = 0
const unknownTokenType Type = 1
const commentTokenType Type = 2                  // //
const incOperatorTokenType Type = 100            // ++
const addOperatorTokenType Type = 101            // +
const subOperatorTokenType Type = 102            // -
const mulOperatorTokenType Type = 103            // *
const divOperatorTokenType Type = 104            // /
const lessTokenType Type = 105                   // <
const greaterTokenType Type = 106                // >
const lessEqualTokenType Type = 116              // <=
const greaterEqualTokenType Type = 117           // >=
const leftParenthesisTokenType Type = 120        // (
const rightParenthesisTokenType Type = 121       // )
const leftBraceTokenType Type = 122              // {
const rightBraceTokenType Type = 123             // }
const commaTokenType Type = 124                  // ,
const semicolonTokenType Type = 125              // ;
const colonTokenType Type = 126                  // :
const periodTokenType Type = 127                 // .
const ifTokenType Type = 230                     //if
const elseTokenType Type = 331                   //else
const funcTokenType Type = 332                   //func
const returnTokenType Type = 333                 //return
const breakTokenType Type = 334                  //break
const forTokenType Type = 335                    //for
const elseifTokenType Type = 336                 //else if
const varTokenType Type = 400                    // var
const assignTokenType Type = 401                 // =
const varAssignTokenType Type = 402              // var x =
const intTokenType Type = 700                    // int
const typeTokenType Type = 999                   // type
const mapObjectType Type = 1001                  // map {}
const arrayObjectType Type = 1002                // array []
const labelType Type = 5000                      // label
const statementType Type = 6001                  // statement
const statementsType Type = 6003                 // statementType
const expressionType Type = 6002                 // expressionType
const callFunctionType Type = 6004               // call function
const nopStatementType Type = 6005               // nop
const assignStatementType Type = 6006            // =
const ObjectType Type = 100000                   // object
const IntObjectType Type = 10000                 //int object
const BoolObjectType Type = 10001                // bool object
const TypeObjectType Type = 10002                // type object
const FuncStatementType Type = 10003             // function object
const typeObjectInitStatementType Type = 11003   // object init statement
const getTypeObjectVarStatementType Type = 10004 // getTypeObjectStatement statement

type Token struct {
	typ  Type
	data string
	line int
}

func (t Token) String() string {
	if t.data == "" {
		return fmt.Sprintf("line:%d type`%s`", t.line, t.typ.String())
	}
	return fmt.Sprintf("line:%d type`%s` data `%s`", t.line, t.typ.String(), t.data)
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
	semicolonToken        = Token{typ: semicolonTokenType}
	colonToken            = Token{typ: colonTokenType}
	periodToken           = Token{typ: periodTokenType}
)

var Keywords = []string{
	"if", "else", "func", "return", "break", "for", "var", "type",
}

var keywordTokenType = map[string]Type{
	"if":     ifTokenType,
	"else":   elseTokenType,
	"func":   funcTokenType,
	"return": returnTokenType,
	"break":  breakTokenType,
	"for":    forTokenType,
	"var":    varTokenType,
	"type":   typeTokenType,
}
