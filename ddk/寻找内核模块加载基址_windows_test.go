package ddk

import (
	"fmt"
	"testing"

	"github.com/ddkwork/golibrary/std/fakeError"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
	"golang.org/x/arch/x86/x86asm"
)

func TestQueryKernelModules(t *testing.T) {
	f := NewKernelModuleFinder()
	modules := mylog.Check2(f.Modules())

	if len(modules) == 0 {
		t.Fatal("no kernel modules found")
	}
	t.Logf("total %d kernel modules loaded", len(modules))
	for _, mod := range modules {
		t.Logf("  %-30s Base=%-18s Size=0x%08X  Path=%s",
			mod.Name,
			mylog.Hex(mod.ImageBase),
			mod.ImageSize,
			mod.FullPath,
		)
	}
}

func TestName(t *testing.T) {
	fakeError.Walk(".", true)
	stream.FmtDir(".")
	stream.FixDir(".")
}

func TestFindCiDllBase(t *testing.T) {
	f := NewKernelModuleFinder()
	mod := f.ModuleByName("ci.dll")
	mylog.Struct(mod)
	if mod.ImageBase == 0 {
		t.Fatal("ci.dll base address should not be 0")
	}
}

func TestFindNtoskrnlBase(t *testing.T) {
	f := NewKernelModuleFinder()
	base := f.ModuleBaseByName("ntoskrnl.exe")

	t.Logf("ntoskrnl.exe base address: %s", mylog.Hex(base))
	if base == 0 {
		t.Fatal("ntoskrnl.exe base address should not be 0")
	}
}

func TestFindKernelModuleNotFound(t *testing.T) {
	f := NewKernelModuleFinder()
	f.ModuleBaseByName("nonexistent_driver.sys")
}

func TestFindExportedSymbolRVA_Ntoskrnl(t *testing.T) {
	f := NewKernelModuleFinder()
	rva := mylog.Check2(f.FindExportedSymbolRVA("ntoskrnl.exe", "NtDeviceIoControlFile"))

	t.Logf("NtDeviceIoControlFile RVA: 0x%06X", rva)
	if rva == 0 {
		t.Fatal("RVA should not be 0")
	}
}

func TestFindExportedSymbolRVA_CiDll(t *testing.T) {
	f := NewKernelModuleFinder()
	rva := mylog.Check2(f.FindExportedSymbolRVA("ci.dll", "CiInitialize"))

	t.Logf("CiInitialize RVA: 0x%06X", rva)
	if rva == 0 {
		t.Fatal("RVA should not be 0")
	}
}

func TestFindExportedSymbolAddress_Ntoskrnl(t *testing.T) {
	f := NewKernelModuleFinder()
	addr := mylog.Check2(f.FindExportedSymbolAddress("ntoskrnl.exe", "ZwDeviceIoControlFile"))

	t.Logf("ZwDeviceIoControlFile address: %s", mylog.Hex(addr))
	if addr == 0 {
		t.Fatal("address should not be 0")
	}
}

func TestFindExportedSymbolAddress_MmFreeNonCachedMemory(t *testing.T) {
	f := NewKernelModuleFinder()
	addr := mylog.Check2(f.FindExportedSymbolAddress("ntoskrnl.exe", "MmFreeNonCachedMemory"))

	t.Logf("MmFreeNonCachedMemory address: %s", mylog.Hex(addr))
}

func TestFindExportedSymbolNotFound(t *testing.T) {
	f := NewKernelModuleFinder()
	mylog.Check2(f.FindExportedSymbolRVA("ntoskrnl.exe", "NonExistentFunction12345"))
}

func TestFindNonExportedSymbol_IopXxxControlFile(t *testing.T) {
	f := NewKernelModuleFinder()
	tracer := func(instructions []x86asm.Inst, baseRVA uint32, data []byte) (uint32, bool) {
		for i, inst := range instructions {
			if inst.Op != x86asm.CALL {
				continue
			}
			next := instructions[i+1]
			next2 := instructions[i+2]
			if next.Op == x86asm.ADD && next2.Op == x86asm.RET {
				for _, arg := range inst.Args {
					if rel, ok := arg.(x86asm.Rel); ok {
						offset := 0
						for j := 0; j <= i; j++ {
							offset += instructions[j].Len
						}
						targetRVA := baseRVA + uint32(offset) + uint32(rel)
						return targetRVA, true
					}
				}
			}
		}
		return 0, false
	}

	addr := mylog.Check2(f.FindNonExportedSymbolAddress("ntoskrnl.exe", "NtDeviceIoControlFile", tracer))

	t.Logf("IopXxxControlFile address: %s", mylog.Hex(addr))
}

