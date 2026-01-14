package notifications

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/wneessen/go-mail"
)

func SendMail(cfg *config.Config, to string, subject string, body string) (err error) {
	// Защита от паник при отправке email
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "PANIC recovered in SendMail: %v\n", r)
			err = fmt.Errorf("panic in SendMail: %v", r)
		}
	}()

	message := mail.NewMsg()
	if err := message.From(cfg.SMTP_FROM_EMAIL); err != nil {
		return err
	}
	if err := message.To(to); err != nil {
		return err
	}
	message.Subject(subject)
	message.SetDate()
	message.SetMessageID()
	message.SetBodyString(mail.TypeTextPlain, body)
	client, err := mail.NewClient(
		cfg.SMTP_HOST,
		mail.WithPort(cfg.SMTP_PORT),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(cfg.SMTP_FROM_EMAIL),
		mail.WithPassword(cfg.SMTP_PASSWORD),
		mail.WithTimeout(3*time.Second),
		mail.WithTLSPolicy(mail.TLSMandatory),
		mail.WithTLSConfig(&tls.Config{
			ServerName:         cfg.SMTP_HOST,
			InsecureSkipVerify: true,
		}),
	)
	if err != nil {
		return err
	}
	if cfg.ENV == "dev" {
		client.SetDebugLog(true)
	}
	if err := client.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
