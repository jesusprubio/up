package internal

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Summary aggregates the reports of a run, per protocol.
//
// It is not safe for concurrent use.
type Summary struct {
	stats []*ProtocolStats
}

// ProtocolStats are the aggregated results of a protocol.
type ProtocolStats struct {
	// Protocol identifier. Example: 'http'.
	ProtocolID string
	// Number of successful requests.
	OK int
	// Number of failed requests.
	Failed int
	// Sum of the response times of the successful requests.
	totalTime time.Duration
	// Last error received, if any.
	LastError string
}

// Average returns the mean response time of the successful requests.
func (s *ProtocolStats) Average() time.Duration {
	if s.OK == 0 {
		return 0
	}
	return s.totalTime / time.Duration(s.OK)
}

// Up reports if at least one request succeeded.
func (s *ProtocolStats) Up() bool {
	return s.OK > 0
}

// Add includes a report in the summary.
func (s *Summary) Add(r *Report) {
	var st *ProtocolStats
	for _, item := range s.stats {
		if item.ProtocolID == r.ProtocolID {
			st = item
			break
		}
	}
	if st == nil {
		st = &ProtocolStats{ProtocolID: r.ProtocolID}
		s.stats = append(s.stats, st)
	}
	if r.Error != "" {
		st.Failed++
		st.LastError = r.Error
		return
	}
	st.OK++
	st.totalTime += r.Time
}

// Stats returns the aggregated results, per protocol, in order of arrival.
func (s *Summary) Stats() []*ProtocolStats {
	return s.stats
}

// Up reports if at least one request succeeded, for any protocol.
func (s *Summary) Up() bool {
	for _, st := range s.stats {
		if st.Up() {
			return true
		}
	}
	return false
}

// String returns the summary ready to be printed.
func (s *Summary) String(format Format) (string, error) {
	switch format {
	case HumanFormat:
		return s.stringHuman(), nil
	case JSONFormat:
		line, err := s.stringJSON()
		if err != nil {
			return "", fmt.Errorf("error generating JSON summary: %w", err)
		}
		return line, nil
	case GrepFormat:
		return s.stringGrep(), nil
	default:
		return "", fmt.Errorf("unsupported format: %v", format)
	}
}

type protocolStatsJSON struct {
	Protocol  string  `json:"protocol"`
	Up        bool    `json:"up"`
	OK        int     `json:"ok"`
	Failed    int     `json:"failed"`
	AvgMs     float64 `json:"avg_ms,omitempty"`
	LastError string  `json:"last_error,omitempty"`
}

type summaryJSON struct {
	Up        bool                `json:"up"`
	Protocols []protocolStatsJSON `json:"protocols"`
}

// Returns the summary in JSON format.
// Example:
// '{"up":true,"protocols":[{"protocol":"http","up":true,"ok":2,"failed":0,"avg_ms":35.2}]}'
func (s *Summary) stringJSON() (string, error) {
	out := summaryJSON{Up: s.Up(), Protocols: []protocolStatsJSON{}}
	for _, st := range s.stats {
		out.Protocols = append(out.Protocols, protocolStatsJSON{
			Protocol:  st.ProtocolID,
			Up:        st.Up(),
			OK:        st.OK,
			Failed:    st.Failed,
			AvgMs:     milliseconds(st.Average()),
			LastError: st.LastError,
		})
	}
	line, err := json.Marshal(out)
	if err != nil {
		return "", fmt.Errorf("marshaling summary: %w", err)
	}
	return string(line), nil
}

// Returns the summary in human readable format.
// Example:
// ✔ http            2/2 ok         avg 35.2ms
// ✘ dns             0/2 ok         (lookup example.com: no such host)
func (s *Summary) stringHuman() string {
	lines := []string{}
	for _, st := range s.stats {
		ratio := fmt.Sprintf("%d/%d ok", st.OK, st.OK+st.Failed)
		line := fmt.Sprintf("%-15s %-14s", bold(st.ProtocolID), ratio)
		if st.Up() {
			avg := fmt.Sprintf("avg %s", st.Average().Round(time.Microsecond))
			lines = append(
				lines,
				fmt.Sprintf("%s %s %s", green("✔"), line, avg),
			)
			continue
		}
		suffix := faint(fmt.Sprintf("(%s)", st.LastError))
		lines = append(lines, fmt.Sprintf("%s %s %s", red("✘"), line, suffix))
	}
	return strings.Join(lines, "\n")
}

// Returns the summary in a grepable format, one line per protocol.
//
// Example: 'http	up	2	0	35.2ms	'
func (s *Summary) stringGrep() string {
	lines := []string{}
	for _, st := range s.stats {
		status := "down"
		if st.Up() {
			status = "up"
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%d\t%d\t%s\t%s",
			st.ProtocolID, status, st.OK, st.Failed,
			st.Average().Round(time.Microsecond), st.LastError,
		))
	}
	return strings.Join(lines, "\n")
}
