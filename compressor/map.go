package compressor

type configs map[string]Compressor

// Map of compressor.
type Map struct {
	configs configs
}

// NewMap for compressor.
func NewMap() *Map {
	f := &Map{
		configs: configs{
			"zstd":   NewZstd(),
			"s2":     NewS2(),
			"snappy": NewSnappy(),
			"none":   NewNone(),
		},
	}

	return f
}

// Register kind and compressor.
func (f *Map) Register(kind string, c Compressor) {
	f.configs[kind] = c
}

// Get from kind.
func (f *Map) Get(kind string) Compressor {
	c, ok := f.configs[kind]
	if !ok {
		return f.configs["none"]
	}

	return c
}
