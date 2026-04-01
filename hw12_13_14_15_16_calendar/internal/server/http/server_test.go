package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockApplication struct {
	mock.Mock
}

func (m *MockApplication) CreateEvent(ctx context.Context, event storage.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockApplication) UpdateEvent(ctx context.Context, event storage.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockApplication) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockApplication) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(storage.Event), args.Error(1)
}

func (m *MockApplication) ListEvents(ctx context.Context, userID string) ([]storage.Event, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockApplication) ListEventsForDay(ctx context.Context, date string) ([]storage.Event, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockApplication) ListEventsForWeek(ctx context.Context, startDate string) ([]storage.Event, error) {
	args := m.Called(ctx, startDate)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockApplication) ListEventsForMonth(ctx context.Context, startDate string) ([]storage.Event, error) {
	args := m.Called(ctx, startDate)
	return args.Get(0).([]storage.Event), args.Error(1)
}

type MockLogger struct{}

func (m *MockLogger) Info(msg string)  {}
func (m *MockLogger) Error(msg string) {}
func (m *MockLogger) Warn(msg string)  {}
func (m *MockLogger) Debug(msg string) {}

func TestHandleCreate(t *testing.T) {
	mockApp := new(MockApplication)
	logger := new(MockLogger)
	server := NewServer(logger, mockApp, "localhost", 8080)

	event := storage.Event{ID: "1", Title: "Test"}
	mockApp.On("CreateEvent", mock.Anything, event).Return(nil)

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/create_event", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	server.handleCreate(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	mockApp.AssertExpectations(t)
}

func TestHandleListForDay(t *testing.T) {
	mockApp := new(MockApplication)
	logger := new(MockLogger)
	server := NewServer(logger, mockApp, "localhost", 8080)

	events := []storage.Event{{ID: "1", Title: "Test"}}
	mockApp.On("ListEventsForDay", mock.Anything, "2024-04-01").Return(events, nil)

	req := httptest.NewRequest("GET", "/events_for_day?date=2024-04-01", nil)
	rr := httptest.NewRecorder()

	server.handleListForDay(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	var resp []storage.Event
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, events, resp)
	mockApp.AssertExpectations(t)
}

func TestHandleDelete(t *testing.T) {
	mockApp := new(MockApplication)
	logger := new(MockLogger)
	server := NewServer(logger, mockApp, "localhost", 8080)

	mockApp.On("DeleteEvent", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest("DELETE", "/delete_event?id=1", nil)
	rr := httptest.NewRecorder()

	server.handleDelete(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockApp.AssertExpectations(t)
}
