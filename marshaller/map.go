package marshaller

type configs map[string]Marshaller

// Map of marshaller.
type Map struct {
	configs configs
}

// NewMap for marshaller.
func NewMap() *Map {
	f := &Map{
		configs: configs{
			"json":     NewJSON(),
			"yaml":     NewYAML(),
			"yml":      NewYAML(),
			"toml":     NewTOML(),
			"proto":    NewProto(),
			"protobuf": NewProto(),
			"gob":      NewGOB(),
		},
	}

	return f
}

// Register kind and marshaller.
func (f *Map) Register(kind string, m Marshaller) {
	f.configs[kind] = m
}

// Get from kind.
func (f *Map) Get(kind string) Marshaller {
	return f.configs[kind]
}
