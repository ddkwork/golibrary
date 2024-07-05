package mylog_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/mylog"
)

func TestCheckM1(t *testing.T) { mylog.Call(m1) }
func TestCheckM2(t *testing.T) { mylog.Call(m2) }
func TestCheckM3(t *testing.T) { mylog.Call(m3) }
func TestCheckM4(t *testing.T) { mylog.Call(m4) }
func TestCheckM5(t *testing.T) { mylog.Call(m5) }
func TestCheckM7(t *testing.T) {
	mylog.Call(func() { m7(nil) })
}

func TestCheckM8(t *testing.T) {
	mylog.Call(func() { m8("") })
}

func m1() {
	atoi := mylog.Check2(strconv.Atoi("123"))
	print(atoi)
	mylog.Check("custom error message")
	println(doNotHers)
}

func m7(file *os.File) {
	mylog.CheckNil(file)
	println(doNotHers)
}

func m8(s string) {
	mylog.Check(len(s))
	println(doNotHers)
}

func m2() {
	atoi := mylog.Check2(strconv.Atoi("123"))
	print(atoi)
	println(mustHers)
}

func m3() {
	m3_()
	println(done)
}

func m3_() string {
	atoi := mylog.Check2(strconv.Atoi("123"))
	print(atoi)
	mylog.Check("custom error message")
	return doNotHers
}

func m4() {
	handleFile()
	println(done)
	mylog.Check(os.Remove("test.txt"))
}

func handleFile() {
	f := mylog.Check2(os.Create("test.txt"))
	defer func() {
		mylog.Check(f.Close())
		mylog.Check(os.Remove("test.txt"))
	}()
	mylog.Check2(f.Write(nil))
	println(doNotHers)
}

func m5() {
	bug()
	println(done)
}

func bug() {
	dir := mylog.Check2(os.ReadDir("2332"))
	dir[1] = nil
	println(mustHers)
	for _, entry := range dir {
		println(entry.Name())
	}
	a := make([]byte, 0)
	a[2] = 0
	println(doNotHers)
}

const (
	mustHers  = "must be hers"
	doNotHers = "do not be hers"
	done      = "done"
)

func TestCheck(t *testing.T) {
	mylog.Call(func() {
		r1 := mylog.Check2(os.Create(""))
		mylog.Check(r1.Close())
	})
}
