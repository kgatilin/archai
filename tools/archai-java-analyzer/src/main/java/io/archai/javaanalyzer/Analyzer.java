package io.archai.javaanalyzer;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.List;
import java.util.Optional;
import java.util.TreeSet;
import java.util.stream.Stream;

import com.github.javaparser.StaticJavaParser;
import com.github.javaparser.ast.CompilationUnit;
import com.github.javaparser.ast.ImportDeclaration;
import com.github.javaparser.ast.Modifier;
import com.github.javaparser.ast.NodeList;
import com.github.javaparser.ast.PackageDeclaration;
import com.github.javaparser.ast.body.AnnotationDeclaration;
import com.github.javaparser.ast.body.AnnotationMemberDeclaration;
import com.github.javaparser.ast.body.BodyDeclaration;
import com.github.javaparser.ast.body.ClassOrInterfaceDeclaration;
import com.github.javaparser.ast.body.ConstructorDeclaration;
import com.github.javaparser.ast.body.EnumConstantDeclaration;
import com.github.javaparser.ast.body.EnumDeclaration;
import com.github.javaparser.ast.body.FieldDeclaration;
import com.github.javaparser.ast.body.MethodDeclaration;
import com.github.javaparser.ast.body.Parameter;
import com.github.javaparser.ast.body.RecordDeclaration;
import com.github.javaparser.ast.body.TypeDeclaration;
import com.github.javaparser.ast.body.VariableDeclarator;
import com.github.javaparser.ast.expr.AnnotationExpr;
import com.github.javaparser.ast.expr.MemberValuePair;
import com.github.javaparser.ast.expr.MethodCallExpr;
import com.github.javaparser.ast.expr.NormalAnnotationExpr;
import com.github.javaparser.ast.expr.SingleMemberAnnotationExpr;
import com.github.javaparser.ast.type.ClassOrInterfaceType;
import com.github.javaparser.ast.type.TypeParameter;
import com.github.javaparser.javadoc.Javadoc;
import com.github.javaparser.resolution.UnsolvedSymbolException;
import com.github.javaparser.symbolsolver.JavaSymbolSolver;
import com.github.javaparser.symbolsolver.resolution.typesolvers.CombinedTypeSolver;
import com.github.javaparser.symbolsolver.resolution.typesolvers.JavaParserTypeSolver;
import com.github.javaparser.symbolsolver.resolution.typesolvers.ReflectionTypeSolver;

import io.archai.javaanalyzer.facts.JavaAnnotation;
import io.archai.javaanalyzer.facts.JavaCall;
import io.archai.javaanalyzer.facts.JavaClass;
import io.archai.javaanalyzer.facts.JavaFacts;
import io.archai.javaanalyzer.facts.JavaField;
import io.archai.javaanalyzer.facts.JavaImport;
import io.archai.javaanalyzer.facts.JavaMethod;
import io.archai.javaanalyzer.facts.JavaParam;

/**
 * Walks one or more Java source roots, parses every {@code .java} file with
 * JavaParser (configured with a reflection + per-root symbol solver) and
 * produces a {@link JavaFacts} document. Resolution failures are recorded as
 * one-line warnings and never abort the run.
 */
public final class Analyzer {

    private final List<Path> srcRoots;
    private final List<String> warnings = new ArrayList<>();

    public Analyzer(List<Path> srcRoots) {
        this.srcRoots = List.copyOf(srcRoots);
    }

    public JavaFacts analyze() throws IOException {
        configureSymbolSolver();

        List<JavaClass> classes = new ArrayList<>();
        List<JavaImport> imports = new ArrayList<>();
        TreeSet<String> packages = new TreeSet<>();

        for (Path root : srcRoots) {
            if (!Files.isDirectory(root)) {
                continue;
            }
            List<Path> javaFiles = collectJavaFiles(root);
            for (Path file : javaFiles) {
                CompilationUnit cu;
                try {
                    cu = StaticJavaParser.parse(file);
                } catch (IOException e) {
                    warnings.add("parse error: " + file + ": " + e.getMessage());
                    continue;
                } catch (RuntimeException e) {
                    warnings.add("parse error: " + file + ": " + e.getMessage());
                    continue;
                }

                String packageName = cu.getPackageDeclaration()
                        .map(PackageDeclaration::getNameAsString)
                        .orElse("");
                if (!packageName.isEmpty()) {
                    packages.add(packageName);
                }

                String relSourceFile = relativizeSource(root, file);

                for (TypeDeclaration<?> type : cu.getTypes()) {
                    JavaClass javaClass = buildClass(type, packageName, relSourceFile);
                    classes.add(javaClass);
                    imports.addAll(buildImports(cu, javaClass.fqn()));
                }
            }
        }

        classes.sort(Comparator.comparing(JavaClass::fqn));
        imports.sort(Comparator
                .comparing(JavaImport::from)
                .thenComparing(JavaImport::toClass)
                .thenComparing(JavaImport::kind));

        List<String> srcRootStrings = new ArrayList<>();
        for (Path root : srcRoots) {
            srcRootStrings.add(root.toString());
        }

        return new JavaFacts(
                JavaFacts.SCHEMA,
                srcRootStrings,
                new ArrayList<>(packages),
                classes,
                imports);
    }

