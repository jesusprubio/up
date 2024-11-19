package internal

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
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
	//URLs (HTTP), host/port strings (TCP) or domains (DNS).
	Input    []string
	Parallel bool
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

// Do the connection requests against the public servers.
//
// The context can be cancelled between different protocol attempts or count
// iterations.
// Returns an error if the setup is invalid.

func (p Probe) Do(ctx context.Context) error {
	err := p.validate()
	if err != nil {
		return fmt.Errorf("invalid setup: %w", err)
	}

	inputs := p.Input
	if len(inputs) == 0 {
		inputs = []string{""}
	}

	p.Logger.Debug("Starting", "setup", p)

	if p.Parallel {
		return p.doParallel(ctx, inputs)
	}
	return p.doSerial(ctx, inputs)
}

func (p Probe) doSerial(ctx context.Context, inputs []string) error {
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
			for _, addr := range inputs {
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
						rhost, extra, err := proto.Probe(addr)
						report := Report{
							ProtocolID: proto.String(),
							Time:       time.Since(start),
							Error:      err,
							RHost:      rhost,
							Extra:      extra,
						}
						p.Logger.Debug(
							"Sending report back",
							"count", count, "report", report,
						)
						p.ReportCh <- &report

						select {
						case <-time.After(p.Delay):
						case <-ctx.Done():
							return nil
						}
					}
				}
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

func (p Probe) doParallel(ctx context.Context, inputs []string) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, len(p.Protocols))

	for count := uint(0); count < p.Count || p.Count == 0; count++ {
		select {
		case <-ctx.Done():
			p.Logger.Debug(
				"Context cancelled between iterations",
				"count", count,
			)
			return nil
		default:
			p.Logger.Debug("New iteration", "count", count)
			for _, addr := range inputs {
				for _, proto := range p.Protocols {
					sem <- struct{}{}
					wg.Add(1)
					go func(addr string, proto Protocol, iterCount uint) {
						defer func() {
							<-sem
							wg.Done()
						}()

						start := time.Now()
						p.Logger.Debug(
							"New protocol",
							"count",
							iterCount,
							"protocol",
							proto,
						)
						rhost, extra, err := proto.Probe(addr)
						report := Report{
							ProtocolID: proto.String(),
							Time:       time.Since(start),
							Error:      err,
							RHost:      rhost,
							Extra:      extra,
						}
						p.Logger.Debug(
							"Sending report back",
							"count", iterCount, "report", report,
						)
						p.ReportCh <- &report

						time.Sleep(p.Delay)
					}(addr, proto, count)
				}
			}

			wg.Wait()
		}
	}

	return nil
}
