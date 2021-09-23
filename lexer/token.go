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
	case CallType:
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
	case VarInitType:
		return ":="
	default:
		panic("unknown token type " + strconv.Itoa(int(t)))
	}
}

const (
	ErrorTokenType               Type = iota
	EOFType                           // EOF
	CommentType                       // //
	StringType                        // string "" ''
	NilType                           // null
	TrueType                          // true
	FalseType                         // false
	NewLineType                       // \n
	NoType                            // !
	OrType                            // ||
	AndType                           // &&
	IncType                           // ++
	AddType                           // +
	SubType                           // -
	MulOpType                         // *
	DivOpType                         // /
	ModOpType                         // %
	LessType                          // <
	GreaterType                       // >
	LessEqualType                     // <=
	GreaterEqualType                  // >=
	EqualType                         // ==
	NoEqualType                       // !=
	LeftParenthesisType               // (
	RightParenthesisType              // )
	LeftBraceType                     // {
	RightBraceType                    // }
	CommaType                         // ,
	SemicolonType                     // ;
	ColonType                         // :
	PeriodType                        // .
	LeftBracketType                   // [
	RightBracketType                  // ]
	IfType                            // if
	ElseType                          // else
	FuncType                          // func
	ReturnType                        // return
	BreakType                         // break
	ForType                           // for
	ElseifType                        // else if
	VarType                           // var
	AssignType                        // =
	VarAssignType                     // var x =
	VarInitType                       // :=
	IntType                           // int
	TypeType                          // type
	MapObjectType                     // map {}
	ArrayObjectType                   // array []
	IDType                            // name
	StatementType                     // statement
	StatementsType                    // StatementType
	ExpressionType                    // ExpressionType
	CallType                          // call function
	NopStatementType                  // nop
	AssignStatementType               // =
	ObjectType                        // objects
	BoolObjectType                    // bool
	TypeObjectType                    // type objects
	FuncStatementType                 // function objects
	TypeObjectInitStatementType       // objects init statement
	PropObjectStatementType           // getTypeObjectStatement statement
	GetObjectObjectStatementType      // getObjectObjectStatement
	FuncCallQueueStatementType        // FuncCallQueueStatement
	DurationObjectType                // DurationObjectType
	TimeObjectType                    // DurationObjectType
	BuiltInFunctionType               // built in function
	CreateObjectStatementType         // createObjectStatement
)

type Token struct {
	Typ  Type
	Val  string
	Line int
}

func (t Token) String() string {
	if t.Val == "" {
		return fmt.Sprintf("line:%d type`%s`", t.Line, t.Typ.String())
	}
	return fmt.Sprintf("line:%d type`%s` var `%s`", t.Line, t.Typ.String(), t.Val)
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
	VarInitToken          = Token{Typ: VarInitType}
	ModToken              = Token{Typ: ModOpType}
	DivToken              = Token{Typ: DivOpType}
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
