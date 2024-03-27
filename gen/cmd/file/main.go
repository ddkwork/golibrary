package main

import (
	"github.com/ddkwork/golibrary/gen"
	"github.com/ddkwork/golibrary/stream/cmd"
)

//go:generate go install github.com/ddkwork/golibrary/gen/cmd/file@8e08f82181f0980837402ee8bdecc3bd3b70a0ae
//go:generate go mod download github.com/ddkwork/golibrary@8e08f82181f0980837402ee8bdecc3bd3b70a0ae

//go:generate go install .

func main() {
	gen.New().FileAction()
	cmd.Run("go mod tidy")
}
