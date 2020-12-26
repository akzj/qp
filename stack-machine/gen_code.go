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

type StackManager struct {
	IP         int
	stack      []string
	stackFrame []struct {
		stack []string
	}
}

func NewStackManager() *StackManager {
	return &StackManager{
	}
}

func (s *StackManager) Store(symbol string) int64 {
	if len(symbol) != 0 {
		for _, label := range s.stack {
			if label == symbol {
				log.Panicln("redefine symbol", symbol, s.stackFrame, s.stack)
			}
		}
	} else {
		s.IP++
	}
	s.stack = append(s.stack, symbol)
	return int64(len(s.stack)) - 1
}

func (s *StackManager) load(label string) (int64, bool) {
	for i := len(s.stack) - 1; i >= 0; i-- {
		if s.stack[i] == label {
			return int64(i), true
		}
	}
	return -1, false
}

func (s *StackManager) pushStackFrame() {
	s.stackFrame = append(s.stackFrame, struct{ stack []string }{stack: s.stack})
	s.stack = nil
}
func (s *StackManager) popStackFrame() {
	s.stack = s.stackFrame[len(s.stackFrame)-1].stack
	s.stackFrame = s.stackFrame[:len(s.stackFrame)-1]
}

func (s *StackManager) SP() int64 {
	return int64(len(s.stack))
}

type GenCode struct {
	symbolTable      *SymbolTable
	builtSymbolTable *SymbolTable
	ins              []Instruction
	toLinks          []toLink
	sm               *StackManager
	funcInstructions map[string]FuncInstruction
}

func NewGenCode() *GenCode {
	gc := &GenCode{
		symbolTable:      NewSymbolTable(),
		builtSymbolTable: NewSymbolTable(),
		ins:              []Instruction{},
		sm:               NewStackManager(),
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
		if statement.Left.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{InstTyp: LoadR, Val: 1})
		}
		genCode.genStatement(statement.Right)
		if statement.Right.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{InstTyp: LoadR, Val: 1})
		}
		genCode.genOpCode(statement.OP)
	case ast.VarAssignStatement:
		genCode.genStatement(statement.Exp)
		if statement.Exp.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{InstTyp: LoadR, Val: 1})
		}
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
		genCode.genCallStatement(statement)
	case ast.NopStatement:
	case *ast.FuncStatement:
		genCode.genFuncStatement(statement)
	case ast.ReturnStatement:
		genCode.genReturnStatement(statement)
	case ast.AssignStatement:
		genCode.genAssignStatement(statement)
	case ast.String:
		genCode.pushIns(Instruction{
			InstTyp: Push,
			ValTyp:  String,
			Str:     string(statement),
		})
	case ast.PeriodStatement:
		genCode.pushIns(Instruction{
			InstTyp: LoadO,
			Str:     statement.Val,
		})
	case ast.ForStatement:
		genCode.genForStatement(statement)
	case ast.IncFieldStatement:
		genCode.genIncFieldStatement(statement)
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
	genCode.sm.Store(label)
	genCode.symbolTable.addSymbol(label)
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

	//if statement
	genCode.genStatement(statement.Statements)

	//fix jump val
	jumpTo := len(genCode.ins) - index + 1
	genCode.ins[index-1].Val = int64(jumpTo)
}

func (genCode *GenCode) genForStatement(statement ast.ForStatement) {

	genCode.genStatement(statement.Pre)

	begin := len(genCode.ins)
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
	jumpStatement := len(genCode.ins)

	//if statement
	genCode.genStatement(statement.Statements)

	genCode.genStatement(statement.Post)

	genCode.pushIns(Instruction{
		InstTyp: Push,
		ValTyp:  Bool,
		Val:     TRUE,
	})
	genCode.pushIns(Instruction{
		InstTyp: Jump,
		JumpTyp: RJump,
		Val:     int64(begin  - len(genCode.ins)),
	})

	//fix jump val
	genCode.ins[jumpStatement-1].Val = int64(len(genCode.ins) - jumpStatement + 1)
}

func (genCode *GenCode) genLoadIns(label string) {
	index, ok := genCode.sm.load(label)
	if ok == false {
		log.Panicln("no find label`" + label + "`")
	}
	genCode.pushIns(Instruction{
		InstTyp: Load,
		Val:     index,
		symbol:  genCode.symbolTable.addSymbol(label),
	})
}

