package flags

// IsStringSet the flag for cmd.
func IsStringSet(s *string) bool {
	return s != nil && *s != ""
}

// IsBoolSet the flag for cmd.
func IsBoolSet(b *bool) bool {
	return b != nil && *b
}
