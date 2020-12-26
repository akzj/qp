package stackmachine

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
)

type InstType byte
type ValType byte
type CmpType byte
type JumpType byte
type ResetStackType byte

const (
	Less      CmpType = iota // <
	LessEQ                   // <=
	Greater                  // >
	GreaterEQ                // >=
	Equal                    // ==

	Push InstType = iota
	Pop
	Add
	Sub
	Load
	Store
	Call
	CallO
	Cmp
	Jump
	Ret
	Label
	Exit
	ResetStack // reset stack
	StoreR     // stack -> register
	LoadR      // register -> stack
	PushS      // set stack base point
	PopS       // pop stack
	LoadO      // load from object

	TRUE  int64 = 1
	FALSE int64 = 0

	Int ValType = iota
	Bool
	Mem
	IP
	String
	Func

	DJump JumpType = 0
	RJump JumpType = 1

	ResetStackD ResetStackType = 0
	ResetStackR ResetStackType = 1
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
	InstTyp       InstType
	ValTyp        ValType
	CmpTyp        CmpType
	JumpTyp       JumpType
	ResetStackTyp ResetStackType
	symbol        int64
	Val           int64
	Str           string
}

func (i Instruction) String(table, builtIn *SymbolTable) string {
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
		if i.ValTyp == IP {
			return "push ip " + strconv.FormatInt(i.Val, 10)
		} else if i.ValTyp == Int {
			return "push " + strconv.FormatInt(i.Val, 10)
		} else if i.ValTyp == Bool {
			if i.Val == TRUE {
				return "push true"
			} else {
				return "push false"
			}
		} else if i.ValTyp == String {
			return "push " + i.Str
		}
	case Exit:
		return "exit"
	case Pop:
		return "pop"
	case Load:
		return "load " + table.symbols[i.symbol] + " " + strconv.FormatInt(i.Val, 10)
	case LoadR:
		return "loadR"
	case LoadO:
		return "LoadO"
	case Store:
		return "store "
	case StoreR:
		return "storeR"
	case Call:
		return "call " + builtIn.symbols[i.symbol]
	case CallO:
		return "CallO"
	case Ret:
		return "return"
	case Label:
		return table.symbols[i.symbol] + ":"
	case ResetStack:
		return "reset " + strconv.FormatInt(i.Val, 10)
	case PushS:
		return "PushS"
	case PopS:
		return "pops"
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
		case Equal:
			return "cmp =="
		}
	}
	log.Panicln("unknown instruction")
	return ""
}

type Object struct {
	VType ValType
	Val   int64
	Str   string
}

func (o Object) String() string {
	if o.VType == Int {
		return strconv.FormatInt(o.Val, 10)
	} else if o.VType == Bool {
		if o.Val == TRUE {
			return "true"
		} else {
			return "false"
		}
	} else if o.VType == String {
		return o.Str
	} else {
		return fmt.Sprintf("{%d %d}", o.VType, o.Val)
	}
}

type StackFrame struct {
	stack []Object
	SP    int64
}

type Machine struct {
	symbolTable        *SymbolTable
	builtInSymbolTable *SymbolTable
	heap               []Object
	stack              []Object
	stackFrames        []StackFrame
	SP                 int64
	instructions       []Instruction
	IP                 int64
	mem                map[string]Object

	R1 [32]Object //return val
}

func New() *Machine {
	return &Machine{
		symbolTable:        NewSymbolTable(),
		builtInSymbolTable: getBuiltInSymbolTable(),
		stack:              make([]Object, 1024*1024),
		SP:                 -1,
		instructions:       nil,
		IP:                 0,
		mem:                map[string]Object{},
	}
}

func getBuiltInSymbolTable() *SymbolTable {
	table := NewSymbolTable()
	for _, builtIn := range BuiltInFunctions {
		table.addSymbol(builtIn.Name)
	}
	return table
}

