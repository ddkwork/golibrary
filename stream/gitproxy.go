package stream

func GitProxy(isSetProxy bool) bool {
	s := New("")
	SetProxy(s, isSetProxy)
	SetNameAndEmail(s)
	SetSafecrlf(s)
	path, ok := JoinHomeFile(".gitconfig")
	if !ok {
		return false
	}
	return WriteTruncate(path, s.String())
}

func SetProxy(s *Stream, isSetProxy bool) {
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

func SetNameAndEmail(s *Stream) {
	s.WriteStringLn(`
[user]
	name = Admin
	email = 2762713521@qq.com
`)
}

func SetSafecrlf(s *Stream) {
	if IsWindows() {
		s.WriteStringLn(`
[core]
	autocrlf = false
`)
	}
}

/*
git提示“warning: LF will be replaced by CRLF”的解决办法

在文件提交时进行safecrlf检查

#拒绝提交包含混合换行符的文件
git config --global core.safecrlf true

#允许提交包含混合换行符的文件
git config --global core.safecrlf false

#提交包含混合换行符的文件时给出警告
git config --global core.safecrlf warn
通俗解释

git 的 Windows 客户端基本都会默认设置 core.autocrlf=true，设置core.autocrlf=true, 只要保持工作区都是纯 CRLF 文件，编辑器用 CRLF 换行，就不会出现警告了；
Linux 最好不要设置 core.autocrlf，因为这个配置算是为 Windows 平台定制；
Windows 上设置 core.autocrlf=false，仓库里也没有配置 .gitattributes，很容易引入 CRLF 或者混合换行符（Mixed Line Endings，一个文件里既有 LF 又有CRLF）到版本库，这样就可能产生各种奇怪的问题。

*/
