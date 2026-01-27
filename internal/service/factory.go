package service

// NewService creates a new diagram service with the given adapters.
func NewService(goReader ModelReader, d2Writer ModelWriter) *Service {
	return &Service{
		goReader: goReader,
		d2Writer: d2Writer,
	}
}
