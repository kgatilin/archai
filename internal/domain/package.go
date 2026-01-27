package domain

// PackageModel is the aggregate root representing a Go package's structure.
// It contains all the symbols (interfaces, structs, functions, type definitions)
// found in the package, along with their dependencies.
type PackageModel struct {
	// Path is the package path relative to the module root, e.g., "internal/service".
	Path string

	// Name is the package name as declared in source files, e.g., "service".
	Name string

	// Interfaces is the list of interface definitions in this package.
	Interfaces []InterfaceDef

	// Structs is the list of struct definitions in this package.
	Structs []StructDef

	// Functions is the list of package-level function definitions.
	Functions []FunctionDef

	// TypeDefs is the list of type definitions (type aliases) in this package.
	TypeDefs []TypeDef

	// Dependencies is the list of dependencies between symbols.
	Dependencies []Dependency
}

// SourceFiles returns a deduplicated list of all source files in the package.
func (p PackageModel) SourceFiles() []string {
	seen := make(map[string]bool)
	var files []string

	addFile := func(f string) {
		if f != "" && !seen[f] {
			seen[f] = true
			files = append(files, f)
		}
	}

	for _, iface := range p.Interfaces {
		addFile(iface.SourceFile)
	}
	for _, s := range p.Structs {
		addFile(s.SourceFile)
	}
	for _, fn := range p.Functions {
		addFile(fn.SourceFile)
	}
	for _, td := range p.TypeDefs {
		addFile(td.SourceFile)
	}

	return files
}

// ExportedInterfaces returns only the exported interfaces.
func (p PackageModel) ExportedInterfaces() []InterfaceDef {
	var result []InterfaceDef
	for _, iface := range p.Interfaces {
		if iface.IsExported {
			result = append(result, iface)
		}
	}
	return result
}

// ExportedStructs returns only the exported structs.
func (p PackageModel) ExportedStructs() []StructDef {
	var result []StructDef
	for _, s := range p.Structs {
		if s.IsExported {
			result = append(result, s)
		}
	}
	return result
}

// ExportedFunctions returns only the exported functions.
func (p PackageModel) ExportedFunctions() []FunctionDef {
	var result []FunctionDef
	for _, fn := range p.Functions {
		if fn.IsExported {
			result = append(result, fn)
		}
	}
	return result
}

// ExportedTypeDefs returns only the exported type definitions.
func (p PackageModel) ExportedTypeDefs() []TypeDef {
	var result []TypeDef
	for _, td := range p.TypeDefs {
		if td.IsExported {
			result = append(result, td)
		}
	}
	return result
}

// HasExportedSymbols returns true if the package has any exported symbols.
func (p PackageModel) HasExportedSymbols() bool {
	for _, iface := range p.Interfaces {
		if iface.IsExported {
			return true
		}
	}
	for _, s := range p.Structs {
		if s.IsExported {
			return true
		}
	}
	for _, fn := range p.Functions {
		if fn.IsExported {
			return true
		}
	}
	for _, td := range p.TypeDefs {
		if td.IsExported {
			return true
		}
	}
	return false
}
