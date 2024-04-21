package pkg

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
	bold  = color.New(color.Bold).SprintFunc()
	faint = color.New(color.Faint).SprintFunc()
)

// Report is the result of a connection attempt.
//
// Depending on the result, only one of the properties 'Response' or 'Error'
// is set.
type Report struct {
	// Protocol used to connect to.
	ProtocolID string `json:"protocol"`
	// Target used to connect to.
	RHost string `json:"rhost"`
	// Response time.
	Time time.Duration `json:"time"`
	// Extra information. Depending on the protocol, it could be:
	// - HTTP: Response code.
	// - TCP: Local address.
	// - DNS: Resolved IP addresses.
	Extra string `json:"extra,omitempty"`
	// Network error.
	Error error `json:"error,omitempty"`
}

// String returns a human-readable representation of the report.
func (r *Report) String() string {
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
