package probex

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/huxulm/liveness/internal/config"
	"k8s.io/kubernetes/pkg/probe"
	htp "k8s.io/kubernetes/pkg/probe/http"
	tcp "k8s.io/kubernetes/pkg/probe/tcp"
)

// Probex instance, includes configuration
type Probex struct {
	conf *config.ProbeConfig
}

// New create a Probex instance
func New(c *config.ProbeConfig) *Probex {
	return (&Probex{conf: c}).applyDefaults()
}

const (
	// DefaultInitialDelaySeconds for set default initialDelaySeconds
	DefaultInitialDelaySeconds = 0
	// DefaultPeriodSeconds for set default periodSeconds
	DefaultPeriodSeconds = 10
	// DefaultSuccessThreshold for set default successThreshold
	DefaultSuccessThreshold = 1
	// DefaultTimeoutSeconds for set default timeoutSeconds
	DefaultTimeoutSeconds = 30
	// DefaultFailureThreshold for set default failureThreshold
	DefaultFailureThreshold = 3
)

func (p *Probex) applyDefaults() *Probex {
	if p.conf == nil {
		panic(fmt.Errorf("probex conf can't be nil"))
	}
	c := p.conf
	if c.FailureThreshold <= 0 {
		c.FailureThreshold = DefaultFailureThreshold
	}
	if c.SuccessThreshold <= 0 {
		c.SuccessThreshold = DefaultSuccessThreshold
	}
	if c.TimeoutSeconds <= 0 {
		c.TimeoutSeconds = DefaultTimeoutSeconds
	}
	if c.PeriodSeconds <= 0 {
		c.PeriodSeconds = DefaultPeriodSeconds
	}
	if c.InitialDelaySeconds < 0 {
		c.InitialDelaySeconds = DefaultInitialDelaySeconds
	}
	return p
}

// ProbeResult compose return values of probe.Probe()
type ProbeResult struct {
	probe.Result
	Reporter *Probex
	Msg      string
}

// Run real probe task with context and send back result to recv
func (p *Probex) Run(ctx context.Context, waiter *sync.WaitGroup, recv chan<- ProbeResult) {
	var (
		isHTTP = p.conf.HTTPGet != nil
		isTCP  = p.conf.TCPSocket != nil
	)

	defer func() {
		if v := recover(); v != nil {
			waiter.Done()
			return
		}
		waiter.Done()
	}()

	for {
		select {
		case <-time.Tick(time.Second * time.Duration(p.conf.PeriodSeconds)):
			if isHTTP {
				req, err := buildRequest(p.conf.HTTPGet)
				if err != nil {
					panic(fmt.Errorf("probe request build failed: %v", err))
				}
				res, msg, _ := htp.New(false).Probe(req, time.Second*time.Duration(p.conf.TimeoutSeconds))
				recv <- ProbeResult{Result: res, Msg: msg, Reporter: p}
			}
			if isTCP {
				res, msg, _ := tcp.New().Probe(p.conf.TCPSocket.Host, p.conf.TCPSocket.Port, time.Second*time.Duration(p.conf.TimeoutSeconds))
				recv <- ProbeResult{Result: res, Msg: msg, Reporter: p}
			}

		case <-ctx.Done():
			log.Printf("probe canceled: %s\n", p.conf.Name)
			return
		}
	}
}

func buildRequest(action *config.HTTPGetAction) (*http.Request, error) {
	return http.NewRequest("GET", fmt.Sprintf("%s://%s%s", action.Scheme, action.Host, action.Path), nil)
}
