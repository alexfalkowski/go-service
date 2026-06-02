package ssh

import (
	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	token "github.com/alexfalkowski/go-service/v2/token/errors"
)

type claims struct {
	Version   string `json:"ver"`
	KeyID     string `json:"kid"`
	Audience  string `json:"aud"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func parseClaims(tkn string) (*claims, []byte, string, error) {
	rawClaims, rawSignature, ok := strings.Cut(tkn, tokenSeparator)
	if !ok || strings.IsEmpty(rawClaims) || strings.IsEmpty(rawSignature) || strings.Contains(rawSignature, tokenSeparator) {
		return nil, nil, strings.Empty, crypto.ErrInvalidMatch
	}

	encoded, err := base64.Decode(rawClaims)
	if err != nil {
		return nil, nil, strings.Empty, crypto.ErrInvalidMatch
	}

	c := &claims{}
	if err := json.Unmarshal(encoded, c); err != nil {
		return nil, nil, strings.Empty, crypto.ErrInvalidMatch
	}

	if strings.IsEmpty(c.KeyID) {
		return nil, nil, strings.Empty, crypto.ErrInvalidMatch
	}

	return c, encoded, rawSignature, nil
}

func validateClaims(c *claims, aud string, now int64, maxLifetime time.Duration) error {
	if c.Audience != aud {
		return token.ErrInvalidAudience
	}
	if c.Version != tokenVersion {
		return crypto.ErrInvalidMatch
	}
	invalidIssuedAt := c.IssuedAt <= 0 || c.IssuedAt > now
	invalidExpiration := c.ExpiresAt <= now || c.ExpiresAt <= c.IssuedAt
	if invalidIssuedAt || invalidExpiration {
		return token.ErrInvalidTime
	}
	if c.ExpiresAt-c.IssuedAt > maxLifetime.Duration().Nanoseconds() {
		return token.ErrInvalidTime
	}

	return nil
}
