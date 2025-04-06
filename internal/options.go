// Package internal provides the core functionality of the application.
package internal

import (
	"errors"
	"flag"
	"time"
)

const targetDesc = "Protocol is required because the format is dependent: URL for HTTP, host:port for TCP, domain for DNS"

// Options are the flags supported by the command line application.
type Options struct {
	// Protocol to use. Example: 'http'.
	Protocol string
	// Where to point the probe.
	// URL (HTTP), host/port string (TCP) or domain (DNS).
	Target string
	// Number of iterations. Zero means infinite.
	Count uint
	// Time to wait for a response.
	Timeout time.Duration
	// Delay between requests.
	Delay time.Duration
	// Stop after the first successful request.
	Stop bool
	// Custom DNS resolver.
	DNSResolver string
	// Output flags.
	// Output in JSON format.
	JSONOutput bool
	// Output in grepable format.
	GrepOutput bool
	// Disable color output.
	NoColor bool
	// Enable debugging.
	Debug bool
	// Show app documentation.
	Help bool
	// Disable stardard input target reading.
	NoStdin bool
}

// Parse fulfills the command line flags provided by the user.
func (opts *Options) Parse() error {
	flag.StringVar(&opts.Protocol, "p", "", "Test only one protocol")
	flag.StringVar(&opts.Target, "tg", "", targetDesc)
	flag.UintVar(&opts.Count, "c", 0, "Number of iterations")
	flag.DurationVar(
		&opts.Timeout, "t", 5*time.Second, "Time to wait for a response",
	)
	flag.DurationVar(
		&opts.Delay, "d", 500*time.Millisecond, "Delay between requests",
	)
	flag.BoolVar(
		&opts.Stop, "s", false, "Stop after the first successful request",
	)
	flag.StringVar(&opts.DNSResolver, "dr", "", "DNS resolution server")
	flag.BoolVar(&opts.JSONOutput, "j", false, "Output in JSON format")
	flag.BoolVar(&opts.GrepOutput, "g", false, "Output in grepable format")
	flag.BoolVar(&opts.NoColor, "nc", false, "Disable color output")
	flag.BoolVar(&opts.Debug, "vv", false, "Verbose output")
	flag.BoolVar(&opts.Help, "h", false, "Show app documentation")
	flag.BoolVar(
		&opts.NoStdin,
		"nstd",
		false,
		"Disable standard input target reading",
	)
	flag.Parse()
	return opts.validate()
}

// Ensures the setup is correct.
func (opts *Options) validate() error {
	if opts.Target != "" && opts.Protocol == "" {
		return errors.New("protocol is required if target is set")
	}
	return nil
}
