package memorystorage

import (
	"context"
	"sync"

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
