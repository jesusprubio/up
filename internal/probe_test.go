package internal

import (
	"context"
	"log/slog"
	"testing"
)

const testHostPort = "127.0.0.1:3355"
const testExtra = "test-extra"

type testProtocol struct{}

func (p *testProtocol) String() string { return "test-proto" }

func (p *testProtocol) Probe(target string) (string, string, error) {
	return testHostPort, testExtra, nil
}

func TestProbeValidate(t *testing.T) {
	protocols := []Protocol{&testProtocol{}}
	t.Run("returns nil with valid setup", func(t *testing.T) {
		reportCh := make(chan *Report)
		defer close(reportCh)
		p := Probe{
			Protocols: protocols, Logger: slog.Default(), ReportCh: reportCh,
		}
		err := p.validate()
		if err != nil {
			t.Fatalf("got %q, want nil", err)
		}
	})
	t.Run("returns an error if 'Protocols' is nil", func(t *testing.T) {
		p := Probe{}
		err := p.validate()
		want := "required property: Protocols"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
	t.Run("returns an error if 'Logger' is nil", func(t *testing.T) {
		p := Probe{Protocols: protocols}
		err := p.validate()
		want := "required property: Logger"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
	t.Run("returns an error if 'ReportCh' is nil", func(t *testing.T) {
		p := Probe{Protocols: protocols, Logger: slog.Default()}
		err := p.validate()
		want := "required property: ReportCh"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
}

func TestProbeDo(t *testing.T) {
	t.Run("returns an error if the setup is invalid", func(t *testing.T) {
		p := Probe{}
		err := p.Do(context.Background())
		want := "invalid setup: required property: Protocols"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
	protocols := []Protocol{&testProtocol{}}
	t.Run("sends back the report in the channel", func(t *testing.T) {
		reportCh := make(chan *Report)
		defer close(reportCh)
		p := Probe{
			Protocols: protocols,
			Count:     2,
			Logger:    slog.Default(),
			ReportCh:  reportCh,
		}
		protoID := protocols[0].String()
		go func(t *testing.T) {
			for report := range p.ReportCh {
				if report.ProtocolID != protoID {
					t.Errorf("got %q, want %q", report.ProtocolID, protoID)
				}
				if report.RHost != testHostPort {
					t.Errorf("got %q, want %q", report.RHost, testHostPort)
				}
				if report.Time == 0 {
					t.Errorf("got %q, want > 0", report.Time)
				}
				if report.Error != "" {
					t.Errorf("got %q, want nil", report.Error)
				}
			}
		}(t)
		err := p.Do(context.Background())
		if err != nil {
			t.Fatalf("got %q, want nil", err)
		}
	})
	t.Run("returns nil if context is canceled", func(t *testing.T) {
		reportCh := make(chan *Report)
		defer close(reportCh)
		p := Probe{
			Protocols: protocols, Logger: slog.Default(), ReportCh: reportCh,
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := p.Do(ctx)
		if err != nil {
			t.Fatalf("got %q, want nil", err)
		}
	})
}