func TestFindNonExportedSymbol_MiGetPteAddress(t *testing.T) {
	f := NewKernelModuleFinder()
	tracer := func(instructions []x86asm.Inst, baseRVA uint32, data []byte) (uint32, bool) {
		for i, inst := range instructions {
			if inst.Op != x86asm.SUB {
				continue
			}
			next := instructions[i+1]
			next2 := instructions[i+2]
			if next.Op == x86asm.MOV && next2.Op == x86asm.CALL {
				for _, arg := range next2.Args {
					if rel, ok := arg.(x86asm.Rel); ok {
						offset := 0
						for j := 0; j <= i; j++ {
							offset += instructions[j].Len
						}
						offset += next.Len
						offset += next2.Len
						targetRVA := baseRVA + uint32(offset) + uint32(rel)
						return targetRVA, true
					}
				}
			}
		}
		return 0, false
	}

	addr := mylog.Check2(f.FindNonExportedSymbolAddress("ntoskrnl.exe", "MmFreeNonCachedMemory", tracer))

	t.Logf("MiGetPteAddress address: %s", mylog.Hex(addr))
}

func TestFindNonExportedSymbol_KeServiceDescriptorTable(t *testing.T) {
	f := NewKernelModuleFinder()
	tracer := func(instructions []x86asm.Inst, baseRVA uint32, data []byte) (uint32, bool) {
		for i, inst := range instructions {
			if inst.Op != x86asm.MOV || i+1 >= len(instructions) {
				continue
			}
			next := instructions[i+1]
			if inst.Op != x86asm.MOV || next.Op != x86asm.JMP {
				continue
			}
			var sysCallNum uint32
			for _, arg := range inst.Args {
				if imm, ok := arg.(x86asm.Imm); ok {
					sysCallNum = uint32(imm)
					break
				}
			}
			if sysCallNum != 0x7 {
				continue
			}
			for _, arg := range next.Args {
				if rel, ok := arg.(x86asm.Rel); ok {
					offset := 0
					for j := 0; j <= i; j++ {
						offset += instructions[j].Len
					}
					offset += next.Len
					targetRVA := baseRVA + uint32(offset) + uint32(rel)
					return targetRVA, true
				}
			}
		}
		return 0, false
	}

	addr := mylog.Check2(f.FindNonExportedSymbolAddress("ntoskrnl.exe", "ZwDeviceIoControlFile", tracer))

	t.Logf("KiServiceInternal address: %s", mylog.Hex(addr))
}

func TestAllKernelModuleBases(t *testing.T) {
	targetModules := []string{
		"ntoskrnl.exe",
		"ci.dll",
		"hal.dll",
		"ndis.sys",
		"ntfs.sys",
		"fltmgr.sys",
		"ksecdd.sys",
		"clipsp.sys",
		"cng.sys",
		"pcw.sys",
	}
	f := NewKernelModuleFinder()
	modules := mylog.Check2(f.Modules())

	found := map[string]uint64{}
	for _, mod := range modules {
		for _, target := range targetModules {
			if containsIgnoreCase(mod.Name, target) {
				found[target] = mod.ImageBase
			}
		}
	}
	for _, target := range targetModules {
		base, ok := found[target]
		if ok {
			t.Logf("  %-20s Base=%s", target, mylog.Hex(base))
		} else {
			t.Logf("  %-20s NOT FOUND", target)
		}
	}
}

func containsIgnoreCase(s, substr string) bool {
	s = fmt.Sprintf("%s", s)
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1 := s[i+j]
			c2 := substr[j]
			if c1 >= 'A' && c1 <= 'Z' {
				c1 += 32
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 += 32
			}
			if c1 != c2 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
