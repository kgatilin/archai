package io.archai.javaanalyzer.facts;

/**
 * A single parameter on a method. Type is the resolved FQN where available,
 * otherwise the string as written in source.
 */
public record JavaParam(String name, String type) {
}
