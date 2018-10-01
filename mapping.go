package surl

// Mapping defines algorithm mapping between decimal and string.
type Mapping interface {
	// Itoa decimal to hex string
	Itoa(int64) string
	// Atoi hex string to decimal
	Atoi(string) (int64, error)
}

var (
	// Hex62 is a basic mapping strategy
	Hex62 = &hex62{}
)
