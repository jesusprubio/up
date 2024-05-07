package internal

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/jesusprubio/up/pkg"
)

const (
	appName = "up"
	AppDesc = `
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

// ProtocolByID returns the protocol implementation whose ID matches the given
// one.
func ProtocolByID(id string) *pkg.Protocol {
	for _, p := range pkg.Protocols {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// Fatal logs the error to the standard output and exits with status 1.
func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
	os.Exit(1)
}

// ReportToLine returns a human-readable representation of the report.
func ReportToLine(r *pkg.Report) string {
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
