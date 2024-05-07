package pkg

import (
	"fmt"
	"net"
	"net/url"
	"testing"
)

func TestRandomCaptivePortal(t *testing.T) {
	got, err := RandomCaptivePortal()
	if err != nil {
		t.Fatal(err)
	}
	_, err = url.Parse(got)
	if err != nil {
		t.Fatalf("invalid URL: %s", got)
	}
}

func TestRandomDNSServer(t *testing.T) {
	got, err := RandomDNSServer()
	if err != nil {
		t.Fatal(err)
	}
	ip := net.ParseIP(got)
	if ip == nil {
		t.Fatalf("invalid IP: %s", got)
	}
}

func TestRandomTCPServer(t *testing.T) {
	got, err := RandomTCPServer()
	if err != nil {
		t.Fatal(err)
	}
	_, port, err := net.SplitHostPort(got)
	if err != nil {
		t.Fatalf("invalid host/port: %s", got)
	}
	if port != "53" {
		t.Fatalf("invalid port: %s", port)
	}
}

func TestRandomDomain(t *testing.T) {
	got, err := RandomDomain()
	if err != nil {
		t.Fatal(err)
	}
	u := fmt.Sprintf("http://%s", got)
	_, err = url.Parse(u)
	if err != nil {
		t.Fatalf("invalid domain: %s", got)
	}
}
