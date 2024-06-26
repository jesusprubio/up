package internal

import (
	"errors"
	"testing"
	"time"

	"github.com/jesusprubio/up/pkg"
)

func TestProtocolByID(t *testing.T) {
	t.Run("returns the protocol if it exists", func(t *testing.T) {
		got := ProtocolByID("http")
		if got == nil {
			t.Fatal("got nil, want a protocol")
		}
	})
	t.Run("returns nil if the protocol doesn't exist", func(t *testing.T) {
		got := ProtocolByID("unknown")
		if got != nil {
			t.Fatalf("got %q, want nil", got)
		}
	})
}

func TestReportToLine(t *testing.T) {
	r := &pkg.Report{
		ProtocolID: "test",
		Time:       1 * time.Second,
		RHost:      "test",
		Extra:      "test",
	}
	t.Run("return success line if no error happened", func(t *testing.T) {
		got := ReportToLine(r)
		want := "✔ test            1s             test (test)"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	t.Run("return error line if an error happened", func(t *testing.T) {
		r.Error = errors.New("probe error")
		got := ReportToLine(r)
		want := "✘ test            1s             test (probe error)"
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}
