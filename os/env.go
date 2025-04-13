package os

import "os"

// GetVariable by key.
func GetVariable(key string) string {
	return os.Getenv(key)
}

// SetVariable of value by key.
func SetVariable(key, value string) error {
	return os.Setenv(key, value)
}

// UnsetVariable by key.
func UnsetVariable(key string) error {
	return os.Unsetenv(key)
}

// VariableExists verifies if the env variable is set.
func VariableExists(key string) bool {
	_, ok := os.LookupEnv(key)

	return ok
}
