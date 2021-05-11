package auth0

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var (
	// ErrInvalidResponse from Auth0.
	ErrInvalidResponse = errors.New("invalid auth0 response")
)

type request struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type response struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type generator struct {
	cfg    *Config
	client *http.Client
}

func (g *generator) Generate() ([]byte, error) {
	req := &request{
		ClientID:     g.cfg.ClientID,
		ClientSecret: g.cfg.ClientSecret,
		Audience:     g.cfg.Audience,
		GrantType:    "client_credentials",
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // nolint:gomnd
	defer cancel()

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

	var resp response
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return []byte(resp.AccessToken), nil
}
