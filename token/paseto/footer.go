package paseto

import (
	"github.com/alexfalkowski/go-service/v3/encoding/json"
	"github.com/alexfalkowski/go-service/v3/strings"
	"github.com/alexfalkowski/go-service/v3/token/errors"
)

type footer struct {
	KeyID string `json:"kid"`
}

func encodeFooter(footer *footer) ([]byte, error) {
	return json.Marshal(footer)
}

func parseFooter(raw []byte) (*footer, error) {
	footer := &footer{}
	if err := json.Unmarshal(raw, footer); err != nil {
		return nil, errors.ErrInvalidKeyID
	}
	if strings.IsEmpty(footer.KeyID) {
		return nil, errors.ErrInvalidKeyID
	}

	return footer, nil
}
