package gbk_test

import (
	"github.com/ddkwork/golibrary/gbk"
	"testing"
)

func TestName(t *testing.T) {
	gbk.Gbk2Utf8All("cmake")
}
