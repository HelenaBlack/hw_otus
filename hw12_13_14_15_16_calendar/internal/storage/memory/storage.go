package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex             // мьютекс для синхронизации доступа к данным
	events map[string]storage.Event // карта событий, ключ - ID события
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range s.events {
		if e.UserID == event.UserID && e.StartTime == event.StartTime {
			return app.ErrDateBusy
		}
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	if ctx.Err() != nil {
		return ctx.Err() // Возвращаем ошибку, если контекст уже отменён
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return context.Canceled
	}

	for id, e := range s.events {
		if id != event.ID && e.UserID == event.UserID && e.StartTime == event.StartTime {
			return app.ErrDateBusy
		}
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return context.Canceled // или custom not found error
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	if ctx.Err() != nil {
		return storage.Event{}, ctx.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[id]
	if !ok {
		return storage.Event{}, context.Canceled // или custom not found error
	}
	return event, nil
}

func (s *Storage) ListEvents(ctx context.Context, userID string) ([]storage.Event, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []storage.Event
	for _, e := range s.events {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (s *Storage) ListEventsForDay(_ context.Context, date string) ([]storage.Event, error) {
	// startDate в формате YYYY-MM-DD
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []storage.Event
	for _, e := range s.events {
		eventDate := time.Unix(e.StartTime, 0).Format("2006-01-02")
		if eventDate == date {
			result = append(result, e)
		}
	}
	return result, nil
}

func (s *Storage) ListEventsForWeek(_ context.Context, startDate string) ([]storage.Event, error) {
	// startDate в формате YYYY-MM-DD
	t, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end := t.AddDate(0, 0, 7)

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []storage.Event
	for _, e := range s.events {
		if e.StartTime >= t.Unix() && e.StartTime < end.Unix() {
			result = append(result, e)
		}
	}
	return result, nil
}

func (s *Storage) ListEventsForMonth(_ context.Context, startDate string) ([]storage.Event, error) {
	// startDate в формате YYYY-MM-DD
	t, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end := t.AddDate(0, 1, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []storage.Event
	for _, e := range s.events {
		if e.StartTime >= t.Unix() && e.StartTime < end.Unix() {
			result = append(result, e)
		}
	}
	return result, nil
}
