//go:generate make -C ../../../ generate
package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/pkg/calendar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEvent(ctx context.Context, id string) (storage.Event, error)
	ListEvents(ctx context.Context, userID string) ([]storage.Event, error)
	ListEventsForDay(ctx context.Context, date string) ([]storage.Event, error)
	ListEventsForWeek(ctx context.Context, startDate string) ([]storage.Event, error)
	ListEventsForMonth(ctx context.Context, startDate string) ([]storage.Event, error)
}

type Server struct {
	calendar.UnimplementedCalendarServiceServer
	logger Logger
	app    Application
	srv    *grpc.Server
}

func NewServer(logger Logger, app Application) *Server {
	s := &Server{
		logger: logger,
		app:    app,
	}

	// Logging interceptor
	loggingInterceptor := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		msg := fmt.Sprintf("method=%s duration=%s", info.FullMethod, duration)
		if err != nil {
			logger.Error(fmt.Sprintf("%s error=%v", msg, err))
		} else {
			logger.Info(msg)
		}

		return resp, err
	}

	s.srv = grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	calendar.RegisterCalendarServiceServer(s.srv, s)
	reflection.Register(s.srv)

	return s
}

func (s *Server) Start(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("grpc server is listening on %s", addr))
	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}

func (s *Server) CreateEvent(
	ctx context.Context,
	req *calendar.CreateEventRequest,
) (*calendar.CreateEventResponse, error) {
	e := storage.Event{
		ID:           req.Event.Id,
		Title:        req.Event.Title,
		Description:  req.Event.Description,
		UserID:       req.Event.UserId,
		StartTime:    req.Event.StartTime,
		EndTime:      req.Event.EndTime,
		NotifyBefore: &req.Event.NotifyBefore,
	}
	err := s.app.CreateEvent(ctx, e)
	if err != nil {
		return &calendar.CreateEventResponse{Error: err.Error()}, err
	}
	return &calendar.CreateEventResponse{Id: e.ID}, nil
}

func (s *Server) UpdateEvent(
	ctx context.Context,
	req *calendar.UpdateEventRequest,
) (*calendar.UpdateEventResponse, error) {
	e := storage.Event{
		ID:           req.Event.Id,
		Title:        req.Event.Title,
		Description:  req.Event.Description,
		UserID:       req.Event.UserId,
		StartTime:    req.Event.StartTime,
		EndTime:      req.Event.EndTime,
		NotifyBefore: &req.Event.NotifyBefore,
	}
	err := s.app.UpdateEvent(ctx, e)
	if err != nil {
		return &calendar.UpdateEventResponse{Error: err.Error()}, err
	}
	return &calendar.UpdateEventResponse{}, nil
}

func (s *Server) DeleteEvent(
	ctx context.Context,
	req *calendar.DeleteEventRequest,
) (*calendar.DeleteEventResponse, error) {
	err := s.app.DeleteEvent(ctx, req.Id)
	if err != nil {
		return &calendar.DeleteEventResponse{Error: err.Error()}, err
	}
	return &calendar.DeleteEventResponse{}, nil
}

func (s *Server) ListEventsForDay(
	ctx context.Context,
	req *calendar.ListEventsForDayRequest,
) (*calendar.ListEventsResponse, error) {
	events, err := s.app.ListEventsForDay(ctx, req.Date)
	return s.toListResponse(events, err)
}

func (s *Server) ListEventsForWeek(
	ctx context.Context,
	req *calendar.ListEventsForWeekRequest,
) (*calendar.ListEventsResponse, error) {
	events, err := s.app.ListEventsForWeek(ctx, req.StartDate)
	return s.toListResponse(events, err)
}

func (s *Server) ListEventsForMonth(
	ctx context.Context,
	req *calendar.ListEventsForMonthRequest,
) (*calendar.ListEventsResponse, error) {
	events, err := s.app.ListEventsForMonth(ctx, req.StartDate)
	return s.toListResponse(events, err)
}

func (s *Server) toListResponse(events []storage.Event, err error) (*calendar.ListEventsResponse, error) {
	if err != nil {
		return &calendar.ListEventsResponse{Error: err.Error()}, err
	}
	protoEvents := make([]*calendar.Event, 0, len(events))
	for _, e := range events {
		var notifyBefore int64
		if e.NotifyBefore != nil {
			notifyBefore = *e.NotifyBefore
		}
		protoEvents = append(protoEvents, &calendar.Event{
			Id:           e.ID,
			Title:        e.Title,
			Description:  e.Description,
			UserId:       e.UserID,
			StartTime:    e.StartTime,
			EndTime:      e.EndTime,
			NotifyBefore: notifyBefore,
		})
	}
	return &calendar.ListEventsResponse{Events: protoEvents}, nil
}
