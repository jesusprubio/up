package internal

import (
	"flag"
	"time"
)

// Options are the flags supported by the command line application.
type Options struct {
	// Input flags.
	// Protocol to use. Example: 'http'.
	Protocol string
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
	// Disable color output.
	NoColor bool
	// Enable debugging.
	Debug bool
	// Show app documentation.
	Help bool
}

// Parse fulfills the command line flags provided by the user.
func (opts *Options) Parse() {
	flag.StringVar(&opts.Protocol, "p", "", "Test only one protocol")
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
	flag.StringVar(&opts.DNSResolver, "r", "", "DNS resolution server")
	flag.BoolVar(&opts.JSONOutput, "j", false, "Output in JSON format")
	flag.BoolVar(&opts.NoColor, "nc", false, "Disable color output")
	flag.BoolVar(&opts.Debug, "dbg", false, "Verbose output")
	flag.BoolVar(&opts.Help, "h", false, "Show app documentation")
	flag.Parse()
}
