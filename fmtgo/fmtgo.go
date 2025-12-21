package main

import "github.com/ddkwork/golibrary/std/stream"

func main() {
	if stream.IsAndroid() {
		return
	}
	stream.RunCommands(`
go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -diff ./...
go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix ./...
go run mvdan.cc/gofumpt@latest -l -w .`)
}
