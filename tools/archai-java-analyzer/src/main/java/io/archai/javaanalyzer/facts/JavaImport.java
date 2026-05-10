package io.archai.javaanalyzer.facts;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * A single import declaration captured per source class. Kind is one of
 * {@code class}, {@code static}, or {@code wildcard}.
 */
public record JavaImport(
        String from,
        @JsonProperty("to_class") String toClass,
        String kind) {
}
