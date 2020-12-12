package qp

import (
	"fmt"
	"strconv"
)

type Type int

func (t Type) String() string {
	switch t {
	case addType:
		return "+"
	case divOpType:
		return "/"
	case subType:
		return "-"
	case mulOpType:
		return "*"
	case intType:
		return "int"
	case leftParenthesisType:
		return "("
	case rightParenthesisType:
		return ")"
	case ifType:
		return "if"
	case elseifType:
		return "else if"
	case elseType:
		return "else"
	case forType:
		return "for"
	case breakType:
		return "break"
	case returnType:
		return "return"
	case leftBraceType:
		return "{"
	case rightBraceType:
		return "}"
	case lessTokenType:
		return "<"
	case lessEqualType:
		return "<="
	case greaterType:
		return ">"
	case greaterEqualType:
		return ">="
	case EOFType:
		return "EOF"
	case IDType:
		return "ID"
	case statementType:
		return "statement"
	case statementsType:
		return "statements"
	case expressionType:
		return "expression"
	case varType:
		return "var"
	case assignType:
		return "="
	case varAssignTokenType:
		return "var ="
	case IntType:
		return "Int"
	case commaType:
		return ","
	case incOperatorTokenType:
		return "++"
	case callFunctionType:
		return "call"
	case semicolonType:
		return ";"
	case assignStatementType:
		return "assignStatement"
	case funcType:
		return "func"
	case ObjectType:
		return "objects"
	case mapObjectType:
		return "map"
	case arrayObjectType:
		return "Array"
	case commentTokenType:
		return "comment"
	case typeType:
		return "type"
	case nopStatementType:
		return "nop"
	case typeObjectInitStatementType:
		return "typeObjectInitStatementType"
	case TypeObjectType:
		return "TypeObjectType"
	case periodType:
		return "."
	case propObjectStatementType:
		return "propObjectStatementType"
	case FuncStatementType:
		return "FuncStatementType"
	case ErrorTokenType:
		return "ErrorTokenType"
	case stringType:
		return "string"
	case nilType:
		return "nil"
	case EqualType:
		return "=="
	case getObjectObjectStatementType:
		return "getObjectObjectStatementType"
	case leftBracketTokenType:
		return "["
	case rightBracketTokenType:
		return "]"
	case funcCallQueueStatementType:
		return "funcCallQueueStatementType"
	case falseType:
		return "false"
	case TrueType:
		return "true"
	case NoEqualTokenType:
		return "!="
	case DurationObjectType:
		return "DurationObjectType"
	case timeObjectType:
		return "timeObjectType"
	case builtInFunctionType:
		return "builtInFunction"
	case NewLineType:
		return `\n`
	case NoType:
		return "!"
	case orType:
		return "||"
	case AndType:
		return "&&"
	default:
		panic("unknown token type " + strconv.Itoa(int(t)))
	}
}

