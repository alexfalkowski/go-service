package id

// Default generator.
var Default = &UUID{}

// Generator to generate an identifier.
type Generator interface {
	// Generate an identifier.
	Generate() string
}

// NewGenerator from config.
func NewGenerator(config *Config) Generator {
	if !IsEnabled(config) {
		return Default
	}

	switch config.Kind {
	case "uuid":
		return Default
	case "ksuid":
		return &KSUID{}
	case "nanoid":
		return &NanoID{}
	case "ulid":
		return &ULID{}
	case "xid":
		return &XID{}
	default:
		return Default
	}
}
