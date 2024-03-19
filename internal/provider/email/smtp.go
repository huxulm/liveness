package email

import (
	"bytes"
	"crypto/tls"
	"html/template"

	"github.com/huxulm/liveness/internal/provider"
	"gopkg.in/gomail.v2"
)

// SMTPConfig defines config properties of SMTPProvider
type SMTPConfig struct {
	Receivers []string `yaml:"receivers"`
	Host      string   `yaml:"host"`
	Port      int      `yaml:"port"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	HTML      string   `yaml:"html"`
}

type smtpProvider struct {
	conf             *SMTPConfig
	dialer           *gomail.Dialer
	messageTmplelate *template.Template
}

// args[0] for {{.Content}}
// args[1] for {{.Reason}}
func (p *smtpProvider) Send(args ...string) error {
	m := gomail.NewMessage(gomail.SetCharset("utf-8"))
	m.SetHeader("From", p.conf.Username)
	m.SetHeader("To", p.conf.Receivers...)
	m.SetHeader("Subject", "线上服务告警")
	html := &bytes.Buffer{}
	if err := p.messageTmplelate.Execute(html, map[string]string{
		"Content": args[0],
		"Reason":  args[1],
	}); err != nil {
		return err
	}
	m.SetBody("text/html", html.String())
	return p.dialer.DialAndSend(m)
}

// NewSMTP create a smtp provider
func NewSMTP(c *SMTPConfig) (provider.Provider, error) {
	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	tmpl, err := template.New("tmpl").Parse(c.HTML)
	if err != nil {
		return nil, err
	}
	return &smtpProvider{dialer: d, conf: c, messageTmplelate: tmpl}, nil
}
