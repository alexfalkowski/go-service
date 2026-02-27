package id

// Generator generates identifiers as strings.
//
// Implementations should return a non-empty identifier suitable for use as a unique key.
// The returned value is expected to be stable (not mutated) and safe to use across goroutines.
//
// Generators generally should not return an error; implementations that can fail typically
// handle failures internally (for example by panicking in exceptional cases) or rely on
// injected dependencies that guarantee success.
type Generator interface {
	// Generate returns a newly generated identifier as a string.
	Generate() string
}
