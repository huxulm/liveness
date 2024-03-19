package config

import (
	"os"

	"github.com/huxulm/liveness/internal/provider/email"
	"github.com/huxulm/liveness/internal/provider/sms"
	"gopkg.in/yaml.v2"
)

// Config provides configuration for probe
type ProbeConfig struct {
	// Length of time before health checking is activated.  In seconds.
	// +optional
	InitialDelaySeconds int32 `yaml:"initialDelaySeconds"`
	// Length of time before health checking times out.  In seconds.
	// +optional
	TimeoutSeconds int32 `yaml:"timeoutSeconds"`
	// How often (in seconds) to perform the probe.
	// +optional
	PeriodSeconds int32 `yaml:"periodSeconds"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Must be 1 for liveness and startup.
	// +optional
	SuccessThreshold int32 `yaml:"successThreshold"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// +optional
	FailureThreshold int32 `yaml:"failureThreshold"`

	Name         string `yaml:"name"`
	ProbeHandler `yaml:",inline"`
}

// HTTPGetAction describes an action based on HTTP Get requests.
type HTTPGetAction struct {
	// Path to access on the HTTP server.
	// +optional
	Path string `yaml:"path,omitempty"`
	// Name or number of the port to access on the container.
	// Number must be in the range 1 to 65535.
	// Name must be an IANA_SVC_NAME.
	Port int `yaml:"port"`
	// Host name to connect to, defaults to the pod IP. You probably want to set
	// "Host" in httpHeaders instead.
	// +optional
	Host string `yaml:"host,omitempty"`
	// Scheme to use for connecting to the host.
	// Defaults to HTTP.
	// +optional
	Scheme string `yaml:"scheme,omitempty"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	// +optional
	HTTPHeaders []HTTPHeader `yaml:"httpHeaders,omitempty"`
}

// HTTPHeader describes a custom header to be used in HTTP probes
type HTTPHeader struct {
	// The header field name.
	// This will be canonicalized upon output, so case-variant names will be understood as the same header.
	Name string `yaml:"name"`
	// The header field value
	Value string `yaml:"value"`
}

// TCPSocketAction describes an action based on opening a socket
type TCPSocketAction struct {
	// Number or name of the port to access on the container.
	// Number must be in the range 1 to 65535.
	// Name must be an IANA_SVC_NAME.
	Port int `yaml:"port"`
	// Optional: Host name to connect to, defaults to the pod IP.
	// +optional
	Host string `yaml:"host,omitempty"`
}

// ProbeHandler ...
type ProbeHandler struct {
	// HTTPGet specifies the http request to perform.
	// +optional
	HTTPGet *HTTPGetAction `yaml:"httpGet,omitempty"`
	// TCPSocket specifies an action involving a TCP port.
	// +optional
	TCPSocket *TCPSocketAction `yaml:"tcpSocket,omitempty"`
}

type ProviderConfig struct {
	sms.SmsConfig    `yaml:",inline"`
	email.SMTPConfig `yaml:",inline"`
}

type Config struct {
	Probes    []*ProbeConfig    `yaml:"probes"`
	Providers []*ProviderConfig `yaml:"providers"`
}

func LoadFromFile(conf string) (*Config, error) {
	var c Config
	r, err := os.Open(conf)
	if err != nil {
		return nil, err
	}

	err = yaml.NewDecoder(r).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
