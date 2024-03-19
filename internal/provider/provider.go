package provider

type Provider interface {
	Send(args ...string) error
}
