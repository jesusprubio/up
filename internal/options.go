// Package internal provides the core functionality of the application.
package internal

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const targetDesc = "Target, format depends on the protocol: URL (http), " +
	"host:port (tcp), domain (dns). Requires a single protocol"

// Options are the flags supported by the command line application.
type Options struct {
	// Protocols to use. Example: 'http'. Empty means all.
	Protocols []string
	// Where to point the probe.
	// URL (HTTP), host/port string (TCP) or domain (DNS).
	Target string
	// Targets passed as positional arguments. Same format as 'Target'.
	Targets []string
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
	// Print only a summary per protocol at the end.
	Summary bool
	// Disable color output.
	NoColor bool
	// Enable debugging.
	Debug bool
	// Show app documentation.
	Help bool
	// Disable stardard input target reading.
	NoStdin bool
}

// Describes a flag registered with a short and a long name.
type flagDef struct {
	short, long, arg, usage string
}

// Parse fulfills the command line flags provided by the user.
func (opts *Options) Parse() error {
	return opts.parse(flag.CommandLine, nil)
}

// Parses the arguments using the given flag set. If 'args' is nil the
// command line arguments ('os.Args') are used.
func (opts *Options) parse(fs *flag.FlagSet, args []string) error {
	var protocol string
	defs := []flagDef{}
	str := func(p *string, short, long, arg, value, usage string) {
		fs.StringVar(p, short, value, usage)
		fs.StringVar(p, long, value, usage)
		defs = append(defs, flagDef{short, long, arg, usage})
	}
	boolean := func(p *bool, short, long, usage string) {
		fs.BoolVar(p, short, false, usage)
		fs.BoolVar(p, long, false, usage)
		defs = append(defs, flagDef{short, long, "", usage})
	}
	str(&protocol, "p", "protocol", "list", "",
		"Protocols to test, comma separated: http,tcp,dns (default all)")
	str(&opts.Target, "tg", "target", "string", "", targetDesc)
	fs.UintVar(&opts.Count, "c", 0, "")
	fs.UintVar(&opts.Count, "count", 0, "")
	defs = append(defs, flagDef{
		"c", "count", "uint", "Number of iterations (default 0, infinite)",
	})
	fs.DurationVar(&opts.Timeout, "t", 5*time.Second, "")
	fs.DurationVar(&opts.Timeout, "timeout", 5*time.Second, "")
	defs = append(defs, flagDef{
		"t", "timeout", "duration", "Time to wait for a response (default 5s)",
	})
	fs.DurationVar(&opts.Delay, "d", 500*time.Millisecond, "")
	fs.DurationVar(&opts.Delay, "delay", 500*time.Millisecond, "")
	defs = append(defs, flagDef{
		"d", "delay", "duration", "Delay between requests (default 500ms)",
	})
	boolean(&opts.Stop, "s", "stop", "Stop after the first successful request")
	str(&opts.DNSResolver, "dr", "dns-resolver", "ip", "",
		"DNS resolution server")
	boolean(&opts.JSONOutput, "j", "json", "Output in JSON lines format")
	boolean(&opts.GrepOutput, "g", "grep", "Output in grepable format")
	boolean(&opts.Summary, "sum", "summary",
		"Print only a summary per protocol at the end")
	boolean(&opts.NoColor, "nc", "no-color", "Disable color output")
	boolean(&opts.Debug, "vv", "verbose", "Verbose output")
	boolean(&opts.Help, "h", "help", "Show app documentation")
	boolean(&opts.NoStdin, "nstd", "no-stdin",
		"Disable standard input target reading")
	fs.Usage = func() { printUsage(fs.Output(), filepath.Base(fs.Name()), defs) }
	if args == nil {
		args = os.Args[1:]
	}
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	opts.Protocols = splitList(protocol)
	opts.Targets = fs.Args()
	return opts.validate()
}

// Prints the usage with the short and long name of each flag.
func printUsage(w io.Writer, name string, defs []flagDef) {
	fmt.Fprintf(w, "Usage: %s [flags] [target ...]\n\nFlags:\n", name)
	for _, d := range defs {
		names := fmt.Sprintf("-%s, --%s", d.short, d.long)
		if d.arg != "" {
			names += " " + d.arg
		}
		fmt.Fprintf(w, "  %-28s %s\n", names, d.usage)
	}
}

// Splits a comma separated list, ignoring empty items and spaces.
func splitList(s string) []string {
	var items []string
	for item := range strings.SplitSeq(s, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}

// Ensures the setup is correct.
func (opts *Options) validate() error {
	hasTargets := opts.Target != "" || len(opts.Targets) > 0
	if hasTargets && len(opts.Protocols) != 1 {
		return errors.New("a single protocol is required if a target is set")
	}
	if opts.JSONOutput && opts.GrepOutput {
		return errors.New("only one output format can be selected")
	}
	return nil
}
