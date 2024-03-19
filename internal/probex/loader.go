package probex

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/huxulm/liveness/internal/provider"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/probe"
)

// AddProbeLoaderFlag add --conf / -c flag
func AddProbeLoaderFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("conf", "c", "probes.yaml", "probes configuration file in yaml format")
}

// ProbeStatus
type ProbeStatus struct {
	ID          uint32
	Status      probe.Result
	LastErrTime *time.Time
}

func idHash(s string) uint32 {
	h := fnv.New32()
	h.Write([]byte(s))
	return h.Sum32()
}

// Scheduler holds the logic of scheduling probe tasks
type Scheduler struct {
	probes    []*Probex
	results   chan ProbeResult
	status    map[uint32]*ProbeStatus
	providers map[string]provider.Provider
}

// Run scheduler all probe tasks in multi threads
func (sh *Scheduler) Run(ctx context.Context) {
	//
	// receive probe results
	go sh.listen(ctx)

	wait := &sync.WaitGroup{}
	for _, pb := range sh.probes {
		wait.Add(1)
		go pb.Run(ctx, wait, sh.results)
	}
	wait.Wait()
}

func (sh *Scheduler) listen(ctx context.Context) {
	for {
		select {
		case r := <-sh.results:
			id := idHash(r.Reporter.conf.Name)
			var status *ProbeStatus
			if old, ok := sh.status[id]; ok {
				status = old
			} else {
				status = &ProbeStatus{ID: id, Status: r.Result, LastErrTime: nil}
			}
			status.Status = r.Result

			if r.Result == probe.Success {
				status.LastErrTime = nil // clear err time last before
			}

			if r.Result == probe.Failure {
				t := time.Now()
				if status.LastErrTime == nil || t.After((*status.LastErrTime).Add(time.Minute*5)) {
					status.LastErrTime = &t
					sh.notifyWithSms(&r)
					sh.notifyWithSMTP(&r)
				}
			}
			sh.status[id] = status
		case <-ctx.Done():
			return
		}
	}
}

func (sh *Scheduler) notifyWithSms(r *ProbeResult) {
	if sms, ok := sh.providers["sms"]; ok {
		args := []string{}
		// http
		if r.Reporter.conf.HTTPGet != nil {
			conf := r.Reporter.conf.HTTPGet
			_ = conf
			args = append(args, fmt.Sprintf("%s://%s%s \n%s", conf.Scheme, conf.Host, conf.Path, r.Msg))
			args = append(args, time.Now().Format("2006-01-02 15:04:05"))
		}

		// tcp
		if r.Reporter.conf.TCPSocket != nil {
			conf := r.Reporter.conf.TCPSocket
			args = append(args, fmt.Sprintf("%s:%d \n%s", conf.Host, conf.Port, r.Msg))
			args = append(args, time.Now().Format("2006-01-02 15:04:05"))
		}
		if len(args) == 0 {
			args = append(args, "Unknown", "Unknown")
		}
		sms.Send(args...)
	}
}

func (sh *Scheduler) notifyWithSMTP(r *ProbeResult) {
	if sms, ok := sh.providers["smtp"]; ok {
		args := []string{}

		// http
		if r.Reporter.conf.HTTPGet != nil {
			conf := r.Reporter.conf.HTTPGet
			args = append(args, fmt.Sprintf("%s://%s%s", conf.Scheme, conf.Host, conf.Path))
			args = append(args, fmt.Sprintf("(%s)[%s]", r.Result, r.Msg))
		}

		// tcp
		if r.Reporter.conf.TCPSocket != nil {
			conf := r.Reporter.conf.TCPSocket
			args = append(args, fmt.Sprintf("%s:%d", conf.Host, conf.Port), r.Msg)
		}
		if len(args) == 0 {
			args = append(args, "Unknown", "Unknown")
		}
		sms.Send(args...)
	}
}
