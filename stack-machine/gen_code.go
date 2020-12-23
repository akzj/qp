package stackmachine

import (
	"bytes"
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
	"log"
	"reflect"
)

type GenCode struct {
	symbolTable *SymbolTable
	ins         []Instruction
	statements  []ast.Statement
}

func NewGenCode(statements []ast.Statement) *GenCode {
	return &GenCode{
		symbolTable: &SymbolTable{
			symbols:   []string{},
			symbolMap: map[string]int64{},
		},
		ins:        nil,
		statements: statements,
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
		genCode.genLoadIns(statement.Label)
	case ast.BinaryOpExpression:
		genCode.genCodeStatement(statement.Right)
		genCode.genCodeStatement(statement.Left)
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
	default:
		log.Panicf("unknown statement %s",reflect.TypeOf(statement).String())
	}
}

func (genCode *GenCode) genOpCode(op lexer.Type) {
	switch op {
	case lexer.AddType:
		genCode.pushIns(Instruction{
			InstTyp: Add,
		})
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
		ValTyp:  0,
		CmpTyp:  0,
		Val:     0,
	})
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
