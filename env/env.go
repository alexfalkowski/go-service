package env

// Environment of the application.
type Environment string

// String representation of the environment.
func (e Environment) String() string {
	return string(e)
}
