package email

import (
	"fmt"
	"log"

	"github.com/wneessen/go-mail"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Mailer struct {
	client *mail.Client
	from   string
}

func NewMailer(cfg Config) (*Mailer, error) {
	if cfg.Password == "" {
		log.Println("smtp pw kosong")
		return &Mailer{
			from: cfg.From,
		}, nil
	}
	client, err := mail.NewClient(cfg.Host,
		mail.WithPort(cfg.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.Username),
		mail.WithPassword(cfg.Password),
	)

	if err != nil {
		return nil, err
	}

	return &Mailer{client: client, from: cfg.From}, nil
}

func (m *Mailer) SendOTP(to string, code string) error {
	if m.client == nil {
		log.Printf("email : %s code : %s", to, code)
		return nil
	}

	msg := mail.NewMsg()
	if err := msg.From(m.from); err != nil {
		return err
	}
	if err := msg.To(to); err != nil {
		return err
	}
	msg.Subject("Kode Verifikasi Akun Esdemy Kamu")
	msg.SetBodyString(mail.TypeTextHTML, otpHTML(code))

	return m.client.DialAndSend(msg)
}

func otpHTML(code string) string {
	return fmt.Sprintf(`
	<div style="font-family: Arial, sans-serif; max-width: 480px; margin: auto; border: 1px solid #e2e2e2; border-radius: 8px; padding: 24px;">
		<h2 style="color: #1a1a1a;">Verifikasi Akun ESdemy</h2>
		<p>Gunakan kode di bawah ini untuk verifikasi akun kamu. Kode berlaku selama 2 menit.</p>
		<div style="font-size: 32px; font-weight: bold; letter-spacing: 8px; text-align: center; padding: 16px; background: #f4f4f4; border-radius: 6px; margin: 16px 0;">
			%s
		</div>
	</div>
	`, code)
}
