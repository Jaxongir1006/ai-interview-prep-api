package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/code19m/errx"
)

type GoogleConfig struct {
	Enabled      bool   `yaml:"enabled"`
	ClientID     string `yaml:"client_id"`
	TokenInfoURL string `yaml:"token_info_url" default:"https://oauth2.googleapis.com/tokeninfo"`
}

type googleProvider struct {
	httpClient   *http.Client
	tokenInfoURL string
	clientID     string
}

func NewGoogleProvider(cfg GoogleConfig) GoogleProvider {
	if !cfg.Enabled {
		return disabledGoogleProvider{}
	}

	tokenInfoURL := cfg.TokenInfoURL
	if tokenInfoURL == "" {
		//nolint:gosec // OAuth provider metadata endpoint, not a credential
		tokenInfoURL = "https://oauth2.googleapis.com/tokeninfo"
	}

	return &googleProvider{
		httpClient:   &http.Client{},
		tokenInfoURL: tokenInfoURL,
		clientID:     cfg.ClientID,
	}
}

type googleTokenInfoResponse struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Audience      string `json:"aud"`
	Error         string `json:"error_description"`
}

func (p *googleProvider) VerifyIDToken(ctx context.Context, idToken string) (*Identity, error) {
	u, err := url.Parse(p.tokenInfoURL)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	query := u.Query()
	query.Set("id_token", idToken)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer resp.Body.Close()

	var payload googleTokenInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, errx.New("google token is invalid")
	}
	if payload.Sub == "" || payload.Email == "" {
		return nil, errx.New("google identity is incomplete")
	}
	if p.clientID != "" && payload.Audience != "" && payload.Audience != p.clientID {
		return nil, errx.New("google token audience mismatch")
	}

	var fullName *string
	if payload.Name != "" {
		fullName = &payload.Name
	}

	return &Identity{
		Provider:       "google",
		ProviderUserID: payload.Sub,
		Email:          strings.ToLower(payload.Email),
		EmailVerified:  strings.EqualFold(payload.EmailVerified, "true"),
		FullName:       fullName,
	}, nil
}

type disabledGoogleProvider struct{}

func (disabledGoogleProvider) VerifyIDToken(context.Context, string) (*Identity, error) {
	return nil, errx.New("google oauth is not configured")
}
