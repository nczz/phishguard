package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/mailgun/mailgun-go/v4"
)

// Message represents an email message.
type Message struct {
	From     string
	FromName string
	To       string
	Subject  string
	HTMLBody string
	TextBody string
	Headers  map[string]string
}

// Mailer is the interface for sending emails.
type Mailer interface {
	Send(ctx context.Context, msg *Message) error
	Name() string
}

// --- SMTP ---

// SMTPMailer sends email via SMTP.
type SMTPMailer struct {
	Host     string
	Port     int
	Username string
	Password string
	TLS      bool
}

func (s *SMTPMailer) Name() string { return "smtp" }

func (s *SMTPMailer) Send(_ context.Context, msg *Message) error {
	addr := net.JoinHostPort(s.Host, strconv.Itoa(s.Port))
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	from := msg.From
	if msg.FromName != "" {
		from = fmt.Sprintf("%s <%s>", msg.FromName, msg.From)
	}

	body := "From: " + from + "\r\n" +
		"To: " + msg.To + "\r\n" +
		"Subject: " + msg.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n"

	for k, v := range msg.Headers {
		body += k + ": " + v + "\r\n"
	}

	if msg.HTMLBody != "" {
		body += "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" + msg.HTMLBody
	} else {
		body += "Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" + msg.TextBody
	}

	if s.TLS {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: s.Host})
		if err != nil {
			return fmt.Errorf("smtp tls dial: %w", err)
		}
		defer conn.Close()

		c, err := smtp.NewClient(conn, s.Host)
		if err != nil {
			return fmt.Errorf("smtp new client: %w", err)
		}
		defer c.Close()

		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
		if err := c.Mail(msg.From); err != nil {
			return fmt.Errorf("smtp mail: %w", err)
		}
		if err := c.Rcpt(msg.To); err != nil {
			return fmt.Errorf("smtp rcpt: %w", err)
		}
		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("smtp data: %w", err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			return fmt.Errorf("smtp write: %w", err)
		}
		return w.Close()
	}

	return smtp.SendMail(addr, auth, msg.From, []string{msg.To}, []byte(body))
}

// --- Mailgun ---

// MailgunMailer sends email via Mailgun API.
type MailgunMailer struct {
	Domain string
	APIKey string
}

func (m *MailgunMailer) Name() string { return "mailgun" }

func (m *MailgunMailer) Send(ctx context.Context, msg *Message) error {
	mg := mailgun.NewMailgun(m.Domain, m.APIKey)

	// Mailgun requires RFC 5322 format: "Name <email@domain>"
	// The email domain MUST match the Mailgun sending domain
	from := msg.From
	if msg.FromName != "" {
		from = fmt.Sprintf("%s <%s>", msg.FromName, msg.From)
	}

	message := mg.NewMessage(from, msg.Subject, msg.TextBody, msg.To)
	if msg.HTMLBody != "" {
		message.SetHtml(msg.HTMLBody)
	}
	// Disable Mailgun's own tracking — we use our own tracking system
	message.SetTracking(false)
	message.SetTrackingClicks(false)
	message.SetTrackingOpens(false)

	for k, v := range msg.Headers {
		message.AddHeader(k, v)
	}

	_, _, err := mg.Send(ctx, message)
	return err
}

// --- SES ---

// SESMailer sends email via AWS SES v2.
type SESMailer struct {
	Region    string
	AccessKey string
	SecretKey string
}

func (s *SESMailer) Name() string { return "ses" }

func (s *SESMailer) Send(ctx context.Context, msg *Message) error {
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(s.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.AccessKey, s.SecretKey, "")),
	)
	if err != nil {
		return fmt.Errorf("ses load config: %w", err)
	}

	client := sesv2.NewFromConfig(cfg)

	from := msg.From
	if msg.FromName != "" {
		from = fmt.Sprintf("%s <%s>", msg.FromName, msg.From)
	}

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(from),
		Destination: &sestypes.Destination{
			ToAddresses: []string{msg.To},
		},
		Content: &sestypes.EmailContent{
			Simple: &sestypes.Message{
				Subject: &sestypes.Content{Data: aws.String(msg.Subject)},
				Body:    &sestypes.Body{},
			},
		},
	}

	if msg.HTMLBody != "" {
		input.Content.Simple.Body.Html = &sestypes.Content{Data: aws.String(msg.HTMLBody)}
	}
	if msg.TextBody != "" {
		input.Content.Simple.Body.Text = &sestypes.Content{Data: aws.String(msg.TextBody)}
	}

	_, err = client.SendEmail(ctx, input)
	return err
}

// --- Factory ---

// NewMailer creates a Mailer by type and config map.
func NewMailer(mailerType string, config map[string]string) (Mailer, error) {
	switch mailerType {
	case "smtp":
		port, _ := strconv.Atoi(config["port"])
		tlsEnabled, _ := strconv.ParseBool(config["tls"])
		return &SMTPMailer{
			Host:     config["host"],
			Port:     port,
			Username: config["username"],
			Password: config["password"],
			TLS:      tlsEnabled,
		}, nil
	case "mailgun":
		return &MailgunMailer{
			Domain: config["domain"],
			APIKey: config["api_key"],
		}, nil
	case "ses":
		return &SESMailer{
			Region:    config["region"],
			AccessKey: config["access_key"],
			SecretKey: config["secret_key"],
		}, nil
	default:
		return nil, fmt.Errorf("unknown mailer type: %s", mailerType)
	}
}
