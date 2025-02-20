package tls

import "crypto/tls"

// NewConfig for tls.
func NewConfig(cfg *Config) (*tls.Config, error) {
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	if !IsEnabled(cfg) || !cfg.HasKeyPair() {
		return config, nil
	}

	cert, err := cfg.GetCert()
	if err != nil {
		return config, err
	}

	key, err := cfg.GetKey()
	if err != nil {
		return config, err
	}

	pair, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return config, err
	}

	config.Certificates = []tls.Certificate{pair}

	return config, nil
}
