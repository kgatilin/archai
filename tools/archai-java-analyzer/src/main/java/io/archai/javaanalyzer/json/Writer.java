package io.archai.javaanalyzer.json;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.MapperFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import io.archai.javaanalyzer.facts.JavaFacts;

import java.io.IOException;
import java.io.OutputStream;

/**
 * Jackson-backed JSON writer. Output is pretty-printed with 2-space indent and
 * preserves declaration order of record components — top-level field order is
 * controlled by the record signatures, not by alphabetical sorting.
 */
public final class Writer {

    private final ObjectMapper mapper;

    public Writer() {
        this.mapper = new ObjectMapper()
                .configure(MapperFeature.SORT_PROPERTIES_ALPHABETICALLY, false)
                .configure(SerializationFeature.INDENT_OUTPUT, true)
                .configure(JsonGenerator.Feature.WRITE_BIGDECIMAL_AS_PLAIN, true);
    }

    public String write(JavaFacts facts) throws IOException {
        return mapper.writerWithDefaultPrettyPrinter().writeValueAsString(facts);
    }

    public void writeTo(JavaFacts facts, OutputStream out) throws IOException {
        mapper.writerWithDefaultPrettyPrinter().writeValue(out, facts);
    }
}
