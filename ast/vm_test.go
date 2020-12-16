package ast

import (
	"fmt"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}
func TestVMMemoryStackTest(t *testing.T) {
	vm := New()

	t.Run("AllocObject", func(t *testing.T) {
		objectA := vm.AllocObject("a")
		objectA1 := vm.GetObject("a")
		if objectA != objectA1 {
			t.Fatal("AllocObject error")
		}
	})

	t.Run("StackFrame", func(t *testing.T) {
		objectB := vm.AllocObject("b")

		//
		vm.PushStackFrame(false)
		objectB1 := vm.AllocObject("b")
		if object := vm.GetObject("b"); object != objectB1 {
			if object == nil {
				t.Fatal("no find objects")
			}
			fmt.Println(object)
			fmt.Println(objectB1)
			t.Fatalf("PushStackFrame failed %s %s", object.String(), objectB1.String())
		}

		//
		vm.PopStackFrame()
		if vm.GetObject("b") == objectB1 {
			t.Fatal("PopStackFrame failed")
		}
		if vm.GetObject("b") != objectB {
			t.Fatal("PopStackFrame failed")
		}
	})

	t.Run("stackFrameIsolate", func(t *testing.T) {
		a := vm.AllocObject("a")

		//isolate
		vm.PushStackFrame(true)

		if vm.GetTypeObject("a") != nil {
			t.Fatal("isolate failed")
		}

		vm.PopStackFrame()
		if vm.GetObject("a") != a {
			t.Fatal("isolate failed")
		}
	})
}
