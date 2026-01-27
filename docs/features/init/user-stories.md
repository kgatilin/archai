# User Stories - Iteration 3: Multi-Package Diagrams, Split & Diff

## Overview

This iteration focuses on **multi-package diagram operations**: generating combined diagrams from multiple packages, splitting them into per-package specs, and comparing diagrams against code or other diagrams.

## User Stories

### US-1: Generate Diagrams - Split Mode (Default)

**As a** developer
**I want to** generate diagrams from multiple packages and save each to its `.arch/` folder
**So that** I can document architecture alongside each package's code

**Acceptance Criteria:**
- Accept one or multiple package paths as input
- Create `.arch/` folder in each package if it doesn't exist
- Generate `pub.d2` (public interfaces only) and/or `internal.d2` (all symbols)
- Default: generate both `pub.d2` and `internal.d2`
- Options: `--pub` (only pub.d2), `--internal` (only internal.d2), or both (default)
- Show cross-package connections (e.g., if `pkg/linter` uses `pkg/archspec.Service`, show that dependency)
- Overwrites existing files

**CLI:**
```bash
# Default: generates both pub.d2 and internal.d2 for each package
go-arch-lint diagram generate pkg/archspec pkg/linter internal/d2

# Public only - generates .arch/pub.d2 for each package
go-arch-lint diagram generate pkg/... --pub

# Internal only - generates .arch/internal.d2 for each package
go-arch-lint diagram generate pkg/... --internal

# Glob pattern support
go-arch-lint diagram generate ./internal/...
```

**Output Structure:**
```
pkg/archspec/
├── .arch/
│   ├── pub.d2       # Public interfaces only
│   └── internal.d2  # All symbols (public + unexported)
pkg/linter/
├── .arch/
│   ├── pub.d2
│   └── internal.d2
```

---

### US-2: Generate Combined Diagram from Code

**As a** developer
**I want to** generate a single combined diagram from multiple packages' **actual code**
**So that** I can visualize the current implementation's cross-package architecture in one file

**Acceptance Criteria:**
- Source: Read actual Go code from packages
- Accept multiple package paths as input
- When `--output=<file.d2>` is specified, generate ONE combined diagram
- Public interfaces only (combined diagrams are for high-level overview)
- Show all packages in one diagram with inter-package connections
- Package relationships visible (e.g., which packages depend on which)

**CLI:**
```bash
# Combined diagram from actual code - single output file
go-arch-lint diagram generate pkg/archspec pkg/linter --output=docs/architecture.d2

# All packages under a path into one diagram
go-arch-lint diagram generate ./pkg/... --output=docs/pkg-overview.d2

# All internal packages
go-arch-lint diagram generate ./internal/... --output=docs/internal-overview.d2
```

**Key Differences:**
- vs US-1: No `--output` flag = Split mode (creates `.arch/` in each package), `--output=<file>` = Combined mode
- vs US-4: Reads **code** (actual state), not saved specs (target state)

---

### US-3: Split Combined Diagram into Per-Package Specs

**As a** developer
**I want to** split a large combined diagram into individual per-package `.arch/` specs
**So that** I can distribute a master architecture plan to live alongside each package's code as target specifications

**Acceptance Criteria:**
- Input: Combined D2 diagram file (e.g., master plan)
- Output: Individual `.arch/pub-spec.d2` and `.arch/internal-spec.d2` files in each package
- Files use `-spec.d2` suffix to indicate **target specification** (vs actual state)
- Preserve type definitions and relationships for each package
- Create `.arch/` folder if it doesn't exist
- Overwrites existing `-spec.d2` files
- Option for dry-run (show what would be created without writing)

**CLI:**
```bash
# Split master diagram into per-package specs
go-arch-lint diagram split docs/master-plan.d2

# Dry run - show what would be created
go-arch-lint diagram split docs/master-plan.d2 --dry-run
```

