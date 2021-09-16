package stackmachine

import (
	"bytes"
	"fmt"
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
	"hash/crc32"
	"log"
	"reflect"
	"strconv"
	"strings"
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
	stack      []string
	stackFrame []struct {
		stack []string
	}
}

func NewStackManager() *StackManager {
	return &StackManager{}
}

func (s *StackManager) Store(symbol string) int64 {
	if len(symbol) != 0 {
		for _, label := range s.stack {
			if label == symbol {
				log.Panicln("redefine symbol", symbol, s.stackFrame, s.stack)
			}
		}
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
	for index, it := range genCode.ins {
		buffer.WriteString(strconv.Itoa(index))
		buffer.WriteString("\t")
		if it.Type != Label {
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

func (genCode *GenCode) Gen(statements []ast.Expression) *GenCode {
	for _, statement := range statements {
		log.Println(statement.String())
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
			Type:   Push,
			ValTyp: Int,
			Val:    int64(statement),
		})
	case ast.Bool:
		var b = FALSE
		if statement {
			b = TRUE
		}
		genCode.pushIns(Instruction{
			Type:   Push,
			ValTyp: Bool,
			Val:    b,
		})
	case *runtime.Object:
		genCode.genObject(statement)
	case ast.BinaryOpExpression:
		genCode.genStatement(statement.Left)
		if statement.Left.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{Type: LoadR, Val: 1})
		}
		genCode.genStatement(statement.Right)
		if statement.Right.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{Type: LoadR, Val: 1})
		}
		genCode.genOpCode(statement.OP)
	case ast.VarAssignStatement:
		genCode.genStatement(statement.Exp)
		if statement.Exp.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{Type: LoadR, Val: 1})
		}
		genCode.genStoreIns(statement.Name)
	case ast.VarStatement:
		genCode.genStatement(statement.Exp)
		if statement.Exp.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{Type: LoadR, Val: 1})
		}
		genCode.genStoreIns(statement.Label)
	case ast.GetVarStatement:
		genCode.genLoadIns(statement.Label)
	case ast.IfExpression:
		genCode.genIfStatement(statement)
	case ast.Expressions:
		for _, next := range statement {
			genCode.genStatement(next)
		}
	case *ast.CallStatement:
		genCode.genCallStatement(statement)
	case ast.NopStatement:
	case *ast.FuncExpression:
		genCode.genFuncStatement(statement)
	case ast.ReturnStatement:
		genCode.genReturnStatement(statement)
	case ast.AssignStatement:
		genCode.genAssignStatement(statement)
	case ast.String:
		genCode.pushIns(Instruction{
			Type:   Push,
			ValTyp: String,
			Str:    string(statement),
		})
	case ast.PeriodStatement:
		genCode.genStatement(statement.Exp)
		genCode.pushIns(Instruction{
			Type: LoadO,
			Str:  statement.Val,
		})
	case ast.ForExpression:
		genCode.genForStatement(statement)
	case ast.IncFieldStatement:
		genCode.genIncFieldStatement(statement)
	case ast.ObjectInitStatement:
		genCode.genObjectInitStatement(statement)
	case *ast.TypeObject:
		genCode.genTypeObject(statement)
	case objectInitStatement:
		genCode.genInitStatement(statement)
	case createObjectStatement:
		genCode.genCreateObjectStatement(statement)
	case ast.NilObject:
		genCode.genNilObject(statement)
	default:
		log.Panicf("unknown statement %s", reflect.TypeOf(statement).String())
	}
}

func (genCode *GenCode) genOpCode(op lexer.Type) {
	switch op {
	case lexer.AddType:
		genCode.pushIns(Instruction{
			Type: Add,
		})
	case lexer.LessType:
		genCode.pushIns(Instruction{
			Type:   Cmp,
			CmpTyp: Less,
		})
	case lexer.GreaterType:
		genCode.pushIns(Instruction{
			Type:   Cmp,
			CmpTyp: Greater,
		})
	case lexer.SubType:
		genCode.pushIns(Instruction{
			Type: Sub,
		})
	case lexer.NoEqualType:
		genCode.pushIns(Instruction{
			Type:   Cmp,
			CmpTyp: NoEqual,
		})
	case lexer.AndType:
		genCode.pushIns(Instruction{
			Type: And,
		})
	case lexer.EqualType:
		genCode.pushIns(Instruction{Type: Cmp,
			CmpTyp: Equal})
	default:
		log.Panicf("unknown instruction %s", op.String())
	}
}

