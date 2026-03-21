# fakeError

## 核心理念

**任何 Go 程序员都不应该忽略错误。**

fakeError 的存在是为了强制执行这一原则。它会自动将所有可能被忽略的错误处理代码转换为明确的 `mylog.Check*` 调用，确保每一个错误都被处理、记录或 panic，而不是被静默忽略。

## 为什么需要 fakeError？

在 Go 中，错误处理是显式的，但开发者经常因为以下原因忽略错误：

1. **懒惰**：`result, _ := someFunc()` - 用 `_` 忽略错误
2. **疏忽**：忘记处理 defer 中的错误
3. **样板代码疲劳**：重复写 `if err != nil { ... }` 让人厌烦
4. **侥幸心理**："这个函数不会失败"

这些行为会导致生产环境中难以调试的问题。fakeError 通过自动化转换，让正确处理错误变得简单，让忽略错误变得困难。

## fakeError 做什么？

fakeError 扫描你的代码，将所有传统的错误处理模式转换为 `mylog.Check*` 函数调用：

### 转换示例

**被忽略的错误 → 被处理的错误**
```go
// 转换前：错误被忽略
result, _ := someFunc()

// 转换后：错误被处理
result := mylog.Check2(someFunc())
```

**冗长的错误检查 → 简洁的错误处理**
```go
// 转换前
result, err := someFunc()
if err != nil {
    log.Fatal(err)
    return
}

// 转换后
result := mylog.Check2(someFunc())
```

**defer 中的错误 → 明确的错误处理**
```go
// 转换前：defer 中的错误经常被忽略
defer func() {
    if err := conn.Close(); err != nil {
        log.Debug(err)
    }
}()

// 转换后：错误被明确处理
defer mylog.Check(conn.Close())
```

## mylog.Check* 的价值

`mylog.Check*` 函数不仅仅是语法糖，它们提供了：

- **错误可见性**：所有错误都会被记录，不会静默消失
- **快速失败**：默认 panic，让问题立即暴露
- **可配置性**：可以通过 mylog 配置调整错误处理策略
- **调用栈追踪**：自动记录错误发生的调用栈
- **统一处理**：整个项目的错误处理逻辑一致

## 使用 fakeError

```go
package main

import "github.com/ddkwork/golibrary/std/fakeError"

func main() {
    // 扫描并转换当前目录的所有 Go 文件
    fakeError.Walk(".", true)
}
```

⚠️ **警告：此工具会原地修改代码，请使用版本控制保留原始代码**

## 转换后的代码风格

转换后的代码更加简洁，同时保持了错误处理的完整性：

```go
func main() {
    config := mylog.Check2(initConfig())
    
    for _, rule := range config.Rules {
        l := mylog.Check2(net.Listen("tcp", rule.ListenAddr))
        
        conn := mylog.Check2(l.Accept())
        defer mylog.Check(conn.Close())
        
        clientHello, clientReader := mylog.Check3(PeekClientHello(conn))
        
        backendConn := mylog.Check2(net.DialTimeout("tcp", target, 5*time.Second))
        defer mylog.Check(backendConn.Close())
        
        mylog.Check2(io.Copy(clientConn, backendConn))
        mylog.Check2(io.Copy(backendConn, clientReader))
    }
}
```

没有冗长的 `if err != nil` 检查，没有被忽略的错误，所有错误都被明确处理。

## 技术细节

fakeError 使用 AST（抽象语法树）分析来识别和转换错误处理模式：

- **AssignStmt**：检测 `result, err := func()` 模式
- **IfStmt**：检测 `if err != nil` 模式
- **DeferStmt**：检测 defer 中的错误处理
- **外部函数检测**：通过 AST 分析确定函数返回类型

## 测试

```bash
go test ./std/fakeError/...
```

## 哲学

> "Errors are values. Treating errors as values is a critical part of Go's design, and ignoring them is a mistake."

fakeError 帮助你实践这一哲学，让错误处理成为代码的一部分，而不是被忽略的负担。
