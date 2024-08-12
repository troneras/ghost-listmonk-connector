package models

import "time"

type RecentActivity struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ActionType  string    `json:"action_type"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}
