# mylog - Go 日志库

一个简洁、高效的 Go 日志库，提供格式化输出和多种日志级别支持。

## 功能特性

- 多种日志级别：Info、Warning、Success、Trace、Error
- 特殊格式支持：Hex、HexDump、Json、Struct
- 自动时间戳和调用位置
- 彩色终端输出
- 自动格式化值
- **Key 自动设置**：自动使用调用者函数名作为 key，用户只需传 value
- **简化错误堆栈**：自动过滤 mylog 内部函数，只显示用户代码调用栈

## 使用方法

### 基本用法

```go
package main

import "github.com/ddkwork/golibrary/std/mylog"

func main() {
    mylog.Info("processing data...")
    mylog.Warning("disk space low")
    mylog.Success("task completed")
}
```

### Key 自动设置

**重要：Key 自动使用调用者函数名，用户只需传 value，不需要也不应该手动设置 key！**

```go
func ProcessData() {
    mylog.Info("processing data...")
    // 输出: ProcessData │ processing data...
}

func LoadConfig() {
    mylog.Warning("config file not found")
    // 输出: LoadConfig │ config file not found
}
```

### 特殊格式

```go
// 十六进制输出
mylog.Hex(uint32(0xDEADBEEF))
mylog.HexDump([]byte{0x01, 0x02, 0x03, 0x04})

// JSON 输出
mylog.Json(`{"key": "value"}`)

// 结构体输出
mylog.Struct(myStruct)
```

## 错误处理

### Check 函数

`Check[T any](result T)` 检查结果，支持三种类型：

| 类型 | 行为 |
|------|------|
| `bool` | 如果为 false 则 panic |
| `string` | 如果不是成功消息则 panic |
| `error` | 如果不是 EOF 且不是成功消息则 panic |

**注意：Check 失败会 panic，需要配合 Call 使用来捕获错误！**

```go
// ❌ 错误用法：Check 失败会 panic 导致程序崩溃
func BadExample() {
    err := someOperation()
    mylog.Check(err)  // 如果 err != nil，程序会 panic 崩溃
    mylog.Info("继续执行")  // 不会执行到这里
}

// ✅ 正确用法：用 Call 包裹，错误时打印堆栈但程序继续
func GoodExample() {
    mylog.Call(func() {
        err := someOperation()
        mylog.Check(err)  // 如果 err != nil，打印错误堆栈后继续
        mylog.Info("继续执行")  // 会执行到这里
    })
}
```

### Check2 函数

`Check2[T any](ret T, err error) T` 同时检查 error 和返回值是否为 nil：

```go
func ReadFile(path string) {
    mylog.Call(func() {
        data := mylog.Check2(os.ReadFile(path))  // 检查 error 和 data 是否为 nil
        mylog.Info("文件大小", len(data))
    })
}
```

### CheckNil 函数

`CheckNil(ptr any)` 检查指针是否为 nil：

```go
mylog.Call(func() {
    ptr := getPointer()
    mylog.CheckNil(ptr)  // 如果 ptr 为 nil 则 panic
})
```

### Call 函数

`Call(f func())` 捕获内部的 panic 并打印错误堆栈，程序不会崩溃：

```go
func SafeOperation() {
    mylog.Call(func() {
        // 这里的任何 panic 都会被捕获并打印堆栈
        f := mylog.Check2(os.Open("config.json"))
        defer f.Close()
        
        data := mylog.Check2(io.ReadAll(f))
        mylog.Info("读取到", len(data), "字节")
    })
    // 即使上面出错，这里也会继续执行
    mylog.Info("操作完成")
}
```

## API 参考

### 基础日志函数

```go
func Info(msg ...any)
func Warning(msg ...any)
func Success(msg ...any)
func Trace(msg ...any)
```

### 特殊格式函数

```go
func Hex[V types.Unsigned](v V) string
func HexDump[V []byte | *bytes.Buffer](buf V)
func Json(msg ...any)
func Struct(object any)
```

### 错误处理函数

```go
func Check[T any](result T) (isEof bool)      // 检查 bool/string/error，失败 panic
func Check2[T any](ret T, err error) T        // 检查 error 和返回值
func CheckNil(ptr any)                         // 检查指针是否为 nil
func Call(f func())                            // 捕获 panic 并打印堆栈
```

## 限制规则

### Value 限制

| 规则 | 说明 |
|------|------|
| 不能为空 | value 必须提供有效内容 |
| 不能包含格式化语法 | 禁止使用 `%s`、`%d` 等格式化符号，mylog 会自动格式化 |

### Key 说明

| 说明 |
|------|
| **Key 完全自动设置**，使用调用者函数名，用户不需要也不应该手动设置 |

## 错误示例

```go
// ❌ 错误：value 为空
mylog.Info()

// ❌ 错误：value 包含格式化语法
mylog.Info("value is %s", "test")

// ❌ 错误：Check 失败会 panic，不用 Call 包裹会导致程序崩溃
func BadCheck() {
    mylog.Check(errors.New("some error"))  // panic!
}
```

## 正确示例

```go
// ✅ 正确：只需传 value，key 自动使用函数名
func ProcessData() {
    mylog.Info("processing data...")
}

// ✅ 正确：多个 value 自动格式化
mylog.Info("user_id", 12345, "ip", "192.168.1.1")

// ✅ 正确：中文放在 value 中
mylog.Info("用户登录成功")

// ✅ 正确：用 Call 包裹 Check，错误时打印堆栈但程序继续
func SafeCheck() {
    mylog.Call(func() {
        mylog.Check(errors.New("some error"))  // 打印错误堆栈，程序继续
    })
    mylog.Info("这里会继续执行")
}
```

## 日志输出格式

### 普通日志

```
2026-03-22 02:47:26    Info -> ProcessData │ processing data... main.go:10
```

格式说明：
- 时间戳
- 日志级别（右对齐）
- 调用者函数名（自动设置）
- value 内容
- 调用位置（文件名:行号）

### 错误堆栈

```
2026-03-22 02:47:26   Error ->           │ open 2332: The system cannot find the file specified.
                                         │ mylog_test.bug check_test.go:82
                                         │ mylog_test.m5 check_test.go:77
                                         │ mylog_test.TestCheckM5 check_test.go:16
```

错误堆栈自动过滤 mylog 内部函数（Check、Check2 等），只显示用户代码调用栈。

## Panic 情况

以下情况会导致 panic：

1. `log value cannot be empty` - value 为空
2. `log value cannot contain format syntax` - value 包含格式化语法
3. `Check` 失败 - 传入 bool(false)、非成功 string、或非 nil error

**重要：使用 `mylog.Call(func(){})` 包裹可以捕获 panic，打印错误堆栈但程序不会崩溃！**

## 许可证

MIT License
