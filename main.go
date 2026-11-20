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
	"slices"
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

	JSON lines (-j), one object per request:
	{"protocol","target","ok","time_ms","error","extra"}

	Summary (-sum), only one line per protocol at the end. With -j a single
	object: {"up","protocols":[{"protocol","up","ok","failed","avg_ms",
	"last_error"}]}

	TARGETS
	Targets can be passed with -tg, as arguments after the flags or through
	the standard input (one per line). They require a single protocol (-p).

	EXIT STATUS
	This utility exits with one of the following values:
	0 At least one response was heard.
	2 The transmission was successful but no responses were received.
	1 Any other error occurred.
	`
	targetConcurrency = 5 // stdin inputs
	exitNoResponse    = 2
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
	err := opts.Parse()
	if err != nil {
		fatal(fmt.Errorf("parsing options: %w", err))
	}
	if opts.Help {
		fmt.Fprintf(os.Stderr, "%s\n", appDesc)
		flag.CommandLine.Usage()
		os.Exit(1)
	}
	if opts.Debug {
		lvl.Set(slog.LevelDebug)
	}
	var stdin string
	if !opts.NoStdin {
		stdin, err = internal.ReadStdin()
		if err != nil {
			fatal(fmt.Errorf("reading stdin: %w", err))
		}
	}
	logger.Debug("Starting ...", "options", opts, "stdin", stdin)
	protocols, err := selectProtocols(opts)
	if err != nil {
		fatal(err)
	}
	targets := opts.Targets
	for line := range strings.SplitSeq(stdin, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			targets = append(targets, line)
		}
	}
	if len(targets) > 0 && len(opts.Protocols) != 1 {
		fatal(errors.New("a single protocol is required to read targets"))
	}
	if opts.Target != "" {
		if len(targets) > 0 {
			logger.Debug("Ignoring target from flag", "target", opts.Target)
		} else {
			targets = []string{opts.Target}
		}
	}
	logger.Info("Starting ...", "protocols", protocols, "count", opts.Count)
	if opts.NoColor {
		color.NoColor = true
	}
	// To wait for termination signals.
	// - 'Interrupt': Ctrl+C from terminal.
	// - 'SIGTERM': Sent from Kubernetes.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		logger.Debug("Listening for termination signals")
		<-sigCh
		logger.Debug("Termination signal received")
		cancel()
	}()
	var format internal.Format
	switch {
	case opts.JSONOutput:
		format = internal.JSONFormat
	case opts.GrepOutput:
		format = internal.GrepFormat
	default:
		format = internal.HumanFormat
	}
	reportCh := make(chan *internal.Report)
	printerDone := make(chan struct{})
	var summary internal.Summary
	go func() {
		defer close(printerDone)
		logger.Debug("Listening for reports ...")
		for report := range reportCh {
			logger.Debug("New report", "report", *report)
			summary.Add(report)
			if !opts.Summary {
				repLine, err := report.String(format)
				if err != nil {
					fatal(err)
				}
				fmt.Println(repLine)
			}
			if report.Error == "" && opts.Stop {
				logger.Debug("Stopping after first successful request")
				cancel()
			}
		}
	}()
	var wg sync.WaitGroup
	run := func(probe internal.Probe) {
		defer wg.Done()
		if ctx.Err() != nil {
			logger.Debug("Context cancelled", "target", probe.Target)
			return
		}
		logger.Debug("Running ...", "setup", probe)
		err := probe.Do(ctx)
		if err != nil {
			fatal(fmt.Errorf(
				"running probe for protocol %s and target %q: %w",
				probe.Proto, probe.Target, err,
			))
		}
	}
	newProbe := func(proto internal.Protocol, target string) internal.Probe {
		return internal.Probe{
			Proto:    proto,
			Count:    opts.Count,
			Delay:    opts.Delay,
			Logger:   logger,
			ReportCh: reportCh,
			Target:   target,
		}
	}
	if len(targets) > 0 {
		logger.Debug("Probing targets", "targets", targets)
		sem := make(chan struct{}, targetConcurrency)
		for _, target := range targets {
			wg.Add(1)
			sem <- struct{}{}
			go func() {
				defer func() { <-sem }()
				run(newProbe(protocols[0], target))
			}()
		}
	} else {
		for _, proto := range protocols {
			wg.Add(1)
			go run(newProbe(proto, ""))
		}
	}
	wg.Wait()
	close(reportCh)
	<-printerDone
	if opts.Summary {
		line, err := summary.String(format)
		if err != nil {
			fatal(err)
		}
		if line != "" {
			fmt.Println(line)
		}
	}
	if !summary.Up() {
		os.Exit(exitNoResponse)
	}
}

// Returns the protocols selected by the user, all of them by default.
func selectProtocols(opts internal.Options) ([]internal.Protocol, error) {
	dnsProtocol := &internal.DNS{Timeout: opts.Timeout}
	if opts.DNSResolver != "" {
		dnsProtocol.Resolver = opts.DNSResolver
	}
	all := []internal.Protocol{
		&internal.HTTP{Timeout: opts.Timeout},
		&internal.TCP{Timeout: opts.Timeout},
		dnsProtocol,
	}
	if len(opts.Protocols) == 0 {
		return all, nil
	}
	selected := []internal.Protocol{}
	for _, id := range opts.Protocols {
		i := slices.IndexFunc(all, func(p internal.Protocol) bool {
			return p.String() == id
		})
		if i == -1 {
			return nil, fmt.Errorf("unknown protocol: %s", id)
		}
		if !slices.Contains(selected, all[i]) {
			selected = append(selected, all[i])
		}
	}
	return selected, nil
}

// Prints the error to the standard output and exits with status 1.
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
	os.Exit(1)
}
