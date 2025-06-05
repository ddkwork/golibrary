package mylog_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func TestBUg(t *testing.T) {
	s := `clang: error: no such file or directory: '/C'
clang: error: no such file or directory: 'clang'

exit status 1`
	mylog.Call(func() {
		mylog.Check(s)
	})
}

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
