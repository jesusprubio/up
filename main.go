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
	"github.com/jesusprubio/up/pkg"
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
	protocols := pkg.Protocols
	if opts.Protocol != "" {
		protocol := internal.ProtocolByID(opts.Protocol)
		if opts.DNSResolver != "" {
			protocol.WithDNSResolver(opts.DNSResolver)
		}
		if protocol == nil {
			internal.Fatal(fmt.Errorf("unknown protocol: %s", opts.Protocol))
		}
		protocols = []*pkg.Protocol{protocol}
	}
	logger.Info("Starting ...", "protocols", protocols, "count", opts.Count)
	if opts.Help {
		fmt.Fprintf(os.Stderr, "%s\n", internal.AppDesc)
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
	reportCh := make(chan *pkg.Report)
	defer close(reportCh)
	probe := pkg.Probe{
		Protocols: protocols,
		Count:     opts.Count,
		Timeout:   opts.Timeout,
		Delay:     opts.Delay,
		Logger:    logger,
		ReportCh:  reportCh,
	}
	go func() {
		logger.Debug("Listening for reports ...")
		for report := range probe.ReportCh {
			logger.Debug("New report", "report", *report)
			var line string
			if opts.JSONOutput {
				reportJSON, err := json.Marshal(report)
				if err != nil {
					internal.Fatal(fmt.Errorf("marshaling report: %w", err))
				}
				line = string(reportJSON)
			} else {
				line = internal.ReportToLine(report)
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
	err := probe.Run(ctx)
	if err != nil {
		internal.Fatal(fmt.Errorf("running probe: %w", err))
	}
}
