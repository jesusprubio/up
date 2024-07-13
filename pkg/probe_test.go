package pkg

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestProtocolValidate(t *testing.T) {
	proto := &Protocol{
		ID:    "test-proto",
		Probe: func(rhost string, timeout time.Duration) (string, error) { return "", nil },
		RHost: func() (string, error) { return "", nil },
	}
	t.Run("returns nil with valid setup", func(t *testing.T) {
		err := proto.validate()
		if err != nil {
			t.Fatalf("got %q, want nil", err)
		}
	})
	t.Run("returns an error if 'Probe' property is nil", func(t *testing.T) {
		p := &Protocol{ID: proto.ID, RHost: proto.RHost}
		err := p.validate()
		want := "required property: Probe"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
	t.Run("returns an error if 'RHost' property is nil", func(t *testing.T) {
		p := &Protocol{ID: proto.ID, Probe: proto.Probe}
		err := p.validate()
		want := "required property: RHost"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
}

func TestProbeValidate(t *testing.T) {
	protocols := []*Protocol{{
		ID: "test-proto",
		Probe: func(rhost string, timeout time.Duration) (string, error) {
			return "", nil
		},
		RHost: func() (string, error) {
			return "", nil
		},
	}}
	t.Run("returns nil with valid setup", func(t *testing.T) {
		reportCh := make(chan *Report)
		defer close(reportCh)
		p := Probe{
			Protocols: protocols,
			Timeout:   1 * time.Second,
			Logger:    slog.Default(),
			ReportCh:  reportCh,
		}
		err := p.validate()
		if err != nil {
			t.Fatalf("got %q, want nil", err)
		}
	})
	t.Run("returns an error if Protocols is nil", func(t *testing.T) {
		p := Probe{}
		err := p.validate()
		want := "required property: Protocols"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
	t.Run("returns an error if a protocol is invalid", func(t *testing.T) {
		p := Probe{Protocols: []*Protocol{{}}}
		err := p.validate()
		want := "invalid protocol: required property: Probe"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
	t.Run("returns an error if Timeout is zero", func(t *testing.T) {
		p := Probe{Protocols: protocols}
		err := p.validate()
		want := "required property: Timeout"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
	t.Run("returns an error if Logger is nil", func(t *testing.T) {
		p := Probe{
			Protocols: protocols,
			Timeout:   1 * time.Second,
		}
		err := p.validate()
		want := "required property: Logger"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
	t.Run("returns an error if ReportCh is nil", func(t *testing.T) {
		p := Probe{
			Protocols: protocols,
			Timeout:   1 * time.Second,
			Logger:    slog.Default(),
		}
		err := p.validate()
		want := "required property: ReportCh"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
}

func TestProbeRun(t *testing.T) {
	t.Run("returns an error if the setup is invalid", func(t *testing.T) {
		p := Probe{}
		err := p.Run(context.Background())
		want := "invalid setup: required property: Protocols"
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err, want)
		}
	})
	hostPort := "192.168.1.1:22"
	localHostPort := "127.0.0.1:3355"
	proto := &Protocol{
		ID: "test-proto",
		Probe: func(rhost string, timeout time.Duration) (string, error) {
			return localHostPort, nil
		},
		RHost: func() (string, error) { return hostPort, nil },
	}
	t.Run("returns nil if 'Count' property is defined", func(t *testing.T) {
		reportCh := make(chan *Report)
		defer close(reportCh)
		p := Probe{
			Protocols: []*Protocol{proto},
			Count:     2,
			Timeout:   1 * time.Second,
			Logger:    slog.Default(),
			ReportCh:  reportCh,
		}
		go func(t *testing.T) {
			for report := range p.ReportCh {
				if report.ProtocolID != proto.ID {
					t.Errorf("got %q, want %q", report.ProtocolID, proto.ID)
				}
				if report.RHost != hostPort {
					t.Errorf("got %q, want %q", report.RHost, hostPort)
				}
				if report.Time == 0 {
					t.Errorf("got %q, want > 0", report.Time)
				}
				if report.Extra != localHostPort {
					t.Errorf("got %q, want %q", report.Extra, localHostPort)
				}
				if report.Error != nil {
					t.Errorf("got %q, want nil", report.Error)
				}
			}
		}(t)
		err := p.Run(context.Background())
		if err != nil {
			t.Fatalf("got %q, want nil", err)
		}
	})
	t.Run("returns nil if context is canceled", func(t *testing.T) {
		reportCh := make(chan *Report)
		defer close(reportCh)
		p := Probe{
			Protocols: []*Protocol{proto},
			Timeout:   1 * time.Second,
			Logger:    slog.Default(),
			ReportCh:  reportCh,
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := p.Run(ctx)
		if err != nil {
			t.Fatalf("got %q, want nil", err)
		}
	})
}
