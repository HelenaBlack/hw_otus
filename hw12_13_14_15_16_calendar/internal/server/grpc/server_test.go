package internalgrpc

import (
	"context"
	"net"
	"testing"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/pkg/calendar"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
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

func (m *MockLogger) Info(string)  {}
func (m *MockLogger) Error(string) {}
func (m *MockLogger) Warn(string)  {}
func (m *MockLogger) Debug(string) {}

func TestCreateEvent(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	mockApp := new(MockApplication)
	logger := new(MockLogger)
	s := NewServer(logger, mockApp)

	go func() {
		if err := s.srv.Serve(lis); err != nil {
			return
		}
	}()
	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := calendar.NewCalendarServiceClient(conn)

	event := &calendar.Event{Id: "1", Title: "Test"}
	mockApp.On("CreateEvent", mock.Anything, mock.Anything).Return(nil)

	resp, err := client.CreateEvent(ctx, &calendar.CreateEventRequest{Event: event})
	require.NoError(t, err)
	require.Equal(t, "1", resp.Id)
	mockApp.AssertExpectations(t)
}