func (m *Machine) Run() {
	for m.IP < int64(len(m.instructions)) {
		ins := m.instructions[m.IP]
		log.Print(ins.String(m.symbolTable, m.builtInSymbolTable), " SP: ", m.SP)
		switch ins.InstTyp {
		case Push:
			m.SP++
			if ins.ValTyp == IP {
				m.stack[m.SP] = Object{
					VType: ins.ValTyp,
					Val:   m.IP + ins.Val, //return addr
				}
			} else if ins.ValTyp == String {
				m.stack[m.SP] = Object{
					VType: ins.ValTyp,
					Str:   ins.Str,
				}
			} else {
				m.stack[m.SP] = Object{
					VType: ins.ValTyp,
					Val:   ins.Val,
				}
			}
		case Pop:
			m.SP--
		case PushS:
			m.stackFrames = append(m.stackFrames, StackFrame{SP: m.SP, stack: m.stack})
			m.stack = m.stack[m.SP+1:]
			m.SP = -1
			//log.Println(m.stackFrames[len(m.stackFrames)-1])
		case PopS:
			//			log.Println(m.SP, m.stack[:m.SP+1])
			//log.Println(m.stackFrames[len(m.stackFrames)-1])
			frame := m.stackFrames[len(m.stackFrames)-1]
			m.stack = frame.stack
			m.SP = frame.SP
			m.stackFrames = m.stackFrames[:len(m.stackFrames)-1]
			//		log.Println(m.SP, m.stack[:m.SP+1])
		case Add, Sub, Cmp:
			operand2 := m.stack[m.SP]
			m.SP--
			operand1 := m.stack[m.SP]
			m.SP--
			var result Object
			switch operand1.VType {
			case Int:
				switch operand1.VType {
				case Int:
					switch ins.InstTyp {
					case Cmp:
						//log.Println(operand1, operand2)
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
						case Equal:
							if operand1.Val == operand2.Val {
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
			//log.Println(operand1, operand2, result)
			m.SP++
			m.stack[m.SP] = result
		case Store:
			val := m.stack[m.SP]
			m.SP--
			m.stack[ins.Val] = val
		case ResetStack:
			if ins.ResetStackTyp == ResetStackR {
				m.SP -= ins.Val
			} else {
				m.SP = ins.Val
			}
		case StoreR:
			object := m.stack[m.SP]
			m.SP--
			m.R1[ins.Val] = object
		case LoadR:
			m.SP++
			m.stack[m.SP] = m.R1[ins.Val]
		case Load:
			sp := ins.Val
			if sp > m.SP {
				log.Panicln("stack error", sp, m.SP)
			}
			m.SP++
			m.stack[m.SP] = m.stack[sp]
		case Ret:
			IP := m.stack[m.SP]
			m.SP--
			m.IP = IP.Val
			continue
		case LoadO:
			object := m.stack[m.SP]
			switch object.VType {
			case String:
				index, ok := BuiltInFunctionsIndex["string."+ins.Str]
				if ok == false {
					log.Panicln("no find string." + ins.Str)
				}
				m.SP++
				m.stack[m.SP] = Object{
					VType: Func,
					Val:   index,
				}
			}
		case Call:
			count := m.R1[0].Val
			objects := m.CallFunc(ins.Val, m.R1[1:count+1]...)
			m.R1[0].Val = int64(len(objects))
			for index, obj := range objects {
				m.R1[index+1] = obj
			}
		case CallO:
			obj := m.stack[m.SP]
			m.SP--
			objects := m.CallFunc(obj.Val, m.stack[m.SP])
			m.SP--

			m.R1[0].Val = int64(len(objects))
			for index, obj := range objects {
				m.R1[index+1] = obj
			}

		case Label:
		case Exit:
			return
		case Jump:
			check := m.stack[m.SP]
			m.SP--
			if check.VType != Bool {
				log.Panicln("expect bool value for check", m.IP, m.SP)
			}
			if check.Val == TRUE {
				//log.Println("true",ins.Val)
				if ins.JumpTyp == DJump {
					m.IP = ins.Val
				} else {
					m.IP += ins.Val
				}
				continue
			}
			//			log.Println("false")
		}
		m.IP++
		log.Println(m.stack[:m.SP+1])
	}
}

func (m *Machine) store(symbol string, val Object) {
	m.mem[symbol] = val
}

func (m *Machine) CallFunc(funcIndex int64, object ...Object) []Object {
	return BuiltInFunctions[funcIndex].Call(object...)
}

func (m *Machine) load(symbol string) Object {
	return m.mem[symbol]
}

func (m *Machine) String() string {
	var buffer bytes.Buffer
	for _, it := range m.instructions {
		buffer.WriteString("\t" + it.String(m.symbolTable, m.builtInSymbolTable))
		buffer.WriteString("\n")
	}
	return buffer.String()
}
