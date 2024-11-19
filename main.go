// Package main implements a simple CLI to use the library.
package main

import (
	"context"
	"encoding/json"
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
	stdin, err := internal.ReadStdin()
	if err != nil {
		fatal(fmt.Errorf("reading stdin: %w", err))
	}
	inputs, err := internal.ProcessInputs(stdin)
	if err != nil {
		fmt.Printf("failed to process the inputs:\n%v\n", err)
	}
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
		Input:     inputs,
	}
	go func() {
		logger.Debug("Listening for reports ...")
		for report := range probe.ReportCh {
			logger.Debug("New report", "report", *report)
			var line string
			if opts.JSONOutput {
				reportJSON, err := json.Marshal(report)
				if err != nil {
					fatal(fmt.Errorf("marshaling report: %w", err))
				}
				line = string(reportJSON)
			} else {
				line = reportToLine(report)
			}
			fmt.Println(line)
			if report.Error == nil {
				if opts.Stop {
					logger.Debug("Stopping after first successful request")
					cancel()
				}
			}
		}
	}()

	logger.Debug("Running ...", "setup", probe)
	err = probe.Do(ctx)
	if err != nil {
		fatal(fmt.Errorf("running probe: %w", err))
	}

	if err != nil {
		fatal(fmt.Errorf("running probe: %w", err))
	}
}

// Logs the error to the standard output and exits with status 1.
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
	os.Exit(1)
}

// Returns a human-readable representation of the report.
func reportToLine(r *internal.Report) string {
	line := fmt.Sprintf("%-15s %-14s %s", bold(r.ProtocolID), r.Time, r.RHost)
	suffix := r.Extra
	prefix := green("✔")
	if r.Error != nil {
		prefix = red("✘")
		suffix = r.Error.Error()
	}
	suffix = fmt.Sprintf("(%s)", suffix)
	return fmt.Sprintf("%s %s %s", prefix, line, faint(suffix))
}

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
	bold  = color.New(color.Bold).SprintFunc()
	faint = color.New(color.Faint).SprintFunc()
)
