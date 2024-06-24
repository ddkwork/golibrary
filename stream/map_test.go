package stream_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/maps"
)

func TestMap(t *testing.T) {
	m := new(maps.SafeMap[string, bool])
	id := uuid.NewString()
	m.Set(id, true)
	mylog.Info("Count", m.Len())
	mylog.Struct(m.Keys())
	mylog.Struct(m.Keys())
	mylog.Info("IsEmpty", m.String())
	mylog.Info("Get", m.Get(id))
	m.Delete(id)
	mylog.Info("Count", m.Len())
}
