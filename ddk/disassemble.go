package ddk

import (
	"fmt"
	"strings"

	"golang.org/x/arch/x86/x86asm"
)

func disassemble(code []byte, stopOnRet bool) []x86asm.Inst {
	var instructions []x86asm.Inst
	for i := 0; i < len(code); {
		inst, err := x86asm.Decode(code[i:], 64)
		if err != nil {
			break
		}
		instructions = append(instructions, inst)
		if stopOnRet && inst.Op == x86asm.RET {
			break
		}
		i += inst.Len
	}
	return instructions
}

func Disassemble(code []byte) []x86asm.Inst {
	return disassemble(code, true)
}

func DisassembleAll(code []byte) []x86asm.Inst {
	return disassemble(code, false)
}

func DisassembleToString(code []byte, addr uint64) string {
	instruction := Disassemble(code)
	var sb strings.Builder
	for i, inst := range instruction {
		syntax := x86asm.IntelSyntax(inst, addr, nil)
		sb.WriteString(fmt.Sprintf("0x%X: %s\n", addr+uint64(i), syntax))
	}
	return sb.String()
}
