package stackmachine

import (
	"fmt"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type InstType byte
type ValType byte
type CmpType byte

const (
	Less   CmpType = 0 + iota // <
	LessEQ                    // <=

	Push InstType = 0 + iota
	Pop
	Add
	Sub
	Load
	Store
	Call
	Cmp
	Jump

	TRUE  = 1
	FALSE = 0
)
const (
	Int ValType = 0 + iota
	Bool
	Mem
)

type SymbolTable struct {
	symbols   []string
	symbolMap map[string]int64
}

func (t *SymbolTable) addSymbol(s string) int64 {
	index, ok := t.symbolMap[s]
	if ok {
		return index
	}
	index = int64(len(t.symbols))
	t.symbols = append(t.symbols, s)
	t.symbolMap[s] = index
	return index
}

type Instruction struct {
	InstTyp InstType
	ValTyp  ValType
	CmpTyp  CmpType
	Val     int64
}

type ValObject struct {
	VType ValType
	Val   int64
}

type Machine struct {
	symbolTable  SymbolTable
	stack        []ValObject
	stackPointer int64
	instructions []Instruction
	pIns         int64
	mem          map[string]ValObject
}

func New() *Machine {
	return &Machine{
		symbolTable: SymbolTable{
			symbols:   []string{},
			symbolMap: map[string]int64{},
		},
		stack:        make([]ValObject, 1024*1024),
		stackPointer: -1,
		instructions: nil,
		pIns:         0,
		mem:          map[string]ValObject{},
	}
}

func (m *Machine) Run() {
	for m.pIns < int64(len(m.instructions)) {
		ins := m.instructions[m.pIns]
		//	log.Println(ins)
		switch ins.InstTyp {
		case Push:
			m.stackPointer++
			m.stack[m.stackPointer] = ValObject{
				VType: ins.ValTyp,
				Val:   ins.Val,
			}
		case Pop:
			m.stackPointer--
		case Add, Sub, Cmp:
			operand2 := m.stack[m.stackPointer]
			m.stackPointer--
			operand1 := m.stack[m.stackPointer]
			m.stackPointer--
			var result ValObject
			switch operand1.VType {
			case Int:
				switch operand1.VType {
				case Int:
					switch ins.InstTyp {
					case Cmp:
						//						log.Println(operand1, operand2)
						result.VType = Bool
						switch ins.CmpTyp {
						case Less:
							result.Val = FALSE
							if operand1.Val < operand2.Val {
								result.Val = TRUE
							}
						case LessEQ:
							result.Val = FALSE
							if operand1.Val <= operand2.Val {
								result.Val = TRUE
							}
						}
					case Add:
						result.VType = Int
						result.Val = operand1.Val + operand2.Val
					case Sub:
						result.VType = Int
						result.Val = operand1.Val - operand2.Val
					}
				}
			}
			//			log.Println(operand1, operand2, result)
			m.stackPointer++
			m.stack[m.stackPointer] = result
		case Store:
			val := m.stack[m.stackPointer]
			m.stackPointer--
			symbol := m.symbolTable.symbols[ins.Val]
			m.store(symbol, val)
		case Load:
			symbol := m.symbolTable.symbols[ins.Val]
			val := m.load(symbol)
			m.stackPointer++
			m.stack[m.stackPointer] = val
		case Call:
			arguments := m.stack[m.stackPointer]
			m.stackPointer--
			//log.Println(arguments.Val)
			m.CallFunc(ins.Val, m.stack[m.stackPointer-arguments.Val+1:m.stackPointer+1]...)
			m.stackPointer -= arguments.Val
		case Jump:
			check := m.stack[m.stackPointer]
			m.stackPointer--
			if check.VType != Bool {
				log.Panic("expect bool value for check")
			}
			if check.Val == TRUE {
				//				log.Println("true",ins.Val)
				m.pIns = ins.Val
				continue
			}
			log.Println("false")
		}
		m.pIns++
		//log.Println(m.stack[:m.stackPointer+1])
	}
}

func (m *Machine) store(symbol string, val ValObject) {
	m.mem[symbol] = val
}

func (m *Machine) CallFunc(funcIndex int64, object ...ValObject) {
	if funcIndex == 0 {
		fmt.Println(object)
	}
}

func (m *Machine) load(symbol string) ValObject {
	return m.mem[symbol]
}
func TestStore(t *testing.T) {
	//var a = 1
	var m = New()
	index := m.symbolTable.addSymbol("a")
	m.instructions = []Instruction{
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1,
		},
		{
			InstTyp: Store,
			ValTyp:  Mem,
			Val:     index,
		},
	}
	m.Run()

}

func TestLoad(t *testing.T) {
	//var a =100
	//print(a)

	var m = New()
	index := m.symbolTable.addSymbol("a")
	m.instructions = []Instruction{
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     100,
		},
		{
			InstTyp: Store,
			ValTyp:  Mem,
			Val:     index,
		},
		{
			InstTyp: Load,
			ValTyp:  Mem,
			Val:     index,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1,
		},
		{
			InstTyp: Call,
		},
	}
	m.Run()
}

func TestAddNum(t *testing.T) {
	var m = New()
	// print(1 +1 +2)
	m.instructions = []Instruction{
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1,
		},
		{
			InstTyp: Add,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     2,
		},
		{
			InstTyp: Add,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1,
		},
		{
			InstTyp: Call,
			Val:     0,
		},
	}
	m.Run()
}

func TestJump(t *testing.T) {
	var m = New()
	//for{print(1)}

	m.instructions = []Instruction{
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1, //argument count 1
		},
		{
			InstTyp: Call,
		},
		{
			InstTyp: Push,
			ValTyp:  Bool,
			Val:     TRUE,
		},
		{
			InstTyp: Jump,
			Val:     0, // jump to begin 0
		},
	}
	m.Run()
}

func TestFor(t *testing.T) {
	var m = New()
	/*
		for a:= 1; a < 10;a++{
			print(a)
		}
	*/

	index := m.symbolTable.addSymbol("a")

	m.instructions = []Instruction{
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     0,
		},
		{
			InstTyp: Store,
			Val:     index, //a := 1
		},
		{
			InstTyp: Load,
			Val:     index,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     10,
		},
		{
			InstTyp: Cmp,
			CmpTyp:  Less, // a < 10
		},
		{
			InstTyp: Jump,
			Val:     8, //todo jump to print
		},
		{
			InstTyp: Push,
			ValTyp:  Bool,
			Val:     TRUE,
		},
		{
			InstTyp: Jump,
			Val:     18, //todo jump to end of for
		},
		{
			InstTyp: Load,
			Val:     index,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1, // argument count
		},
		{
			InstTyp: Call, // call print
		},
		// a++ => a= a+1
		{
			InstTyp: Load,
			Val:     index,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1,
		},
		{
			InstTyp: Add,
		},
		{
			InstTyp: Store,
			Val:     index,
		},
		{
			InstTyp: Push,
			ValTyp:  Bool,
			Val:     TRUE,
		},
		{
			InstTyp: Jump,
			Val:     2, // jump to a < 10
		},
	}
	m.Run()
}
