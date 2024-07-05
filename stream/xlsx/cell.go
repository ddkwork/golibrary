package xlsx

import (
	"strconv"
	"time"

	"github.com/ddkwork/golibrary/mylog"
)

const (
	String CellType = iota
	Number
	Boolean
)

type CellType int

type Cell struct {
	Type  CellType
	Value string
}

func (c *Cell) String() string {
	return c.Value
}

func (c *Cell) Integer() int {
	v := mylog.Check2(strconv.Atoi(c.Value))

	return v
}

func (c *Cell) Float() float64 {
	f := mylog.Check2(strconv.ParseFloat(c.Value, 64))

	return f
}

func (c *Cell) Boolean() bool {
	return c.Value != "0"
}

func (c *Cell) Time() time.Time {
	return timeFromExcelTime(c.Float())
}
