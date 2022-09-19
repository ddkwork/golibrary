package dataBase

import "testing"

type (
	EA interface {
	}
	ea struct {
		maddenMySqldb    Interface
		maddenSqlLite3db Interface
		value            string
		err              error
	}
)

func _TestInterfaceDataBase(t *testing.T) {

}
