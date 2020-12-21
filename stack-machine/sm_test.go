package stackmachine

import (
	"fmt"
	"testing"
)

type InstType byte
type ValType byte

const (
	Push InstType = 0 + iota
	Pop
	Add
	Sub
	Call
	Jump
)
const (
	Int ValType = 0 + iota
)

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
	stack        []ValObject
	instructions []Instruction
	pIns         int
}

func (m *Machine) Run() {
	for m.pIns < len(m.instructions) {
		ins := m.instructions[m.pIns]
		switch ins.InstTyp {
		case Push:
			switch ins.VaType {
			case Int:
				m.stack = append(m.stack, ValObject{
					VType: Int,
					Val:   ins.Val,
				})
			}
		case Pop:
			m.stack = m.stack[:len(m.stack)-1]
		case Add, Sub:
			operand1 := m.stack[len(m.stack)-1]
			operand2 := m.stack[len(m.stack)-2]
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
			m.stack = m.stack[:len(m.stack)-1]
			m.stack = m.stack[:len(m.stack)-1]
			m.stack = append(m.stack, result)
		case Call:
			m.CallFunc(ins.Val)
			m.stack = m.stack[:len(m.stack)-1]
		}
		m.pIns++
		fmt.Println(m.stack)
	}
}

func (m *Machine) CallFunc(val int64) {
	if val == 0 {
		val := m.stack[0]
		fmt.Println(val.Val)
	}
}

func TestName(t *testing.T) {
	var m = Machine{}
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
			InstTyp: Call,
			Val:     0,
		},
	}
	m.Run()
}