    public List<String> warnings() {
        return Collections.unmodifiableList(warnings);
    }

    // ---------------------------------------------------------------------
    // Setup
    // ---------------------------------------------------------------------

    private void configureSymbolSolver() {
        CombinedTypeSolver combined = new CombinedTypeSolver();
        combined.add(new ReflectionTypeSolver());
        for (Path root : srcRoots) {
            if (Files.isDirectory(root)) {
                combined.add(new JavaParserTypeSolver(root));
            }
        }
        JavaSymbolSolver symbolSolver = new JavaSymbolSolver(combined);
        StaticJavaParser.getParserConfiguration().setSymbolResolver(symbolSolver);
    }

    private List<Path> collectJavaFiles(Path root) throws IOException {
        List<Path> result = new ArrayList<>();
        try (Stream<Path> stream = Files.walk(root)) {
            stream.filter(Files::isRegularFile)
                    .filter(p -> p.getFileName().toString().endsWith(".java"))
                    .sorted()
                    .forEach(result::add);
        }
        return result;
    }

    private String relativizeSource(Path root, Path file) {
        Path rel;
        try {
            rel = root.toAbsolutePath().relativize(file.toAbsolutePath());
        } catch (IllegalArgumentException e) {
            rel = file.getFileName();
        }
        return rel.toString().replace('\\', '/');
    }

    // ---------------------------------------------------------------------
    // Class / type extraction
    // ---------------------------------------------------------------------

    private JavaClass buildClass(TypeDeclaration<?> type, String pkg, String sourceFile) {
        String name = type.getNameAsString();
        String fqn = pkg.isEmpty() ? name : pkg + "." + name;
        String kind = kindOf(type);

        List<String> modifiers = sortedModifiers(type.getModifiers());

        List<String> typeParameters = new ArrayList<>();
        if (type instanceof ClassOrInterfaceDeclaration coi) {
            for (TypeParameter tp : coi.getTypeParameters()) {
                typeParameters.add(tp.getNameAsString());
            }
        } else if (type instanceof RecordDeclaration rec) {
            for (TypeParameter tp : rec.getTypeParameters()) {
                typeParameters.add(tp.getNameAsString());
            }
        }

        String extendsClass = null;
        List<String> implementsList = new ArrayList<>();
        List<String> permits = new ArrayList<>();
        if (type instanceof ClassOrInterfaceDeclaration coi) {
            if (coi.isInterface()) {
                if (!coi.getExtendedTypes().isEmpty()) {
                    // interfaces can extend multiple; pick first as "extends" and rest go nowhere
                    // but the schema only has a scalar extends — record first, others as implements? No,
                    // for interfaces the extended types are the parent interfaces. Keep first as extends
                    // and put any extras into implementsList for symmetry.
                    extendsClass = coi.getExtendedTypes(0).getNameAsString();
                    for (int i = 1; i < coi.getExtendedTypes().size(); i++) {
                        implementsList.add(coi.getExtendedTypes(i).getNameAsString());
                    }
                }
            } else {
                if (!coi.getExtendedTypes().isEmpty()) {
                    extendsClass = coi.getExtendedTypes(0).getNameAsString();
                }
                for (ClassOrInterfaceType impl : coi.getImplementedTypes()) {
                    implementsList.add(impl.getNameAsString());
                }
            }
            for (ClassOrInterfaceType permit : coi.getPermittedTypes()) {
                permits.add(permit.getNameAsString());
            }
        } else if (type instanceof RecordDeclaration rec) {
            for (ClassOrInterfaceType impl : rec.getImplementedTypes()) {
                implementsList.add(impl.getNameAsString());
            }
        } else if (type instanceof EnumDeclaration en) {
            for (ClassOrInterfaceType impl : en.getImplementedTypes()) {
                implementsList.add(impl.getNameAsString());
            }
        }

        String doc = type.getJavadoc().map(Javadoc::toText).orElse(null);

        List<JavaField> fields = new ArrayList<>();
        List<JavaMethod> methods = new ArrayList<>();

        // Record components count as fields too.
        if (type instanceof RecordDeclaration rec) {
            for (Parameter component : rec.getParameters()) {
                fields.add(new JavaField(
                        component.getNameAsString(),
                        component.getType().asString(),
                        List.of()));
            }
        }

        // Enum constants render as fields (same shape used in the consumer).
        if (type instanceof EnumDeclaration en) {
            for (EnumConstantDeclaration constant : en.getEntries()) {
                fields.add(new JavaField(
                        constant.getNameAsString(),
                        en.getNameAsString(),
                        List.of()));
            }
        }

        for (BodyDeclaration<?> member : type.getMembers()) {
            if (member instanceof FieldDeclaration fd) {
                List<String> fieldMods = sortedModifiers(fd.getModifiers());
                String typeText = fd.getElementType().asString();
                for (VariableDeclarator var : fd.getVariables()) {
                    fields.add(new JavaField(
                            var.getNameAsString(),
                            typeText,
                            fieldMods));
                }
            } else if (member instanceof MethodDeclaration md) {
                methods.add(buildMethod(md, fqn));
            } else if (member instanceof ConstructorDeclaration cd) {
                methods.add(buildConstructor(cd, fqn));
            } else if (member instanceof AnnotationMemberDeclaration amd) {
                methods.add(buildAnnotationMember(amd, fqn));
            }
        }

        List<JavaAnnotation> annotations = buildAnnotations(type.getAnnotations());

        return new JavaClass(
                fqn,
                pkg,
                name,
                kind,
                modifiers,
                typeParameters,
                extendsClass,
                implementsList,
                permits,
                sourceFile,
                doc,
                fields,
                methods,
                annotations);
    }

