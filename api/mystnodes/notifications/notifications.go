package notifications

import "time"

type Response struct {
	Notifications []Notification `json:"notifications"`
}

type Notification struct {
	ID        string    `json:"id"`
	UserId    string    `json:"userId"`
	Closed    bool      `json:"closed"`
	Icon      string    `json:"icon"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Footer    string    `json:"footer"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Type      string    `json:"type"`
}
