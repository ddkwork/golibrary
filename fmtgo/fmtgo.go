package main

import (
	"github.com/ddkwork/golibrary/std/stream"
)

func main() {
	if stream.IsAndroid() {
		return
	}
	stream.Fix(".")
}
