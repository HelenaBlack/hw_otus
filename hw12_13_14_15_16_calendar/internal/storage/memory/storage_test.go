package memorystorage

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/HelenaBlack/hw_otus/hw12_13_14_15_calendar/internal/storage"
)

func TestStorageCRUD(t *testing.T) {
	s := New()
	ctx := context.Background()
	event := storage.Event{
		ID:        "1",
		Title:     "Test Event",
		UserID:    "user1",
		StartTime: 1000,
		EndTime:   2000,
	}

	if err := s.CreateEvent(ctx, event); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	got, err := s.GetEvent(ctx, event.ID)
	if err != nil || got.ID != event.ID {
		t.Fatalf("GetEvent failed: %v", err)
	}

	event.Title = "Updated"
	if err := s.UpdateEvent(ctx, event); err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}
	got, _ = s.GetEvent(ctx, event.ID)
	if got.Title != "Updated" {
		t.Fatalf("UpdateEvent did not update title")
	}

	list, err := s.ListEvents(ctx, event.UserID)
	if err != nil || len(list) != 1 {
		t.Fatalf("ListEvents failed: %v", err)
	}

	if err := s.DeleteEvent(ctx, event.ID); err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}
	_, err = s.GetEvent(ctx, event.ID)
	if err == nil {
		t.Fatalf("GetEvent should fail after delete")
	}
}

func TestStorageErrDateBusy(t *testing.T) {
	s := New()
	ctx := context.Background()
	e1 := storage.Event{ID: "1", UserID: "u", StartTime: 1000}
	e2 := storage.Event{ID: "2", UserID: "u", StartTime: 1000}

	if err := s.CreateEvent(ctx, e1); err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	err := s.CreateEvent(ctx, e2)
	if !errors.Is(err, app.ErrDateBusy) {
		t.Fatalf("expected ErrDateBusy, got %v", err)
	}
}

func TestStorageConcurrency(t *testing.T) {
	s := New()
	ctx := context.Background()
	wg := sync.WaitGroup{}
	n := 100 // количество одновременных операций

	// Запускаем n горутин, каждая создает одно событие
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e := storage.Event{ID: string(rune(i)), UserID: "u", StartTime: int64(i)}
			_ = s.CreateEvent(ctx, e)
		}(i)
	}
	wg.Wait()

	// Проверяем, что все события были созданы
	list, err := s.ListEvents(ctx, "u")
	if err != nil {
		t.Fatalf("ListEvents failed: %v", err)
	}
	if len(list) != n {
		t.Fatalf("expected %d events, got %d", n, len(list))
	}
}
