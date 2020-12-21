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

const (
	Push InstType = 0 + iota
	Pop
	Add
	Sub
	Load
	Store
	Call
	Jump
)
const (
	Int ValType = 0 + iota
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
	VaType  ValType
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
	pIns         int
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
	for m.pIns < len(m.instructions) {
		ins := m.instructions[m.pIns]
		switch ins.InstTyp {
		case Push:
			switch ins.VaType {
			case Int:
				m.stackPointer++
				m.stack[m.stackPointer] = ValObject{
					VType: Int,
					Val:   ins.Val,
				}
			}
		case Pop:
			m.stackPointer--
		case Add, Sub:
			operand1 := m.stack[m.stackPointer]
			m.stackPointer--
			operand2 := m.stack[m.stackPointer]
			m.stackPointer--
			var result ValObject
			switch operand1.VType {
			case Int:
				switch operand1.VType {
				case Int:
					switch ins.InstTyp {
					case Add:
						result.Val = operand1.Val + operand2.Val
					case Sub:
						result.Val = operand1.Val - operand2.Val
					}
					result.VType = Int
				}
			}
			log.Println(operand1, operand2, result)
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
			log.Println(arguments.Val)
			m.CallFunc(ins.Val, m.stack[m.stackPointer-arguments.Val+1:m.stackPointer+1]...)
			m.stackPointer -= arguments.Val
		}
		m.pIns++
		log.Println(m.stack[:m.stackPointer+1])
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
			VaType:  Int,
			Val:     1,
		},
		{
			InstTyp: Store,
			VaType:  Mem,
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
			VaType:  Int,
			Val:     100,
		},
		{
			InstTyp: Store,
			VaType:  Mem,
			Val:     index,
		},
		{
			InstTyp: Load,
			VaType:  Mem,
			Val:     index,
		},
		{
			InstTyp: Push,
			VaType:  Int,
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
			VaType:  Int,
			Val:     1,
		},
		{
			InstTyp: Push,
			VaType:  Int,
			Val:     1,
		},
		{
			InstTyp: Add,
		},
		{
			InstTyp: Push,
			VaType:  Int,
			Val:     2,
		},
		{
			InstTyp: Add,
		},
		{
			InstTyp: Push,
			VaType:  Int,
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
	/*
	if 2 > 1{
		print(2)
	}else{
		print(1)
	}
	 */
}