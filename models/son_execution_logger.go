// models/execution_log.go
package models

import (
	"time"
)

type SonExecutionLog struct {
	ID           string    `json:"id"`
	SonID        string    `json:"son_id"`
	WebhookLogID string    `json:"webhook_log_id"`
	Status       string    `json:"status"`
	ExecutedAt   time.Time `json:"executed_at"`
	ErrorMessage string    `json:"error_message"`
}

type ActionExecutionLog struct {
	ID             string    `json:"id"`
	ExecutionLogID string    `json:"execution_log_id"`
	ActionType     string    `json:"action_type"`
	Status         string    `json:"status"`
	ExecutedAt     time.Time `json:"executed_at"`
	ErrorMessage   string    `json:"error_message"`
}
