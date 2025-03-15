package csf

import (
	"fmt"
	"strings"
)

// Stringer mimics fmt.Stringer but returns a string for a given value rather
// than being a method on a type. This is used to provide formatting callbacks
// when stringifying a field's value in the template.
type Stringer func(v any) string

// Value is a default Stringer implementation that returns the string
// representation of the input value using fmt.Sprintf.
func Value(v any) string {
	return fmt.Sprintf("%v", v)
}

// Array returns a Stringer conforming function which generates a seperated
// string list from an array input. []string values are joined by the provided
// separator directly, whereas other array value types are string formatted
// using Value and then joined by the separator. If the input is not an array,
// Value(v) is returned.
func Array(sep string) Stringer {
	return func(v any) string {
		switch t := v.(type) {
		case []string:
			return strings.Join(t, sep)
		case []any:
			parts := make([]string, len(t))
			for i, item := range t {
				parts[i] = Value(item)
			}
			return strings.Join(parts, sep)
		default:
			return Value(v)
		}
	}
}

// Const returns a Stringer conforming function that always returns the provided
// string value.
func Const(s string) Stringer {
	return func(_ any) string {
		return s
	}
}

// Eval is a generic interface for evaluating the value of a field in a given
// context map as a string value.
type Eval interface {
	// String evaluates the field's value in the context map and returns its
	// string representation. A non-nil error is returned if the implementation
	// cannot evaluate the value (e.g. the field is required and not present in
	// the context map). An empty string indicates the field was not found or
	// its value is nil.
	String(c map[string]any) (string, error)
}

// Field represents a single named value in the template conforming to Eval.
type Field struct {
	// id is the corresponding key of the field in the context map.
	id string
	// req indicates whether the field is required. If true, the field must be
	// present in the context map and not nil, OR a default value must be set.
	req bool
	// def is the default value for the field. If the field is not present in
	// the context map or its value is nil, the default value will be used.
	def any
	// format is a custom Stringer function to format the field's value.
	// Defaults to Value when Field is constructed via F. Should not be nil.
	format Stringer
}

// Required marks the field as required. If the field is not present in the
// context map and no default value is set, an  error will be returned when the
// template is evaluated.
func (f *Field) Required() *Field {
	f.req = true
	return f
}

// Default sets a default value for the field. If the field is not present in
// the context map or its value is nil, the default value will be used instead.
func (f *Field) Default(v any) *Field {
	f.def = v
	return f
}

// Formatter sets a custom Stringer function to format the field's value. This
// allows for custom formatting logic to be applied to the field's value when
// generating the string representation. If not set, the default Value function
// is used to convert the value to a string.
func (f *Field) Formatter(fmt Stringer) *Field {
	f.format = fmt
	return f
}

// String evaluates the field's value in the context map and returns its
// string representation using its provided Stringer. A non-nil error is
// returned if the field is required and not present in the context map (and no
// acceptable default value is provided). An empty string indicates the field
// was not found or its value is nil.
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

// F creates a new Field instance with the provided id and sets the default
// format to Value. The field instance defaults to being optional (not required)
// and has no default value.
func F(id string) *Field {
	return &Field{
		id:     id,
		format: Value,
	}
}

// FirstMatch is an Eval-confirming implementation that evaluates the given
// fields in order and returns the first successfully evaluated value as a
// string, if any.
type FirstMatch struct {
	fields []Eval
}

// String returns the first non-zero/non-nil value from the list of fields
// provided to First. If no fields are found, an empty string is returned.
// If an error occurs while evaluating a field, it returns the error directly.
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

// First creates an evaluator that returns the first non-nil/non-zero value from
// a list of field values, otherwise returning an empty string. This is useful
// for cases where fields have some form of mutually exclusive relationships.
func First(fields ...Eval) *FirstMatch {
	return &FirstMatch{
		fields: fields,
	}
}

// Constant is an Eval-confirming implementation that always returns the
// provided constant string value. This is useful for cases where a fixed
// string value is needed in the template without evaluating any context map.
type Constant struct {
	v string
}

// String returns the constant value as a string. It does not evaluate the
// context map and always returns the same value.
func (c *Constant) String(_ map[string]any) (string, error) {
	return c.v, nil
}

// C creates a new Constant instance with the provided string value. This
// instance can be used in a template to return a fixed string value without
// evaluating the context map.
func C(v string) *Constant {
	return &Constant{
		v: v,
	}
}

// Template represents a list of ordered Eval values that can be used to
// generate a string representation from an input (a "context map") with
// conditional inclusion behavior and formatting delegates.
type Template struct {
	fields []Eval
}

// String generates a string representation of the template using the provided
// context map inputs. Each template Eval is evaluated in order, returning any
// non-nil errors. Otherwise, the corresponding string, if non-empty, is
// concatenated into a single string with a space separator. If no fields are
// found or all are nil, an empty string is returned.
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

// NewTemplate creates a new Template instance with the provided list of Eval
// values. The fields are stored in the order they are provided and will be
// evaluated in that order when generating the string representation.
func NewTemplate(fields ...Eval) *Template {
	return &Template{
		fields: fields,
	}
}
