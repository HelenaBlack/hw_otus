package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/configs"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.NewConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	l := logger.New(cfg.Logger.Level)
	l.Info("Starting calendar sender...")

	rmqClient, err := rabbitmq.NewClient(
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.Queue,
	)
	if err != nil {
		l.Error("failed to init rabbitmq: " + err.Error())
		return
	}
	defer rmqClient.Close()

	msgs, err := rmqClient.Consume()
	if err != nil {
		l.Error("failed to start consuming: " + err.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	l.Info("Sender is ready to process notifications")

	for i := 0; i < cfg.Sender.WorkerCount; i++ {
		go func(id int) {
			l.Info(fmt.Sprintf("worker %d started", id))
			for {
				select {
				case <-ctx.Done():
					l.Info(fmt.Sprintf("worker %d finished", id))
					return
				case d, ok := <-msgs:
					if !ok {
						return
					}
					var notif storage.Notification
					if err := json.Unmarshal(d.Body, &notif); err != nil {
						l.Error("failed to unmarshal notification: " + err.Error())
						_ = d.Nack(false, false)
						continue
					}

					l.Info(fmt.Sprintf(
						"worker %d got event: %s starting at %d for user: %s",
						id, notif.Title, notif.StartTime, notif.UserID,
					))

					// Report status
					status := storage.NotificationStatus{
						EventID: notif.EventID,
						SentAt:  time.Now().Unix(),
						Success: true,
					}
					if err := rmqClient.PublishToQueue(ctx, "notification_status", status); err != nil {
						l.Error("failed to publish status: " + err.Error())
					}

					_ = d.Ack(false)
				}
			}
		}(i)
	}

	<-ctx.Done()
	l.Info("Sender shutting down...")
}