func (genCode *GenCode) genStoreIns(label string) {
	genCode.sm.Store(label)
	genCode.symbolTable.addSymbol(label)
}

func (genCode *GenCode) genIfStatement(statement ast.IfExpression) {
	genCode.genStatement(statement.Check)
	genCode.pushIns(Instruction{
		Type:    Jump,
		JumpTyp: RJump,
		Val:     3,
	})
	genCode.pushIns(Instruction{
		Type:   Push,
		ValTyp: Bool,
		Val:    TRUE,
	})
	genCode.pushIns(Instruction{
		Type:    Jump,
		JumpTyp: RJump,
	})
	index := len(genCode.ins)

	//if statement
	genCode.genStatement(statement.Statements)

	//fix jump val
	jumpTo := len(genCode.ins) - index + 1
	genCode.ins[index-1].Val = int64(jumpTo)
}

func (genCode *GenCode) genForStatement(statement ast.ForExpression) {

	genCode.genStatement(statement.Pre)

	begin := len(genCode.ins)
	genCode.genStatement(statement.Check)

	baseSP := genCode.sm.SP()
	genCode.pushIns(Instruction{
		Type:    Jump,
		JumpTyp: RJump,
		Val:     3,
	})
	genCode.pushIns(Instruction{
		Type:   Push,
		ValTyp: Bool,
		Val:    TRUE,
	})
	genCode.pushIns(Instruction{
		Type:    Jump,
		JumpTyp: RJump,
	})
	jumpStatement := len(genCode.ins)

	//if statement
	genCode.genStatement(statement.Statements)

	genCode.pushIns(Instruction{
		Type: IncStack,
		Val:  baseSP - genCode.sm.SP(),
	})
	genCode.genStatement(statement.Post)

	genCode.pushIns(Instruction{
		Type:   Push,
		ValTyp: Bool,
		Val:    TRUE,
	})
	genCode.pushIns(Instruction{
		Type:    Jump,
		JumpTyp: RJump,
		Val:     int64(begin - len(genCode.ins)),
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
		Type:   Load,
		Val:    index,
		symbol: genCode.symbolTable.addSymbol(label),
	})
}

func (genCode *GenCode) genArguments(statement *ast.CallStatement) {
	for _, argument := range statement.Arguments {
		genCode.genStatement(argument)
		if argument.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{Type: LoadR, Val: 1})
		}
	}
	genCode.pushIns(Instruction{
		Type:   Push,
		ValTyp: Int,
		Val:    int64(len(statement.Arguments)),
	})
	genCode.pushIns(Instruction{
		Type: StoreR,
		Val:  int64(0),
	})
	var R = int64(len(statement.Arguments))
	for range statement.Arguments {
		genCode.pushIns(Instruction{
			Type: StoreR,
			Val:  int64(R),
		})
		R--
	}
}

func (genCode *GenCode) genCallStatement(statement *ast.CallStatement) {

	//statement.ParentExp todo
	var retIP = int64(len(genCode.ins))
	switch function := statement.Function.(type) {
	case ast.GetVarStatement:
		index, ok := genCode.builtSymbolTable.getSymbol(function.Label)
		if ok == false { // push IP to stack for return
			genCode.pushIns(Instruction{Type: Push, ValTyp: IP})
		}

		genCode.genArguments(statement)

		if ok {
			genCode.pushIns(Instruction{
				Type: Call,
				Val:  index,
			})
		} else if index, ok := genCode.sm.load(function.Label); ok {
			genCode.pushIns(Instruction{
				Type: Load,
				Val:  index,
			})
			genCode.pushIns(Instruction{
				Type: CallO,
			})
		} else {
			genCode.pushIns(Instruction{
				Type:   Push,
				ValTyp: Bool,
				Val:    TRUE,
			})
			genCode.pushIns(Instruction{
				Type:    Jump,
				JumpTyp: DJump,
				Val:     -1, //todo link
			})
			genCode.toLinks = append(genCode.toLinks, toLink{
				label: function.Label,
				IP:    int64(len(genCode.ins)) - 1,
			})
		}
		if ok == false {
			genCode.ins[retIP].Val = int64(len(genCode.ins)) - retIP
		}

	case ast.PeriodStatement:
		genCode.pushIns(Instruction{Type: Push, ValTyp: IP})
		genCode.genStatement(function.Exp)
		genCode.pushIns(Instruction{
			Type:   LoadO,
			ValTyp: String,
			Str:    function.Val,
		})
		statement.Arguments = append(statement.Arguments, statement.ParentExp)
		genCode.genArguments(statement)

		genCode.pushIns(Instruction{
			Type: CallO,
		})
		genCode.ins[retIP].Val = int64(len(genCode.ins)) - retIP
	default:
		log.Panicf("unkown function type %s", reflect.TypeOf(function).String())
	}
}

