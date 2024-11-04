package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Probe is an experiment to measure the connectivity of a network using
// different protocols and public servers.
type Probe struct {
	// Protocols to use.
	Protocols []Protocol
	// Number of iterations. Zero means infinite.
	Count uint
	// Delay between requests.
	Delay time.Duration
	// For debugging purposes.
	Logger *slog.Logger
	// Channel to send back partial results.
	ReportCh chan *Report
}

// Ensures the probe setup is correct.
func (p Probe) validate() error {
	if p.Protocols == nil {
		return newErrorReqProp("Protocols")
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

// Run the connection requests against the public servers.
//
// The context can be cancelled between different protocol attempts or count
// iterations.
// Returns an error if the setup is invalid.
func (p Probe) Run(ctx context.Context) error {
	err := p.validate()
	if err != nil {
		return fmt.Errorf("invalid setup: %w", err)
	}
	p.Logger.Debug("Starting", "setup", p)
	count := uint(0)
	for {
		select {
		case <-ctx.Done():
			p.Logger.Debug(
				"Context cancelled between iterations",
				"count", count,
			)
			return nil
		default:
			p.Logger.Debug("New iteration", "count", count)
			for _, proto := range p.Protocols {
				select {
				case <-ctx.Done():
					p.Logger.Debug(
						"Context cancelled between protocols",
						"count", count, "protocol", proto,
					)
					return nil
				default:
					start := time.Now()
					p.Logger.Debug(
						"New protocol", "count", count, "protocol", proto,
					)
					rhost, err := proto.Probe("")
					report := Report{
						ProtocolID: proto.String(),
						Time:       time.Since(start),
						Error:      err,
						RHost:      rhost,
					}
					p.Logger.Debug(
						"Sending report back",
						"count", count, "report", report,
					)
					p.ReportCh <- &report
					time.Sleep(p.Delay)
				}
			}
			p.Logger.Debug(
				"Iteration finished", "count", count, "p.Count", p.Count,
			)
			count++
			if count == p.Count {
				p.Logger.Debug("Count limit reached", "count", count)
				return nil
			}
		}
	}
}
