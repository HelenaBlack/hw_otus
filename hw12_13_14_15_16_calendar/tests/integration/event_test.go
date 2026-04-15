package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEventCRUD(t *testing.T) {
	userID := uuid.New().String()
	event := storage.Event{
		Title:       "Test Event",
		Description: "Integration Test",
		UserID:      userID,
		StartTime:   time.Now().Add(1 * time.Hour).Unix(),
		EndTime:     time.Now().Add(2 * time.Hour).Unix(),
	}

	// 1. Create Event
	t.Run("Create", func(t *testing.T) {
		body, _ := json.Marshal(event)

		ctx := context.Background() // или t.Context() для Go 1.24+
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, apiURL+"/create_event", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)

		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// 2. List Events (checking listing)
	// Note: Our API doesn't have a direct 'list all for user' in HTTP mux currently,
	// but it has 'events_for_day'. Let's use that.
	t.Run("ListForDay", func(t *testing.T) {
		today := time.Now().Format("2006-01-02")
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			fmt.Sprintf("%s/events_for_day?date=%s",
				apiURL,
				today,
			), nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var events []storage.Event
		err = json.NewDecoder(resp.Body).Decode(&events)
		require.NoError(t, err)

		found := false
		for _, e := range events {
			if e.Title == event.Title && e.UserID == event.UserID {
				found = true
				event.ID = e.ID // capture ID for next tests
				break
			}
		}
		require.True(t, found, "Event not found in listing")
	})

	// 3. Update Event
	t.Run("Update", func(t *testing.T) {
		event.Description = "Updated description"
		body, _ := json.Marshal(event)
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			apiURL+"/update_event",
			bytes.NewBuffer(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 4. Delete Event
	t.Run("Delete", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			fmt.Sprintf("%s/delete_event?id=%s",
				apiURL,
				event.ID,
			), nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestEventBusinessErrors(t *testing.T) {
	t.Run("InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			apiURL+"/create_event",
			bytes.NewBufferString("{invalid json}"),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingIDOnDelete", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			apiURL+"/delete_event",
			nil,
		)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		require.Contains(t, string(body), "missing id")
	})
}
