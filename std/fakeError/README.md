# fakeError

Go 代码错误处理自动转换工具

## 功能说明

fakeError 使用 AST（抽象语法树）分析自动将传统的错误处理模式转换为简化的 `mylog.Check*` 函数调用。

## 支持的转换模式

### 1. 简单错误检查（panic/log.Fatal）

**转换前：**
```go
if err != nil {
    log.Fatal(err)
    return
}
```

**转换后：**
```go
mylog.Check(...)
```

### 2. 错误检查（panic）

**转换前：**
```go
if err != nil {
    panic(err)
}
```

**转换后：**
```go
mylog.Check(...)
```

### 3. 错误检查（continue）

**转换前：**
```go
if err != nil {
    continue
}
```

**转换后：**
```go
mylog.CheckIgnore(err)
continue
```

### 4. Defer 错误处理

**转换前：**
```go
defer func() {
    if err := x.Close(); err != nil {
        log.Debug(err)
    }
}()
```

**转换后：**
```go
defer mylog.Check(x.Close())
```

### 5. 多返回值处理

**转换前：**
```go
result, err := someFunc()
if err != nil {
    return err
}
```

**转换后：**
```go
result := mylog.Check2(someFunc())
```

## 自动功能

- **删除冗余声明**：自动删除 `var err error` 声明
- **自动导入**：自动添加必要的 mylog 包导入
- **保留注释**：保留 `//go:` 构建标签和注释

## 使用方法

```go
package main

import "github.com/ddkwork/golibrary/std/fakeError"

func main() {
    // 遍历当前目录并移除注释
    fakeError.Walk(".", true)
    
    // 仅遍历不移除注释
    fakeError.Walk(".")
}
```

## 注意事项

⚠️ **此工具会原地修改代码，请使用版本控制保留原始代码**

转换设计用于减少样板错误处理代码，同时通过 mylog 包保持错误可见性。

## ⚠️ 忽略错误

fakeError 会转换所有使用 `err` 或 `_` 的错误处理。如果需要忽略错误，可以使用以下方法：

### 方法 1：使用不同的错误变量名

fakeError 只识别 `err` 变量，使用其他变量名可以避免转换：

**不会被转换：**
```go
data, e := os.ReadFile("file.txt")
// 不处理 e，直接忽略
```

### 方法 2：在 if 块中添加其他逻辑

如果 if 块中有其他逻辑（不仅仅是 panic/log/return/continue），fakeError 不会转换：

**不会被转换：**
```go
result, err := someFunc()
if err != nil {
    // 添加注释或其他逻辑，避免被转换
    _ = err
}
```

### 方法 3：避免使用 `_` 接收 error 类型（仅限本文件函数）

**会被转换：**
```go
result, _ := localFuncReturningError()
// 转换为：result := mylog.CheckIgnore(err)
```

**不会被转换：**
```go
result, _ := localFuncReturningNonError()
// 如果最后一个返回类型不是 error，则不会被转换
```

**重要说明**：
- fakeError 通过 AST 分析**本文件中定义的函数**的返回类型
- 对于**外部库函数**，AST 无法获取返回类型，使用 `_` 不会被转换
- 如果外部库函数返回 error，使用 `_` 会导致错误被忽略（不会转换为 CheckIgnore）

**示例：**
```go
// 外部库函数，AST 无法获取返回类型
data, _ := os.ReadFile("file.txt")  // 不会被转换，错误被忽略

// 本文件函数，AST 可以获取返回类型
result, _ := myFuncReturningError()  // 转换为：result := mylog.CheckIgnore(err)
```

## 工作原理

1. 使用 `go/parser` 解析 Go 源代码为 AST
2. 使用 `go/ast` 遍历语法树查找错误处理模式
3. 根据模式应用相应的转换规则
4. 使用 `go/format` 格式化生成的代码

## 转换规则

### AssignStmt 转换

- 检测最后一个赋值变量为 `err` 或 `_`
- 如果是 `_`，检查最后一个返回类型是否为 `error`
- 转换为 `mylog.CheckN()` 调用，N 为返回值数量

### IfStmt 转换

- 检测 `if err != nil` 模式
- 检查 if 块内语句数量
- 如果只有 1 条语句且为 panic/log/return/continue，直接转换

### DeferStmt 转换

- 检测 defer 中的 Close 方法调用
- 转换为 `mylog.Check()` 包装

## 测试

运行测试套件：

```bash
go test ./std/fakeError/...
```

测试覆盖了 10 个不同的转换场景。
