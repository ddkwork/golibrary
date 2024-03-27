package safeType

import (
	"encoding/hex"
	"fmt"

	"github.com/ddkwork/golibrary/mylog"
)

func NewHexString(s HexString) *Data {
	decodeString, err := hex.DecodeString(string(s))
	if !mylog.Error(err) {
		return New(err.Error())
	}
	return New(decodeString)
}

type HexString string

func (d *Data) HexString() HexString      { return HexString(hex.EncodeToString(d.Bytes())) }
func (d *Data) HexStringUpper() HexString { return HexString(fmt.Sprintf("%#X", d.Bytes())[2:]) }
