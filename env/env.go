package env

// Environment represents the runtime environment a service is running in.
//
// This is typically used to drive environment-specific behavior (for example local/dev/stage/prod),
// and is commonly carried in configuration as a simple string value.
type Environment string

// String returns the environment value as a string.
//
// This is a convenience method that preserves the underlying string value without normalization.
func (e Environment) String() string {
	return string(e)
}
