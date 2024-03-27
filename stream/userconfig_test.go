package stream

import (
	"fmt"
	"testing"
)

func Test_getUserConfigDirs(t *testing.T) {
	userConfigDirs, err := GetUserConfigDirs()
	if err != nil {
		panic(err)
	}
	for username, ConfigDir := range userConfigDirs {
		fmt.Println(username + ": " + ConfigDir)
	}
}
