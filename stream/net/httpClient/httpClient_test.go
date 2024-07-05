package httpClient_test

import (
	_ "embed"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/net/httpClient"
)

func TestName(t *testing.T) {
	mylog.Call(func() {
		do()
		println("done")
	})
}

func do() {
	httpClient.New().Get().Url("").Request()
}
