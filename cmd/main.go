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
	"time"

	"github.com/fatih/color"
	"github.com/jesusprubio/up/pkg"
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
	1 Any other value An error occurred.
	`
)

// TODO: Support passing a custom remote server as an argument.

// CLI initialization.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Input options.
	flagProto := flag.String(
		"p",
		"",
		fmt.Sprintf("Use only one protocol: %v", pkg.Protocols),
	)
	flagCount := flag.Uint("c", 0, "Number of iterations. (0 = infinite)")
	flagTimeout := flag.Duration("t",
		5*time.Second,
		"Time to wait for a response",
	)
	flagDelay := flag.Duration(
		"d",
		500*time.Millisecond,
		"Delay between requests",
	)
	flagStop := flag.Bool("s", false, "Stop after the first successful request")
	// Output options.
	flagJSONOutput := flag.Bool("j", false, "Output in JSON format")
	flagNoColor := flag.Bool("nc", false, "Disable color output")
	flagVerbose := flag.Bool("v", false, "Verbose output")
	flagHelp := flag.Bool("h", false, "Show app documentation")
	flag.Parse()
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelError)
	// Only used for debugging.
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	if *flagVerbose {
		lvl.Set(slog.LevelDebug)
	}
	logger.Debug("Running",
		"protocol", *flagProto,
		"count", *flagCount,
		"timeout", *flagTimeout,
		"delay", *flagDelay,
		"stop", *flagStop,
		"json", *flagJSONOutput,
		"no-color", *flagNoColor,
		"verbose", *flagVerbose,
		"help", *flagHelp,
	)
	protocols := pkg.Protocols
	if *flagProto != "" {
		protocol := pkg.ProtocolByID(*flagProto)
		if protocol == nil {
			fatal(fmt.Errorf("unknown protocol: %s", *flagProto))
		}
		protocols = []*pkg.Protocol{protocol}
	}
	logger.Info("Running", "protocols", protocols, "count", *flagCount)
	if *flagHelp {
		fmt.Fprintf(os.Stderr, "%s\n", appDesc)
		flag.Usage()
		os.Exit(1)
	}
	if *flagNoColor {
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
	probe := pkg.Probe{
		Protocols: protocols,
		Count:     *flagCount,
		Delay:     *flagDelay,
		Timeout:   *flagTimeout,
		Logger:    logger,
		ReportCh:  make(chan *pkg.Report),
	}
	go func() {
		logger.Debug("Listening for reports")
		for report := range probe.ReportCh {
			logger.Debug("New report", "report", *report)
			var line string
			if *flagJSONOutput {
				reportJSON, err := json.Marshal(report)
				if err != nil {
					fatal(fmt.Errorf("marshaling report: %w", err))
				}
				line = string(reportJSON)
			} else {
				line = report.String()
			}
			fmt.Println(line)
			if report.Error == nil {
				if *flagStop {
					logger.Debug("Stop requested")
					cancel()
				}
			}
		}
	}()
	logger.Debug("Running", "setup", probe)
	err := probe.Run(ctx)
	if err != nil {
		fatal(err)
	}
	logger.Debug("Bye!")
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
	os.Exit(1)
}
