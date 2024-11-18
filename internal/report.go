package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
)

// Report is the result of a connection attempt.
//
// Only one of the properties 'Response' or 'Error' is set.

type Format int

const (
	JSONFormat Format = iota
	HumanFormat
	GrepFormat
)

type Report struct {
	// Protocol used to connect to.
	ProtocolID string `json:"protocol"`
	// Target used to connect to.
	RHost string `json:"rhost"`
	// Response time.
	Time time.Duration `json:"time"`
	// Network error.
	Error error `json:"error,omitempty"`
	// Extra information. Depends on the protocol.
	Extra  string `json:"extra,omitempty"`
	Format Format
}

func (r *Report) PrintFormatted(f Format) {

	switch f {
	case HumanFormat:
		printHumanFormat(r)
	case JSONFormat:
		printJSONFormat(r)
	case GrepFormat:
		printGrepableFormat(r)
	}

}
func printJSONFormat(r *Report) {

	reportJSON, err := json.Marshal(r)
	if err != nil {
		log.Fatal(fmt.Errorf("marshaling report: %w", err))
	}
	fmt.Println(string(reportJSON))
}

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
	bold  = color.New(color.Bold).SprintFunc()
	faint = color.New(color.Faint).SprintFunc()
)

func printHumanFormat(r *Report) {

	line := fmt.Sprintf("%-15s %-14s %s", bold(r.ProtocolID), r.Time, r.RHost)

	suffix := r.Extra
	prefix := green("✔")
	if r.Error != nil {
		prefix = red("✘")
		suffix = r.Error.Error()
	}
	suffix = fmt.Sprintf("(%s)", suffix)

	fmt.Printf("%s %s %s\n", prefix, line, faint(suffix))
}

// Output: HTTP/1.1    2024-11-18T15:00:00Z    192.168.1.1    success    Request processed successfully
func printGrepableFormat(r *Report) {
	status := "success"
	if r.Error != nil {
		status = "failure"
	}

	line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
		r.ProtocolID,
		r.Time,
		r.RHost,
		status,
		r.Extra,
	)
	fmt.Println(line)
}
