package csf

import (
	"fmt"
	"reflect"
	"strings"
)

// Stringer converts an arbitrary value to its string representation.
type Stringer func(v any) string

// Value is the default [Stringer]. It formats v using fmt.Sprintf("%v", v).
func Value(v any) string {
	return fmt.Sprintf("%v", v)
}

// Const returns a [Stringer] that always returns s, ignoring the input value.
func Const(s string) Stringer {
	return func(_ any) string {
		return s
	}
}

// unpackArray uses reflection to convert a typed slice (e.g. []string) to []any.
func unpackArray(a any) []any {
	vo := reflect.ValueOf(a)
	v := make([]any, vo.Len())
	for i := 0; i < vo.Len(); i++ {
		v[i] = vo.Index(i).Interface()
	}
	return v
}

// Array returns a [Stringer] that treats its input as a slice, formats each
// element with stringer, filters out empty strings, and joins the results with sep.
func Array(sep string, stringer Stringer) Stringer {
	return func(v any) string {
		uv := unpackArray(v)
		strs := make([]string, 0, len(uv))
		for _, item := range uv {
			str := stringer(item)
			if len(str) > 0 {
				strs = append(strs, str)
			}
		}
		return strings.Join(strs, sep)
	}
}

// Eval evaluates a context map and produces a string result.
// An empty string signals that the eval produced no output (e.g. an optional
// field whose key is absent). A non-nil error signals a hard failure.
type Eval interface {
	String(c map[string]any) (string, error)
}

// Field is an [Eval] that looks up a single key in the context map.
// By default a missing or nil key produces an empty string. Use [Field.Required]
// to make it an error, [Field.Default] to supply a fallback, and [Field.Format]
// to control stringification.
type Field struct {
	id     string   // key in the context map
	req    bool     // whether a missing/nil value is an error
	def    any      // fallback when the key is missing or nil
	format Stringer // converts the resolved value to a string
}

// Required marks the field as required. Evaluation returns an error if the key
// is missing or nil and no default has been set.
func (f *Field) Required() *Field {
	f.req = true
	return f
}

// Default sets a fallback value used when the key is missing or nil.
func (f *Field) Default(v any) *Field {
	f.def = v
	return f
}

// Formatter sets a custom [Stringer] for the field's value.
//
// Deprecated: Use [Field.Format].
func (f *Field) Formatter(fmt Stringer) *Field {
	f.format = fmt
	return f
}

// Format sets a custom [Stringer] for the field's value.
// The default is [Value].
func (f *Field) Format(fmt Stringer) *Field {
	f.format = fmt
	return f
}

// String implements [Eval]. It resolves the field's value from the context map,
// falling back to the default if the key is missing or nil.
func (f *Field) String(c map[string]any) (string, error) {
	v := c[f.id]
	if v == nil {
		v = f.def
	}
	if v == nil {
		if f.req {
			return "", fmt.Errorf("context missing required field %q", f.id)
		}
		return "", nil
	}
	return f.format(v), nil
}

// F creates an optional [Field] that looks up id in the context map.
// The field uses [Value] as its default formatter.
func F(id string) *Field {
	return &Field{
		id:     id,
		format: Value,
	}
}

// FirstMatch is an [Eval] that returns the first non-empty result from an
// ordered list of evals, short-circuiting on the first match.
type FirstMatch struct {
	fields []Eval
}

// String implements [Eval]. It returns the first non-empty result, or ""
// if no eval produces output.
func (f *FirstMatch) String(c map[string]any) (string, error) {
	for i, field := range f.fields {
		if s, err := field.String(c); err != nil {
			return "", fmt.Errorf("error evaluating field %d: %w", i, err)
		} else if len(s) > 0 {
			return s, nil
		}
	}
	return "", nil
}

// First creates a [FirstMatch] from the given evals.
func First(fields ...Eval) *FirstMatch {
	return &FirstMatch{
		fields: fields,
	}
}

// Constant is an [Eval] that always returns a fixed string, ignoring the context.
type Constant struct {
	v string
}

// String implements [Eval]. It always returns the constant's value.
func (c *Constant) String(_ map[string]any) (string, error) {
	return c.v, nil
}

// C creates a [Constant] that always returns v.
func C(v string) *Constant {
	return &Constant{
		v: v,
	}
}

// Template is an ordered list of [Eval] values. It implements [Eval] itself,
// so templates can be nested. Evaluating a Template produces a space-joined
// string of all non-empty eval results.
type Template struct {
	fields []Eval
}

// String implements [Eval]. It evaluates each eval in order, collects non-empty
// results, and joins them with a single space.
func (t Template) String(ctx map[string]any) (string, error) {
	sb := make([]string, 0, len(t.fields))
	for _, f := range t.fields {
		if s, err := f.String(ctx); err != nil {
			return "", fmt.Errorf("error evaluating field %q: %w", f, err)
		} else if len(s) > 0 {
			sb = append(sb, s)
		}
	}
	return strings.Join(sb, " "), nil
}

// NewTemplate creates a [Template] from the given evals, evaluated in order.
func NewTemplate(fields ...Eval) *Template {
	return &Template{
		fields: fields,
	}
}
