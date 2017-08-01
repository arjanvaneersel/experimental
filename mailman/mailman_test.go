package mailman

import (
	"net/mail"
	"testing"
)

type testEmailSender struct {
	lastSubject string
	lastBody    string
	lastTo      []*mail.Address
}

// make sure it satisfies the interface
var _ EmailSender = (*testEmailSender)(nil)

func (t *testEmailSender) Send(subject, body string, to ...*mail.Address) {
	t.lastSubject = subject
	t.lastBody = body
	t.lastTo = to
}

func TestSendWelcomeEmail(t *testing.T) {
	sender := &testEmailSender{}
	SendWelcomeEmail(sender, to1, to2)
	if sender.lastSubject != "Welcome" {
		t.Error("Subject line was wrong")
	}
	if sender.To[0] != to1 && sender.To[1] != to2 {
		t.Error("Wrong recipients")
	}
}
