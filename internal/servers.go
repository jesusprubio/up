package internal

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"net/url"
)

const tmplRandom = "creating random %s: %w"

// CaptivePortals are URLs that well-known companies use inspect the network
// connections of their users.
var CaptivePortals []*url.URL = []*url.URL{
	// Google Chrome.
	{
		Scheme: "http",
		Host:   "clients3.google.com:80",
		Path:   "/generate_204",
	},
	// Mozilla Firefox.
	{
		Scheme: "http",
		Host:   "detectportal.firefox.com:80",
		Path:   "/success.txt",
	},
	// Apple.
	{
		Scheme: "http",
		Host:   "www.apple.com:80",
		Path:   "/library/test/success.html",
	},
	// Microsoft.
	{
		Scheme: "http",
		Host:   "www.msftconnecttest.com:80",
		Path:   "/redirect",
	},
	// Android.
	{
		Scheme: "http",
		Host:   "connectivitycheck.android.com:80",
		Path:   "/generate_204",
	},
	// Ubuntu.
	{
		Scheme: "http",
		Host:   "connectivity-check.ubuntu.com:80",
	},
	// Debian.
	{
		Scheme: "http",
		Host:   "network-test.debian.org:80",
	},
}

// Resolvers is a list of public DNS server IP addresses.
var Resolvers = []*net.IP{
	// Cloudflare
	{1, 1, 1, 1},
	{1, 0, 0, 1},
	// Google
	{8, 8, 8, 8},
	{8, 8, 4, 4},
	// OpenDNS
	{208, 67, 222, 222},
	{208, 67, 222, 220},
	// Control D
	{76, 76, 2, 0},
	{76, 76, 10, 0},
	// AdGuard
	{94, 140, 14, 14},
	{94, 140, 15, 15},
	// CleanBrowsing
	{185, 228, 168, 9},
	{185, 228, 169, 9},
	// Verisign
	{64, 6, 64, 6},
	{64, 6, 65, 6},
	// Quad9
	{9, 9, 9, 9},
	{149, 112, 112, 112},
	// Neustar
	{156, 154, 70, 1},
	{156, 154, 71, 1},
	// Yandex
	{77, 88, 8, 8},
	{77, 88, 8, 1},
	// SafeDNS
	{195, 46, 39, 39},
	{195, 46, 39, 40},
	// Norton ConnectSafe
	{199, 85, 126, 10},
	{199, 85, 127, 10},
}

// RandomCaptivePortal returns a randomly selected captive portal URL from the
// list of well-known companies.
//
// Returns an error if the random number generator fails.
func RandomCaptivePortal() (string, error) {
	portal, err := randomItem(CaptivePortals)
	if err != nil {
		return "", err
	}
	return portal.String(), nil
}

// Returns a random item from the given slice.
//
// Returns an error if the random number generator fails.
func randomItem[T any](items []T) (T, error) {
	count := big.NewInt(int64(len(items)))
	index, err := rand.Int(rand.Reader, count)
	if err != nil {
		var zero T
		return zero, fmt.Errorf(tmplRandom, "item", err)
	}
	return items[index.Int64()], nil
}

// RandomDNSServer returns a randomly selected public DNS server address.
//
// Returns an error if the random number generator fails.
func RandomDNSServer() (string, error) {
	resolver, err := randomItem(Resolvers)
	if err != nil {
		return "", err
	}
	return resolver.String(), nil
}

// RandomTCPServer returns a randomly selected public DNS server address.
//
// Returns an error if the random number generator fails.
func RandomTCPServer() (string, error) {
	serverAddr, err := RandomDNSServer()
	if err != nil {
		return "", fmt.Errorf(tmplRandom, "TCP server", err)
	}
	return net.JoinHostPort(serverAddr, "53"), nil
}

// RandomDomain returns a domain selected randomly from the captive portals.
//
// Returns an error if the random number generator fails.
func RandomDomain() (string, error) {
	portalURL, err := RandomCaptivePortal()
	if err != nil {
		return "", fmt.Errorf(tmplRandom, "domain", err)
	}
	u, err := url.Parse(portalURL)
	if err != nil {
		return "", fmt.Errorf("parsing URL %s: %w", portalURL, err)
	}
	return u.Hostname(), nil
}
