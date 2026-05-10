package io.archai.javaanalyzer;

import io.archai.javaanalyzer.facts.JavaFacts;
import io.archai.javaanalyzer.json.Writer;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;

/**
 * CLI entry point for the JavaParser-based archai analyzer. Reads one or more
 * source roots and writes a JavaFacts JSON document to stdout.
 */
public final class Main {

    private static final String USAGE =
            "Usage: java -jar archai-java-analyzer.jar <srcRoot> [<srcRoot> ...]\n" +
                    "\n" +
                    "Walks the given source roots, parses every .java file, and emits a\n" +
                    "JavaFacts JSON document on stdout. Parse warnings are written to\n" +
                    "stderr (one per line) and do not fail the run.\n";

    private Main() {
    }

    public static void main(String[] args) {
        int exit = run(args);
        if (exit != 0) {
            System.exit(exit);
        }
    }

    static int run(String[] args) {
        if (args.length == 0) {
            System.err.print(USAGE);
            return 2;
        }
        for (String arg : args) {
            if (arg.equals("--help") || arg.equals("-h")) {
                System.out.print(USAGE);
                return 0;
            }
        }

        List<Path> roots = new ArrayList<>();
        for (String arg : args) {
            Path p = Path.of(arg);
            if (!Files.isDirectory(p)) {
                System.err.println("error: not a directory: " + arg);
                return 1;
            }
            roots.add(p);
        }

        Analyzer analyzer = new Analyzer(roots);
        JavaFacts facts;
        try {
            facts = analyzer.analyze();
        } catch (IOException e) {
            System.err.println("error: " + e.getMessage());
            return 1;
        }

        Writer writer = new Writer();
        try {
            writer.writeTo(facts, System.out);
            System.out.println();
        } catch (IOException e) {
            System.err.println("error: " + e.getMessage());
            return 1;
        }

        for (String warning : analyzer.warnings()) {
            System.err.println(warning);
        }
        return 0;
    }
}
