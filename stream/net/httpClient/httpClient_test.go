package httpClient_test

import (
	_ "embed"
	"testing"

	"github.com/ddkwork/golibrary/stream/net/httpClient"
)

func TestPb2(t *testing.T) {
	return
	httpClient.MockProtoBufPacket("8080")
}
