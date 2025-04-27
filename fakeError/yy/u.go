package main

import (
	"github.com/ddkwork/golibrary/mylog"
	"regexp"
	"strings"
)

func main() {
	for {
		keyServerName = ".*(" + strings.ReplaceAll(keyServerName[1:], ".", "\\.") + ")$"
		matched, err := regexp.Match(keyServerName, []byte(serverName))
		if err != nil {
			mylog.CheckIgnore(err)
			continue
		}
	}
}
