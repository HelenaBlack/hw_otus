package storage

type Notification struct {
	EventID   string `json:"eventId"`
	Title     string `json:"title"`
	StartTime int64  `json:"startTime"`
	UserID    string `json:"userId"`
}
