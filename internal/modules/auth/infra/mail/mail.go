package mail

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"

	domainmail "github.com/Jaxongir1006/ai-interview-prep-api/internal/modules/auth/domain/mail"

	"github.com/code19m/errx"
)

const (
	ProviderNoop = "noop"
	ProviderSMTP = "smtp"
)

type Config struct {
	Provider string `yaml:"provider" default:"noop" validate:"oneof=noop smtp"`

	FromEmail string `yaml:"from_email" default:"no-reply@hire-ready.local"`
	FromName  string `yaml:"from_name"  default:"HireReady"`

	SMTP SMTPConfig `yaml:"smtp"`
}

type SMTPConfig struct {
	Host string `yaml:"host" default:"localhost"`
	Port int    `yaml:"port" default:"1025"`

	Username string `yaml:"username"`
	Password string `yaml:"password" mask:"true"`
}

func NewSender(cfg Config) domainmail.Sender {
	if cfg.Provider != ProviderSMTP {
		return noopSender{}
	}

	return &smtpSender{
		fromEmail: cfg.FromEmail,
		fromName:  cfg.FromName,
		host:      cfg.SMTP.Host,
		port:      cfg.SMTP.Port,
		username:  cfg.SMTP.Username,
		password:  cfg.SMTP.Password,
	}
}

type noopSender struct{}

func (noopSender) SendVerificationEmail(context.Context, domainmail.VerificationEmail) error {
	return nil
}

type smtpSender struct {
	fromEmail string
	fromName  string
	host      string
	port      int
	username  string
	password  string
}

func (s *smtpSender) SendVerificationEmail(
	_ context.Context,
	msg domainmail.VerificationEmail,
) error {
	addr := net.JoinHostPort(s.host, strconv.Itoa(s.port))
	auth := smtp.Auth(nil)
	if s.username != "" || s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	body := s.message(msg)
	err := smtp.SendMail(addr, auth, s.fromEmail, []string{msg.To}, []byte(body))
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func (s *smtpSender) message(msg domainmail.VerificationEmail) string {
	subject := "Verify your HireReady email"
	textBody := fmt.Sprintf(
		"Verify your email address by opening this link:\r\n\r\n%s\r\n",
		msg.VerificationURL,
	)

	headers := []string{
		"From: " + s.fromHeader(),
		"To: " + msg.To,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + textBody
}

func (s *smtpSender) fromHeader() string {
	if s.fromName == "" {
		return s.fromEmail
	}
	return fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
}
