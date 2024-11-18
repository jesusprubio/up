package internal

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func ReadStdin() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if (info.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("reading from stdin: %w", err)
	}

	return buf.String(), nil
}

func validateAddress(addr string) bool {

	// URL validation
	if _, err := url.ParseRequestURI(addr); err == nil {
		return true
	}

	ipPattern := `^([0-9]{1,3}\.){3}[0-9]{1,3}$`

	// Domain name pattern
	domainPattern := `^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$`

	ipRegex := regexp.MustCompile(ipPattern)
	domainRegex := regexp.MustCompile(domainPattern)

	if ipRegex.MatchString(addr) || domainRegex.MatchString(addr) {
		return true
	}
	return false
}

func ProcessAddrs(s string) []string {
	addrArray := []string{}

	s = strings.TrimSpace(s)
	for _, addr := range strings.Fields(s) {
		addr = strings.TrimSpace(addr)
		if validateAddress(addr) {
			addrArray = append(addrArray, addr)
		} else {
			fmt.Println("Wrong Address format", addr)
		}
	}

	return addrArray
}
