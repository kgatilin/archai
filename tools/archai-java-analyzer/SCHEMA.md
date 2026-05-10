# JavaFacts JSON — v1 schema

Version string: `"javafacts/v1"`. The Java side emits AST-shaped facts only — no archai-domain concepts (no stereotypes, no layers). Domain mapping happens Go-side in `internal/adapter/java/` (issue #102).

## Top-level

```json
{
  "schema": "javafacts/v1",
  "src_roots": ["src/main/java"],
  "packages": ["com.foo.bar"],
  "classes": [ /* JavaClass[] */ ],
  "imports": [ /* JavaImport[] */ ]
}
```

| Field       | Type        | Notes |
|-------------|-------------|-------|
| `schema`    | string      | Always `"javafacts/v1"` for v1 output. |
| `src_roots` | string[]    | The input roots, normalized to absolute paths. |
| `packages`  | string[]    | Sorted, deduplicated FQN packages observed across all parsed classes. |
| `classes`   | JavaClass[] | Sorted alphabetically by `fqn`. |
| `imports`   | JavaImport[]| Sorted by `(from, to_class, kind)`. |

Empty arrays are emitted as `[]`, never omitted.

## JavaClass

```json
{
  "fqn": "com.foo.bar.Foo",
  "package": "com.foo.bar",
  "name": "Foo",
  "kind": "class|interface|enum|record|annotation",
  "modifiers": ["public", "abstract"],
  "type_parameters": ["T extends Number"],
  "extends": "com.foo.bar.Base",
  "implements": ["com.foo.bar.Runnable"],
  "permits": [],
  "source_file": "src/main/java/com/foo/bar/Foo.java",
  "doc": "...",
  "fields": [ /* JavaField[] */ ],
  "methods": [ /* JavaMethod[] */ ],
  "annotations": [ /* JavaAnnotation[] */ ]
}
```

| Field             | Type            | Notes |
|-------------------|-----------------|-------|
| `fqn`             | string          | Fully-qualified type name. |
| `package`         | string          | May be empty for default package. |
| `name`            | string          | Simple type name. |
| `kind`            | enum string     | `class`, `interface`, `enum`, `record`, or `annotation`. |
| `modifiers`       | string[]        | Sorted alphabetically. Includes `public`, `abstract`, `final`, `sealed`, `non-sealed`, `static`, etc. |
| `type_parameters` | string[]        | Source-order. Generic param strings as written in source. |
| `extends`         | string \| null  | Best-effort FQN of the extended type, or `null`. |
| `implements`      | string[]        | Best-effort FQNs of implemented interfaces. |
| `permits`         | string[]        | FQNs from a `permits` clause (sealed types). Empty for non-sealed. |
| `source_file`     | string          | Path relative to src root (forward slashes). |
| `doc`             | string          | Javadoc text, or empty. |
| `fields`          | JavaField[]     | Source order. For records, components surface as fields. For enums, constants surface as field-shaped entries. |
| `methods`         | JavaMethod[]    | Source order. Constructors recorded with `returns: "void"`. |
| `annotations`     | JavaAnnotation[]| Source order. |

## JavaField

```json
{ "name": "x", "type": "int", "modifiers": ["private", "final"] }
```

## JavaMethod

```json
{
  "name": "m",
  "modifiers": ["public"],
  "type_parameters": [],
  "params": [ { "name": "a", "type": "java.lang.String" } ],
  "returns": "void",
  "throws": ["java.io.IOException"],
  "calls": [ /* JavaCall[] */ ]
}
```

## JavaCall

```json
{ "to_class": "com.foo.bar.Bar", "to_method": "baz", "static": false }
```

`to_class` is best-effort: when `MethodCallExpr.resolve()` succeeds, it's the declaring type FQN; otherwise the scope text. `static` reflects resolved info when available, else `false`.

## JavaImport

```json
{ "from": "com.foo.bar.Foo", "to_class": "com.foo.bar.Bar", "kind": "class" }
```

`kind` is `class` (default), `static`, or `wildcard`.

## JavaAnnotation

```json
{ "fqn": "java.lang.Override", "args": [] }
```

`fqn` is best-effort resolved.

## Determinism

Output is fully deterministic given identical inputs:

- Classes sorted by `fqn`.
- Packages sorted, deduplicated.
- Imports sorted by `(from, to_class, kind)`.
- Modifiers and `throws` sorted alphabetically within their lists.
- Fields, methods, annotations preserve source order.

## Build choice

The sub-project uses **Maven**, not Gradle. Reasons: parity with the rest of the Anthropic-internal stack we integrate with; simpler CI integration; smaller dependency footprint than Gradle's Kotlin DSL. Single-module project; the shade plugin produces the fat-jar.
