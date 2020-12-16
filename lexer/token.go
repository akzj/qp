package lexer

import (
	"fmt"
	"strconv"
)

type Type int

func (t Type) String() string {
	switch t {
	case AddType:
		return "+"
	case DivOpType:
		return "/"
	case SubType:
		return "-"
	case MulOpType:
		return "*"
	case IntType:
		return "int"
	case LeftParenthesisType:
		return "("
	case RightParenthesisType:
		return ")"
	case IfType:
		return "if"
	case ElseifType:
		return "else if"
	case ElseType:
		return "else"
	case ForType:
		return "for"
	case BreakType:
		return "break"
	case ReturnType:
		return "return"
	case LeftBraceType:
		return "{"
	case RightBraceType:
		return "}"
	case LessType:
		return "<"
	case LessEqualType:
		return "<="
	case GreaterType:
		return ">"
	case GreaterEqualType:
		return ">="
	case EOFType:
		return "EOF"
	case IDType:
		return "ID"
	case StatementType:
		return "statement"
	case StatementsType:
		return "statements"
	case ExpressionType:
		return "exp"
	case VarType:
		return "var"
	case AssignType:
		return "="
	case VarAssignType:
		return "var ="
	case CommaType:
		return ","
	case IncType:
		return "++"
	case CallFunctionType:
		return "call"
	case SemicolonType:
		return ";"
	case AssignStatementType:
		return "assignStatement"
	case FuncType:
		return "func"
	case ObjectType:
		return "objects"
	case MapObjectType:
		return "map"
	case ArrayObjectType:
		return "Array"
	case CommentType:
		return "comment"
	case TypeType:
		return "type"
	case NopStatementType:
		return "nop"
	case TypeObjectInitStatementType:
		return "TypeObjectInitStatementType"
	case TypeObjectType:
		return "TypeObjectType"
	case PeriodType:
		return "."
	case PropObjectStatementType:
		return "PropObjectStatementType"
	case FuncStatementType:
		return "FuncStatementType"
	case ErrorTokenType:
		return "ErrorTokenType"
	case StringType:
		return "string"
	case NilType:
		return "nil"
	case EqualType:
		return "=="
	case GetObjectObjectStatementType:
		return "GetObjectObjectStatementType"
	case LeftBracketType:
		return "["
	case RightBracketType:
		return "]"
	case FuncCallQueueStatementType:
		return "FuncCallQueueStatementType"
	case FalseType:
		return "false"
	case TrueType:
		return "true"
	case NoEqualType:
		return "!="
	case DurationObjectType:
		return "DurationObjectType"
	case TimeObjectType:
		return "TimeObjectType"
	case BuiltInFunctionType:
		return "builtInFunction"
	case NewLineType:
		return `\n`
	case NoType:
		return "!"
	case OrType:
		return "||"
	case AndType:
		return "&&"
	default:
		panic("unknown token type " + strconv.Itoa(int(t)))
	}
}

