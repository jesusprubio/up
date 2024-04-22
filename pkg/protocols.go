package pkg

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// Protocols included in the library.
var Protocols = []*Protocol{
	{ID: "tcp", Request: requestTCP, RHost: RandomTCPServer},
	{ID: "dns", Request: requestDNS, RHost: RandomDomain},
	{ID: "http", Request: requestHTTP, RHost: RandomCaptivePortal},
}

// ProtocolByID returns a protocol from the list.
func ProtocolByID(id string) *Protocol {
	for _, p := range Protocols {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// Protocol defines a probe attempt.
type Protocol struct {
	ID string
	// Function to make the probe.
	// Returns extra information about the attempt or an error if it failed.
	Request func(rhost string, timeout time.Duration) (string, error)
	// Function to create a random remote
	RHost func() (string, error)
}

// String returns an human-readable representation of the protocol.
func (p *Protocol) String() string {
	return p.ID
}

// Makes an HTTP request.
//
// The extra information is the status code.
func requestHTTP(u string, timeout time.Duration) (string, error) {
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
// The extra information is the local address used to make the request.
func requestTCP(hostPort string, timeout time.Duration) (string, error) {
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
func requestDNS(domain string, timeout time.Duration) (string, error) {
	addrs, err := net.LookupHost(domain)
	if err != nil {
		return "", fmt.Errorf("resolving %s: %w", domain, err)
	}
	return fmt.Sprint(addrs[0]), nil
}
