package internal

import (
	"os"
	"testing"
)

func TestValidateInput(t *testing.T) {
	entries := []struct {
		input    string
		expected bool
	}{
		{"192.180.33.25", true},
		{"example.com", true},
		{"http://192.180.33.25", true},
		{"https://example.org/path?query=123", true},
		{"test-domain.org", true},
		{"256.256.256.256", true},
		{"invalid@domain", false},
		{"not-a-domain", false},
	}

	for _, entry := range entries {
		t.Run(entry.input, func(t *testing.T) {
			result := validateInput(entry.input)
			if result != entry.expected {
				t.Errorf(
					"validateInput(%q) = %v; want %v",
					entry.input,
					result,
					entry.expected,
				)
			}
		})
	}
}

func TestReadAndProcessInputs(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantAddrs []string
	}{
		{
			name:      "single valid url",
			input:     "http://192.168.1.1\n",
			wantAddrs: []string{"http://192.168.1.1"},
		},
		{
			name:      "multiple urls with whitespace",
			input:     "  http://192.168.1.1   https://google.com  \n",
			wantAddrs: []string{"http://192.168.1.1", "https://google.com"},
		},
		{
			name:      "mixed valid and invalid",
			input:     "http://192.168.1.1 not_valid example.com\n",
			wantAddrs: []string{"http://192.168.1.1", "example.com"},
		},
		{
			name:      "empty input",
			input:     "",
			wantAddrs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r, w, _ := os.Pipe()
			defer r.Close()

			os.Stdin = r

			w.Write([]byte(tt.input))
			w.Close()

			got, _ := ReadStdin()
			addrs, _ := ProcessInputs(got)

			if len(addrs) != len(tt.wantAddrs) {
				t.Errorf("got %v, want %v", addrs, tt.wantAddrs)
				return
			}

			for i := range addrs {
				if addrs[i] != tt.wantAddrs[i] {
					t.Errorf("got %v, want %v", addrs, tt.wantAddrs)
					return
				}
			}
		})
	}
}
