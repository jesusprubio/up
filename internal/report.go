package internal

import (
	"time"
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
	Error error `json:"error,omitempty"`
	// Extra information. Depends on the protocol.
	Extra string `json:"extra,omitempty"`
}
