package internal

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"
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
				if report.Error != nil {
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

// Test [Probe.Do] in parallel  and serial  mode
func TestProbeDoParallel(t *testing.T) {
	protocols := []Protocol{&testProtocol{}}
	testCases := []struct {
		name     string
		parallel bool
		inputLen int
		count    uint
	}{
		{"Serial with no inputs", false, 0, 2},
		{"Serial with inputs", false, 2, 2},
		{"Parallel with no inputs", true, 0, 2},
		{"Parallel with inputs", true, 2, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()

			reportCh := make(chan *Report, 100)

			var wg sync.WaitGroup
			wg.Add(1)
			inputs := make([]string, tc.inputLen)
			for i := range inputs {
				inputs[i] = "test-input-" + string('a'+rune(i))
			}

			reports := make([]*Report, 0, 100)

			var reportsMutex sync.Mutex
			go func() {
				defer wg.Done()
				for report := range reportCh {
					reportsMutex.Lock()
					reports = append(reports, report)
					reportsMutex.Unlock()
				}
			}()

			p := Probe{
				Protocols: protocols,
				Count:     tc.count,
				Parallel:  tc.parallel,
				Logger:    slog.Default(),
				ReportCh:  reportCh,
				Input:     inputs,
			}

			err := p.Do(ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			close(reportCh)

			wg.Wait()

			expectedReportCount := int(
				tc.count,
			) * len(
				protocols,
			) * (tc.inputLen)
			if expectedReportCount == 0 {
				expectedReportCount = int(tc.count) * len(protocols)
			}

			if len(reports) != expectedReportCount {
				t.Errorf(
					"incorrect report count. got %d, want %d",
					len(reports),
					expectedReportCount,
				)
			}

			for _, report := range reports {
				if report.ProtocolID != protocols[0].String() {
					t.Errorf(
						"incorrect protocol ID. got %q, want %q",
						report.ProtocolID,
						protocols[0].String(),
					)
				}
				if report.RHost != testHostPort {
					t.Errorf(
						"incorrect remote host. got %q, want %q",
						report.RHost,
						testHostPort,
					)
				}
				if report.Error != nil {
					t.Errorf("unexpected error: %v", report.Error)
				}
			}
		})
	}

	t.Run("returns nil if context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		reportCh := make(chan *Report, 10)
		defer close(reportCh)

		p := Probe{
			Protocols: protocols,
			Logger:    slog.Default(),
			ReportCh:  reportCh,
		}

		err := p.Do(ctx)
		if err != nil {
			t.Fatalf("got %q, want nil", err)
		}
	})
}
