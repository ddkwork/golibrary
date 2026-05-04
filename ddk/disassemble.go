package ddk

import (
	"fmt"
	"iter"
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

func DisassembleToString(code []byte, address uint64) string {
	instruction := Disassemble(code)
	var sb strings.Builder
	for i, inst := range instruction {
		syntax := x86asm.IntelSyntax(inst, address, nil)
		sb.WriteString(fmt.Sprintf("0x%X: %s\n", address+uint64(i), syntax))
	}
	return sb.String()
}

func DisassembleSeq2(code []byte, address uint64) iter.Seq2[string, x86asm.Inst] {
	return func(yield func(string, x86asm.Inst) bool) {
		for i := 0; i < len(code); {
			inst, err := x86asm.Decode(code[i:], 64)
			if err != nil {
				break
			}
			var sb strings.Builder
			syntax := x86asm.IntelSyntax(inst, address, nil)
			sb.WriteString(fmt.Sprintf("0x%X: %s\n", address+uint64(i), syntax))
			if !yield(sb.String(), inst) {
				break
			}
			if inst.Op == x86asm.RET {
				break
			}
			i += inst.Len
		}
	}
}
