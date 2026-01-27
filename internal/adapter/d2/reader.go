package d2

import (
	"context"
	"errors"

	"github.com/kgatilin/archai/internal/domain"
)

// Reader reads D2 diagram files and converts them to domain.PackageModel structures.
// This is a placeholder for US-3+ functionality.
type Reader struct{}

// NewReader creates a new D2 diagram reader.
func NewReader() *Reader {
	return &Reader{}
}

// Read parses D2 diagram files and returns package models.
// This functionality is not yet implemented (planned for US-3).
func (r *Reader) Read(ctx context.Context, paths []string) ([]domain.PackageModel, error) {
	return nil, errors.New("d2.Reader: not implemented (planned for US-3)")
}
