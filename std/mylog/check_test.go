package mylog_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
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
	fmt.Println(atoi)
	mylog.Check("custom error value")
	fmt.Println(doNotHers)
}

func m7(file *os.File) {
	mylog.CheckNil(file)
	fmt.Println(doNotHers)
}

func m8(s string) {
	mylog.Check(len(s))
	fmt.Println(doNotHers)
}

func m2() {
	atoi := mylog.Check2(strconv.Atoi("123"))
	fmt.Println(atoi)
	fmt.Println(mustHers)
}

func m3() {
	m3_()
	fmt.Println(done)
}

func m3_() string {
	atoi := mylog.Check2(strconv.Atoi("123"))
	fmt.Println(atoi)
	mylog.Check("custom error value")
	return doNotHers
}

func m4() {
	handleFile()
	fmt.Println(done)
	mylog.Check(os.Remove("test.txt"))
}

func handleFile() {
	f := mylog.Check2(os.Create("test.txt"))
	defer func() {
		mylog.Check(f.Close())
		mylog.Check(os.Remove("test.txt"))
	}()
	mylog.Check2(f.Write(nil))
	fmt.Println(doNotHers)
}

func m5() {
	bug()
	fmt.Println(done)
}

func bug() {
	dir := mylog.Check2(os.ReadDir("2332"))
	dir[1] = nil
	fmt.Println(mustHers)
	for _, entry := range dir {
		fmt.Println(entry.Name())
	}
	a := make([]byte, 0)
	a[2] = 0
	fmt.Println(doNotHers)
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
