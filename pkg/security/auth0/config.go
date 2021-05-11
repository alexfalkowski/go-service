package auth0

// Config for Auth0.
type Config struct {
	URL          string `envconfig:"AUTH0_URL" required:"true"`
	ClientID     string `envconfig:"AUTH0_CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"AUTH0_CLIENT_SECRET" required:"true"`
	Audience     string `envconfig:"AUTH0_AUDIENCE" required:"true"`
}
