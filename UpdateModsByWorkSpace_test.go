package golibrary

import (
	"testing"

	"go.uber.org/goleak"
)

func TestUpdateModsByWorkSpace(t *testing.T) {
	defer goleak.VerifyNone(t)
	UpdateModsByWorkSpace(false, false,
		"github.com/ddkwork/golibrary@66d1453f648378b7a521cb04d7db47bbf7521e17",
		"cogentcore.org/core@da0f626c53da619895d89587b1b319cd647f665d",
	)
}
