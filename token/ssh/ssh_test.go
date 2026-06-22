package ssh_test

import (
	"testing"

	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	cryptossh "github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"github.com/stretchr/testify/require"
)

func TestConfigIsEnabled(t *testing.T) {
	tests := []struct {
		config  *ssh.Config
		name    string
		enabled bool
	}{
		{name: "nil config"},
		{name: "empty config", config: &ssh.Config{}},
		{name: "with key", config: &ssh.Config{Key: "test"}, enabled: true},
		{name: "with keys", config: &ssh.Config{Keys: ssh.Keys{"test": nil}}, enabled: true},
		{name: "with key and keys", config: test.NewToken("ssh").SSH, enabled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.enabled, tt.config.IsEnabled())
		})
	}
}

func TestConfigRejectsInvalidValues(t *testing.T) {
	valid := test.NewToken("ssh").SSH
	tests := []struct {
		config *ssh.Config
		name   string
	}{
		{
			name: "invalid leeway precision",
			config: &ssh.Config{
				Key:        valid.Key,
				Keys:       valid.Keys,
				Expiration: valid.Expiration,
				Leeway:     time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, test.Validator.Struct(tt.config))
		})
	}

	require.NoError(t, test.Validator.Struct(valid))
}

func TestValid(t *testing.T) {
	token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

	tkn, err := token.Generate(strings.Empty, strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, strings.Empty)
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)

	ssh := ssh.NewToken(nil, nil)
	require.Nil(t, ssh)
}

func TestValidForAudience(t *testing.T) {
	token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

	tkn, err := token.Generate("/service.Method", strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, "/service.Method")
	require.NoError(t, err)
	require.Equal(t, test.UserID.String(), sub)
}

func TestVerifyWithLeeway(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	cfg.Leeway = time.Minute
	token := ssh.NewToken(cfg, test.FS)
	now := time.Now()

	tests := []struct {
		err       error
		issuedAt  time.Time
		expiresAt time.Time
		name      string
	}{
		{
			name:      "issuer clock ahead within leeway",
			issuedAt:  now.Add((30 * time.Second).Duration()),
			expiresAt: now.Add((30 * time.Second).Duration()).Add(cfg.Expiration.Duration()),
		},
		{
			name:      "issuer clock ahead beyond leeway",
			issuedAt:  now.Add((2 * time.Minute).Duration()),
			expiresAt: now.Add((2 * time.Minute).Duration()).Add(cfg.Expiration.Duration()),
			err:       errors.ErrInvalidTime,
		},
		{
			name:      "expired within leeway",
			issuedAt:  now.Add(-(30 * time.Second).Duration()).Add(-cfg.Expiration.Duration()),
			expiresAt: now.Add(-(30 * time.Second).Duration()),
		},
		{
			name:      "expired beyond leeway",
			issuedAt:  now.Add(-(2 * time.Minute).Duration()).Add(-cfg.Expiration.Duration()),
			expiresAt: now.Add(-(2 * time.Minute).Duration()),
			err:       errors.ErrInvalidTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn := signedSSHToken(t, cfg, sshClaims(cfg, tt.issuedAt, tt.expiresAt))
			sub, err := token.Verify(tkn, "/service.Method")
			if tt.err != nil {
				require.Empty(t, sub)
				require.ErrorIs(t, err, tt.err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.UserID.String(), sub)
		})
	}
}

func TestInvalidAudience(t *testing.T) {
	token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

	tkn, err := token.Generate("/service.Method", strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, "/service.Other")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidAudience)
}

func TestInvalidExpired(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	cfg.Expiration = time.Nanosecond
	token := ssh.NewToken(cfg, test.FS)

	tkn, err := token.Generate("/service.Method", strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	time.Sleep(time.Millisecond)

	sub, err := token.Verify(tkn, "/service.Method")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidTime)
}

func TestInvalidLifetimeExceedsConfig(t *testing.T) {
	genCfg := test.NewToken("ssh").SSH
	genCfg.Expiration = time.Hour
	generator := ssh.NewToken(genCfg, test.FS)

	verifyCfg := test.NewToken("ssh").SSH
	verifyCfg.Expiration = time.Minute
	verifier := ssh.NewToken(verifyCfg, test.FS)

	tkn, err := generator.Generate("/service.Method", strings.Empty)
	require.NoError(t, err)

	sub, err := verifier.Verify(tkn, "/service.Method")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidTime)
}

