// Package main implements a simple CLI to use the library.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/fatih/color"
	"github.com/jesusprubio/up/internal"
)

const (
	appName = "up"
	appDesc = `
	Troubleshoot problems with your Internet connection based on different
	protocols and public servers.

	OUTPUT
	Details about each request:
	{Protocol used} {Response time} {Remote server} {Extra info}

	EXIT STATUS
	This utility exits with one of the following values:
	0 At least one response was heard.
	2 The transmission was successful but no responses were received.
	1 Any other error occurred.
	`
	targetConcurrency = 5 // stdin inputs
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Only used for debugging.
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelError)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	stdin, err := internal.ReadStdin()
	if err != nil {
		fatal(fmt.Errorf("reading stdin: %w", err))
	}
	var opts internal.Options
	err = opts.Parse()
	if err != nil {
		fatal(fmt.Errorf("parsing options: %w", err))
	}
	if opts.Debug {
		lvl.Set(slog.LevelDebug)
	}
	logger.Debug("Starting ...", "options", opts, "stdin", stdin)
	dnsProtocol := &internal.DNS{Timeout: opts.Timeout}
	if opts.DNSResolver != "" {
		dnsProtocol.Resolver = opts.DNSResolver
	}
	protocols := []internal.Protocol{
		&internal.HTTP{Timeout: opts.Timeout},
		&internal.TCP{Timeout: opts.Timeout},
		dnsProtocol,
	}
	if opts.Protocol != "" {
		var protocol internal.Protocol
		for _, p := range protocols {
			if p.String() == opts.Protocol {
				protocol = p
				break
			}
		}
		if protocol == nil {
			fatal(fmt.Errorf("unknown protocol: %s", opts.Protocol))
		}
		protocols = []internal.Protocol{protocol}
	}
	logger.Info("Starting ...", "protocols", protocols, "count", opts.Count)
	if opts.Help {
		fmt.Fprintf(os.Stderr, "%s\n", appDesc)
		flag.Usage()
		os.Exit(1)
	}
	if opts.NoColor {
		color.NoColor = true
	}
	// To wait for termination signals.
	// - 'Interrupt': Ctrl+C from terminal.
	// - 'SIGTERM': Sent from Kubernetes.
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		logger.Debug("Listening for termination signals")
		<-sigCh
		logger.Debug("Termination signal received")
		cancel()
	}()
	reportCh := make(chan *internal.Report)
	defer close(reportCh)
	var format internal.Format
	switch {
	case opts.JSONOutput:
		format = internal.JSONFormat
	case opts.GrepOutput:
		format = internal.GrepFormat
	default:
		format = internal.HumanFormat
	}
	go func() {
		logger.Debug("Listening for reports ...")
		for report := range reportCh {
			logger.Debug("New report", "report", *report)
			repLine, err := report.String(format)
			if err != nil {
				fatal(err)
			}
			fmt.Println(repLine)
			if report.Error == "" {
				if opts.Stop {
					logger.Debug("Stopping after first successful request")
					cancel()
				}
			}
		}
	}()
	var wg sync.WaitGroup
	if stdin != "" && !opts.NoStdin {
		logger.Debug("Reading from standard input")
		parts := strings.Split(stdin, "\n")
		logger.Debug("Parts", "parts", parts)
		if opts.Protocol == "" {
			fatal(
				errors.New(
					"protocol is required for standard input target reading",
				),
			)
		}
		if opts.Target != "" {
			logger.Debug(
				"Ignoring target from command line",
				"target",
				opts.Target,
			)
		}
		proto := protocols[0]
		ch := make(chan int, targetConcurrency)
		for _, part := range parts {
			if part == "" {
				logger.Debug("Empty part, skipping")
				continue
			}
			wg.Add(1)
			target := strings.TrimSpace(part)
			go func(tg string) {
				defer func() { wg.Done(); <-ch }()
				select {
				case <-ctx.Done():
					logger.Debug("Context cancelled", "target", tg)
					return
				default:
					probe := internal.Probe{
						Proto:    proto,
						Count:    opts.Count,
						Delay:    opts.Delay,
						Logger:   logger,
						ReportCh: reportCh,
						Target:   target,
					}
					logger.Debug("Running ...", "setup", probe)
					err = probe.Do(ctx)
					if err != nil {
						fatal(
							fmt.Errorf(
								"running probe for target %s: %w",
								target,
								err,
							),
						)
					}
				}
			}(target)
		}
	} else {
		for _, protocol := range protocols {
			wg.Add(1)
			go func(proto internal.Protocol) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					logger.Debug("Context cancelled", "protocol", proto)
					return
				default:
					probe := internal.Probe{
						Proto:    proto,
						Count:    opts.Count,
						Delay:    opts.Delay,
						Logger:   logger,
						ReportCh: reportCh,
						Target:   opts.Target,
					}
					logger.Debug("Running ...", "setup", probe)
					err = probe.Do(ctx)
					if err != nil {
						fatal(fmt.Errorf("running probe for protocol %s: %w", proto, err))
					}
				}
			}(protocol)
		}
	}
	wg.Wait()
}

// Prints the error to the standard output and exits with status 1.
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
	os.Exit(1)
}
