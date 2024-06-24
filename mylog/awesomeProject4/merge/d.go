package main

import (
	"fmt"
	"strings"
)

func main() {
	src := `
package main

import "fmt"

func main() {
	lines := strings.
Split(src, lineBreak)
}
`

	mergedSrc := mergeLines(src, ".", "\n")

	fmt.Println(mergedSrc)
}

func mergeLines(src, mergeToken, lineBreak string) string {
	lines := strings.Split(src, lineBreak)
	var mergedLines []string
	mergeBuffer := ""
	for _, line := range lines {
		if strings.HasSuffix(line, mergeToken) && mergeBuffer == "" {
			mergeBuffer = line
		} else {
			if mergeBuffer != "" {
				mergedLines = append(mergedLines, mergeBuffer+line)
				mergeBuffer = ""
			} else {
				mergedLines = append(mergedLines, line)
			}
		}
	}
	if mergeBuffer != "" {
		mergedLines = append(mergedLines, mergeBuffer)
	}
	return strings.Join(mergedLines, lineBreak)
}
