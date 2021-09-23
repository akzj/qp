package stackmachine

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"log"
	"reflect"
	"strconv"
	"strings"

	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
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

type stackSymbol struct {
	symbol string
	sp     int
}
type stackFrame struct {
	function bool
	stack    []stackSymbol
	sp       int
}
type StackManager struct {
	currStack  stackFrame
	stackFrame []stackFrame
}

func NewStackManager() *StackManager {
	return &StackManager{}
}

func (s *StackManager) Store(symbol string) int64 {
	if len(symbol) != 0 {
		for _, label := range s.currStack.stack {
			if label.symbol == symbol {
				log.Panicln("redefine symbol", symbol, s.stackFrame, s.currStack)
			}
		}
	}
	s.currStack.stack = append(s.currStack.stack, stackSymbol{
		symbol: symbol,
		sp:     s.currStack.sp,
	})
	s.currStack.sp++
	return s.SP()
}

func (s *StackManager) load(label string) (int64, bool) {
	stacks := append([]stackFrame{}, s.currStack)
	stacks = append(stacks, s.stackFrame...)
	for j := len(stacks) - 1; j >= 0; j-- {
		stack := stacks[j]
		for i := len(stack.stack) - 1; i >= 0; i-- {
			if stack.stack[i].symbol == label {
				return int64(stack.stack[i].sp), true
			}
		}
	}
	return -1, false
}

func (s *StackManager) pushStackFrame(funcStack bool) {
	s.stackFrame = append(s.stackFrame, s.currStack)
	if !funcStack {
		s.currStack.stack = nil
	} else {
		s.currStack = stackFrame{function: funcStack}
	}
}
func (s *StackManager) popStackFrame() {
	s.currStack = s.stackFrame[len(s.stackFrame)-1]
	s.stackFrame = s.stackFrame[:len(s.stackFrame)-1]
}

func (s *StackManager) SP() int64 {
	return int64(s.currStack.sp)
}

type CodeGenerator struct {
	symbolTable      *SymbolTable
	builtSymbolTable *SymbolTable
	ins              []Instruction
	toLinks          []toLink
	sm               *StackManager
	funcInstructions map[string]FuncInstruction
}

