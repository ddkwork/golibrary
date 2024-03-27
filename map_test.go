package golibrary

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/ddkwork/golibrary/stream/cmd"
	"github.com/smallnest/safemap"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/google/uuid"
)

// go get  cogentcore.org/core@35f09866eefbd6adceacea7e75707d4882e23460
func TestName(t *testing.T) {
	cmd.CheckLoopvarAndNilPoint()
	c := safemap.New[string, bool]()
	id := uuid.NewString()
	c.Set(id, true)
	mylog.Info("Count", c.Count())
	mylog.Struct(c.Keys())
	mylog.Struct(c.Items())
	mylog.Info("IsEmpty", c.IsEmpty())
	selected, b := c.Get(id)
	if !b {
		return
	}
	mylog.Info("Get", selected)
	c.Remove(id)
	mylog.Info("Count", c.Count())
	mylog.Info("IsEmpty", c.IsEmpty())
}

func ReadLines(fullpath string) ([]string, error) {
	f, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func WriteLines(lines []string, fullpath string) error {
	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}
