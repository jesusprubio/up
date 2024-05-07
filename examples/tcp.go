package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jesusprubio/up/pkg"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	reportCh := make(chan *pkg.Report)
	defer close(reportCh)
	probe := pkg.Probe{
		Protocols: []*pkg.Protocol{pkg.Protocols[1]},
		Count:     3,
		Timeout:   2 * time.Second,
		Logger:    logger,
		ReportCh:  reportCh,
	}
	go func() {
		for report := range probe.ReportCh {
			fmt.Println(report)
		}
	}()
	err := probe.Run(context.Background())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
