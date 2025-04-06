package internal

import (
	"bytes"
	"fmt"
	"os"
)

// ReadStdin reads from the standard input.
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
