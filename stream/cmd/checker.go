package cmd

import (
	"github.com/ddkwork/golibrary/mylog"
)

// todo https://521github.com/uber-go/goleak

//go:generate go install github.com/karamaru-alpha/copyloopvar/cmd/copyloopvar@latest
//go:generate go install github.com/ckaznocha/intrange/cmd/intrange@latest
//go:generate go install go.uber.org/nilaway/cmd/nilaway@latest
//go:generate go install mvdan.cc/gofumpt@latest
//go:generate go install golang.org/x/tools/gopls@latest
//go:generate go install github.com/kisielk/errcheck@latest
//go:generate go install github.com/vektra/mockery/v2@latest
//go:generate mockery --all --with-expecter --inpackage

func CheckLoopvarAndNilPoint() {
	if mylog.IsLinux() {
		Run("go vet -vettool=`which copyloopvar` ./...")
		Run("errcheck -asserts ./...")
		// Run("go vet -vettool=`which intrange` ./...")
		// Run("go vet -vettool=`which nilaway` ./...")
		// Run("gofumpt -l -w .")
		return
	}
	// Run("go vet -vettool=C:\\Users\\Admin\\go\\bin\\copyloopvar.exe ./...")
	Run("errcheck.exe -asserts ./...")
	// Run("go vet -vettool=C:\\Users\\Admin\\go\\bin\\intrange.exe ./...")
	Run("go vet -vettool=C:\\Users\\Admin\\go\\bin\\nilaway.exe ./...")
	// Run("go vet -vettool=C:\\Users\\Admin\\go\\bin\\gopls.exe ./...")
	// Run("gofumpt -l -w .")
}
