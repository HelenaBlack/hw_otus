package storage

type Event struct {
	ID           string // уникальный идентификатор события (UUID)
	Title        string // заголовок события
	Description  string // описание события
	UserID       string // идентификатор пользователя, владельца события
	StartTime    int64  // время начала события (Unix timestamp)
	EndTime      int64  // время окончания события (Unix timestamp)
	NotifyBefore *int64 // количество секунд до события для уведомления (опционально)
}
