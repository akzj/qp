package stackmachine

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"time"
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
	NoEqual                  // !=
)
const (
	Push InstType = iota + 0
	Add
	Sub
	And
	Load
	Store
	Call
	CallO // call object member function
	Cmp
	Jump
	Ret
	Label
	Exit
	IncStack    // update stack
	StoreR      // stack -> register
	LoadR       // register -> stack
	MakeStack   // make stack for function call
	PopStack    // pop stack for return function call
	LoadO       // load from object
	StoreO      // store to object
	MakeArray   // make array
	Append      // append to array
	InitClosure // init lambda closure

	TRUE  int64 = 1
	FALSE int64 = 0
)
const (
	Int ValType = iota
	Bool
	Mem
	IP
	String
	BFunc  // built in function
	OFunc  // object member function
	Lambda //
	Time
	Duration
	Obj
	Array
	Nil // nilObject

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
	Type          InstType
	ValTyp        ValType
	CmpTyp        CmpType
	JumpTyp       JumpType
	ResetStackTyp ResetStackType
	symbol        int64
	Val           int64
	Str           string
}

func (i Instruction) String(table, builtIn *SymbolTable) string {
	switch i.Type {
	case Add:
		return "add"
	case Sub:
		return "sub"
	case And:
		return "&&"
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
			return "push \"" + i.Str + "\""
		} else if i.ValTyp == OFunc || i.ValTyp == Lambda {
			return "push func \"" + i.Str + "\" " + strconv.FormatInt(i.Val, 10)
		} else if i.ValTyp == Obj {
			return "push obj \"" + i.Str + "\""
		} else if i.ValTyp == Nil {
			return "nil"
		} else {
			panic(i.ValTyp)
		}
	case Exit:
		return "exit"
	case InitClosure:
		return "init_closure"
	case StoreO:
		return "StoreO \"" + i.Str + "\""
	case MakeArray:
		return "make_array"
	case Append:
		return "append"
	case Load:
		return "load " + table.symbols[i.symbol] + " " + strconv.FormatInt(i.Val, 10)
	case LoadR:
		return "loadR " + strconv.FormatInt(i.Val, 10)
	case LoadO:
		return "LoadO " + i.Str
	case Store:
		return "store "
	case StoreR:
		return "storeR " + strconv.FormatInt(i.Val, 10)
	case Call:
		return "call " + builtIn.symbols[i.symbol]
	case CallO:
		return "CallO "
	case Ret:
		return "return"
	case Label:
		return table.symbols[i.symbol] + ":"
	case IncStack:
		return "reset " + strconv.FormatInt(i.Val, 10)
	case MakeStack:
		return "make_stack"
	case PopStack:
		return "pop_stack"
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
		case NoEqual:
			return "cmp !="
		}
	}
	log.Panicln("unknown instruction", i)
	return ""
}

type Object struct {
	Type ValType
	Int  int64
	Obj  interface{}
}
type ObjectArray []Object
type objectMap map[string]*Object

func (obj *Object) loadObj(label string) *Object {
	if obj.Obj == nil {
		obj.Obj = make(objectMap)
	}
	o, ok := obj.Obj.(objectMap)[label]
	if ok {
		return o
	}
	o = &Object{}
	obj.Obj.(objectMap)[label] = o
	return o
}

func (obj *Object) Store(str string, ele Object) {
	if obj.Obj == nil {
		obj.Obj = make(objectMap)
	}
	obj.Obj.(objectMap)[str] = &ele
}

func (obj Object) String() string {
	if obj.Type == Int {
		return strconv.FormatInt(obj.Int, 10)
	} else if obj.Type == Bool {
		if obj.Int == TRUE {
			return "true"
		} else {
			return "false"
		}
	} else if obj.Type == String {
		return obj.Obj.(string)
	} else if obj.Type == Time {
		return obj.Obj.(time.Time).String()
	} else if obj.Type == Duration {
		return time.Duration(obj.Int).String()
	} else if obj.Type == Obj {
		return "object"
	} else if obj.Type == OFunc || obj.Type == BFunc {
		return "{ function " + strconv.FormatInt(obj.Int, 10) + " }"
	} else if obj.Type == Lambda {
		return "{ lambda " + strconv.FormatInt(obj.Int, 10) + " }"
	} else if obj.Type == Array {
		return fmt.Sprintf("%+v", obj.Obj)
	} else {
		return fmt.Sprintf("{%d %d}", obj.Type, obj.Int)
	}
}

type StackFrame struct {
	stack []Object
	SP    int64
}

type Machine struct {
	symbolTable        *SymbolTable
	builtInSymbolTable *SymbolTable
	stack              []Object
	stackFrames        []StackFrame
	SP                 int64
	instructions       []Instruction
	IP                 int64
	R                  [32]Object //register
	closure            ObjectArray
}

