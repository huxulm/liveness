package probex

import (
	"context"

	"github.com/huxulm/liveness/internal/config"
	"github.com/huxulm/liveness/internal/provider"
	"github.com/huxulm/liveness/internal/provider/email"
	"github.com/huxulm/liveness/internal/provider/sms"
)

// Run probex deamon in background
func Run(ctx context.Context, conf *config.Config) error {
	sch := Scheduler{results: make(chan ProbeResult, 1000), status: make(map[uint32]*ProbeStatus)}

	var probes []*Probex
	for _, c := range conf.Probes {
		probes = append(probes, New(c))
	}

	sch.probes = probes
	sch.providers = map[string]provider.Provider{}

	for _, provider := range conf.Providers {
		switch provider.Type {
		case "sms":
			sch.providers["sms"], _ = sms.NewSms(&sms.SmsConfig{
				Key:          provider.Key,
				Secret:       provider.Secret,
				Phones:       provider.Phones,
				Provider:     "sms",
				TemplateCode: provider.TemplateCode,
				SignName:     provider.SignName,
			})
		case "smtp":
			sch.providers["smtp"], _ = email.NewSMTP(&email.SMTPConfig{
				Username:  provider.Username,
				Password:  provider.Password,
				Host:      provider.Host,
				Port:      provider.Port,
				HTML:      provider.HTML,
				Receivers: provider.Receivers,
			})
		}
	}

	sch.Run(ctx)
	return nil
}