    private String kindOf(TypeDeclaration<?> type) {
        if (type instanceof ClassOrInterfaceDeclaration coi) {
            return coi.isInterface() ? "interface" : "class";
        }
        if (type instanceof EnumDeclaration) {
            return "enum";
        }
        if (type instanceof RecordDeclaration) {
            return "record";
        }
        if (type instanceof AnnotationDeclaration) {
            return "annotation";
        }
        return "class";
    }

    // ---------------------------------------------------------------------
    // Method extraction
    // ---------------------------------------------------------------------

    private JavaMethod buildMethod(MethodDeclaration md, String ownerFqn) {
        List<String> modifiers = sortedModifiers(md.getModifiers());
        List<String> typeParameters = new ArrayList<>();
        for (TypeParameter tp : md.getTypeParameters()) {
            typeParameters.add(tp.getNameAsString());
        }
        List<JavaParam> params = buildParams(md.getParameters());
        String returns = md.getType().isVoidType() ? "void" : md.getType().asString();
        List<String> throwsList = sortedThrows(md.getThrownExceptions());
        List<JavaCall> calls = collectCalls(md, ownerFqn, md.getNameAsString());
        List<JavaAnnotation> annotations = buildAnnotations(md.getAnnotations());

        return new JavaMethod(
                md.getNameAsString(),
                modifiers,
                typeParameters,
                params,
                returns,
                throwsList,
                calls,
                annotations);
    }

    private JavaMethod buildConstructor(ConstructorDeclaration cd, String ownerFqn) {
        List<String> modifiers = sortedModifiers(cd.getModifiers());
        List<String> typeParameters = new ArrayList<>();
        for (TypeParameter tp : cd.getTypeParameters()) {
            typeParameters.add(tp.getNameAsString());
        }
        List<JavaParam> params = buildParams(cd.getParameters());
        List<String> throwsList = sortedThrows(cd.getThrownExceptions());
        List<JavaCall> calls = collectCalls(cd, ownerFqn, cd.getNameAsString());
        List<JavaAnnotation> annotations = buildAnnotations(cd.getAnnotations());

        return new JavaMethod(
                cd.getNameAsString(),
                modifiers,
                typeParameters,
                params,
                "void",
                throwsList,
                calls,
                annotations);
    }

