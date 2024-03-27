package stream

import (
	"bufio"
	"os"
	"os/user"
	"runtime"
	"strings"
)

func GetUserConfigDirs() (UserConfigDirs map[string]string, err error) {
	UserConfigDirs = make(map[string]string)
	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		u, err := user.Current()
		if err != nil {
			return nil, err
		}
		UserConfigDirs[u.Username] = dir
	} else if IsTermux() {
		dir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		u, err := user.Current()
		if err != nil {
			return nil, err
		}
		UserConfigDirs[u.Username] = dir
	} else {
		file, err := os.Open("/etc/passwd")
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				username := parts[0]
				u, err := user.Lookup(username)
				if err != nil {
					continue
				}
				dir := u.HomeDir + "/.config"
				if strings.Contains(dir, "root") || strings.Contains(dir, "home") {
					UserConfigDirs[username] = dir
				}
			}
		}
	}
	return UserConfigDirs, nil
}
