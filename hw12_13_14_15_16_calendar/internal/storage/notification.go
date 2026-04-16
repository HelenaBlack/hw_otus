package storage

type Notification struct {
	EventID   string `json:"eventId"`
	Title     string `json:"title"`
	StartTime int64  `json:"startTime"`
	UserID    string `json:"userId"`
}

type NotificationStatus struct {
	EventID string `json:"eventId"`
	SentAt  int64  `json:"sentAt"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