    private JavaMethod buildAnnotationMember(AnnotationMemberDeclaration amd, String ownerFqn) {
        List<String> modifiers = sortedModifiers(amd.getModifiers());
        String returns = amd.getType().asString();
        List<JavaAnnotation> annotations = buildAnnotations(amd.getAnnotations());

        return new JavaMethod(
                amd.getNameAsString(),
                modifiers,
                List.of(),
                List.of(),
                returns,
                List.of(),
                List.of(),
                annotations);
    }

    private List<JavaParam> buildParams(NodeList<Parameter> source) {
        List<JavaParam> result = new ArrayList<>();
        for (Parameter p : source) {
            result.add(new JavaParam(p.getNameAsString(), p.getType().asString()));
        }
        return result;
    }

    private List<String> sortedThrows(NodeList<com.github.javaparser.ast.type.ReferenceType> source) {
        TreeSet<String> set = new TreeSet<>();
        for (com.github.javaparser.ast.type.ReferenceType ref : source) {
            set.add(ref.asString());
        }
        return new ArrayList<>(set);
    }

    private List<JavaCall> collectCalls(com.github.javaparser.ast.Node node, String ownerFqn, String methodName) {
        List<JavaCall> calls = new ArrayList<>();
        node.findAll(MethodCallExpr.class).forEach(call -> {
            String toClass;
            String toMethod = call.getNameAsString();
            boolean staticCall = false;
            try {
                var resolved = call.resolve();
                toClass = resolved.declaringType().getQualifiedName();
                staticCall = resolved.isStatic();
            } catch (UnsolvedSymbolException e) {
                toClass = scopeFallback(call);
                warnings.add("unresolved call: " + ownerFqn + "#" + methodName
                        + " -> " + toMethod + ": " + e.getMessage());
            } catch (RuntimeException e) {
                toClass = scopeFallback(call);
                warnings.add("unresolved call: " + ownerFqn + "#" + methodName
                        + " -> " + toMethod + ": " + e.getMessage());
            }
            calls.add(new JavaCall(toClass, toMethod, staticCall));
        });
        return calls;
    }

    private String scopeFallback(MethodCallExpr call) {
        Optional<com.github.javaparser.ast.expr.Expression> scope = call.getScope();
        if (scope.isPresent()) {
            return scope.get().toString();
        }
        return "";
    }

    // ---------------------------------------------------------------------
    // Imports & annotations
    // ---------------------------------------------------------------------

    private List<JavaImport> buildImports(CompilationUnit cu, String fromFqn) {
        List<JavaImport> result = new ArrayList<>();
        for (ImportDeclaration imp : cu.getImports()) {
            String to = imp.getNameAsString();
            String kind;
            if (imp.isStatic()) {
                kind = "static";
            } else if (imp.isAsterisk()) {
                kind = "wildcard";
            } else {
                kind = "class";
            }
            result.add(new JavaImport(fromFqn, to, kind));
        }
        return result;
    }

    private List<JavaAnnotation> buildAnnotations(NodeList<AnnotationExpr> source) {
        List<JavaAnnotation> result = new ArrayList<>();
        for (AnnotationExpr ann : source) {
            String fqn = resolveAnnotationFqn(ann);
            List<String> args = new ArrayList<>();
            if (ann instanceof SingleMemberAnnotationExpr single) {
                args.add(single.getMemberValue().toString());
            } else if (ann instanceof NormalAnnotationExpr normal) {
                for (MemberValuePair pair : normal.getPairs()) {
                    args.add(pair.getNameAsString() + "=" + pair.getValue().toString());
                }
            }
            result.add(new JavaAnnotation(fqn, args));
        }
        return result;
    }

    private String resolveAnnotationFqn(AnnotationExpr ann) {
        try {
            return ann.resolve().getQualifiedName();
        } catch (UnsolvedSymbolException e) {
            warnings.add("unresolved annotation: " + ann.getNameAsString() + ": " + e.getMessage());
            return ann.getNameAsString();
        } catch (RuntimeException e) {
            warnings.add("unresolved annotation: " + ann.getNameAsString() + ": " + e.getMessage());
            return ann.getNameAsString();
        }
    }

    // ---------------------------------------------------------------------
    // Helpers
    // ---------------------------------------------------------------------

    private List<String> sortedModifiers(NodeList<Modifier> source) {
        TreeSet<String> set = new TreeSet<>();
        for (Modifier m : source) {
            set.add(m.getKeyword().asString());
        }
        return new ArrayList<>(set);
    }
}
