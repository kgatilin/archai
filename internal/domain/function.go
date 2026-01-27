package domain

// FunctionDef represents a Go function definition (package-level, no receiver).
type FunctionDef struct {
	// Name is the function name.
	Name string

	// Params is the list of function parameters.
	Params []ParamDef

	// Returns is the list of return types.
	Returns []TypeRef

	// IsExported indicates if this function is exported (starts with uppercase).
	IsExported bool

	// SourceFile is the filename where this function is defined, e.g., "factory.go".
	SourceFile string

	// Doc is the documentation comment for this function.
	Doc string

	// Stereotype is the DDD classification, either detected via heuristics
	// or explicitly set via an archspec annotation.
	Stereotype Stereotype
}

// Signature returns a formatted function signature string.
// Format: "FunctionName(param1 Type1, param2 Type2) (ReturnType1, ReturnType2)"
func (f FunctionDef) Signature() string {
	// Reuse MethodDef.Signature() by creating a temporary MethodDef
	m := MethodDef{
		Name:    f.Name,
		Params:  f.Params,
		Returns: f.Returns,
	}
	return m.Signature()
}

// SignatureWithVisibility returns the function signature with a visibility prefix.
// Exported functions get "+" prefix, unexported get "-" prefix.
func (f FunctionDef) SignatureWithVisibility() string {
	prefix := "-"
	if f.IsExported {
		prefix = "+"
	}
	return prefix + f.Signature()
}
