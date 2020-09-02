package go_saas_mailer_mailgun

import (
	"context"
	"sync"

	"github.com/go-saas/go-saas/mailer"
	"github.com/mailgun/mailgun-go/v4"
)

type Mailer struct {
	Handlers map[string]func(data map[string]interface{}) (go_saas_mailer.Mail, error)
	Mailgun  *mailgun.MailgunImpl
	ApiBase  string
	Domain   string
	ApiKey   string

	*sync.RWMutex
}

func (mailer *Mailer) GetHandlers() map[string]func(data map[string]interface{}) (go_saas_mailer.Mail, error) {
	mailer.RLock()
	defer mailer.RUnlock()

	return mailer.Handlers
}

func (mailer *Mailer) GetHandler(name string) (func(data map[string]interface{}) (go_saas_mailer.Mail, error), error) {
	mailer.RLock()
	defer mailer.RUnlock()

	handler, exists := mailer.Handlers[name]

	if !exists {
		return nil, go_saas_mailer.MailNotExistsErr
	}

	return handler, nil
}

func (mailer *Mailer) SetHandler(name string, handler func(data map[string]interface{}) (go_saas_mailer.Mail, error)) {
	mailer.Lock()
	defer mailer.Unlock()

	mailer.Handlers[name] = handler
}

func (mailer *Mailer) GetMail(name string, data map[string]interface{}) (go_saas_mailer.Mail, error) {
	handler, err := mailer.GetHandler(name)

	if err != nil {
		return nil, err
	}

	return handler(data)
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
