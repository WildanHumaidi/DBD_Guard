package cron

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dbd-guard/backend/internal/usecase"
	robfig "github.com/robfig/cron/v3"
)

type SLACron struct {
	cron          *robfig.Cron
	reportUsecase usecase.ReportUsecase
}

func NewSLACron(reportUsecase usecase.ReportUsecase) *SLACron {
	// Setup panic recovery and basic audit logging
	logger := robfig.PrintfLogger(log.New(os.Stdout, "[SLA Cron] ", log.LstdFlags))
	c := robfig.New(
		robfig.WithChain(
			robfig.Recover(logger),
		),
	)
	return &SLACron{
		cron:          c,
		reportUsecase: reportUsecase,
	}
}

func (c *SLACron) Start() {
	// Run every hour to check for SLA violations
	_, err := c.cron.AddFunc("@hourly", func() {
		startTime := time.Now()
		log.Println("[SLA Cron] [AUDIT] START: Running Auto-Release check for expired tasks...")
		ctx := context.Background()
		count, err := c.reportUsecase.AutoReleaseExpiredTasks(ctx)
		duration := time.Since(startTime)
		if err != nil {
			log.Printf("[SLA Cron] [AUDIT] FAIL: Error auto-releasing tasks (Duration: %v): %v\n", duration, err)
		} else {
			log.Printf("[SLA Cron] [AUDIT] SUCCESS: Auto-released %d tasks due to SLA violation (Duration: %v)\n", count, duration)
		}
	})

	if err != nil {
		log.Fatalf("[SLA Cron] Failed to add cron job: %v", err)
	}

	c.cron.Start()
	log.Println("[SLA Cron] Started successfully")
}

func (c *SLACron) Stop() {
	c.cron.Stop()
	log.Println("[SLA Cron] Stopped")
}
