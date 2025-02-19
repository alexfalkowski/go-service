package env

// Environment of the application.
type Environment string

// Development environment.
const Development = Environment("development")

// IsDevelopment environment.
func (e Environment) IsDevelopment() bool {
	return e == "development" || e == "dev"
}

// String representation of the environment.
func (e Environment) String() string {
	return string(e)
}
