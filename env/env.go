package env

// Development environment.
const Development = Environment("development")

// Environment of the application.
type Environment string

// IsDevelopment environment.
func (e Environment) IsDevelopment() bool {
	return e == "development" || e == "dev"
}

// String representation of the environment.
func (e Environment) String() string {
	return string(e)
}
