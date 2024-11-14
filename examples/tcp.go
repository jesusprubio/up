package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jesusprubio/up/internal"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	reportCh := make(chan *internal.Report)
	defer close(reportCh)
	probe := internal.Probe{
		Protocols: []internal.Protocol{&internal.TCP{Timeout: 2 * time.Second}},
		Count:     3,
		Logger:    logger,
		ReportCh:  reportCh,
	}
	go func() {
		for report := range probe.ReportCh {
			fmt.Println(report)
		}
	}()
	err := probe.Do(context.Background())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
