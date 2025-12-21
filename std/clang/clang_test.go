package clang

import (
	"testing"
)

func TestClangFormat(t *testing.T) {
}

func TestFixFunctionSignatures(t *testing.T) {
	println(fixFunctionSignatures(code))
}

var code = `

PUSERMODE_DEBUGGING_PROCESS_DETAILS
ThreadHolderGetProcessDebuggingDetailsByThreadId(UINT32 ThreadId) {
  PLIST_ENTRY TempList = 0;
  PLIST_ENTRY TempList2 = 0;
  TempList = &g_ProcessDebuggingDetailsListHead;
  while (&g_ProcessDebuggingDetailsListHead != TempList->Flink) {
    TempList = TempList->Flink;
    PUSERMODE_DEBUGGING_PROCESS_DETAILS ProcessDebuggingDetails =
        CONTAINING_RECORD(TempList, USERMODE_DEBUGGING_PROCESS_DETAILS,
                          AttachedProcessList);
    TempList2 = &ProcessDebuggingDetails->ThreadsListHead;
    while (&ProcessDebuggingDetails->ThreadsListHead != TempList2->Flink) {
      TempList2 = TempList2->Flink;
      PUSERMODE_DEBUGGING_THREAD_HOLDER ThreadHolder = CONTAINING_RECORD(
          TempList2, USERMODE_DEBUGGING_THREAD_HOLDER, ThreadHolderList);
      for (size_t i = 0; i < MAX_THREADS_IN_A_PROCESS_HOLDER; i++) {
        if (ThreadHolder->Threads[i].ThreadId == ThreadId) {
          return ProcessDebuggingDetails;
        }
      }
    }
  }
  return NULL;
}
`