type createObjectStatement struct {
	label string
}

func (c createObjectStatement) Invoke() runtime.Invokable {
	panic("implement me")
}

func (c createObjectStatement) GetType() lexer.Type {
	return lexer.CreateObjectStatementType
}

func (c createObjectStatement) String() string {
	panic("implement me")
}

func (genCode *GenCode) genCreateObjectStatement(statement createObjectStatement) {
	genCode.pushIns(Instruction{
		Type:   Push,
		ValTyp: Obj,
		Str:    statement.label,
		Val:    genCode.symbolTable.addSymbol(statement.label),
	})
}

func (genCode *GenCode) genObjectInitStatement(statement ast.ObjectInitStatement) {
	switch obj := statement.Exp.(type) {
	case ast.GetVarStatement:
		genCode.genCallStatement(&ast.CallStatement{
			ParentExp: nil,
			Function:  ast.GetVarStatement{Label: obj.Label + "." + objectInitFunctionName},
			Arguments: []ast.Expression{createObjectStatement{label: obj.Label}},
		})
		genCode.pushIns(Instruction{Type: LoadR, Val: 1})
	default:
		log.Panicln(reflect.TypeOf(obj).String())
	}

}

/*

translate object member function .add init function
*/

func (genCode *GenCode) genObject(label *runtime.Object) {
	genCode.genStatement(label.Pointer)
}

const closureLabel = "__Closure__"

func (genCode *GenCode) genFuncStatement(statement *ast.FuncExpression) {
	if statement.Closure {

		//generate function label
		hash := crc32.NewIEEE()
		hash.Write([]byte(statement.String()))
		statement.Label = fmt.Sprintf("lambda_%d", hash.Sum32())

		// link
		defer func() {
			genCode.toLinks = append(genCode.toLinks, toLink{
				label: statement.Label,
				IP:    int64(len(genCode.ins)),
			})
			genCode.pushIns(Instruction{
				Type:   Push,
				ValTyp: Lambda,
				Str:    statement.Label,
				Val:    -1, //to link
			})

			genCode.pushIns(Instruction{Type: MakeArray})
			// store closure val to function
			for _, label := range statement.ClosureLabel {
				index, ok := genCode.sm.load(label)
				if ok == false {
					log.Panicf("no find label `%s`", label)
				}
				genCode.pushIns(Instruction{
					Type: Load,
					Val:  index,
				})
				genCode.pushIns(Instruction{Type: Append})
			}
			genCode.pushIns(Instruction{
				Type: Load,
				Val:  -2, //top stack
			})
			genCode.pushIns(Instruction{
				Type: StoreO,
				Str:  closureLabel,
			})
		}()
	}

	genCode.sm.pushStackFrame()
	defer genCode.sm.popStackFrame()

	label := statement.Label
	if len(label) == 0 {
		label = strings.Join(statement.Labels, ".")
	}
	done := genCode.prepareGenFunction(label)
	defer done()

	genCode.pushIns(Instruction{
		Type:   Label,
		symbol: genCode.symbolTable.addSymbol(label),
	})
	genCode.pushIns(Instruction{
		Type: MakeStack,
	})

	if len(statement.Parameters) > 1 {
		if statement.Parameters[0] == "this" {
			statement.Parameters = append(statement.Parameters[1:], "this")
		}
	}
	// arguments
	for i := 0; i < len(statement.Parameters); i++ {
		genCode.pushIns(Instruction{
			Type: LoadR,
			Val:  int64(i + 1),
		})
		genCode.symbolTable.addSymbol(statement.Parameters[i])
		genCode.sm.Store(statement.Parameters[i])
	}

	for _, it := range statement.ClosureLabel {
		genCode.symbolTable.addSymbol(it)
		genCode.sm.Store(it)
	}

	if statement.Closure {
		genCode.pushIns(Instruction{Type: InitClosure})
	}

	// closure

	var last runtime.Invokable
	for _, statement := range statement.Statements {
		genCode.genStatement(statement)
		last = statement
	}
	switch last.(type) {
	case ast.ReturnStatement:
	default:
		genCode.pushIns(Instruction{Type: PopStack})
		genCode.pushIns(Instruction{Type: Ret})
	}
}

