package oauth

import "context"

type Identity struct {
	Provider       string
	ProviderUserID string
	Email          string
	EmailVerified  bool
	FullName       *string
}

type GoogleProvider interface {
	VerifyIDToken(ctx context.Context, idToken string) (*Identity, error)
}

type GitHubProvider interface {
	AuthenticateCode(ctx context.Context, code string) (*Identity, error)
}
