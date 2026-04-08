package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/code19m/errx"
)

type GitHubConfig struct {
	Enabled       bool   `yaml:"enabled"`
	ClientID      string `yaml:"client_id"`
	ClientSecret  string `yaml:"client_secret"`
	TokenURL      string `yaml:"token_url"       default:"https://github.com/login/oauth/access_token"`
	UserURL       string `yaml:"user_url"        default:"https://api.github.com/user"`
	UserEmailsURL string `yaml:"user_emails_url" default:"https://api.github.com/user/emails"`
}

type gitHubProvider struct {
	httpClient    *http.Client
	clientID      string
	clientSecret  string
	tokenURL      string
	userURL       string
	userEmailsURL string
}

func NewGitHubProvider(cfg GitHubConfig) GitHubProvider {
	if !cfg.Enabled {
		return disabledGitHubProvider{}
	}

	return &gitHubProvider{
		httpClient:    &http.Client{},
		clientID:      cfg.ClientID,
		clientSecret:  cfg.ClientSecret,
		tokenURL:      defaultString(cfg.TokenURL, "https://github.com/login/oauth/access_token"),
		userURL:       defaultString(cfg.UserURL, "https://api.github.com/user"),
		userEmailsURL: defaultString(cfg.UserEmailsURL, "https://api.github.com/user/emails"),
	}
}

type gitHubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type gitHubUserResponse struct {
	ID    int64   `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type gitHubEmailResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (p *gitHubProvider) AuthenticateCode(ctx context.Context, code string) (*Identity, error) {
	accessToken, err := p.exchangeCode(ctx, code)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	userResp, err := p.fetchUser(ctx, accessToken)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	email, verified, err := p.resolveEmail(ctx, accessToken, userResp.Email)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if email == "" {
		return nil, errx.New("github account does not expose a usable email")
	}

	return &Identity{
		Provider:       "github",
		ProviderUserID: strconv.FormatInt(userResp.ID, 10),
		Email:          strings.ToLower(email),
		EmailVerified:  verified,
		FullName:       userResp.Name,
	}, nil
}

func (p *gitHubProvider) exchangeCode(ctx context.Context, code string) (string, error) {
	values := url.Values{}
	values.Set("client_id", p.clientID)
	values.Set("client_secret", p.clientSecret)
	values.Set("code", code)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.tokenURL,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", errx.Wrap(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", errx.Wrap(err)
	}
	defer resp.Body.Close()

	var payload gitHubAccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return "", errx.Wrap(err)
	}
	if resp.StatusCode >= http.StatusBadRequest || payload.AccessToken == "" {
		return "", errx.New("github code exchange failed")
	}

	return payload.AccessToken, nil
}

func (p *gitHubProvider) fetchUser(ctx context.Context, accessToken string) (*gitHubUserResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.userURL, nil)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	defer resp.Body.Close()

	var payload gitHubUserResponse
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	if resp.StatusCode >= http.StatusBadRequest || payload.ID == 0 {
		return nil, errx.New("github user fetch failed")
	}

	return &payload, nil
}

func (p *gitHubProvider) resolveEmail(
	ctx context.Context,
	accessToken string,
	directEmail *string,
) (string, bool, error) {
	if directEmail != nil && *directEmail != "" {
		return *directEmail, false, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.userEmailsURL, nil)
	if err != nil {
		return "", false, errx.Wrap(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", false, errx.Wrap(err)
	}
	defer resp.Body.Close()

	var payload []gitHubEmailResponse
	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return "", false, errx.Wrap(err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return "", false, errx.New("github email fetch failed")
	}

	for _, email := range payload {
		if email.Primary && email.Email != "" {
			return email.Email, email.Verified, nil
		}
	}
	for _, email := range payload {
		if email.Email != "" {
			return email.Email, email.Verified, nil
		}
	}

	return "", false, nil
}

func defaultString(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}

type disabledGitHubProvider struct{}

func (disabledGitHubProvider) AuthenticateCode(context.Context, string) (*Identity, error) {
	return nil, errx.New("github oauth is not configured")
}
