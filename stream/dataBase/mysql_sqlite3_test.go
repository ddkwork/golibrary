package dataBase

import (
	"testing"
)

type (
	EA interface{}
	ea struct {
		maddenMySqldb    Interface
		maddenSqlLite3db Interface
		value            string
	}
)

func TestName(t *testing.T) {
}
