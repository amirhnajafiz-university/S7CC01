package mail

import "github.com/mailgun/mailgun-go"

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