const ErrorTokenType Type = -1
const EOFType Type = 0
const CommentType Type = 2                     // //
const StringType Type = 3                      // string "" ''
const NilType Type = 4                         // null
const TrueType Type = 5                        // true
const FalseType Type = 6                       // false
const NewLineType Type = 7                     // \n
const NoType Type = 8                          // !
const OrType Type = 9                          // ||
const AndType Type = 10                        // &&
const IncType Type = 100                       // ++
const AddType Type = 101                       // +
const SubType Type = 102                       // -
const MulOpType Type = 103                     // *
const DivOpType Type = 104                     // /
const LessType Type = 105                      // <
const GreaterType Type = 106                   // >
const LessEqualType Type = 116                 // <=
const GreaterEqualType Type = 117              // >=
const EqualType Type = 118                     // ==
const NoEqualType Type = 119                   // !=
const LeftParenthesisType Type = 120           // (
const RightParenthesisType Type = 121          // )
const LeftBraceType Type = 122                 // {
const RightBraceType Type = 123                // }
const CommaType Type = 124                     // ,
const SemicolonType Type = 125                 // ;
const ColonType Type = 126                     // :
const PeriodType Type = 127                    // .
const LeftBracketType Type = 128               // [
const RightBracketType Type = 129              // ]
const IfType Type = 230                        // if
const ElseType Type = 331                      // else
const FuncType Type = 332                      // func
const ReturnType Type = 333                    // return
const BreakType Type = 334                     // break
const ForType Type = 335                       // for
const ElseifType Type = 336                    // else if
const VarType Type = 400                       // var
const AssignType Type = 401                    // =
const VarAssignType Type = 402                 // var x =
const IntType Type = 700                       // int
const TypeType Type = 999                      // type
const MapObjectType Type = 1001                // map {}
const ArrayObjectType Type = 1002              // array []
const IDType Type = 5000                       // name
const StatementType Type = 6001                // statement
const StatementsType Type = 6003               // StatementType
const ExpressionType Type = 6002               // ExpressionType
const CallFunctionType Type = 6004             // call function
const NopStatementType Type = 6005             // nop
const AssignStatementType Type = 6006          // =
const ObjectType Type = 100000                 // objects
const BoolObjectType Type = 10001              // bool
const TypeObjectType Type = 10002              // type objects
const FuncStatementType Type = 10003           // function objects
const TypeObjectInitStatementType Type = 11003 // objects init statement
const PropObjectStatementType Type = 10004     // getTypeObjectStatement statement
const GetObjectObjectStatementType Type = 1005 // getObjectObjectStatement
const FuncCallQueueStatementType Type = 10006  // FuncCallQueueStatement
const DurationObjectType Type = 10007          // DurationObjectType
const TimeObjectType Type = 10008              // DurationObjectType
const BuiltInFunctionType = 10009              // built in function

type Token struct {
	Typ  Type
	Val  string
	Line int
}

func (t Token) String() string {
	if t.Val == "" {
		return fmt.Sprintf("Line:%d type`%s`", t.Line, t.Typ.String())
	}
	return fmt.Sprintf("Line:%d type`%s` Val `%s`", t.Line, t.Typ.String(), t.Val)
}

var (
	EmptyToken            = Token{Typ: EOFType}
	AddOperatorToken      = Token{Typ: AddType}
	MulOperatorToken      = Token{Typ: MulOpType}
	LeftParenthesisToken  = Token{Typ: LeftParenthesisType}
	RightParenthesisToken = Token{Typ: RightParenthesisType}
	LeftBraceToken        = Token{Typ: LeftBraceType}
	RightBraceToken       = Token{Typ: RightBraceType}
	LessToken             = Token{Typ: LessType}
	LessEqualToken        = Token{Typ: LessEqualType}
	GreaterToken          = Token{Typ: GreaterType}
	GreaterEqualToken     = Token{Typ: GreaterEqualType}
	AssignToken           = Token{Typ: AssignType}
	CommaToken            = Token{Typ: CommaType}
	IncOperatorToken      = Token{Typ: IncType}
	SemicolonToken        = Token{Typ: SemicolonType}
	ColonToken            = Token{Typ: ColonType}
	PeriodToken           = Token{Typ: PeriodType}
	EqualToken            = Token{Typ: EqualType}
	LeftBracketToken      = Token{Typ: LeftBracketType}
	RightBracketToken     = Token{Typ: RightBracketType}
	NoEqualToken          = Token{Typ: NoEqualType}
	SubOperatorToken      = Token{Typ: SubType}
	OrToken               = Token{Typ: OrType}
	AndToken              = Token{Typ: AndType}
)

var Keywords = []string{
	"if", "else", "func", "return", "break", "for", "var", "type", "nil", "true", "false",
}

var KeywordType = map[string]Type{
	"if":     IfType,
	"else":   ElseType,
	"func":   FuncType,
	"return": ReturnType,
	"break":  BreakType,
	"for":    ForType,
	"var":    VarType,
	"type":   TypeType,
	"nil":    NilType,
	"true":   TrueType,
	"false":  FalseType,
}
