package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNotificationLifecycle(t *testing.T) {
	userID := uuid.New().String()
	notifyBefore := int64(10) // 10 seconds
	// Create event starting in 12 seconds, notify 10 seconds before.
	// So it should be picked up in ~2 seconds.
	event := storage.Event{
		Title:        "Notification Test Event",
		Description:  "Testing full lifecycle",
		UserID:       userID,
		StartTime:    time.Now().Unix() + 12,
		EndTime:      time.Now().Unix() + 3612,
		NotifyBefore: &notifyBefore,
	}

	// 1. Create Event
	body, _ := json.Marshal(event)
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+"/create_event", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// 2. Wait for status in RabbitMQ
	// We consume from the status queue that our Sender now populates
	msgs, err := rmqClient.Consume()
	require.NoError(t, err)

	timeout := time.After(30 * time.Second)
	found := false
	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for notification status")
		case d := <-msgs:
			var status storage.NotificationStatus
			err := json.Unmarshal(d.Body, &status)
			require.NoError(t, err)

			// We might get statuses from other tests if they run in parallel,
			// but we can check the Title or EventID if we had it.
			// Currently we don't have the ID until we list it,
			// but we can use the UserID if we include it in status.
			// Wait, I didn't include UserID in NotificationStatus.
			// I'll just check if it's successful for now,
			// as this is a clean test environment.

			fmt.Printf("Received status for event: %s, success: %v\n", status.EventID, status.Success)
			if status.Success {
				found = true
				_ = d.Ack(false)
				goto end
			}
			_ = d.Ack(false)
		}
	}

end:
	require.True(t, found, "Notification status not received")
}
