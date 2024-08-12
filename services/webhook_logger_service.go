package services

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/troneras/ghost-listmonk-connector/database"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type WebhookLogger struct {
	db *sql.DB
}

func NewWebhookLogger() *WebhookLogger {
	return &WebhookLogger{db: database.GetDB()}
}

func (l *WebhookLogger) CreateWebhookLog(userID string, req *http.Request, body []byte) (string, error) {
	// Convert headers to JSON
	headerMap := make(map[string]string)
	for k, v := range req.Header {
		headerMap[k] = v[0]
	}
	headersJSON, err := json.Marshal(headerMap)
	if err != nil {
		return "", err
	}

	logID := utils.GenerateUUID()

	// Insert initial log into database
	_, err = l.db.Exec(`
		INSERT INTO webhook_logs (id, user_id, timestamp, method, path, headers, body, status_code, duration)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, logID, userID, time.Now(), req.Method, req.URL.Path, string(headersJSON), string(body), 200, 0)

	if err != nil {
		utils.ErrorLogger.Errorf("Failed to insert initial webhook log: %v", err)
		return "", err
	}

	return logID, nil
}

func (l *WebhookLogger) UpdateWebhookLog(logID string, statusCode int, response interface{}, duration time.Duration) error {
	// Convert response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	// Update log in database
	_, err = l.db.Exec(`
		UPDATE webhook_logs
		SET status_code = ?, response_body = ?, duration = ?
		WHERE id = ?
	`, statusCode, string(responseJSON), int(duration.Milliseconds()), logID)

	if err != nil {
		utils.ErrorLogger.Errorf("Failed to update webhook log: %v", err)
		return err
	}

	return nil
}

func (l *WebhookLogger) GetWebhookLogs(userID string, limit, offset int) ([]WebhookLog, int, error) {
	// First, get the total count of logs for this user
	var total int
	err := l.db.QueryRow("SELECT COUNT(*) FROM webhook_logs WHERE user_id = ?", userID).Scan(&total)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get total log count: %v", err)
		return nil, 0, err
	}

	// Now, get the paginated logs
	rows, err := l.db.Query(`
		SELECT id, timestamp, method, path, status_code, duration
		FROM webhook_logs
		WHERE user_id = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, userID, limit, offset)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to query webhook logs: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var logs []WebhookLog
	for rows.Next() {
		var log WebhookLog
		err := rows.Scan(&log.ID, &log.Timestamp, &log.Method, &log.Path, &log.StatusCode, &log.Duration)
		if err != nil {
			utils.ErrorLogger.Printf("Failed to scan webhook log: %v", err)
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		utils.ErrorLogger.Printf("Error iterating over webhook logs: %v", err)
		return nil, 0, err
	}

	utils.InfoLogger.Printf("Retrieved %d webhook logs for user %s (total: %d)", len(logs), userID, total)
	return logs, total, nil
}

func (l *WebhookLogger) GetWebhookLogDetails(id string) (*WebhookLogDetails, error) {
	var log WebhookLogDetails
	err := l.db.QueryRow(`
		SELECT id, user_id, timestamp, method, path, headers, body, status_code, response_body, duration
		FROM webhook_logs
		WHERE id = ?
	`, id).Scan(&log.ID, &log.UserID, &log.Timestamp, &log.Method, &log.Path, &log.Headers, &log.Body, &log.StatusCode, &log.ResponseBody, &log.Duration)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (l *WebhookLogger) GetWebhookLogForReplay(id string) (*WebhookLogDetails, error) {
	var log WebhookLogDetails
	err := l.db.QueryRow(`
		SELECT id, user_id, timestamp, method, path, headers, body, status_code, response_body, duration
		FROM webhook_logs
		WHERE id = ?
	`, id).Scan(&log.ID, &log.UserID, &log.Timestamp, &log.Method, &log.Path, &log.Headers, &log.Body, &log.StatusCode, &log.ResponseBody, &log.Duration)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

type WebhookLog struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Timestamp  time.Time `json:"timestamp"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status_code"`
	Duration   int       `json:"duration"`
}

type WebhookLogDetails struct {
	WebhookLog
	Headers      string `json:"headers"`
	Body         string `json:"body"`
	ResponseBody string `json:"response_body"`
}
