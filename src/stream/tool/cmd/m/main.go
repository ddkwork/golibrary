package main

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/stream/tool/cmd"
)

func main() {
	//b, err := cmd.Run("C:\\Windows\\System32\\PING.EXE www.baidu.com -t ")
	b, err := cmd.Run("ping www.baidu.com -t ")
	if !mylog.Error(err) {
		return
	}
	mylog.Json("ast", b)
	select {}
}
