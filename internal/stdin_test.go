package internal

import (
	"os"
	"testing"
)

// Tests [validateInput] func
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

// Tests  [ReadStdin] and [ProcessInputs] func
func TestReadAndProcessInputs(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantInputs []string
		wantErr    bool
	}{
		{
			name:       "single valid url",
			input:      "http://192.168.1.1\n",
			wantInputs: []string{"http://192.168.1.1"},
			wantErr:    false,
		},
		{
			name:       "multiple urls with whitespace",
			input:      "  http://192.168.1.1   https://google.com  \n",
			wantInputs: []string{"http://192.168.1.1", "https://google.com"},
			wantErr:    false,
		},
		{
			name:       "mixed valid and invalid",
			input:      "http://192.168.1.1 not_valid example.com\n",
			wantInputs: []string{"http://192.168.1.1"},
			wantErr:    true,
		},
		{
			name:       "empty input",
			input:      "",
			wantInputs: []string{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			defer r.Close()
			os.Stdin = r
			w.Write([]byte(tt.input))
			w.Close()
			got, readErr := ReadStdin()
			if readErr != nil {
				t.Errorf("ReadStdin() error = %v", readErr)
				return
			}
			inputs, processErr := ProcessInputs(got)
			if (processErr != nil) != tt.wantErr {
				t.Errorf(
					"ProcessInputs() error = %v, wantErr %v",
					processErr,
					tt.wantErr,
				)
				return
			}
			if len(inputs) != len(tt.wantInputs) {
				t.Errorf("got %v, want %v", inputs, tt.wantInputs)
				return
			}
			for i := range inputs {
				if inputs[i] != tt.wantInputs[i] {
					t.Errorf("got %v, want %v", inputs, tt.wantInputs)
					return
				}
			}
		})
	}
}
