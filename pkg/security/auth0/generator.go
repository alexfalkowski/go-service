package auth0

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	sjwt "github.com/alexfalkowski/go-service/pkg/security/jwt"
	"github.com/dgraph-io/ristretto"
	"github.com/form3tech-oss/jwt-go"
)

// ErrInvalidResponse from Auth0.
var ErrInvalidResponse = errors.New("invalid response")

type generatorRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type generatorResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type generator struct {
	cfg    *Config
	client *http.Client
}

func (g *generator) Generate(ctx context.Context) ([]byte, error) {
	req := &generatorRequest{
		ClientID:     g.cfg.ClientID,
		ClientSecret: g.cfg.ClientSecret,
		Audience:     g.cfg.Audience,
		GrantType:    "client_credentials",
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.cfg.URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Content-Type", "application/json")

	httpResp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 { // nolint:gomnd
		return nil, ErrInvalidResponse
	}

	var resp generatorResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return []byte(resp.AccessToken), nil
}

type cachedGenerator struct {
	cfg   *Config
	cache *ristretto.Cache

	sjwt.Generator
}

// nolint:forcetypeassert
func (g *cachedGenerator) Generate(ctx context.Context) ([]byte, error) {
	cacheKey := g.cfg.CacheKey("generate")

	v, ok := g.cache.Get(cacheKey)
	if ok {
		return v.([]byte), nil
	}

	key, err := g.Generator.Generate(ctx)
	if err != nil {
		return nil, err
	}

	parser := &jwt.Parser{}

	token, _, err := parser.ParseUnverified(string(key), jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	exp := claims["exp"].(float64)
	ttl := time.Unix(int64(exp), 0).Add(-30 * time.Second)

	g.cache.SetWithTTL(cacheKey, key, 0, time.Until(ttl))

	return key, nil
}
