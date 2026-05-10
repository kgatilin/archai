package io.archai.javaanalyzer.facts;

import java.util.List;

/**
 * A declared field on a class, interface, enum, or record component. Modifiers
 * are alphabetically sorted; {@code null} lists are normalised to empty.
 */
public record JavaField(String name, String type, List<String> modifiers) {

    public JavaField {
        modifiers = modifiers == null ? List.of() : List.copyOf(modifiers);
    }
}
