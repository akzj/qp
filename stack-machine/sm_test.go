package stackmachine

import (
	"fmt"
	"log"
	"math"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
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

func TestIf(t *testing.T) {
	/*
		var a = 2
		if a > 1{
			print(a)
		}
	*/

	var m = New()
	index := m.symbolTable.addSymbol("a")
	m.instructions = []Instruction{
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     2,
		},
		{
			InstTyp: Store,
			Val:     index,
		},
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
			InstTyp: Cmp,
			CmpTyp:  Greater,
		},
		{
			InstTyp: Jump,
			Val:     3,
			JumpTyp: RJump,
		},
		{
			InstTyp: Push,
			ValTyp:  Bool,
			Val:     TRUE,
		},
		{
			InstTyp: Jump,
			Val:     4,
			JumpTyp: RJump,
		},
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
			InstTyp: Call,
			Val:     1,
		},
	}
	fmt.Println(m.String())
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

func TestFunc(t *testing.T) {
	/*
		func fib(val1,val2) {
			print(val1,val2)
		}

		fib(1,2)
	*/

	var m = New()
	val1 := m.symbolTable.addSymbol("val1")
	val2 := m.symbolTable.addSymbol("val2")
	m.instructions = []Instruction{

		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     5, //return IP
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     2,
		},
		{
			InstTyp: Push,
			ValTyp:  Bool,
			Val:     TRUE,
		},
		{
			InstTyp: Jump, // call function
			Val:     7,
		},

		//end of call

		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     math.MaxInt64, //to end
		},
		{
			InstTyp: Ret,
		},


		//function instruction
		{
			InstTyp: Store,
			Val:     val1, //store `val` to mem
		},
		{
			InstTyp: Store,
			Val:     val2, //store `val` to mem
		},
		{
			InstTyp: Load,
			Val:     val1,
		},
		{
			InstTyp: Load,
			Val:     val2,
		},
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     2,
		},
		{
			InstTyp: Call,
		},
		{
			InstTyp: Ret,
		},
	}
	m.Run()
}

func TestReturnVal(t *testing.T) {
	/*
		func fib() {
			return 1
		}
		var a = fib()
		println(a)
	*/

	var m = New()
	a1 := m.symbolTable.addSymbol("a")
	m.instructions = []Instruction{
		{
			InstTyp: Push,
			ValTyp:  Int,
			Val:     1,
		},
		{
			InstTyp: Store,
			Val: a1,
		},
		{
			InstTyp: Load,
			Val: a1,
		},


	}
	m.Run()

}
