package d2

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kgatilin/archai/internal/domain"
)

// Writer writes domain.PackageModel structures as D2 diagram files.
type Writer struct{}

// NewWriter creates a new D2 diagram writer.
func NewWriter() *Writer {
	return &Writer{}
}

// Write generates a D2 diagram from a package model and writes it to the output.
func (w *Writer) Write(ctx context.Context, model domain.PackageModel, opts domain.WriteOptions) error {
	// Check context cancellation before starting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Build D2 content
	builder := newD2TextBuilder()
	content := builder.Build(model, opts.PublicOnly)

	// Write to stdout if requested
	if opts.ToStdout {
		fmt.Print(content)
		return nil
	}

	// Determine output path
	outputPath := opts.OutputPath
	if outputPath == "" {
		return fmt.Errorf("output path is required when not writing to stdout")
	}

	// Ensure parent directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating output directory %s: %w", dir, err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing D2 file %s: %w", outputPath, err)
	}

	return nil
}
