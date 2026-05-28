package csf

import (
	"testing"
)

func TestValue(t *testing.T) {
	if s := Value("hello"); s != "hello" {
		t.Fatalf("expected %q, got %q", "hello", s)
	}
	if s := Value(42); s != "42" {
		t.Fatalf("expected %q, got %q", "42", s)
	}
	if s := Value(true); s != "true" {
		t.Fatalf("expected %q, got %q", "true", s)
	}
}

func TestConst(t *testing.T) {
	fn := Const("fixed")
	if s := fn("anything"); s != "fixed" {
		t.Fatalf("expected %q, got %q", "fixed", s)
	}
	if s := fn(nil); s != "fixed" {
		t.Fatalf("expected %q, got %q", "fixed", s)
	}
}

func TestArray(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		fn := Array(", ", Value)
		s := fn([]string{"foo", "bar"})
		if s != "foo, bar" {
			t.Fatalf("expected %q, got %q", "foo, bar", s)
		}
	})

	t.Run("ints", func(t *testing.T) {
		fn := Array(", ", Value)
		s := fn([]int{10, 20})
		if s != "10, 20" {
			t.Fatalf("expected %q, got %q", "10, 20", s)
		}
	})

	t.Run("empty array", func(t *testing.T) {
		fn := Array(", ", Value)
		s := fn([]string{})
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("filters empty strings", func(t *testing.T) {
		fn := Array(", ", func(v any) string {
			if v.(string) == "skip" {
				return ""
			}
			return v.(string)
		})
		s := fn([]string{"a", "skip", "b"})
		if s != "a, b" {
			t.Fatalf("expected %q, got %q", "a, b", s)
		}
	})

	t.Run("custom separator", func(t *testing.T) {
		fn := Array(" | ", Value)
		s := fn([]string{"x", "y", "z"})
		if s != "x | y | z" {
			t.Fatalf("expected %q, got %q", "x | y | z", s)
		}
	})

	t.Run("structs", func(t *testing.T) {
		type msg struct{ text string }
		fn := Array(", ", func(v any) string {
			return v.(msg).text
		})
		s := fn([]msg{{"hello"}, {"world"}})
		if s != "hello, world" {
			t.Fatalf("expected %q, got %q", "hello, world", s)
		}
	})
}

func TestField(t *testing.T) {
	t.Run("optional present", func(t *testing.T) {
		s, err := F("a").String(map[string]any{"a": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		if s != "foo" {
			t.Fatalf("expected %q, got %q", "foo", s)
		}
	})

	t.Run("optional missing", func(t *testing.T) {
		s, err := F("a").String(map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("optional nil value", func(t *testing.T) {
		s, err := F("a").String(map[string]any{"a": nil})
		if err != nil {
			t.Fatal(err)
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("required present", func(t *testing.T) {
		s, err := F("a").Required().String(map[string]any{"a": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		if s != "foo" {
			t.Fatalf("expected %q, got %q", "foo", s)
		}
	})

	t.Run("required missing", func(t *testing.T) {
		s, err := F("a").Required().String(map[string]any{})
		if err == nil {
			t.Fatal("expected error for missing required field")
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("required nil value", func(t *testing.T) {
		s, err := F("a").Required().String(map[string]any{"a": nil})
		if err == nil {
			t.Fatal("expected error for nil required field")
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("default used when missing", func(t *testing.T) {
		s, err := F("a").Default("fallback").String(map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
		if s != "fallback" {
			t.Fatalf("expected %q, got %q", "fallback", s)
		}
	})

	t.Run("default not used when present", func(t *testing.T) {
		s, err := F("a").Default("fallback").String(map[string]any{"a": "actual"})
		if err != nil {
			t.Fatal(err)
		}
		if s != "actual" {
			t.Fatalf("expected %q, got %q", "actual", s)
		}
	})

	t.Run("default satisfies required", func(t *testing.T) {
		s, err := F("a").Required().Default("fallback").String(map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
		if s != "fallback" {
			t.Fatalf("expected %q, got %q", "fallback", s)
		}
	})

	t.Run("format", func(t *testing.T) {
		s, err := F("a").Format(Const("always")).String(map[string]any{"a": "ignored"})
		if err != nil {
			t.Fatal(err)
		}
		if s != "always" {
			t.Fatalf("expected %q, got %q", "always", s)
		}
	})

	t.Run("format with array", func(t *testing.T) {
		s, err := F("a").Format(Array(", ", Value)).String(map[string]any{
			"a": []string{"x", "y"},
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "x, y" {
			t.Fatalf("expected %q, got %q", "x, y", s)
		}
	})
}

func TestConstant(t *testing.T) {
	t.Run("returns fixed value", func(t *testing.T) {
		s, err := C("hello").String(map[string]any{"a": "ignored"})
		if err != nil {
			t.Fatal(err)
		}
		if s != "hello" {
			t.Fatalf("expected %q, got %q", "hello", s)
		}
	})

	t.Run("returns fixed value with nil context", func(t *testing.T) {
		s, err := C("hello").String(nil)
		if err != nil {
			t.Fatal(err)
		}
		if s != "hello" {
			t.Fatalf("expected %q, got %q", "hello", s)
		}
	})
}

func TestFirst(t *testing.T) {
	t.Run("returns first present", func(t *testing.T) {
		s, err := First(F("a"), F("b")).String(map[string]any{
			"a": "foo",
			"b": "bar",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "foo" {
			t.Fatalf("expected %q, got %q", "foo", s)
		}
	})

	t.Run("skips missing", func(t *testing.T) {
		s, err := First(F("a"), F("b")).String(map[string]any{
			"b": "bar",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "bar" {
			t.Fatalf("expected %q, got %q", "bar", s)
		}
	})

	t.Run("no matches", func(t *testing.T) {
		s, err := First(F("a"), F("b")).String(map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("required field missing errors", func(t *testing.T) {
		s, err := First(F("a").Required(), F("b")).String(map[string]any{
			"b": "bar",
		})
		if err == nil {
			t.Fatal("expected error for missing required field")
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("required field not reached", func(t *testing.T) {
		s, err := First(F("a"), F("b").Required()).String(map[string]any{
			"a": "foo",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "foo" {
			t.Fatalf("expected %q, got %q", "foo", s)
		}
	})
}

func TestTemplate(t *testing.T) {
	t.Run("empty template", func(t *testing.T) {
		s, err := NewTemplate().String(map[string]any{"a": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("nil context", func(t *testing.T) {
		s, err := NewTemplate(F("a")).String(nil)
		if err != nil {
			t.Fatal(err)
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("concatenates fields", func(t *testing.T) {
		s, err := NewTemplate(F("a"), F("b")).String(map[string]any{
			"a": "foo",
			"b": "bar",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "foo bar" {
			t.Fatalf("expected %q, got %q", "foo bar", s)
		}
	})

	t.Run("skips missing optional", func(t *testing.T) {
		s, err := NewTemplate(F("a"), F("b"), F("c")).String(map[string]any{
			"a": "foo",
			"c": "baz",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "foo baz" {
			t.Fatalf("expected %q, got %q", "foo baz", s)
		}
	})

	t.Run("errors on missing required", func(t *testing.T) {
		s, err := NewTemplate(F("a").Required()).String(map[string]any{})
		if err == nil {
			t.Fatal("expected error for missing required field")
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})

	t.Run("mixed constants and fields", func(t *testing.T) {
		s, err := NewTemplate(C("hello"), F("name"), C("!")).String(map[string]any{
			"name": "world",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "hello world !" {
			t.Fatalf("expected %q, got %q", "hello world !", s)
		}
	})

	t.Run("first inside template", func(t *testing.T) {
		s, err := NewTemplate(
			C("val:"),
			First(F("primary"), F("fallback")),
		).String(map[string]any{
			"fallback": "backup",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "val: backup" {
			t.Fatalf("expected %q, got %q", "val: backup", s)
		}
	})

	t.Run("nested template", func(t *testing.T) {
		outer := NewTemplate(C("name:"), NewTemplate(F("first"), F("last")))
		s, err := outer.String(map[string]any{
			"first": "John",
			"last":  "Doe",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "name: John Doe" {
			t.Fatalf("expected %q, got %q", "name: John Doe", s)
		}
	})

	t.Run("nested template with missing optional", func(t *testing.T) {
		inner := NewTemplate(F("first"), F("last"))
		outer := NewTemplate(C("name:"), inner)
		s, err := outer.String(map[string]any{
			"first": "John",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "name: John" {
			t.Fatalf("expected %q, got %q", "name: John", s)
		}
	})

	t.Run("nested template all missing collapses", func(t *testing.T) {
		inner := NewTemplate(F("a"), F("b"))
		outer := NewTemplate(C("prefix"), inner, C("suffix"))
		s, err := outer.String(map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
		if s != "prefix suffix" {
			t.Fatalf("expected %q, got %q", "prefix suffix", s)
		}
	})

	t.Run("deeply nested templates", func(t *testing.T) {
		innermost := NewTemplate(F("x"))
		middle := NewTemplate(C("["), innermost, C("]"))
		outer := NewTemplate(C("result:"), middle)
		s, err := outer.String(map[string]any{
			"x": "val",
		})
		if err != nil {
			t.Fatal(err)
		}
		if s != "result: [ val ]" {
			t.Fatalf("expected %q, got %q", "result: [ val ]", s)
		}
	})

	t.Run("nested template propagates error", func(t *testing.T) {
		inner := NewTemplate(F("a").Required())
		outer := NewTemplate(C("start"), inner)
		s, err := outer.String(map[string]any{})
		if err == nil {
			t.Fatal("expected error from nested required field")
		}
		if s != "" {
			t.Fatalf("expected empty string, got %q", s)
		}
	})
}
