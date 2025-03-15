package csf

import (
	"testing"
)

func Test_EmptyTemplate(t *testing.T) {
	st := NewTemplate()
	s, err := st.String(map[string]any{
		"a": "foo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "" {
		t.Fatalf("expected empty string, got %q", s)
	}
}

func Test_EmptyContext(t *testing.T) {
	st := NewTemplate(
		F("a"),
	)
	s, err := st.String(nil)
	if err != nil {
		t.Fatal(err)
	}
	if s != "" {
		t.Fatalf("expected empty string, got %q", s)
	}
}

func Test_Concat(t *testing.T) {
	st := NewTemplate(
		F("a"),
		F("b"),
	)
	s, err := st.String(map[string]any{
		"a": "foo",
		"b": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo bar" {
		t.Fatalf("expected 'foo bar', got %q", s)
	}
}

func Test_ConcatMissingOptional(t *testing.T) {
	st := NewTemplate(
		F("a"),
		F("b"),
	)
	s, err := st.String(map[string]any{
		"a": "foo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatalf("expected 'foo', got %q", s)
	}
}

func Test_ConcatRequired(t *testing.T) {
	st := NewTemplate(
		F("a").Required(),
		F("b"),
	)
	s, err := st.String(map[string]any{
		"a": "foo",
		"b": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo bar" {
		t.Fatalf("expected 'foo bar', got %q", s)
	}
}

func Test_ConcatMissingRequired(t *testing.T) {
	st := NewTemplate(
		F("a").Required(),
		F("b"),
	)
	s, err := st.String(map[string]any{})
	if err == nil {
		t.Fatalf("expected error for missing required field, got nil")
	}
	if s != "" {
		t.Fatalf("expected empty string, got %q", s)
	}
}

func Test_DefaultValue(t *testing.T) {
	st := NewTemplate(
		F("a").Required(),
		F("b").Default("bar"),
		F("c").Required(),
	)
	s, err := st.String(map[string]any{
		"a": "foo",
		"c": "baz",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo bar baz" {
		t.Fatalf("expected 'foo bar baz', got %q", s)
	}
}

func Test_CustomFieldFormat(t *testing.T) {
	st := NewTemplate(
		F("a").Formatter(Const("foo")),
	)
	s, err := st.String(map[string]any{
		"a": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatalf("expected 'foo', got %q", s)
	}
}

func Test_ArrayFieldFormat(t *testing.T) {
	st := NewTemplate(
		F("a").Required().Formatter(Array(", ")),
	)
	s, err := st.String(map[string]any{
		"a": []string{"foo", "bar"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo, bar" {
		t.Fatalf("expected 'foo, bar', got %q", s)
	}
}

func Test_First(t *testing.T) {
	st := NewTemplate(
		First(
			F("a"),
			F("b"),
		),
	)
	s, err := st.String(map[string]any{
		"a": "foo",
		"b": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatalf("expected 'foo', got %q", s)
	}
}

func Test_FirstReverse(t *testing.T) {
	st := NewTemplate(
		First(
			F("b"),
			F("a"),
		),
	)
	s, err := st.String(map[string]any{
		"a": "foo",
		"b": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "bar" {
		t.Fatalf("expected 'bar', got %q", s)
	}
}

func Test_FirstOffset(t *testing.T) {
	st := NewTemplate(
		First(F("a"), F("b")),
	)
	s, err := st.String(map[string]any{
		"b": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "bar" {
		t.Fatalf("expected 'bar', got %q", s)
	}
}

func Test_FirstRequired(t *testing.T) {
	st := NewTemplate(
		First(
			F("a"),
			F("b").Required(),
		),
	)
	s, err := st.String(map[string]any{
		"a": "foo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatalf("expected 'foo', got %q", s)
	}
}

func Test_FirstRequiredBad(t *testing.T) {
	st := NewTemplate(
		First(
			F("a").Required(),
			F("b"),
		),
	)
	s, err := st.String(map[string]any{
		"b": "foo",
	})
	if err == nil {
		t.Fatalf("expected error for missing required field, got nil")
	}
	if s != "" {
		t.Fatalf("expected empty string, got %q", s)
	}
}

func Test_Constant(t *testing.T) {
	st := NewTemplate(
		C("foo"),
		F("a"),
		C("baz"),
	)
	s, err := st.String(map[string]any{
		"a": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo bar baz" {
		t.Fatalf("expected 'foo bar baz', got %q", s)
	}
}
