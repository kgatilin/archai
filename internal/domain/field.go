package domain

import "fmt"

// FieldDef represents a struct field definition.
type FieldDef struct {
	// Name is the field name.
	Name string

	// Type is the type reference for this field.
	Type TypeRef

	// IsExported indicates if this field is exported (starts with uppercase).
	IsExported bool

	// Tag is the struct tag string, e.g., `json:"name,omitempty"`.
	Tag string
}

// String returns a human-readable representation of the field.
func (f FieldDef) String() string {
	result := fmt.Sprintf("%s %s", f.Name, f.Type.String())
	if f.Tag != "" {
		result += " " + f.Tag
	}
	return result
}

// StringWithVisibility returns the field string with a visibility prefix.
// Exported fields get "+" prefix, unexported get "-" prefix.
func (f FieldDef) StringWithVisibility() string {
	prefix := "-"
	if f.IsExported {
		prefix = "+"
	}
	return prefix + f.Name + " " + f.Type.String()
}
