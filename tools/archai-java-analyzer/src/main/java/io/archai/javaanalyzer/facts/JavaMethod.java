package io.archai.javaanalyzer.facts;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

/**
 * A declared method (including constructors). {@code throwsList} marshals as
 * {@code "throws"} in JSON to match the schema and avoid the Java keyword
 * collision in source.
 */
public record JavaMethod(
        String name,
        List<String> modifiers,
        @JsonProperty("type_parameters") List<String> typeParameters,
        List<JavaParam> params,
        String returns,
        @JsonProperty("throws") List<String> throwsList,
        List<JavaCall> calls,
        List<JavaAnnotation> annotations) {

    public JavaMethod {
        modifiers = modifiers == null ? List.of() : List.copyOf(modifiers);
        typeParameters = typeParameters == null ? List.of() : List.copyOf(typeParameters);
        params = params == null ? List.of() : List.copyOf(params);
        throwsList = throwsList == null ? List.of() : List.copyOf(throwsList);
        calls = calls == null ? List.of() : List.copyOf(calls);
        annotations = annotations == null ? List.of() : List.copyOf(annotations);
    }
}
