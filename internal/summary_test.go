package internal

import (
	"testing"
	"time"

	"github.com/fatih/color"
)

func newTestSummary() *Summary {
	var s Summary
	s.Add(&Report{ProtocolID: "http", Time: 30 * time.Millisecond})
	s.Add(&Report{ProtocolID: "dns", Error: "error-0"})
	s.Add(&Report{ProtocolID: "http", Time: 40 * time.Millisecond})
	s.Add(&Report{ProtocolID: "http", Error: "error-1"})
	return &s
}

func TestSummaryAdd(t *testing.T) {
	stats := newTestSummary().Stats()
	if len(stats) != 2 {
		t.Fatalf("got %d protocols, want 2", len(stats))
	}
	http, dns := stats[0], stats[1]
	if http.ProtocolID != "http" || dns.ProtocolID != "dns" {
		t.Fatalf("got %q and %q, want arrival order", http.ProtocolID,
			dns.ProtocolID)
	}
	if http.OK != 2 || http.Failed != 1 {
		t.Fatalf("got %d ok and %d failed, want 2 and 1", http.OK, http.Failed)
	}
	if http.Average() != 35*time.Millisecond {
		t.Fatalf("got %s, want 35ms", http.Average())
	}
	if http.LastError != "error-1" {
		t.Fatalf("got %q, want %q", http.LastError, "error-1")
	}
	if dns.Up() || dns.Average() != 0 {
		t.Fatal("got dns up, want down with no average")
	}
}

func TestSummaryUp(t *testing.T) {
	t.Run("returns true if any protocol is up", func(t *testing.T) {
		if !newTestSummary().Up() {
			t.Fatal("got false, want true")
		}
	})
	t.Run("returns false if all the requests failed", func(t *testing.T) {
		var s Summary
		s.Add(&Report{ProtocolID: "dns", Error: "error-0"})
		if s.Up() {
			t.Fatal("got true, want false")
		}
	})
	t.Run("returns false without reports", func(t *testing.T) {
		var s Summary
		if s.Up() {
			t.Fatal("got true, want false")
		}
	})
}

func TestSummaryString(t *testing.T) {
	color.NoColor = true
	s := newTestSummary()
	tests := []struct {
		name   string
		format Format
		want   string
	}{
		{
			"returns a summary using JSON format",
			JSONFormat,
			`{"up":true,"protocols":[` +
				`{"protocol":"http","up":true,"ok":2,"failed":1,` +
				`"avg_ms":35,"last_error":"error-1"},` +
				`{"protocol":"dns","up":false,"ok":0,"failed":1,` +
				`"last_error":"error-0"}]}`,
		},
		{
			"returns a summary using human format",
			HumanFormat,
			"✔ http            2/3 ok         avg 35ms\n" +
				"✘ dns             0/1 ok         (error-0)",
		},
		{
			"returns a summary using a grepable format",
			GrepFormat,
			"http\tup\t2\t1\t35ms\terror-1\ndns\tdown\t0\t1\t0s\terror-0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.String(tt.format)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
	t.Run("returns an empty protocol list without reports", func(t *testing.T) {
		var empty Summary
		got, err := empty.String(JSONFormat)
		if err != nil {
			t.Fatal(err)
		}
		want := `{"up":false,"protocols":[]}`
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}
