package internal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
)

// Format is the output format of the report.
type Format int

const (
	HumanFormat Format = iota
	JSONFormat
	GrepFormat
)

// Report is the result of a connection attempt.
//
// Only one of the properties 'Response' or 'Error' is set.
type Report struct {
	// Protocol used to connect to.
	ProtocolID string `json:"protocol"`
	// Target used to connect to.
	RHost string `json:"rhost"`
	// Response time.
	Time time.Duration `json:"time"`
	// Network error.
	Error string `json:"error,omitempty"`
	// Extra information. Depends on the protocol.
	Extra string `json:"extra,omitempty"`
}

// String returns a string representation of the report.
func (r *Report) String(format Format) (string, error) {
	switch format {
	case HumanFormat:
		return r.stringHuman(), nil
	case JSONFormat:
		line, err := r.stringJSON()
		if err != nil {
			return "", fmt.Errorf("error generating JSON report: %w", err)
		}
		return line, nil
	case GrepFormat:
		return r.stringGrep(), nil
	default:
		return "", fmt.Errorf("unsupported format: %v", format)
	}
}

// Returns the report in JSON format.
// Example:
// '{"protocol":"tcp","rhost":"64.6.65.6:53","time":13433165,"extra":"192.168.1.177:39384"}'
func (r *Report) stringJSON() (string, error) {
	reportJSON, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("marshaling report: %w", err)
	}
	return string(reportJSON), nil
}

// Returns the report in human readable format.
// Example: '✔ tcp    100.077875ms   77.88.8.8:53 (192.168.1.177:43586)
func (r *Report) stringHuman() string {
	line := fmt.Sprintf("%-15s %-14s %s", bold(r.ProtocolID), r.Time, r.RHost)
	suffix := r.Extra
	prefix := green("✔")
	if r.Error != "" {
		prefix = red("✘")
		suffix = r.Error
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

// Returns the report in a grepable format.
//
// Example: 'tcp     13.944825ms     195.46.39.40:53 success 192.168.1.177:43296
func (r *Report) stringGrep() string {
	status := "ok"
	if r.Error != "" {
		status = "error"
	}
	suffix := r.Extra
	if r.Error != "" {
		suffix = r.Error
	}
	line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
		r.ProtocolID, r.Time, r.RHost, status, suffix,
	)
	return line
}
