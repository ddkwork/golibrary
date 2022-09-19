package caseconv

import (
	"testing"
)

func TestName(t *testing.T) {
	for _, s := range name {
		println(ToCamelUpper(s, true))
	}
}

var name = []string{

	"	HIDDEN_HOOK_READ_AND_WRITE                                                 ",
	"	HIDDEN_HOOK_READ                                                           ",
	"	HIDDEN_HOOK_WRITE                                                          ",
	"	HIDDEN_HOOK_EXEC_DETOURS                                                   ",
	"	HIDDEN_HOOK_EXEC_CC                                                        ",
	"	SYSCALL_HOOK_EFER_SYSCALL                                                  ",
	"	SYSCALL_HOOK_EFER_SYSRET                                                   ",
	"	CPUID_INSTRUCTION_EXECUTION                                                ",
	"	RDMSR_INSTRUCTION_EXECUTION                                                ",
	"	WRMSR_INSTRUCTION_EXECUTION                                                ",
	"	IN_INSTRUCTION_EXECUTION                                                   ",
	"	OUT_INSTRUCTION_EXECUTION                                                  ",
	"	EXCEPTION_OCCURRED                                                         ",
	"	EXTERNAL_INTERRUPT_OCCURRED                                                ",
	"	DEBUG_REGISTERS_ACCESSED                                                   ",
	"	TSC_INSTRUCTION_EXECUTION                                                  ",
	"	PMC_INSTRUCTION_EXECUTION                                                  ",
	"	VMCALL_INSTRUCTION_EXECUTION                                               ",
	"	CONTROL_REGISTER_MODIFIED                                                  ",
	"	DEBUGGER_EVENT_TYPE_ENUM                                                   ",
	"	BREAK_TO_DEBUGGER                                                          ",
	"	RUN_SCRIPT                                                                 ",
	"	RUN_CUSTOM_CODE                                                            ",
	"	DEBUGGER_EVENT_SYSCALL_SYSRET_SAFE_ACCESS_MEMORY                           ",
	"	DEBUGGER_EVENT_SYSCALL_SYSRET_HANDLE_ALL_UD                                ",
	"	DEBUGGER_MODIFY_EVENTS_QUERY_STATE                                         ",
	"	DEBUGGER_MODIFY_EVENTS_ENABLE                                              ",
	"	DEBUGGER_MODIFY_EVENTS_DISABLE                                             ",
	"	DEBUGGER_MODIFY_EVENTS_CLEAR                                               ",
	"	DEBUGGER_MODIFY_EVENTS_TYPE                                                ",
	"	struct__DEBUGGER_MODIFY_EVENTS                                             ",
	"	VirtualAddress                                                             ",
	"	ProcessId                                                                  ",
	"	Pml4eVirtualAddress                                                        ",
	"	Pml4eValue                                                                 ",
	"	PdpteVirtualAddress                                                        ",
	"	PdpteValue                                                                 ",
	"	PdeVirtualAddress                                                          ",
	"	PdeValue                                                                   ",
	"	PteVirtualAddress                                                          ",
	"	PteValue                                                                   ",
	"	KernelStatus                                                               ",
	"	DEBUGGER_READ_PAGE_TABLE_ENTRIES_DETAILS                                   ",
	"	PDEBUGGER_READ_PAGE_TABLE_ENTRIES_DETAILS                                  ",
	"	struct__DEBUGGER_VA2PA_AND_PA2VA_COMMANDS                                  ",
	"	DEBUGGER_VA2PA_AND_PA2VA_COMMANDS                                          ",
	"	PDEBUGGER_VA2PA_AND_PA2VA_COMMANDS                                         ",
	"	struct__DEBUGGER_DT_COMMAND_OPTIONS                                        ",
	"	ypeName                                                                    ",
	"	DEBUGGER_SHOW_COMMAND_DT                                                   ",
	"	DEBUGGER_SHOW_COMMAND_DISASSEMBLE64                                        ",
	"	DEBUGGER_SHOW_COMMAND_DISASSEMBLE32                                        ",
	"	DEBUGGER_SHOW_COMMAND_DB                                                   ",
	"	DEBUGGER_SHOW_COMMAND_DC                                                   ",
	"	DEBUGGER_SHOW_COMMAND_DQ                                                   ",
	"	DEBUGGER_SHOW_COMMAND_DD                                                   ",
}
