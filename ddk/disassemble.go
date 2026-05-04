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

func DisassembleSeq2(code []byte, baseAddr int, sym x86asm.SymLookup) iter.Seq2[string, x86asm.Inst] {
	return func(yield func(string, x86asm.Inst) bool) {
		mode := 32
		if uint64(baseAddr) > 0xFFFFFFFF {
			mode = 64
		}

		off := 0
		for off < len(code) {
			rip := uint64(baseAddr + off)
			inst, err := x86asm.Decode(code[off:], mode)

			var opHex, asm string
			if err != nil || inst.Len <= 0 {
				b := code[off]
				opHex = fmt.Sprintf("%02X ", b)
				asm = fmt.Sprintf(".byte %02X", b)
				off++
			} else {
				// x86原生 Sym 符号格式化函数
				for _, b := range code[off : off+inst.Len] {
					opHex += fmt.Sprintf("%02X ", b)
				}
				asm = x86asm.IntelSyntax(inst, rip, sym)
				off += inst.Len
			}

			// x86最长15字节 opcode 固定对齐，Maximum instruction size is 15 bytes.
			opAlign := fmt.Sprintf("%-45s", opHex)

			var addrFmt string
			if uint64(baseAddr) <= 0xFFFFFFFF {
				addrFmt = "%08X"
			} else {
				addrFmt = "%016X"
			}

			line := fmt.Sprintf(addrFmt+": %s%s", rip, opAlign, asm)
			if !yield(line, inst) {
				return
			}
			if inst.Op == x86asm.RET {
				return
			}
		}
	}
}

