// Package internalhttp предоставляет HTTP-сервер для приложения календаря.
// Включает middleware для логирования запросов и graceful shutdown.
package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	logger  Logger      // логгер для записи событий сервера
	app     Application // интерфейс к бизнес-логике приложения
	httpSrv *http.Server
}

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

func NewServer(logger Logger, app Application, host string, port int) *Server {
	mux := http.NewServeMux()

	s := &Server{
		logger: logger,
		app:    app,
	}

	mux.HandleFunc("/create_event", s.handleCreate)
	mux.HandleFunc("/update_event", s.handleUpdate)
	mux.HandleFunc("/delete_event", s.handleDelete)
	mux.HandleFunc("/events_for_day", s.handleListForDay)
	mux.HandleFunc("/events_for_week", s.handleListForWeek)
	mux.HandleFunc("/events_for_month", s.handleListForMonth)

	h := loggingMiddleware(logger)(mux)

	addr := fmt.Sprintf("%s:%d", host, port)
	s.httpSrv = &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	var e storage.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.app.CreateEvent(r.Context(), e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var e storage.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.app.UpdateEvent(r.Context(), e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if err := s.app.DeleteEvent(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleListForDay(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	events, err := s.app.ListEventsForDay(r.Context(), date)
	s.writeJSON(w, events, err)
}

func (s *Server) handleListForWeek(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	events, err := s.app.ListEventsForWeek(r.Context(), date)
	s.writeJSON(w, events, err)
}

func (s *Server) handleListForMonth(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	events, err := s.app.ListEventsForMonth(r.Context(), date)
	s.writeJSON(w, events, err)
}

func (s *Server) writeJSON(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) Start(ctx context.Context) error {
	// Горутина для graceful shutdown при отмене контекста
	go func() {
		<-ctx.Done()
		_ = s.Stop(context.Background())
	}()
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.httpSrv.Shutdown(ctxTimeout)
}
