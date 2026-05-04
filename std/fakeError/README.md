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

## mylog.Check* 是安全的 — panic 回收机制

很多人担心 `mylog.Check` 内部调用 `panic` 不安全。**这是误解。**

`mylog.Check*` 通过 `recover()` 机制回收 panic：

```go
// mylog/call.go 核心实现
func Call(f func()) {
    callWithHandler(f, func(err error) { l.errorCall(err) })
}

func callWithHandler(f func(), errHandler func(error error)) {
    defer recovery(errHandler)  // ← recover 捕获 panic
    f()                          // ← 如果 Check 遇到 error 会 panic
}

func recovery(handler recoveryHandler) {
    if recovered := recover(); recovered != nil && handler != nil {
        // 将 panic 转为 error 处理（记录日志 + 调用栈）
        // 然后从 panic 点返回，不继续执行后续代码
        handler(e)
    }
}
```

**关键点：**
- `mylog.Check(err)` 在 err != nil 时 **panic**
- 外层 `defer recovery()` **立即捕获**这个 panic
- 记录完整错误信息和调用栈后，**从 panic 发生处返回**
- **不会崩溃程序，不会继续执行后续代码**

这等价于 `if err != nil { return }` 但更简洁，还附带日志和调用栈。

## 什么情况应该替换为 mylog.Check*？

### ✅ 应该替换

**1. 简单错误判断 + 直接返回**
```go
// 转换前
f, err := os.Open("")
if err != nil {
    return
}

// 转换后
f := mylog.Check2(os.Open(""))
```

**2. 错误判断 + 返回 bool/error**
```go
// 转换前
token, err := OpenProcessToken(...)
if err != nil {
    return false
}

// 转换后
token := mylog.Check2(OpenProcessToken(...))
```

**3. 接口方法返回 error**
```go
// 转换前
n, err := w.Write([]byte("hello"))
if err != nil {
    return
}

// 转换后
n := mylog.Check2(w.Write([]byte("hello")))
```

**4. defer 中的 Close/清理操作**
```go
// 转换前
defer f.Close()

// 转换后（必须用 func() 包装，否则参数立即求值）
defer func() { mylog.Check(f.Close()) }()
```

> ⚠️ **为什么 defer 必须用 `func()` 包装？**
>
> Go 的 `defer expr` 规则：**表达式在 defer 声明时立即求值**。
>
> ```go
> defer mylog.Check(token.Close())  // ❌ token.Close() 立即执行！句柄当场关闭
> defer func() { mylog.Check(token.Close()) }()  // ✅ 延迟到函数返回时执行
> ```

---

### ❌ 不应该替换

**1. 根据 error 类型分支处理不同逻辑**
```go
// 不替换：不同的 error 有不同的处理策略
schService, e := CreateService(sc)
if e != nil {
    if e == ERROR_SERVICE_EXISTS {
        return false          // 已存在，直接返回
    }
    if e == ERROR_SERVICE_MARKED_FOR_DELETE {
        return false          // 标记删除中，提示用户重试
    }
    return false              // 其他错误
}
```

**2. 需要 Warning 提示用户处理后重试**
```go
// 不替换：需要输出提示信息让用户手动干预
h, e := CreateFile(namePtr, ...)
if e != nil {
    mylog.Warning("CreateFile failed", "error", e)  // 用户看到后需要处理环境问题
    r.driver.Stop()
    r.driver.Remove()
    return false
}
```

**3. 循环内的 break（不是函数返回）**
```go
// 不替换：break 只是跳出循环，不是退出函数
for off := 0; off < len(data); {
    inst, e := x86asm.Decode(data[off:], 64)
    if e != nil || inst.Len == 0 {
        break    // mylog.Check 会 panic 退出整个函数，而不是只跳出循环
    }
    off += inst.Len
}
```

> **循环 break 的正确处理方案**：使用 `mylog.CheckIgnore` 或保持原样

**4. if err 块内有业务逻辑**
```go
// 不替换：err 块内有赋值、计算等业务逻辑
if err != nil {
    count++
    result = defaultValue
    log.Printf("fallback triggered: %v", err)
    return
}
```

## 转换示例总览

| 模式 | 转换前 | 转换后 |
|------|--------|--------|
| 忽略错误 | `r, _ := f()` | `r := mylog.Check2(f())` |
| 简单返回 | `if err != nil { return }` | `mylog.Check(f())` |
| 返回 bool | `if err != nil { return false }` | `x := mylog.Check2(f())` |
| log.Fatal | `if err != nil { log.Fatal(err) }` | `mylog.Check(f())` |
| panic | `if err != nil { panic(err) }` | `mylog.Check(f())` |
| defer Close | `defer f.Close()` | `defer func() { mylog.Check(f.Close()) }()` |

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

## 技术细节

fakeError 使用 AST（抽象语法树）分析来识别和转换错误处理模式：

- **AssignStmt**：检测 `result, err := func()` 模式
- **IfStmt**：检测 `if err != nil` 模式（仅简单返回情况）
- **DeferStmt**：检测 defer 中的错误处理，自动添加 `func()` 包装
- **外部函数检测**：通过 AST 分析确定函数返回类型
- **接口方法检测**：识别接口方法调用的 error 返回值
- **智能跳过**：根据 if 块内语句复杂度判断是否应该跳过转换

## 测试

```bash
go test ./std/fakeError/...
```

测试覆盖 22 个场景，包括：
- 基础错误处理转换（test1~test16）
- defer Close 自动包装（test17）
- 简单返回替换（test18）
- 复杂 err 分支不替换（test19）
- Warning+return 不替换（test20）
- 循环 break 不替换（test21）
- 接口方法 error 替换（test22）

## 哲学

> "Errors are values. Treating errors as values is a critical part of Go's design, and ignoring them is a mistake."

fakeError 帮助你实践这一哲学，让错误处理成为代码的一部分，而不是被忽略的负担。
