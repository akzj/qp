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

type FuncInstruction struct {
	toLinks []toLink
	ins     []Instruction
	label   string
}

type GenCode struct {
	symbolTable      *SymbolTable
	builtSymbolTable *SymbolTable
	ins              []Instruction
	toLinks          []toLink
	funcInstructions map[string]FuncInstruction
}

func NewGenCode() *GenCode {
	gc := &GenCode{
		symbolTable:      NewSymbolTable(),
		builtSymbolTable: NewSymbolTable(),
		ins:              []Instruction{},
		funcInstructions: map[string]FuncInstruction{},
	}
	for _, function := range BuiltInFunctions {
		gc.builtSymbolTable.addSymbol(function.Name)
	}
	return gc
}

func (genCode *GenCode) String() string {
	var buffer bytes.Buffer
	for _, it := range genCode.ins {
		if it.InstTyp != Label {
			buffer.WriteString("\t")
		}
		buffer.WriteString(it.String(genCode.symbolTable, genCode.builtSymbolTable))
		buffer.WriteString("\n")
	}
	/*for _, function := range genCode.funcInstructions {
		buffer.WriteString(function.label + ":")
		buffer.WriteString("\n")
		for _, it := range function.ins {
			buffer.WriteString("\t" + it.String(genCode.symbolTable))
			buffer.WriteString("\n")
		}
	}*/
	return buffer.String()
}

func (genCode *GenCode) pushIns(instruction Instruction) {
	genCode.ins = append(genCode.ins, instruction)
}

func (genCode *GenCode) Gen(statements []ast.Statement) *GenCode {
	for _, statement := range statements {
		genCode.genStatement(statement)
	}
	genCode.GenExit()
	linker := NewLinker(genCode.funcInstructions,
		genCode.ins,
		genCode.toLinks,
		genCode.symbolTable,
		genCode.builtSymbolTable)
	genCode.ins = linker.link()
	return genCode
}

func (genCode *GenCode) genStatement(statement runtime.Invokable) {
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
		genCode.genStatement(statement.Left)
		genCode.genStatement(statement.Right)
		genCode.genOpCode(statement.OP)
	case ast.VarAssignStatement:
		genCode.genStatement(statement.Exp)
		if statement.Exp.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{InstTyp: LoadR})
		}
		log.Println(statement.Exp.GetType())
		genCode.genStoreIns(statement.Name)
	case ast.VarStatement:
		genCode.genStatement(statement.Exp)
		genCode.genStoreIns(statement.Label)
	case ast.GetVarStatement:
		genCode.genLoadIns(statement.Label)
	case ast.IfStatement:
		genCode.genIfStatement(statement)
	case ast.Statements:
		for _, next := range statement {
			genCode.genStatement(next)
		}
	case *ast.CallStatement:
		genCode.genFuncCallStatement(statement)
	case ast.NopStatement:
	case *ast.FuncStatement:
		genCode.genFuncStatement(statement)
	case ast.ReturnStatement:
		genCode.genReturnStatement(statement)
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
	case lexer.SubType:
		genCode.pushIns(Instruction{
			InstTyp: Sub,
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
	genCode.genStatement(statement.Check)
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
	genCode.genStatement(statement.Statements)
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

func (genCode *GenCode) genFuncCallStatement(statement *ast.CallStatement) {
	//statement.ParentExp todo
	switch function := statement.Function.(type) {
	case ast.GetVarStatement:
		var II = int64(len(genCode.ins))
		index, ok := genCode.builtSymbolTable.getSymbol(function.Label)
		if ok == false { // push IP to stack for return
			genCode.pushIns(Instruction{InstTyp: Push, ValTyp: IP})
		}
		for _, argument := range statement.Arguments {
			genCode.genStatement(argument)
		}
		genCode.pushIns(Instruction{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     int64(len(statement.Arguments)),
		})
		if ok {
			genCode.pushIns(Instruction{
				InstTyp: Call,
				Val:     index,
			})
		} else {
			genCode.pushIns(Instruction{
				InstTyp: Push,
				ValTyp:  Bool,
				Val:     TRUE,
			})
			genCode.pushIns(Instruction{
				InstTyp: Jump,
				JumpTyp: DJump,
				Val:     -1, //todo link
			})
			genCode.toLinks = append(genCode.toLinks, toLink{
				label: function.Label,
				IP:    int64(len(genCode.ins)) - 1,
			})
			genCode.ins[II].Val = int64(len(genCode.ins)) - II
		}
	default:
		log.Panicf("unkown function type %s", reflect.TypeOf(function).String())
	}
}

func (genCode *GenCode) genObject(label *runtime.Object) {
	genCode.genStatement(label.Pointer)
}

func (genCode *GenCode) genFuncStatement(statement *ast.FuncStatement) {
	done := genCode.prepareGenFunction(statement.Label)
	defer done()
	genCode.pushIns(Instruction{
		InstTyp: Label,
		Val:     genCode.symbolTable.addSymbol(statement.Label),
	})
	//check argument count
	genCode.pushIns(Instruction{
		InstTyp: Push,
		ValTyp:  Int,
		Val:     int64(len(statement.Parameters)),
	})
	genCode.pushIns(Instruction{
		InstTyp: Cmp,
		CmpTyp:  Equal,
	})
	genCode.pushIns(Instruction{
		InstTyp: Jump,
		JumpTyp: RJump,
		Val:     2,
	})
	genCode.pushIns(Instruction{
		InstTyp: Call,
		Val:     genCode.symbolTable.addSymbol("panic"),
	})
	for i := len(statement.Parameters) - 1; i >= 0; i-- {
		genCode.pushIns(Instruction{
			InstTyp: Store,
			Val:     genCode.symbolTable.addSymbol(statement.Parameters[i]),
		})
	}
	for _, statement := range statement.Statements {
		genCode.genStatement(statement)
	}
	genCode.pushIns(Instruction{InstTyp: Ret})
}

func (genCode *GenCode) prepareGenFunction(label string) func() {
	ins := genCode.ins
	toLink := genCode.toLinks
	genCode.ins = nil
	genCode.toLinks = nil
	return func() {
		genCode.funcInstructions[label] = FuncInstruction{
			toLinks: genCode.toLinks,
			ins:     genCode.ins,
			label:   label,
		}
		genCode.ins = ins
		genCode.toLinks = toLink
	}
}

func (genCode *GenCode) genReturnStatement(statement ast.ReturnStatement) {
	genCode.genStatement(statement.Exp)
	genCode.pushIns(Instruction{
		InstTyp: StoreR,
	})
	genCode.pushIns(Instruction{InstTyp: Ret})
}

func (genCode *GenCode) GenExit() {
	genCode.pushIns(Instruction{
		InstTyp: Exit,
	})
}
