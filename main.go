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
	"time"

	"github.com/fatih/color"
	"github.com/jesusprubio/up/pkg"
)

// TODO(#39): STDIN piped input.

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
	var opts options
	opts.parse()
	if opts.verbose {
		// TODO(#37): Debug and verbose should not be the same thing.
		lvl.Set(slog.LevelDebug)
	}
	logger.Debug("Starting", "options", opts)
	protocols := pkg.Protocols
	if opts.protocol != "" {
		protocol := pkg.ProtocolByID(opts.protocol)
		if protocol == nil {
			fatal(fmt.Errorf("unknown protocol: %s", opts.protocol))
		}
		protocols = []*pkg.Protocol{protocol}
	}
	logger.Info("Running", "protocols", protocols, "count", opts.count)
	if opts.help {
		fmt.Fprintf(os.Stderr, "%s\n", appDesc)
		flag.Usage()
		os.Exit(1)
	}
	if opts.noColor {
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
		Count:     opts.count,
		Timeout:   opts.timeout,
		Delay:     opts.delay,
		Logger:    logger,
		ReportCh:  make(chan *pkg.Report),
	}
	go func() {
		logger.Debug("Listening for reports")
		for report := range probe.ReportCh {
			logger.Debug("New report", "report", *report)
			var line string
			if opts.jsonOutput {
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
				if opts.stop {
					logger.Debug("Stopping after first successful request")
					cancel()
				}
			}
		}
	}()
	logger.Debug("Running", "setup", probe)
	err := probe.Run(ctx)
	if err != nil {
		fatal(fmt.Errorf("running probe: %w", err))
	}
	logger.Debug("Bye!")
}

// Flags passed by the user.
type options struct {
	// Input flags.
	// Protocol to use.
	protocol string
	// Number of iterations. Zero means infinite.
	count uint
	// Time to wait for a response.
	timeout time.Duration
	// Delay between requests.
	delay time.Duration
	// Stop after the first successful request.
	stop bool
	// Output flags.
	// Output in JSON format.
	jsonOutput bool
	// Disable color output.
	noColor bool
	// Verbose output.
	verbose bool
	// Show app documentation.
	help bool
}

// Parses the command line flags provided by the user.
func (opts *options) parse() {
	flag.StringVar(&opts.protocol, "p", "", "Test only one protocol")
	flag.UintVar(&opts.count, "c", 0, "Number of iterations")
	flag.DurationVar(
		&opts.timeout,
		"t",
		5*time.Second,
		"Time to wait for a response",
	)
	flag.DurationVar(
		&opts.delay,
		"d",
		500*time.Millisecond,
		"Delay between requests",
	)
	flag.BoolVar(
		&opts.stop,
		"s",
		false,
		"Stop after the first successful request",
	)
	flag.BoolVar(&opts.jsonOutput, "j", false, "Output in JSON format")
	flag.BoolVar(&opts.noColor, "nc", false, "Disable color output")
	flag.BoolVar(&opts.verbose, "v", false, "Verbose output")
	flag.BoolVar(&opts.help, "h", false, "Show app documentation")
	flag.Parse()
}

// Logs the error to the standard output and exits with status 1.
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
	os.Exit(1)
}

// String returns a human-readable representation of the report.
func reportToLine(r *pkg.Report) string {
	// TODO(#40): Use Go string padding.
	line := fmt.Sprintf("%s\t%s\t%s", bold(r.ProtocolID), r.Time, r.RHost)
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