func NewCodeGenerator() *CodeGenerator {
	gc := &CodeGenerator{
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

func (genCode *CodeGenerator) String() string {
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

func (genCode *CodeGenerator) pushIns(instruction Instruction) {
	genCode.ins = append(genCode.ins, instruction)
}

func (genCode *CodeGenerator) Gen(statements []ast.Expression) *CodeGenerator {
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

func (genCode *CodeGenerator) genStatement(statement runtime.Invokable) int {
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
		return 1
	case ast.VarInitExpression:
		genCode.genStatement(statement.Exp)
		if statement.Exp.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{Type: LoadR, Val: 1})
		}
		genCode.genStoreIns(statement.Name)
		return 1
	case ast.VarStatement:
		genCode.genStatement(statement.Exp)
		if statement.Exp.GetType() == lexer.CallType {
			genCode.pushIns(Instruction{Type: LoadR, Val: 1})
		}
		genCode.genStoreIns(statement.Label)
		return 1
	case ast.GetVarStatement:
		genCode.genLoadIns(statement.Label)
	case ast.IfExpression:
		genCode.genIfStatement(statement)
	case ast.Expressions:
		var stackSize int
		for _, next := range statement {
			stackSize += genCode.genStatement(next)
		}
		return stackSize
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
	case ast.ParenthesisExpression:
		genCode.genStatement(statement.Exp)
	default:
		log.Panicf("unknown statement %s", reflect.TypeOf(statement).String())
	}
	return 0
}

func (genCode *CodeGenerator) genOpCode(op lexer.Type) {
	switch op {
	case lexer.MulOpType:
		genCode.pushIns(Instruction{
			Type: Mul,
		})
	case lexer.AddType:
		genCode.pushIns(Instruction{
			Type: Add,
		})
	case lexer.ModOpType:
		genCode.pushIns(Instruction{
			Type: Mod,
		})
	case lexer.DivOpType:
		genCode.pushIns(Instruction{
			Type: Div,
		})
	case lexer.LessType:
		genCode.pushIns(Instruction{
			Type:   Cmp,
			CmpTyp: Less,
		})
	case lexer.LessEqualType:
		genCode.pushIns(Instruction{
			Type:   Cmp,
			CmpTyp: LessEQ,
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

func (genCode *CodeGenerator) genStoreIns(label string) {
	genCode.sm.Store(label)
	genCode.symbolTable.addSymbol(label)
}

func (genCode *CodeGenerator) genIfStatement(statement ast.IfExpression) {
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
	genCode.sm.pushStackFrame(false)
	stackSize := genCode.genStatement(statement.Statements)
	genCode.sm.popStackFrame()

	//clear if expression stack
	genCode.pushIns(Instruction{
		Type: MoveStack,
		Val:  -int64(stackSize),
	})

	//fix jump val
	jumpTo := len(genCode.ins) - index + 1
	genCode.ins[index-1].Val = int64(jumpTo)
}

func (genCode *CodeGenerator) genForStatement(statement ast.ForExpression) {

	genCode.sm.pushStackFrame(false)
	preStackSize := genCode.genStatement(statement.Pre)

	begin := len(genCode.ins)
	genCode.genStatement(statement.Check)

	if preStackSize > 0 {
		genCode.pushIns(Instruction{
			Type:    Jump,
			JumpTyp: RJump,
			Val:     4,
		})
		genCode.pushIns(Instruction{
			Type: MoveStack,
			Val:  -int64(preStackSize),
		})
	} else {
		genCode.pushIns(Instruction{
			Type:    Jump,
			JumpTyp: RJump,
			Val:     3,
		})
	}
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

	// statement

	stackSize := genCode.genStatement(statement.Statements)

	//reset for expression Stack
	if stackSize > 0 {
		genCode.pushIns(Instruction{
			Type: MoveStack,
			Val:  -int64(stackSize),
		})
	}

	genCode.genStatement(statement.Post)
	genCode.sm.pushStackFrame(false)

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

func (genCode *CodeGenerator) genLoadIns(label string) {
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

func (genCode *CodeGenerator) genArguments(statement *ast.CallStatement) {
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

func (genCode *CodeGenerator) genCallStatement(statement *ast.CallStatement) {

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

func (genCode *CodeGenerator) genCreateObjectStatement(statement createObjectStatement) {
	genCode.pushIns(Instruction{
		Type:   Push,
		ValTyp: Obj,
		Str:    statement.label,
		Val:    genCode.symbolTable.addSymbol(statement.label),
	})
}

func (genCode *CodeGenerator) genObjectInitStatement(statement ast.ObjectInitStatement) {
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

func (genCode *CodeGenerator) genObject(label *runtime.Object) {
	genCode.genStatement(label.Pointer)
}

const closureLabel = "__Closure__"

func (genCode *CodeGenerator) genFuncStatement(statement *ast.FuncExpression) {
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

	genCode.sm.pushStackFrame(true)
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

func (genCode *CodeGenerator) genInitStatement(statement objectInitStatement) {
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

func (genCode *CodeGenerator) genTypeObject(statement *ast.TypeObject) {
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

func (genCode *CodeGenerator) prepareGenFunction(label string) func() {
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

func (genCode *CodeGenerator) genReturnStatement(statement ast.ReturnStatement) {
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

func (genCode *CodeGenerator) GenExit() {
	genCode.pushIns(Instruction{
		Type: Exit,
	})
}

func (genCode *CodeGenerator) genAssignStatement(statement ast.AssignStatement) {
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

func (genCode *CodeGenerator) genIncFieldStatement(statement ast.IncFieldStatement) {
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

func (genCode *CodeGenerator) genNilObject(statement ast.NilObject) {
	genCode.pushIns(Instruction{
		Type:   Push,
		ValTyp: Nil,
		Str:    "nil",
	})
}
