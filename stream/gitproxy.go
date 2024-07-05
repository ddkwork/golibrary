package stream

import (
	"github.com/ddkwork/golibrary/mylog"
)

func GitProxy(isSetProxy bool) {
	mylog.Call(func() {
		s := NewBuffer("")
		SetProxy(s, isSetProxy)
		SetNameAndEmail(s)
		SetSafecrlf(s)
		path := JoinHomeFile(".gitconfig")
		WriteTruncate(path, s.String())
	})
}

func SetProxy(s *Buffer, isSetProxy bool) {
	if !isSetProxy {
		return
	}
	s.WriteStringLn(`
[http]
    proxy = socks5://127.0.0.1:7890
[https]
    proxy = socks5://127.0.0.1:7890
`)
}

func SetNameAndEmail(s *Buffer) {
	s.WriteStringLn(`
[user]
	name = Admin
	email = 2762713521@qq.com
`)
}

func SetSafecrlf(s *Buffer) {
	if IsWindows() {
		s.WriteStringLn(`
[core]
	autocrlf = false

[safe]
	directory = *

`)
	}
}
