package cron

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dbd-guard/backend/internal/usecase"
	robfig "github.com/robfig/cron/v3"
)

type FeedbackCron struct {
	cron          *robfig.Cron
	reportUsecase usecase.ReportUsecase
}

func NewFeedbackCron(reportUsecase usecase.ReportUsecase) *FeedbackCron {
	// Setup panic recovery and basic audit logging
	logger := robfig.PrintfLogger(log.New(os.Stdout, "[Feedback Cron] ", log.LstdFlags))
	c := robfig.New(
		robfig.WithChain(
			robfig.Recover(logger),
		),
	)
	return &FeedbackCron{
		cron:          c,
		reportUsecase: reportUsecase,
	}
}

func (c *FeedbackCron) Start() {
	// Run every day at 10:00 AM to check for 14-days feedback reminder
	_, err := c.cron.AddFunc("0 10 * * *", func() {
		startTime := time.Now()
		log.Println("[Feedback Cron] [AUDIT] START: Running 14-days evaluation feedback check...")
		ctx := context.Background()
		count, err := c.reportUsecase.SendFeedbackReminders(ctx)
		duration := time.Since(startTime)
		if err != nil {
			log.Printf("[Feedback Cron] [AUDIT] FAIL: Error sending feedback reminders (Duration: %v): %v\n", duration, err)
		} else {
			log.Printf("[Feedback Cron] [AUDIT] SUCCESS: Sent %d feedback reminders (Duration: %v)\n", count, duration)
		}
	})

	if err != nil {
		log.Fatalf("[Feedback Cron] Failed to add cron job: %v", err)
	}

	c.cron.Start()
	log.Println("[Feedback Cron] Started successfully")
}

func (c *FeedbackCron) Stop() {
	c.cron.Stop()
	log.Println("[Feedback Cron] Stopped")
}
