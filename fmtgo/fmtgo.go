package main

import "github.com/ddkwork/golibrary/std/stream"

func main() {
	if stream.IsAndroid() {
		return
	}
	stream.RunCommands(`
go fix ./...
go run mvdan.cc/gofumpt@latest -l -w .`)
}