const ErrorTokenType Type = -1
const EOFType Type = 0
const commentTokenType Type = 2                // //
const stringType Type = 3                      // string "" ''
const nilType Type = 4                         // null
const TrueType Type = 5                        // true
const falseType Type = 6                       // false
const NewLineType Type = 7                     // \n
const NoType Type = 8                          // !
const orType Type = 9                          // ||
const AndType Type = 10                        // &&
const incOperatorTokenType Type = 100          // ++
const addType Type = 101                       // +
const subType Type = 102                       // -
const mulOpType Type = 103                     // *
const divOpType Type = 104                     // /
const lessTokenType Type = 105                 // <
const greaterType Type = 106                   // >
const lessEqualType Type = 116                 // <=
const greaterEqualType Type = 117              // >=
const EqualType Type = 118                     // ==
const NoEqualTokenType Type = 119              // !=
const leftParenthesisType Type = 120           // (
const rightParenthesisType Type = 121          // )
const leftBraceType Type = 122                 // {
const rightBraceType Type = 123                // }
const commaType Type = 124                     // ,
const semicolonType Type = 125                 // ;
const colonTokenType Type = 126                // :
const periodType Type = 127                    // .
const leftBracketTokenType Type = 128          // [
const rightBracketTokenType Type = 129         // ]
const ifType Type = 230                        // if
const elseType Type = 331                      // else
const funcType Type = 332                      // func
const returnType Type = 333                    // return
const breakType Type = 334                     // break
const forType Type = 335                       // for
const elseifType Type = 336                    // else if
const varType Type = 400                       // var
const assignType Type = 401                    // =
const varAssignTokenType Type = 402            // var x =
const intType Type = 700                       // int
const typeType Type = 999                      // type
const mapObjectType Type = 1001                // map {}
const arrayObjectType Type = 1002              // array []
const IDType Type = 5000                       // label
const statementType Type = 6001                // statement
const statementsType Type = 6003               // statementType
const expressionType Type = 6002               // expressionType
const callFunctionType Type = 6004             // call function
const nopStatementType Type = 6005             // nop
const assignStatementType Type = 6006          // =
const ObjectType Type = 100000                 // objects
const IntType Type = 10000                     // int
const BoolObjectType Type = 10001              // bool
const TypeObjectType Type = 10002              // type objects
const FuncStatementType Type = 10003           // function objects
const typeObjectInitStatementType Type = 11003 // objects init statement
const propObjectStatementType Type = 10004     // getTypeObjectStatement statement
const getObjectObjectStatementType Type = 1005 // getObjectObjectStatement
const funcCallQueueStatementType Type = 10006  // FuncCallQueueStatement
const DurationObjectType Type = 10007          //DurationObjectType
const timeObjectType Type = 10008              //DurationObjectType
const builtInFunctionType = 10009              // built in function

type Token struct {
	typ  Type
	val  string
	line int
}

func (t Token) String() string {
	if t.val == "" {
		return fmt.Sprintf("line:%d type`%s`", t.line, t.typ.String())
	}
	return fmt.Sprintf("line:%d type`%s` val `%s`", t.line, t.typ.String(), t.val)
}

var (
	emptyToken            = Token{typ: EOFType}
	addOperatorToken      = Token{typ: addType}
	mulOperatorToken      = Token{typ: mulOpType}
	leftParenthesisToken  = Token{typ: leftParenthesisType}
	rightParenthesisToken = Token{typ: rightParenthesisType}
	leftBraceToken        = Token{typ: leftBraceType}
	rightBraceToken       = Token{typ: rightBraceType}
	lessToken             = Token{typ: lessTokenType}
	lessEqualToken        = Token{typ: lessEqualType}
	greaterToken          = Token{typ: greaterType}
	greaterEqualToken     = Token{typ: greaterEqualType}
	assignToken           = Token{typ: assignType}
	commaToken            = Token{typ: commaType}
	incOperatorToken      = Token{typ: incOperatorTokenType}
	semicolonToken        = Token{typ: semicolonType}
	colonToken            = Token{typ: colonTokenType}
	periodToken           = Token{typ: periodType}
	equalToken            = Token{typ: EqualType}
	leftBracketToken      = Token{typ: leftBracketTokenType}
	rightBracketToken     = Token{typ: rightBracketTokenType}
	NoEqualToken          = Token{typ: NoEqualTokenType}
	subOperatorToken      = Token{typ: subType}
	orToken               = Token{typ: orType}
	andToken              = Token{typ: AndType}
)

var Keywords = []string{
	"if", "else", "func", "return", "break", "for", "var", "type", "nil", "true", "false",
}

var keywordTokenType = map[string]Type{
	"if":     ifType,
	"else":   elseType,
	"func":   funcType,
	"return": returnType,
	"break":  breakType,
	"for":    forType,
	"var":    varType,
	"type":   typeType,
	"nil":    nilType,
	"true":   TrueType,
	"false":  falseType,
}