func (genCode *GenCode) genCallStatement(statement *ast.CallStatement) {
	//statement.ParentExp todo
	switch function := statement.Function.(type) {
	case ast.GetVarStatement:
		var II = int64(len(genCode.ins))
		index, ok := genCode.builtSymbolTable.getSymbol(function.Label)
		if ok == false { // push IP to stack for return
			//genCode.sm.Store("")
			genCode.pushIns(Instruction{InstTyp: Push, ValTyp: IP})
		}
		var R int64
		genCode.pushIns(Instruction{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     int64(len(statement.Arguments)),
		})
		genCode.pushIns(Instruction{
			InstTyp: StoreR,
			Val:     int64(R),
		})
		for _, argument := range statement.Arguments {
			R++
			genCode.genStatement(argument)
			genCode.pushIns(Instruction{
				InstTyp: StoreR,
				Val:     int64(R),
			})
		}
		if ok == false {
			genCode.pushIns(Instruction{
				InstTyp: PushS,
			})
		}
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
	case ast.PeriodStatement:
		//log.Println(reflect.TypeOf(statement.ParentExp).String())
		genCode.genStatement(statement.ParentExp)
		genCode.genStatement(statement.Function)
		genCode.pushIns(Instruction{
			InstTyp: CallO,
		})
	default:
		log.Panicf("unkown function type %s", reflect.TypeOf(function).String())
	}
}

func (genCode *GenCode) genObject(label *runtime.Object) {
	genCode.genStatement(label.Pointer)
}

func (genCode *GenCode) genFuncStatement(statement *ast.FuncStatement) {

	genCode.sm.pushStackFrame()
	defer genCode.sm.popStackFrame()

	done := genCode.prepareGenFunction(statement.Label)
	defer done()

	genCode.pushIns(Instruction{
		InstTyp: Label,
		symbol:  genCode.symbolTable.addSymbol(statement.Label),
	})
	//check argument count

	genCode.pushIns(Instruction{
		InstTyp: LoadR,
		Val:     int64(0),
	})

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
	for i := 0; i < len(statement.Parameters); i++ {
		genCode.symbolTable.addSymbol(statement.Parameters[i])
		index := genCode.sm.Store(statement.Parameters[i])
		genCode.pushIns(Instruction{
			InstTyp: LoadR,
			Val:     int64(i + 1),
		})
		log.Println(statement.Parameters[i], index)
	}

	var last runtime.Invokable
	for _, statement := range statement.Statements {
		genCode.genStatement(statement)
		last = statement
	}
	switch last.(type) {
	case ast.ReturnStatement:
	default:
		genCode.pushIns(Instruction{InstTyp: PopS})
		genCode.pushIns(Instruction{InstTyp: Ret})
	}
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
	if statement.Exp.GetType() == lexer.CallType {
		genCode.pushIns(Instruction{InstTyp: LoadR, Val: 1})
	}
	genCode.pushIns(Instruction{
		InstTyp: StoreR,
		Val:     1,
	})
	genCode.pushIns(Instruction{InstTyp: PopS})
	genCode.pushIns(Instruction{InstTyp: Ret})
}

func (genCode *GenCode) GenExit() {
	genCode.pushIns(Instruction{
		InstTyp: Exit,
	})
}

func (genCode *GenCode) genAssignStatement(statement ast.AssignStatement) {
	genCode.genStatement(statement.Exp)
	switch obj := statement.Left.(type) {
	case ast.GetVarStatement:
		if index, ok := genCode.sm.load(obj.Label); ok {
			genCode.pushIns(Instruction{
				InstTyp: Store,
				Val:     index,
			})
		}
	default:
		log.Panicln(reflect.TypeOf(obj).String())
	}
}

func (genCode *GenCode) genIncFieldStatement(statement ast.IncFieldStatement) {
	object := statement.Exp.(ast.GetVarStatement)
	index, ok := genCode.sm.load(object.Label)
	if ok == false {
		log.Panicln("no find label`" + object.Label + "`")
	}
	genCode.pushIns(Instruction{
		InstTyp: Load,
		Val:     index,
	})
	genCode.pushIns(Instruction{
		InstTyp: Push,
		ValTyp:  Int,
		Val:     1,
	})
	genCode.pushIns(Instruction{
		InstTyp: Add,
	})
	genCode.pushIns(Instruction{
		InstTyp: Store,
		Val:     index,
	})

}
