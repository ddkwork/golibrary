//go:build !windows

package i18n

import "os"

func Locale() string {
	locale := os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LANG")
		if locale == "" {
			locale = "en_US.UTF-8"
		}
	}
	return locale
}
