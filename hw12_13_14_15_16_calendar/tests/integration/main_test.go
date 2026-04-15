package integration

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/rabbitmq"
)

var (
	apiURL    = "http://test-calendar:8080"
	rmqClient *rabbitmq.Client
)

func TestMain(m *testing.M) {
	// Wait for API to be ready
	if err := waitForAPI(); err != nil {
		log.Fatalf("API not ready: %v", err)
	}

	// Wait for RabbitMQ and init client for status checking
	var err error
	for i := 0; i < 10; i++ {
		rmqClient, err = rabbitmq.NewClient("guest", "guest", "test-rabbitmq", 5672, "notification_status")

		if err == nil {
			break
		}
		fmt.Printf("Waiting for RabbitMQ... %v\n", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("RabbitMQ not ready: %v", err)
	}

	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	defer rmqClient.Close()
	return m.Run()
}

func waitForAPI() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for API")
		default:
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL+"/events", nil)
			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
					return nil
				}
			}
			fmt.Println("Waiting for API...")
			time.Sleep(2 * time.Second)
		}
	}
}