**Output Structure:**
```
pkg/archspec/.arch/
├── pub-spec.d2       # Target specification (from split)
├── internal-spec.d2  # Target specification (from split)

# Later, after running 'diagram generate':
├── pub.d2            # Actual state from code
└── internal.d2       # Actual state from code
```

**Workflow Example:**
```bash
# 1. Create master architecture plan
vim docs/feature-plan.d2

# 2. Split into target specs
go-arch-lint diagram split docs/feature-plan.d2

# 3. Implement the feature
# ... write code ...

# 4. Generate actual state from code
go-arch-lint diagram generate pkg/archspec pkg/linter

# 5. Compare spec (target) vs actual
go-arch-lint diagram diff pkg/archspec pkg/linter
```

---

### US-4: Compose Diagram from Saved Specs

**As a** developer
**I want to** compose a single combined diagram from all saved `.arch/` specification files
**So that** I can see the full target architecture (not current code) across multiple packages

**Acceptance Criteria:**
- Source: Read saved `.arch/*-spec.d2` or `.arch/*.d2` files
- Input: List of package paths (looks for `.arch/pub-spec.d2` or `.arch/pub.d2` in each)
- Output: Single combined D2 diagram
- Compose from saved specs, **NOT from code** (unlike US-2)
- Handle missing `.arch/` folders gracefully (warn and skip)
- Option to specify which files: `--spec` (use `-spec.d2` files), `--code` (use `.d2` files), or auto-detect

**CLI:**
```bash
# Compose from saved specs (target state)
go-arch-lint diagram compose pkg/archspec pkg/linter --output=docs/target-architecture.d2

# Compose from code-generated diagrams (if generated previously)
go-arch-lint diagram compose pkg/... --code --output=docs/current-architecture.d2

# Auto-detect (prefer -spec.d2 if exists, fallback to .d2)
go-arch-lint diagram compose pkg/... --output=docs/architecture.d2
```

**Key Differences:**
- vs US-2: Reads **saved diagrams** (from `.arch/` folders), not Go source code
- `--spec` flag: Use `-spec.d2` files (target architecture)
- `--code` flag: Use `.d2` files (previously generated from code)
- Useful for: Visualizing target architecture, comparing against old versions, documentation

---

### US-5: Diff Spec vs Code

**As a** developer
**I want to** compare saved specifications against actual code
**So that** I can see what has changed or diverged between the two

