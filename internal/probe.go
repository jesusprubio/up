package internal

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Probe is an experiment to measure the connectivity of a network.
type Probe struct {
	// Protocol to use.
	Proto Protocol
	// Number of iterations. Zero means infinite.
	Count uint
	// Delay between requests.
	Delay time.Duration
	// For debugging purposes.
	Logger *slog.Logger
	// Channel to send back partial results.
	ReportCh chan *Report
	// Optional. Where to point the probe.
	// URL (HTTP), host/port string (TCP) or domain (DNS).
	Target string
}

// Ensures the probe setup is correct.
func (p Probe) validate() error {
	if p.Proto == nil {
		return newErrorReqProp("Proto")
	}
	// 'Delay' could be zero.
	if p.Logger == nil {
		return newErrorReqProp("Logger")
	}
	if p.ReportCh == nil {
		return newErrorReqProp("ReportCh")
	}
	return nil
}

func newErrorReqProp(prop string) error {
	return fmt.Errorf("required property: %s", prop)
}

// Do makes the connection requests against the public servers.
//
// The context can be cancelled between different protocol attempts or count
// iterations.
// Returns an error if the setup is invalid.
func (p Probe) Do(ctx context.Context) error {
	err := p.validate()
	if err != nil {
		return fmt.Errorf("invalid setup: %w", err)
	}
	p.Logger.Debug("Starting", "setup", p)
	count := uint(0)
	for {
		select {
		case <-ctx.Done():
			p.Logger.Debug("Context cancelled", "count", count)
			return nil
		default:
			p.Logger.Debug(
				"New iteration",
				"count",
				count,
				"protocol",
				p.Proto,
				"target",
				p.Target,
			)
			start := time.Now()
			target, extra, err := p.Proto.Probe(p.Target)
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			report := Report{
				ProtocolID: p.Proto.String(),
				Time:       time.Since(start),
				Error:      errMsg,
				Target:     target,
				Extra:      extra,
			}
			p.Logger.Debug("Sending report back", "report", report)
			p.ReportCh <- &report
			time.Sleep(p.Delay)
			count++
			if p.Count > 0 && count >= p.Count {
				p.Logger.Debug("Count limit reached", "count", count)
				return nil
			}
		}
	}
}
