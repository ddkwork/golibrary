# mylog - Go 日志库

一个简洁、高效的 Go 日志库，提供格式化输出和多种日志级别支持。

## 功能特性

- 多种日志级别：Info、Warning、Success、Trace、Error
- 特殊格式支持：Hex、HexDump、Json、Struct
- 自动时间戳和调用位置
- 彩色终端输出
- 自动格式化值
- **自动填充 key**：自动使用调用者函数名作为 key，无需手动指定
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

### 自动填充 Key

Key 自动使用调用者函数名，无需手动指定：

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

### 错误处理

```go
func ReadFile(path string) {
    data, err := os.ReadFile(path)
    mylog.Check(err)  // 错误时打印堆栈并继续
    mylog.Info("file size", len(data))
}

func SafeOperation() {
    mylog.Call(func() {
        // 错误会自动捕获并打印堆栈
        f := mylog.Check2(os.Open("config.json"))
        defer f.Close()
    })
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
func Check(err error) bool
func Check2[T any](v T, err error) T
func CheckNil(ptr any)
func Call(f func())
```

## 限制规则

### Value 限制

| 规则 | 说明 |
|------|------|
| 不能为空 | value 必须提供有效内容 |
| 不能包含格式化语法 | 禁止使用 `%s`、`%d` 等格式化符号，因为 mylog 会自动格式化 |

### Key 限制

| 规则 | 说明 |
|------|------|
| 自动填充 | 使用调用者函数名作为 key |
| 最大长度 10 | 超长函数名截断为 `xxx...` 格式 |

## 错误示例

```go
// ❌ 错误：value 为空
mylog.Info()

// ❌ 错误：value 包含格式化语法
mylog.Info("value is %s", "test")
```

## 正确示例

```go
// ✅ 正确：自动使用函数名作为 key
func ProcessData() {
    mylog.Info("processing data...")
}

// ✅ 正确：多个 value 自动格式化
mylog.Info("user_id", 12345, "ip", "192.168.1.1")

// ✅ 正确：中文放在 value 中
mylog.Info("用户登录成功")
```

## 日志输出格式

### 普通日志

```
2026-03-22 02:47:26    Info -> ProcessData │ processing data... main.go:10
```

格式说明：
- 时间戳
- 日志级别（右对齐）
- 调用者函数名（自动填充）
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

## 许可证

MIT License
