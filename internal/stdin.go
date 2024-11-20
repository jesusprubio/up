package internal

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	domainPattern = `^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$`
	ipPattern     = `^(\d{1,3}\.){3}\d{1,3}$`
)

func ReadStdin() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve stdin : %w", err)
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

var (
	ipRegex     = regexp.MustCompile(ipPattern)
	domainRegex = regexp.MustCompile(domainPattern)
)

// Ensures the inputs are correct
func validateInput(addr string) bool {
	_, err := url.ParseRequestURI(addr)
	if err == nil {
		return true
	}
	if ipRegex.MatchString(addr) || domainRegex.MatchString(addr) {
		return true
	}
	return false
}
func ProcessInputs(s string) ([]string, error) {
	var inputs []string
	for _, input := range strings.Fields(s) {
		if validateInput(input) {
			inputs = append(inputs, input)
		} else {
			return inputs, fmt.Errorf("invalid address format: %s", input)
		}
	}
	return inputs, nil
}
