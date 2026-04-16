package storage

type Event struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	UserID       string `json:"userId"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	NotifyBefore *int64 `json:"notifyBefore,omitempty"`
	NotifySent   bool   `json:"notifySent"`
}
