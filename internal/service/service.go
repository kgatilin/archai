package service

// Service orchestrates diagram generation operations.
// It uses adapters to read from various sources (Go code, D2 files)
// and write to various destinations (D2 files, stdout).
type Service struct {
	// goReader reads package models from Go source code.
	goReader ModelReader

	// d2Reader reads package models from D2 diagram files.
	// This will be nil until US-3 (diagram split operation).
	d2Reader ModelReader

	// d2Writer writes package models as D2 diagram files.
	d2Writer ModelWriter
}
