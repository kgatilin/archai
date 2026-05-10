# archai-java-analyzer

JavaParser-based source walker that emits **JavaFacts** JSON for archai. The Java side is intentionally semantically dumb — it dumps what the AST tells it. archai's domain mapping (stereotypes, layers, etc.) lives Go-side in `internal/adapter/java/` (see issue #102).

## Build

Requires JDK 17+ and Maven.

```
mvn -B -ntp package
```

Produces a fat-jar at `target/archai-java-analyzer-0.1.0.jar`.

From the archai repo root, the convenience target also copies it to `dist/`:

```
make build-java-analyzer
```

## Run

```
java -jar target/archai-java-analyzer-0.1.0.jar <srcRoot> [<srcRoot> ...] > facts.json
```

- Walks each source root recursively for `.java` files.
- Emits a single JavaFacts JSON document on stdout.
- Parse warnings (e.g. unresolved external symbols) are written to stderr, one per line — non-fatal.

Exit codes:
- `0` — success
- `1` — fatal error (unreadable directory, IO error)
- `2` — invalid arguments

## Schema

See [`SCHEMA.md`](./SCHEMA.md) for the JavaFacts JSON v1 shape.

## License

Apache-2.0.
