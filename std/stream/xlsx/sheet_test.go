package xlsx

import (
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
)

func TestLoad(t *testing.T) {
	t.Skip()
	mylog.Check2(Load("D:\\clone\\HyperDbg\\hyperdbg\\demo\\宁夏装车计件表.xlsx"))
}
