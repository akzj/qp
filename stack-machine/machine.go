package stackmachine

import (
	"bytes"
	"log"
	"strconv"
)

type InstType byte
type ValType byte
type CmpType byte
type JumpType byte

const (
	Less      CmpType = iota // <
	LessEQ                   // <=
	Greater                  // >
	GreaterEQ                // >=

	Push InstType = iota
	Pop
	Add
	Sub
	Load
	Store
	Call
	Cmp
	Jump
	Ret

	TRUE  int64 = 1
	FALSE int64 = 0

	Int ValType = iota
	Bool
	Mem

	DJump JumpType = 0
	RJump JumpType = 1
)

type SymbolTable struct {
	symbols   []string
	symbolMap map[string]int64
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{
		symbols:   []string{},
		symbolMap: map[string]int64{},
	}
	for _, function := range BuiltInFunctions {
		st.addSymbol(function.Name)
	}
	return st
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

func (t *SymbolTable) getSymbol(s string) (int64, bool) {
	index, ok := t.symbolMap[s]
	return index, ok
}

type Instruction struct {
	InstTyp InstType
	ValTyp  ValType
	CmpTyp  CmpType
	JumpTyp JumpType
	Val     int64
}

func (i Instruction) String(table *SymbolTable) string {
	switch i.InstTyp {
	case Add:
		return "add"
	case Sub:
		return "sub"
	case Jump:
		if i.JumpTyp == DJump {
			return "jump D " + strconv.FormatInt(i.Val, 10)
		} else {
			return "jump R " + strconv.FormatInt(i.Val, 10)
		}
	case Push:
		return "push " + strconv.FormatInt(i.Val, 10)
	case Pop:
		return "pop"
	case Load:
		return "load " + table.symbols[i.Val]
	case Store:
		return "store " + table.symbols[i.Val]
	case Call:
		return "call " + table.symbols[i.Val]
	case Ret:
		return "return"
	case Cmp:
		switch i.CmpTyp {
		case Less:
			return "cmp <"
		case LessEQ:
			return "cmp <="
		case Greater:
			return "cmp >"
		case GreaterEQ:
			return "cmp >="
		}
	}
	panic("unknown instruction")
}

type Object struct {
	VType ValType
	Val   int64
}

type Machine struct {
	symbolTable  *SymbolTable
	heap         []Object
	stack        []Object
	stackPointer int64
	instructions []Instruction
	IP           int64
	mem          map[string]Object
}

func New() *Machine {
	return &Machine{
		symbolTable:  NewSymbolTable(),
		stack:        make([]Object, 1024*1024),
		stackPointer: -1,
		instructions: nil,
		IP:           0,
		mem:          map[string]Object{},
	}
}

func (m *Machine) Run() {
	for m.IP < int64(len(m.instructions)) {
		ins := m.instructions[m.IP]
		//	log.Println(ins)
		switch ins.InstTyp {
		case Push:
			m.stackPointer++
			m.stack[m.stackPointer] = Object{
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
			var result Object
			switch operand1.VType {
			case Int:
				switch operand1.VType {
				case Int:
					switch ins.InstTyp {
					case Cmp:
						//						log.Println(operand1, operand2)
						result.VType = Bool
						result.Val = FALSE
						switch ins.CmpTyp {
						case Less:
							if operand1.Val < operand2.Val {
								result.Val = TRUE
							}
						case LessEQ:
							if operand1.Val <= operand2.Val {
								result.Val = TRUE
							}
						case Greater:
							if operand1.Val > operand2.Val {
								result.Val = TRUE
							}
						case GreaterEQ:
							if operand1.Val >= operand2.Val {
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
		case Ret:
			IP := m.stack[m.stackPointer]
			m.stackPointer--
			m.IP = IP.Val
			continue
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
				if ins.JumpTyp == DJump {
					m.IP = ins.Val
				} else {
					m.IP += ins.Val
				}
				continue
			}
			log.Println("false")
		}
		m.IP++
		//log.Println(m.stack[:m.stackPointer+1])
	}
}

func (m *Machine) store(symbol string, val Object) {
	m.mem[symbol] = val
}

func (m *Machine) CallFunc(funcIndex int64, object ...Object) {
	BuiltInFunctions[funcIndex].Call(object...)
}

func (m *Machine) load(symbol string) Object {
	return m.mem[symbol]
}

func (m *Machine) String() string {
	var buffer bytes.Buffer
	for _, it := range m.instructions {
		buffer.WriteString(it.String(m.symbolTable))
		buffer.WriteString("\n")
	}
	return buffer.String()
}
