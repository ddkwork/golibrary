package hardwareIndo

import (
	os "os"
	unsafe "unsafe"
)

type DWORD = uint32
type BOOL = int32
type BYTE = uint8
type WORD = uint16
type FLOAT = float32
type PFLOAT = *float32
type INT = int32
type UINT = uint32
type PUINT = *uint32
type PBOOL = *int32
type LPBOOL = *int32
type PBYTE = *uint8
type LPBYTE = *uint8
type PINT = *int32
type LPINT = *int32
type PWORD = *uint16
type LPWORD = *uint16
type LPLONG = *int32
type PDWORD = *uint32
type LPDWORD = *uint32
type LPVOID = unsafe.Pointer
type LPCVOID = unsafe.Pointer
type ULONG = uint32
type PULONG = *uint32
type USHORT = uint16
type PUSHORT = *uint16
type UCHAR = uint8
type PUCHAR = *uint8
type CHAR = int8
type SHORT = int16
type LONG = int32
type _cgoa_1_windef struct {
	Xbf_0 uint16
}
type _cgoa_2_windef struct {
	Xbf_0 uint16
}
type _cgoa_3_windef struct {
	Xbf_0 uint16
}
type _cgoa_4_windef struct {
	Xbf_0 uint16
}
type _cgoa_5_windef struct {
	Xbf_0 uint16
}
type _cgoa_6_windef struct {
	Xbf_0 uint16
}
type _cgoa_7_windef struct {
	Xbf_0 uint16
}
type struct__IDINFO struct {
	GenConfig             uint16
	NumCyls               uint16
	Reserved2             uint16
	NumHeads              uint16
	Reserved4             uint16
	Reserved5             uint16
	NumSectorsPerTrack    uint16
	VendorUnique          [3]uint16
	SerialNumber          [20]uint8
	BufferType            uint16
	BufferSize            uint16
	ECCSize               uint16
	FirmwareRev           [8]uint8
	ModelNumber           [40]uint8
	MoreVendorUnique      uint16
	Reserved48            uint16
	Capabilities          _cgoa_1_windef
	Reserved1             uint16
	PIOTiming             uint16
	DMATiming             uint16
	FieldValidity         _cgoa_2_windef
	NumCurCyls            uint16
	NumCurHeads           uint16
	NumCurSectorsPerTrack uint16
	CurSectorsLow         uint16
	CurSectorsHigh        uint16
	MultSectorStuff       _cgoa_3_windef
	dwTotalSectors        uint32
	SingleWordDMA         uint16
	MultiWordDMA          _cgoa_4_windef
	PIOCapacity           _cgoa_5_windef
	MinMultiWordDMACycle  uint16
	RecMultiWordDMACycle  uint16
	MinPIONoFlowCycle     uint16
	MinPOIFlowCycle       uint16
	Reserved69            [11]uint16
	MajorVersion          _cgoa_6_windef
	MinorVersion          uint16
	Reserved82            [6]uint16
	UltraDMA              _cgoa_7_windef
	Reserved89            [167]uint16
}
type IDINFO = struct__IDINFO
type PIDINFO = *struct__IDINFO
type struct__DRIVERSTATUS struct {
	DriverError uint8
	IDEError    uint8
	Reserved    [2]uint8
	dwReserved  [2]uint32
}
type DRIVERSTATUS = struct__DRIVERSTATUS
type PDRIVERSTATUS = *struct__DRIVERSTATUS
type LPDRIVERSTATUS = *struct__DRIVERSTATUS
type struct__SENDCMDOUTPARAMS struct {
	BufferSize   uint32
	DriverStatus struct__DRIVERSTATUS
	Buffer       [1]uint8
}
type SENDCMDOUTPARAMS = struct__SENDCMDOUTPARAMS
type PSENDCMDOUTPARAMS = *struct__SENDCMDOUTPARAMS
type LPSENDCMDOUTPARAMS = *struct__SENDCMDOUTPARAMS
type struct__GETVERSIONINPARAMS struct {
	Version       uint8
	Revision      uint8
	Reserved1     uint8
	IDEDeviceMap  uint8
	fCapabilities uint32
	Reserved2     [4]uint32
}
type GETVERSIONINPARAMS = struct__GETVERSIONINPARAMS
type PGETVERSIONINPARAMS = *struct__GETVERSIONINPARAMS
type LPGETVERSIONINPARAMS = *struct__GETVERSIONINPARAMS
type struct__IDEREGS struct {
	FeaturesReg     uint8
	SectorCountReg  uint8
	SectorNumberReg uint8
	CylLowReg       uint8
	CylHighReg      uint8
	DriveHeadReg    uint8
	CommandReg      uint8
	Reserved        uint8
}
type IDEREGS = struct__IDEREGS
type PIDEREGS = *struct__IDEREGS
type LPIDEREGS = *struct__IDEREGS
type struct__SENDCMDINPARAMS struct {
	BufferSize  uint32
	irDriveRegs struct__IDEREGS
	DriveNumber uint8
	Reserved1   [3]uint8
	Reserved2   [4]uint32
	Buffer      [1]uint8
}
type SENDCMDINPARAMS = struct__SENDCMDINPARAMS
type PSENDCMDINPARAMS = *struct__SENDCMDINPARAMS
type LPSENDCMDINPARAMS = *struct__SENDCMDINPARAMS

func _cgo_main() int32 {
	return int32(0)
}
func main() {
	os.Exit(int(_cgo_main()))
}
