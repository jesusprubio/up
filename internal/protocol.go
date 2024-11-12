package internal

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Default timeout for the probes.
const timeout = 5 * time.Second

// Protocol defines a probe attempt.
type Protocol interface {
	// Returns the identifier.
	// Example: "http".
	String() string
	// Attempt to check the connectivity to the target.
	// The target depends on the protocol. For example, for HTTP it's a URL.
	// Returns the used target or error if the attempt failed. Some protocols
	// include an additional string with extra information. For example, the
	// HTTP protocol returns the status code.
	Probe(target string) (string, string, error)
}

// HTTP protocol implementation.
type HTTP struct {
	Timeout time.Duration
}

// String returns the identifier of the protocol.
func (h *HTTP) String() string {
	return "http"
}

// Probe makes an HTTP request to a random captive portal.
//
// The target is a URL.
// The extra data is the status code.
func (h *HTTP) Probe(target string) (string, string, error) {
	cli := &http.Client{Timeout: h.Timeout}
	url := target
	if url == "" {
		var err error
		url, err = RandomCaptivePortal()
		if err != nil {
			return "", "", fmt.Errorf("selecting captive portal: %w", err)
		}
	}
	resp, err := cli.Get(url)
	if err != nil {
		return "", "", err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", "", fmt.Errorf("closing response body: %w", err)
	}
	return url, resp.Status, nil
}

// TCP protocol implementation.
type TCP struct {
	Timeout time.Duration
}

// String returns the identifier of the protocol.
func (t *TCP) String() string {
	return "tcp"
}

// Probe makes a TCP request to a random server.
//
// The target is a host:port.
// The extra data is the local interface.
func (t *TCP) Probe(target string) (string, string, error) {
	hostPort := target
	if hostPort == "" {
		var err error
		hostPort, err = RandomTCPServer()
		if err != nil {
			return "", "", fmt.Errorf("selecting TCP server: %w", err)
		}
	}
	conn, err := net.DialTimeout("tcp", hostPort, t.Timeout)
	if err != nil {
		return "", "", err
	}
	err = conn.Close()
	if err != nil {
		return "", "", fmt.Errorf("closing connection: %w", err)
	}
	return hostPort, conn.LocalAddr().String(), nil
}

// DNS protocol implementation.
type DNS struct {
	Timeout time.Duration
	// Custom DNS resolver.
	Resolver string
}

// String returns the identifier of the protocol.
func (d *DNS) String() string {
	return "dns"
}

// Probe resolves a random domain name.
//
// The target is a domain name.
// The extra data is the first resolved IP address.
func (d *DNS) Probe(target string) (string, string, error) {
	var r net.Resolver
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if d.Resolver != "" {
		r.PreferGo = true
		r.Dial = func(ctx context.Context, network, address string) (
			net.Conn, error,
		) {
			nd := net.Dialer{Timeout: d.Timeout}
			return nd.DialContext(ctx, network, fmt.Sprintf(
				"%s:%s", d.Resolver, "53",
			))
		}
	}
	domain := target
	if domain == "" {
		var err error
		domain, err = RandomDomain()
		if err != nil {
			return "", "", fmt.Errorf("selecting domain: %w", err)
		}
	}
	addrs, err := r.LookupHost(ctx, domain)
	if err != nil {
		return "", "", err
	}
	return domain, addrs[0], nil
}
