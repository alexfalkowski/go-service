package env

// Environment of the application.
type Environment string

// IsDevelopment environment.
func (e Environment) IsDevelopment() bool {
	return e == "development" || e == "dev"
}
