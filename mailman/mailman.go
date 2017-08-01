package mailman

import (
	"fmt"
	"net/mail"
)

type MailMan struct{}

type EmailSender interface {
	Send(subject, body string, to ...*mail.Address)
}

func (m *MailMan) Send(subject, body string, to ...*mail.Address) {
	fmt.Println("Sending email...")
}

func New() *MailMan {
	return &MailMan{}
}
