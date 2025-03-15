# csf

csf (context-based string formatting) is a utility library for generating templated strings using named key-value pairs provided by a map. I have most commonly used it when I have some "blob" of data (e.g. deserialized JSON) and I want to generate a human-facing string representation depending on the data included in the blob.

## Example 

```go
input := map[string]any{ // usually from deserialized files or database
    "legal_name":     "John Smith",
    "preferred_name": "Johnny Apple",
    "email":          "john-smith@example.com",
}

st := csf.NewTemplate(
    csf.C("Author:"), // Constant string to add label
    csf.First( // First will return the first non-empty value in the list
        csf.F("preferred_name"),        // Preferred name has priority, but is optional
        csf.F("legal_name").Required(), // Required will generate error if neither is present
    ),
    csf.F("email").Formatter(func(v any) string { // Add email if present, but not required
        return fmt.Sprintf("<%s>", v) // Custom formatter to wrap email in angle brackets
    }),
)

s, err := st.String(input)
if err != nil {
    panic(err) // handle error
}
fmt.Println(s)

// Outputs `Author: Johnny Apple <john-smith@example.com>`
```

Try editing the input map to remove various fields. When neither name is present, an error is generated. Providing `preferred_name` will always override `legal_name`. Note how removing the email will result in only the name returned, without any trailing whitespace or empty brackets.

## Features

- Fields may be optional (e.g. `csf.F("email")`) or required (e.g. `csf.F("legal_name").Required()`). Required fields will generate an error when missing, whereas optional fields are simply ignored.
- Custom format functions are supported for stringify various types of fields (e.g. `csf.F("email").Formatter(func(v any) string { ... })`). csf additionally provides a few basic ones (`csf.Value`, `csf.Array` and `csf.Const`) out of the box.
- Conditional field logic can be implemented using `csf.First` to return the first non-empty value in the list. This is useful for providing a fallback value when the primary value is not present, or when fields have a mutually exclusive/overriding relationship.
- Constants can be placed in the template using `csf.C("constant")`.

## Motivation

For better or worse I routinely find myself needing to generate human-readable strings from deserialized map blobs. Usually I reach for text/template, or a bunch of `if` statements and parsing code in a `String()` method. This library is an attempt to provide a more structured way to do this, particularly when the data comes in various "shapes" (i.e. variable field inclusion).

## Limitations

Is it fast? Not particularly, and I doubt any faster than the equivalent code rolled out. Is it type-safe? Not really. It is also "one shot" meaning you cannot progressively rewrite the string as you go. Don't expect this to replace "deeply" nested logic.

## Installation

`go get github.com/Cryptkeeper/csf`

See the example included above. For more examples, see the [test cases](csf_test.go).
