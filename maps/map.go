package maps

// StringAny a common map with a string key and a value of any.
type StringAny map[string]any

// IsEmpty verifies if the map has an empty length.
func (m StringAny) IsEmpty() bool {
	return len(m) == 0
}