**Acceptance Criteria:**
- Input: Package paths to compare
- Process:
  1. Read saved `.arch/*-spec.d2` files (specification)
  2. Generate in-memory model from actual code (don't save to file)
  3. Compare spec vs code, showing differences neutrally
- Output: Git-style diff showing what's in spec, what's in code, what's changed
- If no `-spec.d2` files exist, show error
- Exit code: 0 if match, 1 if differences, 2 if error

**CLI:**
```bash
# Compare spec vs code for specific packages
go-arch-lint diagram diff pkg/archspec pkg/linter

# Compare all packages that have .arch/*-spec.d2 files
go-arch-lint diagram diff --all

# Output as markdown
go-arch-lint diagram diff pkg/archspec --format=markdown
```

**Example Output (Git-style):**
```
pkg/archspec:

  In spec, not in code:
    - Service.Split(path string) error
    - Service.Compare(paths []string) CompareResult
    - Validator.ValidateAll() []ValidationError

  In code, not in spec:
    + FileDef.Imports() []PackageModel
    + NewComposer(reader ModelReader, writer ModelWriter) Composer

  Different (exists in both):
    ~ Generator.Generate()
      spec: Generate(path string) error
      code: Generate(paths []string, opts GenerateOptions) ([]GenerateResult, error)

    ~ Service interface
      spec: 3 methods
      code: 4 methods (added Compose)

Summary: 3 in spec only, 2 in code only, 2 modified
Exit code: 1
```

**Files Compared:**
- `.arch/pub-spec.d2` vs actual code (public API)
- `.arch/internal-spec.d2` vs actual code (if exists, includes unexported symbols)

---

### US-6: Diff Between Two Diagrams

**As a** developer
**I want to** compare any two diagrams against each other
**So that** I can see differences between two architecture versions or sources

**Acceptance Criteria:**
- Input: Two diagram sources (can be files or generated from code/specs)
- Sources can be:
  - A D2 file path
  - `code:<packages>` (generate from code in memory)
  - `specs:<packages>` (compose from saved specs in memory)
- Output: Symmetric diff (what's in source, what's in target, what's different)
- Git-style diff format (consistent with US-5)
- Exit code: 0 if match, 1 if differences, 2 if error

**CLI:**
```bash
# Compare two diagram files
go-arch-lint diagram diff --source=docs/old-arch.d2 --target=docs/new-arch.d2

# Compare spec vs code (shorthand for US-5)
go-arch-lint diagram diff --source=specs:pkg/archspec --target=code:pkg/archspec

# Compare saved specs against a master plan
go-arch-lint diagram diff --source=docs/master-plan.d2 --target=specs:pkg/...

# Compare two different code snapshots
go-arch-lint diagram diff --source=code:pkg/... --target=code:internal/...
```

**Example Output:**
```
Comparing source (docs/v1-plan.d2) vs target (docs/v2-plan.d2):

  In source, not in target:
    - pkg/archspec.Generator interface
    - pkg/archspec.ValidateOptions struct

  In target, not in source:
    + pkg/archspec.Service interface
    + pkg/archspec.Composer interface
    + internal/compose package

  Different (exists in both):
    ~ pkg/archspec.ModelReader
      source: 2 methods (Read, ReadAll)
      target: 1 method (Read)

Summary: 2 removed, 3 added, 1 modified
```

---

### US-7: Diagram Diff Output

**As a** developer
**I want to** generate a D2 diagram that visually highlights differences between two sources
**So that** I can see architectural changes at a glance with color-coded additions, removals, and modifications

**Acceptance Criteria:**
- Input: Same as US-6 (two diagram sources to compare)
- Output: D2 diagram file with color-coded differences (alternative to text output)
- Color scheme (git-style):
  - **Red** - In source, not in target (removed)
  - **Green** - In target, not in source (added)
  - **Yellow** - Exists in both but different (modified)
  - **Gray/Default** - Unchanged (optional, for context)
- Visual indicators:
  - Types/interfaces/structs marked with colors
  - Methods/fields within types can have individual colors
  - Legend showing what colors mean
  - Annotations for significant changes (e.g., signature differences)
- Option to show only differences (`--changes-only`) or include unchanged items for context
- Requires `--output=<file>` to save the diagram

**CLI:**
```bash
# Generate diagram diff (instead of text output)
go-arch-lint diagram diff --source=docs/v1.d2 --target=docs/v2.d2 --diagram --output=docs/diff.d2

# Diff spec vs code, as diagram
go-arch-lint diagram diff pkg/archspec --diagram --output=docs/spec-vs-code.d2

# Show only changes (hide unchanged items)
go-arch-lint diagram diff --source=specs:pkg/... --target=code:pkg/... --diagram --changes-only --output=docs/changes.d2

# Include unchanged for context (default when using --diagram)
go-arch-lint diagram diff --source=docs/old.d2 --target=docs/new.d2 --diagram --output=docs/diff.d2

# Text output (default, no --diagram flag)
go-arch-lint diagram diff pkg/archspec
```

**Visual Example (D2 syntax):**
```d2
# Legend
legend: {
  removed: { style.fill: "#ffcccc"; label: "Removed (in source, not in target)" }
  added: { style.fill: "#ccffcc"; label: "Added (in target, not in source)" }
  modified: { style.fill: "#ffffcc"; label: "Modified (exists in both, different)" }
}

pkg.archspec: {
  # Removed interface (red)
  Generator: {
    shape: class
    stereotype: "<<interface>>"
    style.fill: "#ffcccc"
    style.stroke: "#ff0000"

    "+Generate(path string)": error
  }

  # Added interface (green)
  Service: {
    shape: class
    stereotype: "<<interface>>"
    style.fill: "#ccffcc"
    style.stroke: "#00aa00"

    "+Generate(paths []string, opts GenerateOptions)": "([]GenerateResult, error)"
    "+Validate(path string)": ValidationResult
    "+Split(file string)": error
    "+Compare(paths []string)": ComparisonResult
  }

  # Modified interface (yellow)
  ModelReader: {
    shape: class
    stereotype: "<<interface>>"
    style.fill: "#ffffcc"
    style.stroke: "#aaaa00"

    "+Read(path string)": "(PackageModel, error)"  # unchanged
    "+ReadAll(paths []string)": "([]PackageModel, error)" {
      # Method removed in target
      style.fill: "#ffcccc"
      style.stroke: "#ff0000"
    }
  }
}
```

**D2 Color Palette:**
- Removed: `style.fill: "#ffcccc"`, `style.stroke: "#ff0000"` (light red background, red border)
- Added: `style.fill: "#ccffcc"`, `style.stroke: "#00aa00"` (light green background, green border)
- Modified: `style.fill: "#ffffcc"`, `style.stroke: "#aaaa00"` (light yellow background, yellow border)
- Unchanged: Default styling (no special colors)

**Use Cases:**
- Code review: Generate visual diff before/after a refactor
- Documentation: Show what changed in an architecture update
- Planning: Visualize gap between current state and target architecture
- Communication: Share visual diffs with team members

---

### US-8: Documentation and Style Guide

**As a** developer
**I want to** access built-in documentation about D2 diagram conventions and workflow
**So that** I understand how to work with architecture diagrams effectively

**Acceptance Criteria:**
- Command outputs comprehensive guide to stdout or file
- Documentation includes:
  - Overview and ideology (diagram-driven development)
  - Workflow examples (plan → split → implement → diff)
  - File structure and naming conventions (`.arch/`, `pub.d2`, `internal.d2`, `-spec.d2`)
  - D2 syntax conventions specific to this tool
  - Stereotype reference (what each stereotype means)
  - Color scheme and DDD semantics
  - Factory conventions (dependencies as fields, return type)
  - File container conventions
  - Method signature syntax
  - Common patterns and examples
- Option to output to file for offline reference
- Option to output specific sections (e.g., just stereotypes, just colors)

**CLI:**
```bash
# Output full guide to stdout
go-arch-lint diagram docs

# Save to file
go-arch-lint diagram docs --output=docs/DIAGRAM_GUIDE.md

# Show specific section
go-arch-lint diagram docs --section=stereotypes
go-arch-lint diagram docs --section=colors
go-arch-lint diagram docs --section=workflow
```

**Documentation Sections:**

1. **Overview & Ideology**
   - Diagram-driven development principles
   - Diagrams as source of truth
   - Why D2 for architecture specs

2. **Workflow Guide**
   - Complete workflow examples
   - When to use each command
   - Integration with development process

3. **File Structure**
   - `.arch/` folder purpose
   - `pub.d2` vs `internal.d2`
   - `-spec.d2` suffix for target specifications
   - Where files live (with code, not separate docs folder)

4. **D2 Conventions**
   - File containers (grouping by source file)
   - Package containers
   - Method signature syntax
   - Factory conventions

5. **Stereotypes Reference**
   - Complete list with descriptions
   - When to use each stereotype
   - Examples for each type

6. **Color Scheme (DDD-Inspired)**
   - Color palette table with hex codes
   - DDD concept mapping
   - When to use each color
   - Consistency guidelines

7. **Common Patterns**
   - Interface definitions
   - Factory functions
   - Dependency injection patterns
   - Cross-package relationships

**Example Output (excerpt):**
```markdown
# D2 Architecture Diagrams - Style Guide

## Overview

go-arch-lint supports diagram-driven development through D2 architecture
specifications. Diagrams become the source of truth for package architecture.

## Workflow

1. **Plan** - Create master architecture plan (docs/feature-plan.d2)
2. **Split** - Distribute into per-package specs (.arch/pub-spec.d2)
3. **Implement** - Write Go code matching the spec
4. **Generate** - Create diagrams from actual code (.arch/pub.d2)
5. **Diff** - Compare spec vs code to track progress

## Stereotypes Reference

| Stereotype | Type | Usage | Example |
|------------|------|-------|---------|
| `<<interface>>` | Interface | Public interface | `Service <<interface>>` |
| `<<struct>>` | Struct | General struct | `config <<struct>>` |
| `<<factory>>` | Function | Creation function | `NewService <<factory>>` |
| `<<service>>` | Struct | Service implementation | `service <<service>>` |
| `<<aggregate>>` | Struct | DDD aggregate root | `PackageModel <<aggregate>>` |
| `<<value>>` | Struct | Value object | `GenerateOptions <<value>>` |

## Color Scheme (DDD-Inspired)

| Color | Hex | DDD Concept | Use Case |
|-------|-----|-------------|----------|
| Blue | #e8f4fc | Aggregate/Entity | Core domain models |
| Purple | #f0e8fc | Service/Repository | Behavioral contracts |
| Green | #e8fce8 | Factory | Creation points |
| Gray | #f8f8f8 | Value Object | Options, results |

...
```

---

## Non-Functional Requirements

### NFR-1: In-Memory Model Generation

For diff operations, diagrams should be generated in-memory without writing to files:
- Generate `PackageModel` from Go code
- Generate `PackageModel` from D2 specs
- Compare models directly, not file contents

### NFR-2: Consistent Output Format

All diff commands should support two output modes:
- **Text** (default) - Plain text to stdout with git-style diff format
- **Diagram** (with `--diagram --output=<file>`) - D2 file with color-coded differences

### NFR-3: Graceful Handling

- Missing `.arch/` folders: warn and continue
- Invalid D2 syntax: error with clear message and location
- Partial matches: show what could be compared

---

## Command Summary

| Command | Source | Description |
|---------|--------|-------------|
| `diagram generate <paths>` | Code | Split mode: Generate `.arch/pub.d2` + `.arch/internal.d2` per package from code |
| `diagram generate <paths> --pub` | Code | Split mode: Generate `.arch/pub.d2` only from code |
| `diagram generate <paths> --internal` | Code | Split mode: Generate `.arch/internal.d2` only from code |
| `diagram generate <paths> --output=<file>` | Code | Combined mode: Generate single diagram file from code |
| `diagram split <file>` | D2 file | Split combined diagram into per-package `-spec.d2` files |
| `diagram compose <paths> --output=<file>` | Specs | Compose single diagram from saved `.arch/` spec files |
| `diagram diff <paths>` | Specs vs Code | Compare saved specs vs actual code (text output) |
| `diagram diff --source=X --target=Y` | Any | Compare two diagram sources (text output) |
| `diagram diff --source=X --target=Y --diagram --output=<file>` | Any | Generate diff as color-coded D2 diagram |
| `diagram docs` | Built-in | Output complete style guide and workflow documentation |

---

## Implementation Notes

### Model Layer (Internal)

The diff operations require in-memory model generation:

```go
// Generate model from code (don't write file)
codeModel := reader.ReadPackages(paths...)

// Generate model from saved specs
specModel := compose.FromSpecs(paths...)

// Compare models
diff := compare.Models(specModel, codeModel)
```

### Diagram Sources Abstraction

```go
type DiagramSource interface {
    Load() ([]PackageModel, error)
}

type FileSource struct { Path string }
type CodeSource struct { Paths []string }
type SpecsSource struct { Paths []string }
```

This allows uniform handling of diff operations regardless of where the diagrams come from.
