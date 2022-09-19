package windef

import (
	"github.com/ddkwork/golibrary/src/stream"
	"github.com/ddkwork/golibrary/src/stream/tool"
)

func CTL_CODE(deviceType, function, method, access uint32) uint32 {
	return ((deviceType) << 16) | ((access) << 14) | ((function) << 2) | (method)
}

var (
	SMART_GET_VERSION    uint32 = CTL_CODE(IOCTL_DISK_BASE, 0x0020, METHOD_BUFFERED, FILE_READ_ACCESS)
	SMART_RCV_DRIVE_DATA uint32 = CTL_CODE(IOCTL_DISK_BASE, 0x0022, METHOD_BUFFERED, FILE_READ_ACCESS|FILE_WRITE_ACCESS)
)

func Creat(structBody string, isWriteType bool) (ok bool) {
	s := stream.New()
	if isWriteType {
		s.WriteStringLn(defBody)
	}
	s.WriteStringLn(structBody)
	s.WriteStringLn(`int main() { return 0; }`)
	if !tool.File().WriteTruncate("./windefTest/windef.c", s.Bytes()) {
		return
	}
	return tool.File().WriteTruncate("./windefTest/windef.cmd", "c2go -v windef windef.c")
}

var defBody = `
typedef unsigned long DWORD;
typedef int BOOL;
typedef unsigned char BYTE;
typedef unsigned short WORD;
typedef float FLOAT;
typedef FLOAT *PFLOAT;
typedef int INT;
typedef unsigned int UINT;
typedef unsigned int *PUINT;
typedef int BOOL;
typedef unsigned char BYTE;
typedef unsigned short WORD;
typedef float FLOAT;
typedef FLOAT *PFLOAT;
typedef BOOL *PBOOL;
typedef BOOL *LPBOOL;
typedef BYTE *PBYTE;
typedef BYTE *LPBYTE;
typedef int *PINT;
typedef int *LPINT;
typedef WORD *PWORD;
typedef WORD *LPWORD;
typedef long *LPLONG;
typedef DWORD *PDWORD;
typedef DWORD *LPDWORD;
typedef void *LPVOID;
typedef const void *LPCVOID;
typedef int INT;
typedef unsigned int UINT;
typedef unsigned int *PUINT;
typedef unsigned long ULONG;
typedef ULONG *PULONG;
typedef unsigned short USHORT;
typedef USHORT *PUSHORT;
typedef unsigned char UCHAR;
typedef UCHAR *PUCHAR;
typedef char CHAR;
typedef short SHORT;
typedef long LONG;
typedef int INT;
`
