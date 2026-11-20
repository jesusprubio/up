package internal

import (
	"flag"
	"io"
	"slices"
	"testing"
	"time"
)

func parseTestOptions(args ...string) (Options, error) {
	fs := flag.NewFlagSet("up", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	if args == nil {
		args = []string{} // nil means 'os.Args'.
	}
	var opts Options
	err := opts.parse(fs, args)
	return opts, err
}

func TestOptionsParse(t *testing.T) {
	t.Run("supports short flags", func(t *testing.T) {
		opts, err := parseTestOptions("-p", "http", "-c", "2", "-j", "-sum")
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(opts.Protocols, []string{"http"}) || opts.Count != 2 ||
			!opts.JSONOutput || !opts.Summary {
			t.Fatalf("unexpected options: %+v", opts)
		}
	})
	t.Run("supports long flags", func(t *testing.T) {
		opts, err := parseTestOptions(
			"--protocol", "dns", "--count", "3", "--timeout", "1s",
			"--grep", "--summary", "--no-stdin", "--dns-resolver", "1.1.1.1",
		)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(opts.Protocols, []string{"dns"}) || opts.Count != 3 ||
			opts.Timeout != time.Second || !opts.GrepOutput || !opts.Summary ||
			!opts.NoStdin || opts.DNSResolver != "1.1.1.1" {
			t.Fatalf("unexpected options: %+v", opts)
		}
	})
	t.Run("uses the defaults", func(t *testing.T) {
		opts, err := parseTestOptions()
		if err != nil {
			t.Fatal(err)
		}
		if opts.Protocols != nil || opts.Count != 0 ||
			opts.Timeout != 5*time.Second ||
			opts.Delay != 500*time.Millisecond {
			t.Fatalf("unexpected options: %+v", opts)
		}
	})
	t.Run("splits a protocol list", func(t *testing.T) {
		opts, err := parseTestOptions("-p", "http, dns,")
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"http", "dns"}
		if !slices.Equal(opts.Protocols, want) {
			t.Fatalf("got %q, want %q", opts.Protocols, want)
		}
	})
	t.Run("reads targets from the arguments", func(t *testing.T) {
		opts, err := parseTestOptions("-p", "dns", "a.com", "b.com")
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"a.com", "b.com"}
		if !slices.Equal(opts.Targets, want) {
			t.Fatalf("got %q, want %q", opts.Targets, want)
		}
	})
	errTests := []struct {
		name string
		args []string
		want string
	}{
		{
			"returns an error for a target without protocol",
			[]string{"-tg", "a.com"},
			"a single protocol is required if a target is set",
		},
		{
			"returns an error for targets with many protocols",
			[]string{"-p", "http,dns", "a.com"},
			"a single protocol is required if a target is set",
		},
		{
			"returns an error for many output formats",
			[]string{"-j", "-g"},
			"only one output format can be selected",
		},
	}
	for _, tt := range errTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseTestOptions(tt.args...)
			if err == nil || err.Error() != tt.want {
				t.Fatalf("got %v, want %q", err, tt.want)
			}
		})
	}
}