func New() *Machine {
	return &Machine{
		symbolTable:        NewSymbolTable(),
		builtInSymbolTable: getBuiltInSymbolTable(),
		stack:              make([]Object, 1024*1024),
		SP:                 -1,
		instructions:       nil,
		IP:                 0,
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
	var tick int64
	defer func() {
		log.Println("tick", tick)
	}()
	for m.IP < int64(len(m.instructions)) {
		tick++
		ins := m.instructions[m.IP]
		log.Print(ins.String(m.symbolTable, m.builtInSymbolTable), " SP: ", m.SP)
		switch ins.Type {
		case Push:
			m.SP++
			m.stack[m.SP].Type = ins.ValTyp
			if ins.ValTyp == IP {
				m.stack[m.SP].Int = m.IP + ins.Val //return addr
			} else if ins.ValTyp == String {
				m.stack[m.SP].Obj = ins.Str
			} else if ins.ValTyp == Obj {
				m.stack[m.SP].Obj = make(objectMap)
			} else if ins.ValTyp == Lambda {
				m.stack[m.SP].Obj = make(objectMap)
				m.stack[m.SP].Int = ins.Val
			} else {
				m.stack[m.SP].Int = ins.Val
			}
		case MakeStack:
			m.stackFrames = append(m.stackFrames, StackFrame{SP: m.SP, stack: m.stack})
			m.stack = m.stack[m.SP+1:]
			m.SP = -1
		case PopStack:
			frame := m.stackFrames[len(m.stackFrames)-1]
			m.stack = frame.stack
			m.SP = frame.SP
			m.stackFrames = m.stackFrames[:len(m.stackFrames)-1]

		case Add, Sub, Cmp:
			operand2 := &m.stack[m.SP]
			m.SP--
			operand1 := &m.stack[m.SP]
			m.SP--
			var result Object
			switch operand1.Type {
			case Int:
				switch operand2.Type {
				case Int:
					switch ins.Type {
					case Cmp:
						//log.Println(operand1, operand2)
						result.Type = Bool
						result.Int = FALSE
						switch ins.CmpTyp {
						case Less:
							if operand1.Int < operand2.Int {
								result.Int = TRUE
							}
						case LessEQ:
							if operand1.Int <= operand2.Int {
								result.Int = TRUE
							}
						case Greater:
							if operand1.Int > operand2.Int {
								result.Int = TRUE
							}
						case GreaterEQ:
							if operand1.Int >= operand2.Int {
								result.Int = TRUE
							}
						case Equal:
							if operand1.Int == operand2.Int {
								result.Int = TRUE
							}
						case NoEqual:
							if operand1.Int != operand2.Int {
								result.Int = TRUE
							}
						default:
							log.Panicln("unknown instruction", ins, m.IP)
						}
					case Add:
						result.Type = Int
						result.Int = operand1.Int + operand2.Int
					case Sub:
						result.Type = Int
						result.Int = operand1.Int - operand2.Int
					default:
						log.Panicln("unknown instruction", ins, m.IP)
					}
				case Nil:
					result.Type = Bool
					result.Int = FALSE
					switch ins.CmpTyp {
					case Equal:
					case NoEqual:
						result.Int = TRUE
					}
				default:
					log.Panicln("unknown instruction", ins.String(m.symbolTable, m.builtInSymbolTable), m.IP, operand2.String())
				}
			case Time:
				switch operand2.Type {
				case Time:
					switch ins.Type {
					case Sub:
						result.Type = Duration
						result.Int = int64(operand1.Obj.(time.Time).Sub(operand2.Obj.(time.Time)))
						operand1.Obj = nil
						operand2.Obj = nil
					default:
						log.Panicln("unknown instruction", ins, m.IP)
					}
				default:
					log.Panicln("unknown instruction", ins, m.IP)
				}
			case Obj:
				switch operand2.Type {
				case Nil:
					switch ins.Type {
					case Cmp:
						result.Type = Bool
						result.Int = FALSE
						switch ins.CmpTyp {
						case NoEqual:
							result.Int = TRUE
						case Equal:
							result.Int = FALSE
						default:
							log.Panicln("unknown instruction", ins, m.IP)
						}
					default:
						log.Panicln("unknown instruction", ins, m.IP)
					}
				default:
					log.Panicln("unknown instruction", ins, m.IP)
				}
			case Bool:
				switch operand2.Type {
				case Bool:
					result.Type = Bool
					result.Int = FALSE
					switch ins.Type {
					case And:
						if operand1.Int == TRUE && operand2.Int == TRUE {
							result.Int = TRUE
						}
					default:
						log.Panicln("unknown instruction", ins, m.IP)
					}
				}
			case Nil:
				switch operand2.Type {
				case Nil:
					result.Type = Bool
					result.Int = FALSE
					switch ins.CmpTyp {
					case Equal:
						result.Int = TRUE
					case NoEqual:
					default:
						log.Panicln("unknown instruction",
							ins.String(m.symbolTable, m.builtInSymbolTable),
							m.IP,
							operand2.String())
					}
				}
			default:
				log.Panicln("unknown instruction", ins, m.IP)
			}
			//log.Println(operand1, operand2, result)
			m.SP++
			m.stack[m.SP] = result
		case Store:
			val := m.stack[m.SP]
			m.SP--
			m.stack[ins.Val] = val
		case IncStack:
			m.SP += ins.Val
		case StoreR:
			m.R[ins.Val] = m.stack[m.SP]
			m.stack[m.SP].Obj = nil
			m.SP--
		case LoadR:
			m.SP++
			m.stack[m.SP] = m.R[ins.Val]
		case Load:
			SP := ins.Val
			if ins.Val < 0 {
				SP = m.SP + ins.Val + 1 // -1 top
			}
			if SP > m.SP || SP < 0 {
				log.Panicln("stack error", SP, m.SP)
			}
			m.SP++
			m.stack[m.SP] = m.stack[SP]
		case Ret:
			IP := m.stack[m.SP]
			m.SP--
			m.IP = IP.Int
			continue
		case LoadO:
			obj := m.stack[m.SP]
			switch obj.Type {
			case String:
				index, ok := BuiltInFunctionsIndex["string."+ins.Str]
				if ok == false {
					log.Panicln("no find string." + ins.Str)
				}
				m.SP++
				m.stack[m.SP] = Object{
					Type: BFunc,
					Int:  index,
				}
			case Obj, Lambda:
				m.stack[m.SP] = *obj.loadObj(ins.Str)
			default:
				log.Panicln("unknown obj type", obj, m.SP)
			}
		case StoreO:
			obj := m.stack[m.SP]
			m.SP--
			ele := m.stack[m.SP]
			m.SP--
			switch obj.Type {
			case Obj, Lambda:
				obj.Store(ins.Str, ele)
				obj.Obj = nil
				ele.Obj = nil
			default:
				log.Panicln("unknown obj type", obj.Type)
			}
		case Call:
			count := m.R[0].Int
			objects := m.CallFunc(ins.Val, m.R[1:count+1]...)
			for i := range m.R[1 : count+1] {
				m.R[i+1].Obj = nil
			}
			m.R[0].Int = int64(len(objects))
			for index, obj := range objects {
				m.R[index+1] = obj
			}
		case CallO:
			f := &m.stack[m.SP]
			m.SP--
			switch f.Type {
			case BFunc:
				objects := m.CallFunc(f.Int, m.stack[m.SP])
				f.Obj = nil
				m.stack[m.SP].Obj = nil
				m.SP--
				m.SP-- //pop IP on the stack
				m.R[0].Int = int64(len(objects))
				for index, obj := range objects {
					m.R[index+1] = obj
				}
			case OFunc:
				m.IP = f.Int
			case Lambda:
				m.IP = f.Int
				m.closure = f.loadObj(closureLabel).Obj.(ObjectArray)
				continue
			default:
				log.Panicln("no function type", f.Type, m.IP)
			}

		case Label:
		case Exit:
			return
		case Jump:
			check := m.stack[m.SP]
			m.SP--
			if check.Type != Bool {
				log.Panicln("expect bool value for check", m.IP, m.SP)
			}
			if check.Int == TRUE {
				//log.Println("true",ins.Int)
				if ins.JumpTyp == DJump {
					m.IP = ins.Val
				} else {
					m.IP += ins.Val
				}
				continue
			}
		//			log.Println("false")

		case MakeArray:
			m.SP++
			m.stack[m.SP] = Object{
				Type: Array,
				Obj:  make(ObjectArray, 0, 8),
			}
		case Append:
			m.stack[m.SP-1].Obj = append(m.stack[m.SP-1].Obj.(ObjectArray), m.stack[m.SP])
			m.stack[m.SP].Obj = nil
			m.SP--
		case InitClosure:
			for _, obj := range m.closure {
				m.SP++
				m.stack[m.SP] = obj
			}
			m.closure = nil
		}
		m.IP++
		log.Println(m.stack[:m.SP+1])
	}
}

func (m *Machine) CallFunc(funcIndex int64, object ...Object) []Object {
	return BuiltInFunctions[funcIndex].Call(object...)
}

func (m *Machine) String() string {
	var buffer bytes.Buffer
	for _, it := range m.instructions {
		buffer.WriteString("\t" + it.String(m.symbolTable, m.builtInSymbolTable))
		buffer.WriteString("\n")
	}
	return buffer.String()
}
