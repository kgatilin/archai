package io.archai.javaanalyzer.facts;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

/**
 * A top-level type declaration: class, interface, enum, record, or annotation.
 * The Jackson property name {@code implements} is used for {@code implementsList}
 * to keep the JSON schema natural while avoiding the Java keyword.
 */
public record JavaClass(
        String fqn,
        @JsonProperty("package") String pkg,
        String name,
        String kind,
        List<String> modifiers,
        @JsonProperty("type_parameters") List<String> typeParameters,
        @JsonProperty("extends") String extendsClass,
        @JsonProperty("implements") List<String> implementsList,
        List<String> permits,
        @JsonProperty("source_file") String sourceFile,
        String doc,
        List<JavaField> fields,
        List<JavaMethod> methods,
        List<JavaAnnotation> annotations) {

    public JavaClass {
        modifiers = modifiers == null ? List.of() : List.copyOf(modifiers);
        typeParameters = typeParameters == null ? List.of() : List.copyOf(typeParameters);
        implementsList = implementsList == null ? List.of() : List.copyOf(implementsList);
        permits = permits == null ? List.of() : List.copyOf(permits);
        fields = fields == null ? List.of() : List.copyOf(fields);
        methods = methods == null ? List.of() : List.copyOf(methods);
        annotations = annotations == null ? List.of() : List.copyOf(annotations);
    }
}
