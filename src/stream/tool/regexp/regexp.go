package regexp

import (
	"fmt"
	"regexp"
)

type (
	Interface interface {
		RegexpWebBodyBlocks(tagName string) string
		IntegerToIP(ip int64) string
	}
	object struct {
	}
)

var Default = New()

func New() Interface { return &object{} }

var RegexpCenter = `(.+?)`

func (*object) RegexpWebBodyBlocks(tagName string) string {
	return `<` + tagName + `[^>]*?>[\w\W]*?<\/` + tagName + `>`
}

func (*object) IntegerToIP(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

var (
	RegexpIp     = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))`) //todo panic
	RegexpIpPort = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))):([0-9]+)`)
)
