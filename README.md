# structdisplay

[English](README.md) | [简体中文](README.zh-CN.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/1Vewton/structdisplay.svg)](https://pkg.go.dev/github.com/1Vewton/structdisplay)
[![Go Version](https://img.shields.io/badge/Go-1.27%2B-00ADD8?logo=go)](https://go.dev/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`structdisplay` is a small, zero-dependency Go library that renders any value — especially nested structs and pointers — as readable, indented, Markdown-like text.

```text
- Profile:
    - Name: Alice
    - Age: 30
    - Address:
        - City: Shanghai
        - Zip: 200000
```

## Features

- **Reflection based** — no struct tags, extra interfaces, or code generation.
- **Pointer aware** — follows multi-level pointers and detects `nil` values.
- **Controllable field order** — render fields in the order you define.
- **Field and type filtering** — skip sensitive or noisy data.
- **Custom `nil` output** — decide how missing values are displayed.
- **Pluggable templates** — replace the default Markdown renderer.
- **Recursion guard** — put an upper bound on rendering depth.
- **Concurrency-safe helpers** — and no third-party dependencies.

## Requirements

- Go 1.27 or newer.

## Installation

```bash
go get github.com/1Vewton/structdisplay
```

## Quick start

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

Output:

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

| Parameter | Type | Description |
| --- | --- | --- |
| `title` | `string` | Label rendered for the current value. |
| `baseString` | `string` | Prefix the rendered text is appended to. Use `""` to start fresh. |
| `data` | `any` | Value to render, usually a struct or a pointer to a struct. |
| `layer` | `int` | Starting nesting depth. Use `0` for a top-level value. |

Every `With...` option mutates and returns the same converter, so options can be chained.

## How rendering works

1. Pointers are followed until a concrete value or `nil` is reached.
2. `nil` values are rendered by the configured `NilToolset`, or omitted when no handler matches.
3. The value is checked against the configured ignore terms.
4. A struct is rendered as a title, one line per exported field, and then a tail. Each field is rendered recursively with `layer + 1`.
5. Any other value is passed to the template's `DisplayField` method.
6. Unexported struct fields are skipped.

The default template produces Markdown-style bullets with four spaces per level. `map[any]any` values are expanded recursively; other values use Go's `%v` formatting.

## Field order

Field order is defined per struct type:

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

Use the struct type (`reflect.TypeOf(Profile{})`), not the pointer type. Register a sequence for every struct type, including nested ones, when output needs to be stable. Fields that are not listed are still rendered after the listed ones, but their relative order is not guaranteed.

## Ignoring fields and types

The snippets below reuse `Address` from the quick start and add these types. They also assume `strings` and the `ignore` and `process` packages are imported:

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

Use `ignore` to keep data out of the result:

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
// AddTypeByType(reflect.TypeOf(Secret{})) does the same thing.

converter := display.NewConverter().
	WithFieldSequence(order).
	WithIgnoredTerms(ignored).
	WithRecursionLimit(10)
```

Given an `Account` with `Name`, `Age`, `Password`, `Address`, and `Token *Secret`, the result is:

```text
- Account:
    - Name: Alice
    - Age: 30
    - Address:
        - City: Shanghai
        - Zip: 200000
```

- `AddField` matches field names exactly and applies to all struct types.
- `AddType` matches `reflect.Type.Name()`; `AddTypeByType` derives the name from a `reflect.Type`.
- `NewTermsIgnoredSlice` is an alternative implementation of `ignore.TermsIgnored` with the same behavior.

## Handling `nil` values

A `nil` field is omitted by default. Register a handler to render it explicitly:

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

With an `Account{Name: "Bob", Age: 18}` whose `Address` is `nil`, the result is:

```text
- Account:
    - Name: Bob
    - Age: 18
    - Address: <nil>
```

Handlers are keyed by field name, and the callback receives the current `layer` so it can produce matching indentation. The returned string is appended exactly as returned.

## Custom templates

A template controls every piece of text. Implement `template.Template`:

```go
type Template interface {
	DisplayField(layer int, fieldName any, value any) string
	DisplayTitle(layer int, title string) string
	DisplayTail(layer int) string
}
```

- `DisplayTitle` is called before a struct's fields.
- `DisplayField` is called for every non-struct value, including maps.
- `DisplayTail` is called after a struct's fields.
- `layer` starts at the value passed to `ConvertAnyData` and increases with each nesting level.

Example — a plain `key = value` renderer:

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

Use it like any other template:

```go
converter := display.NewConverter().
	WithTemplate(PlainTemplate{}).
	WithFieldSequence(order).
	WithRecursionLimit(10)
```

Output:

```text
Profile:
  Name = Alice
  Age = 30
  Address:
    City = Shanghai
    Zip = 200000
```

## Recursion limit

`WithRecursionLimit` bounds how deep rendering may go:

```go
converter := display.NewConverter().WithRecursionLimit(10)
```

When the limit is reached, rendering stops and `ConvertAnyData` returns an error. Pointer dereferences also advance the internal depth, so leave headroom for long pointer chains. Without a limit, no upper bound is applied.

## API overview

### `pkg/display`

- `NewConverter() *Converter[template.Template]` — creates a converter with the default template.
- `WithTemplate(T)` — uses a custom template.
- `WithFieldSequence(*sequence.FieldSequence)` — sets field order.
- `WithIgnoredTerms(terms)` — sets ignore rules; accepts any `ignore.TermsIgnored` implementation.
- `WithProcessNilToolset(*process.NilToolset)` — sets `nil` handlers.
- `WithRecursionLimit(int)` — sets the maximum depth.
- `ConvertAnyData(title, baseString string, data any, layer int) (string, error)` — renders a value.

### `pkg/sequence`

- `NewFieldSequence()`
- `AddSequence(reflect.Type, []string)`
- `AddSequenceByTypeName(string, []string)`
- `GetSequence(reflect.Type) []string`
- `GetSequenceByTypeName(string) []string`

### `pkg/ignore`

- `NewTermsIgnoredSlice()`, `NewTermsIgnoredMap()`
- `AddField(...string)`, `AddType(...string)`, `AddTypeByType(...reflect.Type)`
- `IsFieldIgnored(string) bool`, `IsTypeIgnored(string) bool`, `CheckIsTypeIgnoredByType(reflect.Type) bool`

### `pkg/process`

- `NewNilToolset()`
- `AddProcessMethod(fieldName string, method NilMethod)`, where `NilMethod` is `func(layer int) string`

### `pkg/template`

- `Template` — interface with `DisplayField`, `DisplayTitle`, and `DisplayTail`.
- `templatedefault.NewDefaultTemplate(nil)` — the built-in Markdown-style template.

## Concurrency

`FieldSequence`, `TermsIgnoredSlice`, `TermsIgnoredMap`, and `NilToolset` are safe for concurrent use. A `Converter` only reads its configuration during conversion, so configure it once and reuse it. If you supply a custom template, make sure it is safe for concurrent use when the converter is shared.

## Project layout

```text
.
├── internal/operations/          Reflection helpers (FetchData, GetFields)
├── pkg/display/                  Converter and conversion flow
├── pkg/ignore/                   Field and type filtering
├── pkg/process/                  Custom nil handling
├── pkg/sequence/                 Field order definitions
├── pkg/template/                 Template interface
└── pkg/template/templatedefault/ Default Markdown-style template
```

## License

MIT — see [LICENSE](LICENSE).