func TestInvalidMissingIssuedAt(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	claims := map[string]any{
		"ver": "v1",
		"kid": cfg.Key,
		"sub": cfg.Key,
		"aud": "/service.Method",
		"exp": time.Now().Add(time.Hour.Duration()).UnixNano(),
	}

	token := ssh.NewToken(cfg, test.FS)
	tkn := signedSSHToken(t, cfg, claims)

	sub, err := token.Verify(tkn, "/service.Method")
	require.Empty(t, sub)
	require.ErrorIs(t, err, errors.ErrInvalidTime)
}

func TestInvalidSignedClaims(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	token := ssh.NewToken(cfg, test.FS)
	now := time.Now()

	tests := []struct {
		err    error
		claims map[string]any
		name   string
	}{
		{
			name: "invalid version",
			err:  crypto.ErrInvalidMatch,
			claims: map[string]any{
				"ver": "v2",
				"kid": cfg.Key,
				"sub": cfg.Key,
				"aud": "/service.Method",
				"iat": now.UnixNano(),
				"exp": now.Add(time.Hour.Duration()).UnixNano(),
			},
		},
		{
			name: "issued at in future",
			err:  errors.ErrInvalidTime,
			claims: map[string]any{
				"ver": "v1",
				"kid": cfg.Key,
				"sub": cfg.Key,
				"aud": "/service.Method",
				"iat": now.Add(time.Minute.Duration()).UnixNano(),
				"exp": now.Add(time.Hour.Duration()).UnixNano(),
			},
		},
		{
			name: "expiration before issued at",
			err:  errors.ErrInvalidTime,
			claims: map[string]any{
				"ver": "v1",
				"kid": cfg.Key,
				"sub": cfg.Key,
				"aud": "/service.Method",
				"iat": now.UnixNano(),
				"exp": now.UnixNano(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn := signedSSHToken(t, cfg, tt.claims)

			sub, err := token.Verify(tkn, "/service.Method")
			require.Empty(t, sub)
			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestInvalidSubjectDiffersFromKeyID(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	token := ssh.NewToken(cfg, test.FS)
	now := time.Now()
	claims := map[string]any{
		"ver": "v1",
		"kid": cfg.Key,
		"sub": "other",
		"aud": "/service.Method",
		"iat": now.UnixNano(),
		"exp": now.Add(time.Hour.Duration()).UnixNano(),
	}

	tkn := signedSSHToken(t, cfg, claims)

	sub, err := token.Verify(tkn, "/service.Method")
	require.Empty(t, sub)
	require.ErrorIs(t, err, crypto.ErrInvalidMatch)
}

func TestInvalidMalformedTokens(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	token := ssh.NewToken(cfg, test.FS)
	now := time.Now()

	tests := []struct {
		err   error
		name  string
		token string
	}{
		{
			name:  "invalid claims base64",
			token: "%%%." + base64.Encode([]byte("signature")),
			err:   crypto.ErrInvalidMatch,
		},
		{
			name:  "invalid claims json",
			token: base64.Encode([]byte("{")) + "." + base64.Encode([]byte("signature")),
			err:   crypto.ErrInvalidMatch,
		},
		{
			name: "missing key id",
			token: encodedSSHToken(t, map[string]any{
				"ver": "v1",
				"sub": cfg.Key,
				"aud": "/service.Method",
				"iat": now.UnixNano(),
				"exp": now.Add(time.Hour.Duration()).UnixNano(),
			}, base64.Encode([]byte("signature"))),
			err: crypto.ErrInvalidMatch,
		},
		{
			name: "invalid signature base64",
			token: encodedSSHToken(t, map[string]any{
				"ver": "v1",
				"kid": cfg.Key,
				"aud": "/service.Method",
				"iat": now.UnixNano(),
				"exp": now.Add(time.Hour.Duration()).UnixNano(),
			}, "%%%"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := token.Verify(tt.token, "/service.Method")
			require.Empty(t, sub)
			require.Error(t, err)

			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}

func TestValidNameWithDash(t *testing.T) {
	cfg := test.NewToken("ssh").SSH
	cfg.Keys["test-user"] = cfg.Keys.Get(cfg.Key)
	cfg.Key = "test-user"

	token := ssh.NewToken(cfg, test.FS)

	tkn, err := token.Generate(strings.Empty, strings.Empty)
	require.NoError(t, err)
	require.NotEmpty(t, tkn)

	sub, err := token.Verify(tkn, strings.Empty)
	require.NoError(t, err)
	require.Equal(t, "test-user", sub)
}

func TestInvalidPrivateKey(t *testing.T) {
	token := ssh.NewToken(&ssh.Config{
		Expiration: time.Hour,
		Key:        "test",
		Keys: ssh.Keys{
			"test": &ssh.Key{
				Config: test.NewSSH("secrets/ssh_public", "secrets/none"),
			},
		},
	}, test.FS)

	_, err := token.Generate(strings.Empty, strings.Empty)
	require.Error(t, err)
}

func TestInvalidTokenShapes(t *testing.T) {
	for _, tkn := range []string{strings.Empty, "none-", "test-", "test-bob"} {
		t.Run(tkn, func(t *testing.T) {
			token := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)

			sub, err := token.Verify(tkn, strings.Empty)
			require.Error(t, err)
			require.Empty(t, sub)
		})
	}
}

func TestInvalidPublicKey(t *testing.T) {
	valid := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
	tkn, err := valid.Generate(strings.Empty, strings.Empty)
	require.NoError(t, err)

	token := ssh.NewToken(&ssh.Config{
		Expiration: time.Hour,
		Keys: ssh.Keys{
			test.UserID.String(): &ssh.Key{
				Config: test.NewSSH("secrets/none", "secrets/ssh_private"),
			},
		},
	}, test.FS)

	sub, err := token.Verify(tkn, strings.Empty)
	require.Error(t, err)
	require.Empty(t, sub)
}

func TestInvalidUnknownKey(t *testing.T) {
	valid := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
	tkn, err := valid.Generate(strings.Empty, strings.Empty)
	require.NoError(t, err)

	token := ssh.NewToken(&ssh.Config{
		Expiration: time.Hour,
		Keys: ssh.Keys{
			"other": &ssh.Key{
				Config: test.NewSSH("secrets/ssh_public", "secrets/ssh_private"),
			},
		},
	}, test.FS)

	sub, err := token.Verify(tkn, strings.Empty)
	require.Empty(t, sub)
	require.ErrorIs(t, err, crypto.ErrInvalidMatch)
}

func TestInvalidSignature(t *testing.T) {
	valid := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
	tkn, err := valid.Generate(strings.Empty, strings.Empty)
	require.NoError(t, err)

	encoded, _, ok := strings.Cut(tkn, ".")
	require.True(t, ok)

	sub, err := valid.Verify(encoded+"."+base64.Encode([]byte("bad")), "other")
	require.Empty(t, sub)
	require.ErrorIs(t, err, crypto.ErrInvalidMatch)
}

func TestInvalidConfig(t *testing.T) {
	token := ssh.NewToken(nil, test.FS)
	require.Nil(t, token)
}

func TestInvalidConfigDoesNotPanic(t *testing.T) {
	t.Run("generate with verification only config", func(t *testing.T) {
		token := ssh.NewToken(&ssh.Config{
			Keys: ssh.Keys{
				"test": &ssh.Key{
					Config: test.NewSSH("secrets/ssh_public", "secrets/ssh_private"),
				},
			},
		}, test.FS)

		tkn, err := token.Generate(strings.Empty, strings.Empty)
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate with empty key name", func(t *testing.T) {
		cfg := test.NewToken("ssh").SSH
		cfg.Key = strings.Empty

		token := ssh.NewToken(cfg, test.FS)

		tkn, err := token.Generate(strings.Empty, strings.Empty)
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate with empty expiration", func(t *testing.T) {
		cfg := test.NewToken("ssh").SSH
		cfg.Expiration = 0

		token := ssh.NewToken(cfg, test.FS)

		tkn, err := token.Generate(strings.Empty, strings.Empty)
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})
}

func TestInvalidGenerateKeyConfigDoesNotPanic(t *testing.T) {
	t.Run("generate with missing active key", func(t *testing.T) {
		cfg := test.NewToken("ssh").SSH
		cfg.Key = "missing"

		token := ssh.NewToken(cfg, test.FS)

		tkn, err := token.Generate(strings.Empty, strings.Empty)
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate with nil active key", func(t *testing.T) {
		cfg := test.NewToken("ssh").SSH
		cfg.Keys[cfg.Key] = nil

		token := ssh.NewToken(cfg, test.FS)

		tkn, err := token.Generate(strings.Empty, strings.Empty)
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("generate with active key missing config", func(t *testing.T) {
		cfg := test.NewToken("ssh").SSH
		cfg.Keys[cfg.Key] = &ssh.Key{}

		token := ssh.NewToken(cfg, test.FS)

		tkn, err := token.Generate(strings.Empty, strings.Empty)
		require.Empty(t, tkn)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})
}

func TestInvalidVerifyConfigDoesNotPanic(t *testing.T) {
	t.Run("verify with matching key missing config", func(t *testing.T) {
		valid := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
		tkn, err := valid.Generate(strings.Empty, strings.Empty)
		require.NoError(t, err)

		token := ssh.NewToken(&ssh.Config{
			Expiration: time.Hour,
			Keys: ssh.Keys{
				test.UserID.String(): &ssh.Key{},
			},
		}, test.FS)

		sub, err := token.Verify(tkn, strings.Empty)
		require.Empty(t, sub)
		require.ErrorIs(t, err, errors.ErrInvalidConfig)
	})

	t.Run("verify with invalid matching key config", func(t *testing.T) {
		valid := ssh.NewToken(test.NewToken("ssh").SSH, test.FS)
		tkn, err := valid.Generate(strings.Empty, strings.Empty)
		require.NoError(t, err)

		token := ssh.NewToken(&ssh.Config{
			Expiration: time.Hour,
			Keys: ssh.Keys{
				test.UserID.String(): &ssh.Key{
					Config: test.NewSSH("secrets/none", "secrets/ssh_private"),
				},
			},
		}, test.FS)

		sub, err := token.Verify(tkn, strings.Empty)
		require.Empty(t, sub)
		require.Error(t, err)
	})
}

func TestKeysGet(t *testing.T) {
	keys := ssh.Keys{
		"other": nil,
		"test":  &ssh.Key{Config: test.NewSSH("secrets/ssh_public", "secrets/ssh_private")},
	}

	key := keys.Get("test")
	require.NotNil(t, key)
	require.Nil(t, keys.Get("other"))
	require.Nil(t, keys.Get("missing"))
	require.Nil(t, ssh.Keys{}.Get("missing"))
	require.Nil(t, ssh.Keys(nil).Get("missing"))
}

func sshClaims(cfg *ssh.Config, issuedAt, expiresAt time.Time) map[string]any {
	return map[string]any{
		"ver": "v1",
		"kid": cfg.Key,
		"sub": cfg.Key,
		"aud": "/service.Method",
		"iat": issuedAt.UnixNano(),
		"exp": expiresAt.UnixNano(),
	}
}

func signedSSHToken(t *testing.T, cfg *ssh.Config, claims map[string]any) string {
	t.Helper()

	signer, err := cryptossh.NewSigner(test.FS, cfg.Keys.Get(cfg.Key).Config)
	require.NoError(t, err)

	encoded, err := json.Marshal(claims)
	require.NoError(t, err)

	signature, err := signer.Sign(encoded)
	require.NoError(t, err)

	return strings.Join(".", base64.Encode(encoded), base64.Encode(signature))
}

func encodedSSHToken(t *testing.T, claims map[string]any, signature string) string {
	t.Helper()

	encoded, err := json.Marshal(claims)
	require.NoError(t, err)

	return strings.Join(".", base64.Encode(encoded), signature)
}
