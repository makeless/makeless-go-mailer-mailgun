package go_saas_mailer_mailgun

import (
	"context"
	"sync"

	"github.com/go-saas/go-saas/mailer"
	"github.com/mailgun/mailgun-go/v4"
)

type Mailer struct {
	Mailgun *mailgun.MailgunImpl
	ApiBase string
	Domain  string
	ApiKey  string

	*sync.RWMutex
}

func (mailer *Mailer) GetMailgun() *mailgun.MailgunImpl {
	mailer.RLock()
	defer mailer.RUnlock()

	return mailer.Mailgun
}

func (mailer *Mailer) SetMailgun(mailgun *mailgun.MailgunImpl) {
	mailer.Lock()
	defer mailer.Unlock()

	mailer.Mailgun = mailgun
}

func (mailer *Mailer) GetApiBase() string {
	mailer.RLock()
	defer mailer.RUnlock()

	return mailer.ApiBase
}

func (mailer *Mailer) GetDomain() string {
	mailer.RLock()
	defer mailer.RUnlock()

	return mailer.Domain
}

func (mailer *Mailer) GetApiKey() string {
	mailer.RLock()
	defer mailer.RUnlock()

	return mailer.ApiKey
}

func (mailer *Mailer) Init() error {
	mg := mailgun.NewMailgun(mailer.GetDomain(), mailer.GetApiKey())
	mg.SetAPIBase(mailer.GetApiBase())

	mailer.SetMailgun(mg)
	return nil
}

func (mailer *Mailer) Send(ctx context.Context, mail go_saas_mailer.Mail) error {
	var message = mailer.GetMailgun().NewMessage(
		mail.GetFrom(),
		mail.GetSubject(),
		string(mail.GetMessage()),
		mail.GetTo()...,
	)

	message.SetHtml(string(mail.GetHtmlMessage()))

	for _, cc := range mail.GetCc() {
		message.AddCC(cc)
	}

	for _, bcc := range mail.GetBcc() {
		message.AddCC(bcc)
	}

	for key := range mail.GetHeaders() {
		message.AddHeader(key, mail.GetHeaders().Get(key))
	}

	for _, attachment := range mail.GetAttachments() {
		message.AddBufferAttachment(attachment.GetFilename(), attachment.GetData())
	}

	_, _, err := mailer.GetMailgun().Send(ctx, message)
	return err
}
