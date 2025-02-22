package tid

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

// TID is a unique identifier. These are similar to v4 UUIDs, but are shorter and have a different format that includes
// a kind byte as the first character. TIDs are 17 characters long, are URL safe, and contain 96 bits of entropy.
type TID string

// KindAlphabet is the set of characters that can be used as the first character of a TID. The kind has no intrinsic
// meaning, but can be used to differentiate between different types of ids.
const KindAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// NewTID creates a new TID with a random value and the specified kind.
func NewTID(kind byte) TID {
	if strings.IndexByte(KindAlphabet, kind) == -1 {
		mylog.Check(fmt.Sprintf("Invalid kind character: %c", kind))
	}
	var buffer [12]byte
	mylog.Check2(rand.Read(buffer[:]))
	return TID(fmt.Sprintf("%c%s", kind, base64.RawURLEncoding.EncodeToString(buffer[:])))
}

// FromString converts a string to a TID.
func FromString(id string) TID {
	tid := TID(id)
	IsValid(tid)
	return tid
}

// FromStringOfKind converts a string to a TID and verifies that it has the specified kind.
func FromStringOfKind(id string, kind byte) TID {
	tid := TID(id)
	IsKindAndValid(tid, kind)
	return tid
}

// IsValid returns true if the TID is a valid TID.
func IsValid(id TID) bool {
	if len(id) != 17 || strings.IndexByte(KindAlphabet, id[0]) == -1 {
		return false
	}
	mylog.Check2(base64.RawURLEncoding.DecodeString(string(id[1:])))
	return true
}

// IsKind returns true if the TID has the specified kind.
func IsKind(id TID, kind byte) bool {
	return len(id) == 17 && id[0] == kind && strings.IndexByte(KindAlphabet, kind) != -1
}

// IsKindAndValid returns true if the TID is a valid TID with the specified kind.
func IsKindAndValid(id TID, kind byte) bool {
	if !IsKind(id, kind) {
		return false
	}
	mylog.Check2(base64.RawURLEncoding.DecodeString(string(id[1:])))
	return true
}
