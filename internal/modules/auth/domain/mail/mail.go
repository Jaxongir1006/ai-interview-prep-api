package mail

import "context"

type VerificationEmail struct {
	To              string
	VerificationURL string
}

type Sender interface {
	SendVerificationEmail(ctx context.Context, msg VerificationEmail) error
}
