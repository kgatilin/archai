package io.archai.javaanalyzer.facts;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * A method invocation observed inside a method body. {@code staticCall}
 * marshals as the JSON key {@code "static"}.
 */
public record JavaCall(
        @JsonProperty("to_class") String toClass,
        @JsonProperty("to_method") String toMethod,
        @JsonProperty("static") boolean staticCall) {
}
