package internal

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Example_validateAddress() {
	entries := []string{
		"192.180.33.25",
		"example.com",
		"http://192.180.33.25",
		"https://example.org/path?query=123",
		"test-domain.org",
		"256.256.256.256",
		"invalid@domain",
		"not-a-domain",
	}

	for _, addr := range entries {
		fmt.Println(addr+":", validateAddress(addr))
	}
	// Output:
	// 192.180.33.25: true
	// example.com: true
	// http://192.180.33.25: true
	// https://example.org/path?query=123: true
	// test-domain.org: true
	// 256.256.256.256: true
	// invalid@domain: false
	// not-a-domain: false
}

func TestReadAndProcessAddrs(t *testing.T) {
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
			os.Stdin = r
			defer func() {
				r.Close()
			}()

			w.Write([]byte(tt.input))
			w.Close()

			got, _ := ReadStdin()
			addrs := ProcessAddrs(got)

			if !reflect.DeepEqual(addrs, tt.wantAddrs) {
				t.Errorf("got %v, want %v", addrs, tt.wantAddrs)
			}
		})
	}
}