type objectInitStatement struct {
	functions []ast.FuncExpression
}

func (i objectInitStatement) Invoke() runtime.Invokable {
	panic("implement me")
}

func (i objectInitStatement) GetType() lexer.Type {
	panic("implement me")
}

func (i objectInitStatement) String() string {
	panic("implement me")
}

func (genCode *GenCode) genInitStatement(statement objectInitStatement) {
	for _, function := range statement.functions {
		genCode.pushIns(Instruction{
			Type:   Push,
			ValTyp: OFunc,
			Str:    function.Labels[1],
			Val:    -1, // to link
		})
		genCode.toLinks = append(genCode.toLinks, toLink{
			label: strings.Join(function.Labels, "."),
			IP:    int64(len(genCode.ins)) - 1,
		})
		genCode.pushIns(Instruction{
			Type: Load,
			Val:  0,
		})
		genCode.pushIns(Instruction{
			Type: StoreO,
			Str:  function.Labels[1],
		})
	}
}

const objectInitFunctionName = "__init__"

func (genCode *GenCode) genTypeObject(statement *ast.TypeObject) {
	var initStatement objectInitStatement
	for _, object := range statement.GetObjects() {
		switch obj := object.Pointer.(type) {
		case *ast.FuncExpression:
			initStatement.functions = append(initStatement.functions, *obj)
		default:
			log.Println("type", reflect.TypeOf(obj).String())
		}
	}
	for _, function := range initStatement.functions {
		genCode.genFuncStatement(&function)
	}
	//generate init function for object
	genCode.genFuncStatement(&ast.FuncExpression{
		Labels:     []string{statement.Label, objectInitFunctionName},
		Parameters: []string{"this"},
		Statements: []ast.Expression{initStatement,
			ast.ReturnStatement{
				Exp: ast.GetVarStatement{Label: "this"},
			}},
	})
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
		genCode.pushIns(Instruction{Type: LoadR, Val: 1})
	}
	genCode.pushIns(Instruction{
		Type: StoreR,
		Val:  1,
	})
	genCode.pushIns(Instruction{Type: PopStack})
	genCode.pushIns(Instruction{Type: Ret})
}

func (genCode *GenCode) GenExit() {
	genCode.pushIns(Instruction{
		Type: Exit,
	})
}

func (genCode *GenCode) genAssignStatement(statement ast.AssignStatement) {
	genCode.genStatement(statement.Exp)
	switch obj := statement.Left.(type) {
	case ast.GetVarStatement:
		if index, ok := genCode.sm.load(obj.Label); ok {
			genCode.pushIns(Instruction{
				Type: Store,
				Val:  index,
			})
		}
	case ast.PeriodStatement:
		genCode.genStatement(obj.Exp)
		genCode.pushIns(Instruction{
			Type: StoreO,
			Str:  obj.Val,
		})
	default:
		log.Panicln(reflect.TypeOf(obj).String())
	}
}

func (genCode *GenCode) genIncFieldStatement(statement ast.IncFieldStatement) {
	switch obj := statement.Exp.(type) {
	case ast.GetVarStatement:
		index, ok := genCode.sm.load(obj.Label)
		if ok == false {
			log.Panicln("no find label`" + obj.Label + "`")
		}
		genCode.pushIns(Instruction{
			Type: Load,
			Val:  index,
		})
		genCode.pushIns(Instruction{
			Type:   Push,
			ValTyp: Int,
			Val:    1,
		})
		genCode.pushIns(Instruction{
			Type: Add,
		})
		genCode.pushIns(Instruction{
			Type: Store,
			Val:  index,
		})
	case ast.PeriodStatement:
		genCode.genStatement(obj.Exp)
		genCode.pushIns(Instruction{
			Type: LoadO,
			Val:  0,
			Str:  obj.Val,
		})
		genCode.pushIns(Instruction{
			Type:   Push,
			ValTyp: Int,
			Val:    1,
		})
		genCode.pushIns(Instruction{
			Type: Add,
		})
		genCode.genStatement(obj.Exp)
		genCode.pushIns(Instruction{
			Type: StoreO,
			Str:  obj.Val,
		})
	}
}

func (genCode *GenCode) genNilObject(statement ast.NilObject) {
	genCode.pushIns(Instruction{
		Type:   Push,
		ValTyp: Nil,
		Str:    "nil",
	})
}
