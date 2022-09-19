package golibrary

//go:generate  go get -v -d github.com/ddkwork/librarygo@main
//go:generate  go get -v -d github.com/ddkwork/librarygo/mylog@latest
//go:generate  go build .
//go:generate  go  clean -modcache
//go:generate  go work init
//go:generate  go work use -r .
//go:generate  go work sync
//go:generate  go env -w GOPROXY=https://goproxy.cn
//go:generate  go env -w GOPRIVATE=gitee.com
//git tag -d v1.1.0
//git tag -l
//git push --delete origin v1.1.0
//https://blog.csdn.net/qq_39545674/article/details/120632719
//go:generate go get github.com/libsgh/PanIndex@v3.1.0
//go:generate git tag -d v1.1.0
//go:generate git push --delete origin v1.1.0
//import (
//	"go.uber.org/goleak"
//	"testing"
//)
//
//func leak() {
//	f := make(chan struct{})
//	go func() {
//		f <- struct{}{}
//	}()
//}
//
//func TestName(t *testing.T) {
//	defer goleak.VerifyNone(t)
//	leak()
//}
