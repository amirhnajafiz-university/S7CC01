package mail

import (
	"github.com/mailgun/mailgun-go"
)

const (
	sender = "admin@aut.ac.ir"
)

// Mailgun
// contains variables for connection to mailgun.
type Mailgun struct {
	APIKEY string
	Client *mailgun.MailgunImpl
}

// NewConnection
// opens a new connection for mailgun service.
func NewConnection(cfg Config) *Mailgun {
	return &Mailgun{
		APIKEY: cfg.APIKEY,
		Client: mailgun.NewMailgun(cfg.Domain, cfg.APIKEY),
	}
}

// Send
// emails with mailgun client.
func (m *Mailgun) Send(content, subject, receiver string) error {
	message := m.Client.NewMessage(sender, subject, content, receiver)

	_, _, err := m.Client.Send(message)

	return err
}
