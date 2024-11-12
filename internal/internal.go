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

// Fatal logs the error to the standard output and exits with status 1.
func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", appName, err)
	os.Exit(1)
}

// ReportToLine returns a human-readable representation of the report.
func ReportToLine(r *pkg.Report) string {
	symbol := green("✔")
	suffix := r.RHost
	if r.Error != nil {
		symbol = red("✘")
		suffix = r.Error.Error()
	}
	return fmt.Sprintf("%s %s", symbol, fmt.Sprintf(
		"%-15s %-14s %-15s", bold(r.ProtocolID), r.Time, faint(suffix),
	))
}

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
	bold  = color.New(color.Bold).SprintFunc()
	faint = color.New(color.Faint).SprintFunc()
)
