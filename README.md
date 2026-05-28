# csf

csf (context-based string formatting) is a Go library for building human-readable strings from map data. Define a template of named fields, pass in a `map[string]any`, and csf assembles the output - automatically handling missing fields, fallback logic, and custom formatting.

## Installation

```
go get github.com/Cryptkeeper/csf
```

## Example

```go
st := csf.NewTemplate(
    csf.C("Author:"),
    csf.First(
        csf.F("preferred_name"),
        csf.F("legal_name").Required(),
    ),
    csf.F("company").Format(func(v any) string {
        return fmt.Sprintf("@ %s", v)
    }),
    csf.F("email").Format(func(v any) string {
        return fmt.Sprintf("<%s>", v)
    }),
)
```

The same template produces different output depending on which fields are present in the input:

```go
st.String(map[string]any{
    "legal_name": "John Smith",
    "email":      "john-smith@example.com",
})
// "Author: John Smith <john-smith@example.com>"

st.String(map[string]any{
    "legal_name":     "John Smith",
    "preferred_name": "Johnny Apple",
    "email":          "john-smith@example.com",
})
// "Author: Johnny Apple <john-smith@example.com>"

st.String(map[string]any{
    "preferred_name": "Johnny Apple",
    "company":        "ExampleCo",
    "email":          "john-smith@example.com",
})
// "Author: Johnny Apple @ ExampleCo <john-smith@example.com>"
```

Each field is evaluated against the input map, and non-empty results are joined with spaces. Optional fields that are missing or nil are silently omitted - no trailing spaces or empty brackets. `First` picks `preferred_name` when available, falling back to `legal_name`. Since `legal_name` is marked required, an error is returned when neither name is present.

## Features

- **Optional and required fields** - `csf.F("email")` is optional (omitted when missing), `csf.F("name").Required()` returns an error when missing.
- **Default values** - `csf.F("role").Default("member")` provides a fallback when the key is absent or nil.
- **Custom formatters** - `csf.F("tags").Format(csf.Array(", ", csf.Value))` controls how a field's value is stringified. Built-in stringers include `Value`, `Array`, and `Const`.
- **Conditional selection** - `csf.First(fields...)` returns the first non-empty value, useful for fallbacks or mutually exclusive fields.
- **Constants** - `csf.C("literal")` injects a fixed string into the template.
- **Nesting** - Templates implement the `Eval` interface, so they can be composed within other templates.

## Extending

Both core interfaces - `Eval` and `Stringer` - are small enough to implement yourself. `Eval` is a single method (`String(map[string]any) (string, error)`), and `Stringer` is just `func(any) string`. This makes it straightforward to add behavior that the library doesn't provide out of the box.

For example, a custom `Eval` that wraps a value in parentheses only when present:

```go
type Parenthesized struct {
    field csf.Eval
}

func (p *Parenthesized) String(c map[string]any) (string, error) {
    s, err := p.field.String(c)
    if err != nil || s == "" {
        return s, err
    }
    return "(" + s + ")", nil
}

// csf.NewTemplate(csf.F("name"), &Parenthesized{csf.F("nickname")})
// with {"name": "John", "nickname": "Johnny"} -> "John (Johnny)"
```

Or a reusable `Stringer` that truncates long values:

```go
func Truncate(max int) csf.Stringer {
    return func(v any) string {
        s := csf.Value(v)
        if len(s) > max {
            return s[:max] + "..."
        }
        return s
    }
}

// csf.F("bio").Format(Truncate(50))
```

## Motivation

Generating human-readable strings from map data usually means reaching for `text/template` or writing a pile of `if` statements in a `String()` method. csf is a more declarative alternative for data that comes in variable shapes with optional field inclusion.

## Limitations

csf is not especially fast, not type-safe, and is "one shot" - you cannot progressively build the string. It is likely not a replacement for deeply nested formatting logic, but it may help.

For more examples, see the [test cases](csf_test.go).
