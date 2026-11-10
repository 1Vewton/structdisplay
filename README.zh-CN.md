# structdisplay

[English](README.md) | 简体中文

[![Go Reference](https://pkg.go.dev/badge/github.com/1Vewton/structdisplay.svg)](https://pkg.go.dev/github.com/1Vewton/structdisplay)
[![Go Version](https://img.shields.io/badge/Go-1.27%2B-00ADD8?logo=go)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`structdisplay` 是一个零第三方依赖的轻量 Go 库，可以把任意值（尤其是嵌套结构体和多级指针）渲染成带缩进、类似 Markdown 的可读文本。

```text
- Profile:
    - Name: Alice
    - Age: 30
    - Address:
        - City: Shanghai
        - Zip: 200000
```

## 特性

- **基于反射**：无需结构体标签、额外接口或代码生成。
- **支持指针**：自动解引用多级指针并识别 `nil`。
- **可控的字段顺序**：按自定义顺序输出字段。
- **字段与类型过滤**：跳过敏感或无关数据。
- **自定义 `nil` 输出**：自行决定缺失值的展示方式。
- **可插拔模板**：替换默认的 Markdown 渲染器。
- **递归保护**：为渲染深度设置上限。
- **并发安全的辅助类型**：且零第三方依赖。

## 环境要求

- Go 1.27 或更高版本。

## 安装

```bash
go get github.com/1Vewton/structdisplay
```

## 快速开始

```go
package main

import (
	"fmt"
	"reflect"

	"github.com/1Vewton/structdisplay/pkg/display"
	"github.com/1Vewton/structdisplay/pkg/sequence"
)

type Address struct {
	City string
	Zip  string
}

type Profile struct {
	Name    string
	Age     int
	Address *Address
}

func main() {
	data := &Profile{
		Name:    "Alice",
		Age:     30,
		Address: &Address{City: "Shanghai", Zip: "200000"},
	}

	order := sequence.NewFieldSequence().
		AddSequence(
			reflect.TypeOf(Profile{}),
			[]string{"Name", "Age", "Address"},
		).
		AddSequence(
			reflect.TypeOf(Address{}),
			[]string{"City", "Zip"},
		)

	converter := display.NewConverter().
		WithFieldSequence(order).
		WithRecursionLimit(10)

	out, err := converter.ConvertAnyData("Profile", "", data, 0)
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
```

输出：

```text
- Profile:
    - Name: Alice
    - Age: 30
    - Address:
        - City: Shanghai
        - Zip: 200000
```

### `ConvertAnyData`

```go
ConvertAnyData(title string, baseString string, data any, layer int) (string, error)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `title` | `string` | 为当前值渲染的标题。 |
| `baseString` | `string` | 渲染结果会追加到该前缀之后；传 `""` 表示从空字符串开始。 |
| `data` | `any` | 要渲染的值，通常是结构体或结构体指针。 |
| `layer` | `int` | 起始嵌套层级，顶层传 `0`。 |

所有 `With...` 方法都会修改并返回同一个转换器，因此可以链式调用。

## 渲染流程

1. 逐层解引用指针，直到得到具体值或 `nil`。
2. `nil` 值交由已配置的 `NilToolset` 处理；没有匹配的处理器时直接省略。
3. 使用已配置的 ignore 规则检查该值。
4. 结构体先输出标题，再逐行输出每个导出字段，最后输出尾部；每个字段以 `layer + 1` 递归渲染。
5. 其他类型的值交给模板的 `DisplayField` 方法处理。
6. 未导出的结构体字段会被跳过。

默认模板使用 Markdown 风格的项目符号，每层缩进 4 个空格。`map[any]any` 会递归展开，其他值使用 Go 的 `%v` 格式化。

## 字段顺序

字段顺序按结构体类型分别定义：

```go
order := sequence.NewFieldSequence().
	AddSequence(
		reflect.TypeOf(Profile{}),
		[]string{"Name", "Age", "Address"},
	).
	AddSequence(
		reflect.TypeOf(Address{}),
		[]string{"City", "Zip"},
	)
```

请使用结构体类型（`reflect.TypeOf(Profile{})`），而不是指针类型。如果需要稳定的输出，请为每个结构体类型（包括嵌套类型）都注册顺序。未列出的字段仍会排在已列字段之后输出，但它们之间的相对顺序不作保证。

## 忽略字段和类型

下面的片段复用快速开始中的 `Address`，并新增以下类型；同时假定已导入 `strings` 以及 `ignore`、`process` 包：

```go
type Secret struct {
	Value string
}

type Account struct {
	Name     string
	Age      int
	Password string
	Address  *Address
	Token    *Secret
}
```

使用 `ignore` 将数据排除在输出之外：

```go
order := sequence.NewFieldSequence().
	AddSequence(
		reflect.TypeOf(Account{}),
		[]string{"Name", "Age", "Password", "Address", "Token"},
	).
	AddSequence(
		reflect.TypeOf(Address{}),
		[]string{"City", "Zip"},
	)

ignored := ignore.NewTermsIgnoredMap()
ignored.AddField("Password")
ignored.AddType("Secret")
// AddTypeByType(reflect.TypeOf(Secret{})) 的效果相同。

converter := display.NewConverter().
	WithFieldSequence(order).
	WithIgnoredTerms(ignored).
	WithRecursionLimit(10)
```

对于包含 `Name`、`Age`、`Password`、`Address` 和 `Token *Secret` 的 `Account`，结果为：

```text
- Account:
    - Name: Alice
    - Age: 30
    - Address:
        - City: Shanghai
        - Zip: 200000
```

- `AddField` 按字段名精确匹配（区分大小写），并对所有结构体类型生效。
- `AddType` 按 `reflect.Type.Name()` 匹配；`AddTypeByType` 从 `reflect.Type` 推导类型名。
- `NewTermsIgnoredSlice` 是 `ignore.TermsIgnored` 的另一种实现，行为相同。

## 处理 `nil` 值

默认情况下，`nil` 字段会被省略。注册处理器即可显式渲染：

```go
nilHandler := process.NewNilToolset().AddProcessMethod(
	"Address",
	func(layer int) string {
		return strings.Repeat("    ", layer) + "- Address: <nil>\n"
	},
)

converter := display.NewConverter().
	WithFieldSequence(order).
	WithIgnoredTerms(ignored).
	WithProcessNilToolset(nilHandler).
	WithRecursionLimit(10)
```

对于 `Address` 为 `nil` 的 `Account{Name: "Bob", Age: 18}`，结果为：

```text
- Account:
    - Name: Bob
    - Age: 18
    - Address: <nil>
```

处理器按字段名匹配，回调会收到当前 `layer`，便于生成对应的缩进。返回的字符串会原样追加到结果中。

## 自定义模板

模板负责所有文本输出，实现 `template.Template` 即可：

```go
type Template interface {
	DisplayField(layer int, fieldName any, value any) string
	DisplayTitle(layer int, title string) string
	DisplayTail(layer int) string
}
```

- `DisplayTitle` 在结构体字段之前调用。
- `DisplayField` 对所有非结构体值调用（包括 map）。
- `DisplayTail` 在结构体字段之后调用。
- `layer` 从 `ConvertAnyData` 传入的值开始，并随嵌套层级递增。

示例：一个简单的 `key = value` 渲染器：

```go
type PlainTemplate struct{}

func (PlainTemplate) DisplayField(layer int, fieldName any, value any) string {
	return fmt.Sprintf("%s%v = %v\n", strings.Repeat("  ", layer), fieldName, value)
}

func (PlainTemplate) DisplayTitle(layer int, title string) string {
	return fmt.Sprintf("%s%s:\n", strings.Repeat("  ", layer), title)
}

func (PlainTemplate) DisplayTail(layer int) string {
	return ""
}
```

使用方式：

```go
converter := display.NewConverter().
	WithTemplate(PlainTemplate{}).
	WithFieldSequence(order).
	WithRecursionLimit(10)
```

输出：

```text
Profile:
  Name = Alice
  Age = 30
  Address:
    City = Shanghai
    Zip = 200000
```

## 递归限制

`WithRecursionLimit` 用于限制渲染深度：

```go
converter := display.NewConverter().WithRecursionLimit(10)
```

达到上限时，渲染会停止，`ConvertAnyData` 返回错误。指针解引用也会推进内部层级，因此对于较长的指针链应预留足够余量。未设置上限时，不施加深度限制。

## API 概览

### `pkg/display`

- `NewConverter() *Converter[template.Template]`：使用默认模板创建转换器。
- `WithTemplate(T)`：使用自定义模板。
- `WithFieldSequence(*sequence.FieldSequence)`：设置字段顺序。
- `WithIgnoredTerms(terms)`：设置忽略规则，接受任意 `ignore.TermsIgnored` 实现。
- `WithProcessNilToolset(*process.NilToolset)`：设置 `nil` 处理器。
- `WithRecursionLimit(int)`：设置最大深度。
- `ConvertAnyData(title, baseString string, data any, layer int) (string, error)`：渲染一个值。

### `pkg/sequence`

- `NewFieldSequence()`
- `AddSequence(reflect.Type, []string)`
- `AddSequenceByTypeName(string, []string)`
- `GetSequence(reflect.Type) []string`
- `GetSequenceByTypeName(string) []string`

### `pkg/ignore`

- `NewTermsIgnoredSlice()`、`NewTermsIgnoredMap()`
- `AddField(...string)`、`AddType(...string)`、`AddTypeByType(...reflect.Type)`
- `IsFieldIgnored(string) bool`、`IsTypeIgnored(string) bool`、`CheckIsTypeIgnoredByType(reflect.Type) bool`

### `pkg/process`

- `NewNilToolset()`
- `AddProcessMethod(fieldName string, method NilMethod)`，其中 `NilMethod` 为 `func(layer int) string`

### `pkg/template`

- `Template`：包含 `DisplayField`、`DisplayTitle` 和 `DisplayTail` 的接口。
- `templatedefault.NewDefaultTemplate(nil)`：内置的 Markdown 风格模板。

## 并发

`FieldSequence`、`TermsIgnoredSlice`、`TermsIgnoredMap` 和 `NilToolset` 均可安全地并发使用。转换期间 `Converter` 只读取自身配置，因此建议配置一次后重复使用。如果传入自定义模板，在共享转换器时请确保模板本身也是并发安全的。

## 项目结构

```text
.
├── internal/operations/          反射辅助函数（FetchData、GetFields）
├── pkg/display/                  转换器与转换流程
├── pkg/ignore/                   字段与类型过滤
├── pkg/process/                  自定义 nil 处理
├── pkg/sequence/                 字段顺序定义
├── pkg/template/                 模板接口
└── pkg/template/templatedefault/ 默认 Markdown 风格模板
```

## 许可证

MIT，详见 [LICENSE](LICENSE)。
