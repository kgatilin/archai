// Package service provides business operations for diagram generation.
package service

import (
	"context"

	"github.com/kgatilin/archai/internal/domain"
)

// ModelReader reads package models from a source (code or diagrams).
// This interface is implemented by:
// - adapter/golang.Reader (reads Go source code)
// - adapter/d2.Reader (reads D2 files, planned for US-3)
type ModelReader interface {
	Read(ctx context.Context, paths []string) ([]domain.PackageModel, error)
}

// ModelWriter writes package models to a destination.
// This interface is implemented by:
// - adapter/d2.Writer (writes D2 diagram files)
type ModelWriter interface {
	Write(ctx context.Context, model domain.PackageModel, opts domain.WriteOptions) error
}
