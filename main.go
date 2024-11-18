// Package main implements a simple CLI to use the library.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
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
)

// TODO(#39): STDIN piped input.

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Only used for debugging.
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelError)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	var opts internal.Options
	opts.Parse()
	if opts.Debug {
		lvl.Set(slog.LevelDebug)
	}
	logger.Debug("Starting ...", "options", opts)
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
	probe := internal.Probe{
		Protocols: protocols,
		Count:     opts.Count,
		Delay:     opts.Delay,
		Logger:    logger,
		ReportCh:  reportCh,
	}
	go func() {
		logger.Debug("Listening for reports ...")
		// Print Report Lines
		for report := range probe.ReportCh {
			logger.Debug("New report", "report", *report)

			format := internal.HumanFormat // Default format
			if opts.JSONOutput {
				format = internal.JSONFormat
			} else if opts.GrepFormat {
				format = internal.GrepFormat
			}
			report.PrintFormatted(format)

			if report.Error == nil {
				if opts.Stop {
					logger.Debug("Stopping after first successful request")
					cancel()
				}
			}
		}
	}()
	logger.Debug("Running ...", "setup", probe)
	err := probe.Do(ctx)
	if err != nil {
		fatal(fmt.Errorf("running probe: %w", err))
	}
}

// Logs the error to the standard output and exits with status 1.
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
	os.Exit(1)
}
