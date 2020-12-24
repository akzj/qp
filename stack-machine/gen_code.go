package stackmachine

import (
	"bytes"
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
	"log"
	"reflect"
)

type toLink struct {
	label string
	IP    int64
}
type GenCode struct {
	symbolTable *SymbolTable
	ins         []Instruction
	statements  []ast.Statement
	toLinks     []toLink
}

func NewGenCode(statements []ast.Statement) *GenCode {
	return &GenCode{
		ins:         []Instruction{},
		symbolTable: NewSymbolTable(),
		statements:  statements,
	}
}

func (genCode *GenCode) pushIns(instruction Instruction) {
	genCode.ins = append(genCode.ins, instruction)
}

func (genCode *GenCode) Gen() *GenCode {
	for _, statement := range genCode.statements {
		genCode.genCodeStatement(statement)
	}
	return genCode
}

func (genCode *GenCode) genCodeStatement(statement runtime.Invokable) {
	switch statement := statement.(type) {
	case ast.Int:
		genCode.pushIns(Instruction{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     int64(statement),
		})
	case ast.Bool:
		var b = FALSE
		if statement {
			b = TRUE
		}
		genCode.pushIns(Instruction{
			InstTyp: Push,
			ValTyp:  Bool,
			Val:     b,
		})
	case *runtime.Object:
		genCode.genObject(statement)
	case ast.BinaryOpExpression:
		genCode.genCodeStatement(statement.Left)
		genCode.genCodeStatement(statement.Right)
		genCode.genOpCode(statement.OP)
	case ast.VarAssignStatement:
		genCode.genCodeStatement(statement.Exp)
		genCode.genStoreIns(statement.Name)
	case ast.VarStatement:
		genCode.genCodeStatement(statement.Exp)
		genCode.genStoreIns(statement.Label)
	case ast.GetVarStatement:
		genCode.genLoadIns(statement.Label)
	case ast.IfStatement:
		genCode.genIfStatement(statement)
	case ast.Statements:
		for _, next := range statement {
			genCode.genCodeStatement(next)
		}
	case *ast.FuncCallStatement:
		genCode.genFuncCallStatement(statement)
	case ast.NopStatement:
	case *ast.FuncStatement:
		genCode.genFuncStatement(statement)
	default:
		log.Panicf("unknown statement %s", reflect.TypeOf(statement).String())
	}
}

func (genCode *GenCode) genOpCode(op lexer.Type) {
	switch op {
	case lexer.AddType:
		genCode.pushIns(Instruction{
			InstTyp: Add,
		})
	case lexer.LessType:
		genCode.pushIns(Instruction{
			InstTyp: Cmp,
			CmpTyp:  Less,
		})
	case lexer.GreaterType:
		genCode.pushIns(Instruction{
			InstTyp: Cmp,
			CmpTyp:  Greater,
		})
	default:
		log.Panicf("unknown instruction %s", op.String())
	}
}

func (genCode *GenCode) genStoreIns(label string) {
	genCode.pushIns(Instruction{
		InstTyp: Store,
		Val:     genCode.symbolTable.addSymbol(label),
	})
}

func (genCode *GenCode) genIfStatement(statement ast.IfStatement) {
	genCode.genCodeStatement(statement.Check)
	genCode.pushIns(Instruction{
		InstTyp: Jump,
		JumpTyp: RJump,
		Val:     3,
	})
	genCode.pushIns(Instruction{
		InstTyp: Push,
		ValTyp:  Bool,
		Val:     TRUE,
	})
	genCode.pushIns(Instruction{
		InstTyp: Jump,
		JumpTyp: RJump,
	})
	index := len(genCode.ins)
	genCode.genCodeStatement(statement.Statements)
	jumpTo := len(genCode.ins) - index + 1
	genCode.ins[index-1].Val = int64(jumpTo)
}

func (genCode *GenCode) genLoadIns(label string) {
	index := genCode.symbolTable.addSymbol(label)
	genCode.pushIns(Instruction{
		InstTyp: Load,
		Val:     index,
	})
}

func (genCode *GenCode) String() string {
	var buffer bytes.Buffer
	for _, it := range genCode.ins {
		buffer.WriteString(it.String(genCode.symbolTable))
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (genCode *GenCode) genFuncCallStatement(statement *ast.FuncCallStatement) {
	//statement.ParentExp todo
	for _, argument := range statement.Arguments {
		genCode.genCodeStatement(argument)
	}
	genCode.pushIns(Instruction{
		InstTyp: Push,
		ValTyp:  Int,
		Val:     int64(len(statement.Arguments)),
	})
	switch function := statement.Function.(type) {
	case ast.GetVarStatement:
		if index, ok := genCode.symbolTable.getSymbol(function.Label); ok {
			genCode.pushIns(Instruction{
				InstTyp: Call,
				Val:     index,
			})
		} else {
			genCode.pushIns(Instruction{
				InstTyp: Jump,
				JumpTyp: DJump,
				Val:     0,
			})
			genCode.toLinks = append(genCode.toLinks, toLink{
				label: function.Label,
				IP:    int64(len(genCode.ins)),
			})
		}
	default:
		log.Panicf("unkown function type %s", reflect.TypeOf(function).String())
	}
}

func (genCode *GenCode) genObject(label *runtime.Object) {
	genCode.genCodeStatement(label.Pointer)
}

func (genCode *GenCode) genFuncStatement(statement *ast.FuncStatement) {

}
