package internal

import (
	"errors"
	"testing"
	"time"

	"github.com/jesusprubio/up/pkg"
)

func TestReportToLine(t *testing.T) {
	r := &pkg.Report{
		ProtocolID: "test",
		RHost:      "test",
		Time:       1 * time.Second,
	}
	t.Run("return success line if no error happened", func(t *testing.T) {
		got := ReportToLine(r)
		want := "✔ test            1s             test           "
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
	t.Run("return error line if an error happened", func(t *testing.T) {
		r.Error = errors.New("probe error")
		got := ReportToLine(r)
		want := "✘ test            1s             probe error    "
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}
