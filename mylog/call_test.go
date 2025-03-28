package mylog_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func Test_log_printAndWrite2(t *testing.T) {
	mylog.Call(func() {
		mylog.Check(errors.New("this is a err value"))
	})
}

func TestName(t *testing.T) {
	mylog.Call(func() {
		atoi := mylog.Check2(strconv.Atoi("123"))
		fmt.Println(atoi)
		mylog.Check("custom error value")
		fmt.Println(doNotHers)
	})
}
