package stackmachine

import "log"

type Linker struct {
	functionPoint      map[string]int64
	functions          map[string]FuncInstruction
	ins                []Instruction
	toLink             []toLink
	symbolTable        *SymbolTable
	builtInSymbolTable *SymbolTable
}

func NewLinker(functions map[string]FuncInstruction,
	ins []Instruction,
	toLink []toLink,
	symbolTable *SymbolTable,
	builtInSymbolTable *SymbolTable) Linker {
	return Linker{
		functionPoint:      map[string]int64{},
		functions:          functions,
		ins:                ins,
		toLink:             toLink,
		symbolTable:        symbolTable,
		builtInSymbolTable: builtInSymbolTable,
	}
}

func (linker *Linker) link() []Instruction {
	for _, link := range linker.toLink {
		ins := linker.ins[link.IP]
		if ins.InstTyp != Jump && ins.ValTyp != OFunc {
			log.Println("toLink error", ins.String(linker.symbolTable, linker.builtInSymbolTable), link.IP)
		}
		ins.Val = linker.getFunctionPoint(link.label)
		linker.ins[link.IP] = ins
		log.Println(ins)
	}
	return linker.ins
}

func (linker *Linker) getFunctionPoint(label string) int64 {
	if point, ok := linker.functionPoint[label]; ok {
		return point
	}
	function, ok := linker.functions[label]
	if ok == false {
		log.Panicf("no find function `%s`\n", label)
	}
	index := int64(len(linker.ins))
	linker.functionPoint[label] = index
	linker.ins = append(linker.ins, function.ins...)
	for _, link := range function.toLinks {
		linker.ins[link.IP+index].Val = linker.getFunctionPoint(link.label)
	}
	return index
}
