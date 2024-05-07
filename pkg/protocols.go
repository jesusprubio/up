package pkg

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// Protocols included in the library.
var Protocols = []*Protocol{
	{ID: "http", Probe: httpProbe, RHost: RandomCaptivePortal},
	{ID: "tcp", Probe: tcpProbe, RHost: RandomTCPServer},
	{ID: "dns", Probe: dnsProbe, RHost: RandomDomain},
}

// Protocol defines a probe attempt.
type Protocol struct {
	ID string
	// Probe implementation for this protocol.
	// Returns extra information about the attempt or an error if it failed.
	Probe func(rhost string, timeout time.Duration) (string, error)
	// Function to create a random remote
	RHost func() (string, error)
}

// String returns an human-readable representation of the protocol.
func (p *Protocol) String() string {
	return p.ID
}

// Ensures the required properties are set.
func (p *Protocol) validate() error {
	if p.Probe == nil {
		return fmt.Errorf(tmplRequiredProp, "Probe")
	}
	if p.RHost == nil {
		return fmt.Errorf(tmplRequiredProp, "RHost")
	}
	return nil
}

// Makes an HTTP request.
//
// The extra information is the status code.
func httpProbe(u string, timeout time.Duration) (string, error) {
	cli := &http.Client{Timeout: timeout}
	resp, err := cli.Get(u)
	if err != nil {
		return "", fmt.Errorf("making request to %s: %w", u, err)
	}
	err = resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("closing response body: %w", err)
	}
	return resp.Status, nil
}

// Makes a TCP request.
//
// The extra information is the local host/port.
func tcpProbe(hostPort string, timeout time.Duration) (string, error) {
	conn, err := net.DialTimeout("tcp", hostPort, timeout)
	if err != nil {
		return "", fmt.Errorf("making request to %s: %w", hostPort, err)
	}
	err = conn.Close()
	if err != nil {
		return "", fmt.Errorf("closing connection: %w", err)
	}
	return conn.LocalAddr().String(), nil
}

// Resolves a domain name.
//
// The extra information is the first resolved IP address.
// TODO(#31)
func dnsProbe(domain string, timeout time.Duration) (string, error) {
	addrs, err := net.LookupHost(domain)
	if err != nil {
		return "", fmt.Errorf("resolving %s: %w", domain, err)
	}
	return fmt.Sprint(addrs[0]), nil
}
