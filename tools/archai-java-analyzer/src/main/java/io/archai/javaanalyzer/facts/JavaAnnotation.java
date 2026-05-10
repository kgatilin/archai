package io.archai.javaanalyzer.facts;

import java.util.List;

/**
 * A single annotation. {@code args} contains the source-text representation of
 * each argument expression (or pair string for {@code key=value}).
 */
public record JavaAnnotation(String fqn, List<String> args) {

    public JavaAnnotation {
        args = args == null ? List.of() : List.copyOf(args);
    }
}
