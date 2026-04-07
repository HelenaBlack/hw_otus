package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/app"
	config "github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/configs"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage/sql"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.NewConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	l := logger.New(cfg.Logger.Level)
	l.Info("Starting calendar scheduler...")

	var s app.Storage
	if cfg.Storage.Type == "sql" {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)
		s, err = sqlstorage.New(dsn)
	} else {
		s = memorystorage.New()
	}
	if err != nil {
		l.Error("failed to init storage: " + err.Error())
		os.Exit(1)
	}

	rmqClient, err := rabbitmq.NewClient(
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.Queue,
	)
	if err != nil {
		l.Error("failed to init rabbitmq: " + err.Error())
		os.Exit(1)
	}
	defer rmqClient.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	ticker := time.NewTicker(time.Duration(cfg.Scheduler.ScanInterval) * time.Second)
	defer ticker.Stop()

	// Cleanup ticker - once per day
	cleanupTicker := time.NewTicker(24 * time.Hour)
	defer cleanupTicker.Stop()

	l.Info("Scheduler is running")

	for {
		select {
		case <-ctx.Done():
			l.Info("Scheduler shutting down...")
			return
		case <-ticker.C:
			now := time.Now().Unix()
			l.Debug("Scanning for events to notify...")
			events, err := s.GetEventsForNotification(ctx, now)
			if err != nil {
				l.Error("failed to get events for notification: " + err.Error())
				continue
			}

			for _, e := range events {
				notif := storage.Notification{
					EventID:   e.ID,
					Title:     e.Title,
					StartTime: e.StartTime,
					UserID:    e.UserID,
				}
				if err := rmqClient.Publish(ctx, notif); err != nil {
					l.Error("failed to publish notification: " + err.Error())
					continue
				}
				if err := s.UpdateNotificationSent(ctx, e.ID); err != nil {
					l.Error("failed to update notification sent flag: " + err.Error())
				}
				l.Info("Notification sent for event: " + e.ID)
			}
		case <-cleanupTicker.C:
			olderThan := time.Now().AddDate(-1, 0, 0).Unix()
			l.Info("Cleaning up old events...")
			if err := s.DeleteOldEvents(ctx, olderThan); err != nil {
				l.Error("failed to delete old events: " + err.Error())
			}
		}
	}
}
