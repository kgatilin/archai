package d2

import "github.com/kgatilin/archai/internal/domain"

// D2TextBuilderForTest wraps d2TextBuilder for testing.
type D2TextBuilderForTest struct {
	builder *d2TextBuilder
}

// NewD2TextBuilderForTest creates a builder for testing purposes.
func NewD2TextBuilderForTest() *D2TextBuilderForTest {
	return &D2TextBuilderForTest{
		builder: newD2TextBuilder(),
	}
}

// Build generates D2 content from a package model.
func (b *D2TextBuilderForTest) Build(pkg domain.PackageModel, publicOnly bool) string {
	return b.builder.Build(pkg, publicOnly)
}
