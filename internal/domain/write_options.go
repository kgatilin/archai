package domain

// WriteOptions configures how models are written to output.
type WriteOptions struct {
	// OutputPath is the file path to write to.
	OutputPath string

	// PublicOnly includes only exported symbols when true.
	PublicOnly bool

	// ToStdout writes to stdout instead of file when true.
	ToStdout bool
}
