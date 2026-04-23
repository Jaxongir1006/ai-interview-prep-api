package mail

import "context"

type VerificationEmail struct {
	To              string
	VerificationURL string
}

type PasswordResetEmail struct {
	To       string
	ResetURL string
}

type Sender interface {
	SendVerificationEmail(ctx context.Context, msg VerificationEmail) error
	SendPasswordResetEmail(ctx context.Context, msg PasswordResetEmail) error
}
