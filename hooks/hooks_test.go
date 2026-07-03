package hooks_test

import (
	"strconv"
	"testing"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	webhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	gen := hooks.NewGenerator(rand.NewGenerator(rand.NewReader()))

	secret, err := gen.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, secret)
}

func TestNewHookReturnsSourceError(t *testing.T) {
	_, err := hooks.NewHook(test.FS, newConfig(test.FilePath("secrets/none")))
	require.Error(t, err)
}

func TestNewHookReturnsInvalidSecretError(t *testing.T) {
	_, err := hooks.NewHook(test.FS, newConfig(test.FilePath("secrets/redis")))
	require.Error(t, err)
}

func TestNewHookSignsWithActiveSecret(t *testing.T) {
	hook := newRotatingHook(t)
	payload := []byte("body")
	timestamp := time.Now()

	signature, err := hook.Sign("id-1", timestamp, payload)
	require.NoError(t, err)

	headers := signedHeaders("id-1", timestamp, signature)
	require.NoError(t, newSingleHook(t, webhookSecret("dGVzdA==")).Verify(payload, headers))
	require.Error(t, newSingleHook(t, webhookSecret("b3RoZXI=")).Verify(payload, headers))
}

func TestNewHookVerifiesTrustedSecrets(t *testing.T) {
	hook := newRotatingHook(t)
	payload := []byte("body")

	tests := map[string]string{
		"active":   webhookSecret("dGVzdA=="),
		"previous": webhookSecret("b3RoZXI="),
	}
	for name, secret := range tests {
		t.Run(name, func(t *testing.T) {
			signature, timestamp := sign(t, newSingleHook(t, secret), payload)

			require.NoError(t, hook.Verify(payload, signedHeaders("id-1", timestamp, signature)))
		})
	}

	signature, timestamp := sign(t, newSingleHook(t, webhookSecret("YmFk")), payload)
	require.Error(t, hook.Verify(payload, signedHeaders("id-1", timestamp, signature)))
}

func TestNewHookReturnsInvalidConfigForMissingActiveSecret(t *testing.T) {
	_, err := hooks.NewHook(test.FS, &hooks.Config{
		Key: "missing",
		Secrets: hooks.Secrets{
			"current": webhookSecret("dGVzdA=="),
		},
	})
	require.ErrorIs(t, err, hooks.ErrInvalidConfig)
}

func TestNewHookRejectsEmptySecret(t *testing.T) {
	_, err := hooks.NewHook(test.FS, newConfig(""))
	require.ErrorIs(t, err, hooks.ErrEmptySecret)
}

func TestNewHookRejectsEmptyEnvSecret(t *testing.T) {
	t.Setenv("EMPTY_WEBHOOK_VALUE", "")

	source := "env:" + "EMPTY_WEBHOOK_VALUE"
	_, err := hooks.NewHook(test.FS, newConfig(source))
	require.ErrorIs(t, err, hooks.ErrEmptySecret)
}

func TestNewHookReturnsNilWhenDisabled(t *testing.T) {
	h, err := hooks.NewHook(nil, nil)
	require.NoError(t, err)
	require.Nil(t, h)
}

func newRotatingHook(t *testing.T) *hooks.Hook {
	t.Helper()

	hook, err := hooks.NewHook(test.FS, &hooks.Config{
		Key: "current",
		Secrets: hooks.Secrets{
			"current":  webhookSecret("dGVzdA=="),
			"previous": webhookSecret("b3RoZXI="),
		},
	})
	require.NoError(t, err)

	return hook
}

func newSingleHook(t *testing.T, secret string) *hooks.Hook {
	t.Helper()

	hook, err := hooks.NewHook(test.FS, newConfig(secret))
	require.NoError(t, err)

	return hook
}

func newConfig(secret string) *hooks.Config {
	return &hooks.Config{
		Key: "current",
		Secrets: hooks.Secrets{
			"current": secret,
		},
	}
}

func sign(t *testing.T, hook *hooks.Hook, payload []byte) (string, time.Time) {
	t.Helper()

	timestamp := time.Now()
	signature, err := hook.Sign("id-1", timestamp, payload)
	require.NoError(t, err)

	return signature, timestamp
}

func signedHeaders(id string, timestamp time.Time, signature string) http.Header {
	header := http.Header{}
	header.Set(webhooks.HeaderWebhookID, id)
	header.Set(webhooks.HeaderWebhookSignature, signature)
	header.Set(webhooks.HeaderWebhookTimestamp, strconv.FormatInt(timestamp.Unix(), 10))

	return header
}

func webhookSecret(secret string) string {
	return "whsec_" + secret
}
