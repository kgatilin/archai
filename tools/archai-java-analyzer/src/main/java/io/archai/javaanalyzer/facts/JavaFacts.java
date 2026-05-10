package io.archai.javaanalyzer.facts;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

/**
 * Top-level JavaFacts document. Schema identifier is fixed at
 * {@code javafacts/v1}; consumers gate on this string.
 */
public record JavaFacts(
        String schema,
        @JsonProperty("src_roots") List<String> srcRoots,
        List<String> packages,
        List<JavaClass> classes,
        List<JavaImport> imports) {

    public static final String SCHEMA = "javafacts/v1";

    public JavaFacts {
        srcRoots = srcRoots == null ? List.of() : List.copyOf(srcRoots);
        packages = packages == null ? List.of() : List.copyOf(packages);
        classes = classes == null ? List.of() : List.copyOf(classes);
        imports = imports == null ? List.of() : List.copyOf(imports);
    }
}
