// services/son_execution_logger.go
package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/troneras/ghost-listmonk-connector/database"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type SonExecutionLogger struct {
	db    *sql.DB
	redis *redis.Client
}

func NewSonExecutionLogger(redisAddr string) *SonExecutionLogger {
	return &SonExecutionLogger{
		db:    database.GetDB(),
		redis: redis.NewClient(&redis.Options{Addr: redisAddr}),
	}
}

type SonStats struct {
	Name       string `json:"name"`
	Executions int    `json:"executions"`
	Success    int    `json:"success"`
	Failure    int    `json:"failure"`
}

func (l *SonExecutionLogger) LogSonExecution(sonID, webhookLogID string, status string, errorMessage string) (string, error) {
	executionID := utils.GenerateUUID()
	_, err := l.db.Exec(`
		INSERT INTO son_execution_logs (id, son_id, webhook_log_id, execution_status, error_message)
		VALUES (?, ?, ?, ?, ?)
	`, executionID, sonID, webhookLogID, status, errorMessage)
	if err != nil {
		return "", err
	}
	return executionID, nil
}

func (l *SonExecutionLogger) LogActionExecution(executionID string, actionType string, status string, errorMessage string) error {
	_, err := l.db.Exec(`
		INSERT INTO son_execution_action_logs (id, son_execution_log_id, action_type, action_status, error_message)
		VALUES (?, ?, ?, ?, ?)
	`, utils.GenerateUUID(), executionID, actionType, status, errorMessage)

	if status == "failure" {
		l.UpdateSonExecutionStatus(executionID, "failure", errorMessage)
	}

	return err
}

func (l *SonExecutionLogger) GetSonExecutionLogs(userID string, limit, offset int) ([]models.SonExecutionLog, int, error) {
	var total int
	err := l.db.QueryRow(`
		SELECT COUNT(*) 
		FROM son_execution_logs sel
		JOIN sons s ON sel.son_id = s.id
		WHERE s.user_id = ?
	`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := l.db.Query(`
		SELECT sel.id, sel.son_id, sel.webhook_log_id, sel.execution_status, sel.executed_at, sel.error_message
		FROM son_execution_logs sel
		JOIN sons s ON sel.son_id = s.id
		WHERE s.user_id = ?
		ORDER BY sel.executed_at DESC
		LIMIT ? OFFSET ?
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.SonExecutionLog
	for rows.Next() {
		var log models.SonExecutionLog
		err := rows.Scan(&log.ID, &log.SonID, &log.WebhookLogID, &log.Status, &log.ExecutedAt, &log.ErrorMessage)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

func (l *SonExecutionLogger) GetActionExecutionLogs(executionID string) ([]models.ActionExecutionLog, error) {
	rows, err := l.db.Query(`
		SELECT id, son_execution_log_id, action_type, action_status, executed_at, error_message
		FROM son_execution_action_logs
		WHERE son_execution_log_id = ?
		ORDER BY executed_at ASC
	`, executionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ActionExecutionLog
	for rows.Next() {
		var log models.ActionExecutionLog
		err := rows.Scan(&log.ID, &log.ExecutionLogID, &log.ActionType, &log.Status, &log.ExecutedAt, &log.ErrorMessage)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (l *SonExecutionLogger) UpdateSonExecutionStatus(executionID string, status string, errorMessage string) error {
	_, err := l.db.Exec(`
        UPDATE son_execution_logs
        SET execution_status = ?, error_message = ?
        WHERE id = ?
    `, status, errorMessage, executionID)
	return err
}

func (l *SonExecutionLogger) GetSonStats(ctx context.Context, userID string, timeframe string) ([]SonStats, error) {
	cacheKey := fmt.Sprintf("son_stats:%s:%s", userID, timeframe)

	// Try to get from cache
	cachedStats, err := l.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var stats []SonStats
		err = json.Unmarshal([]byte(cachedStats), &stats)
		if err == nil {
			return stats, nil
		}
	}

	// If not in cache or error, fetch from database
	duration, err := utils.ParseDuration(timeframe)
	if err != nil {
		return nil, err
	}

	query := `
        SELECT s.id, s.name, 
               COUNT(*) as executions,
               SUM(CASE WHEN sel.execution_status = 'success' THEN 1 ELSE 0 END) as success,
               SUM(CASE WHEN sel.execution_status = 'failure' THEN 1 ELSE 0 END) as failure
        FROM sons s
        LEFT JOIN son_execution_logs sel ON s.id = sel.son_id
        WHERE s.user_id = ? AND sel.executed_at >= ?
        GROUP BY s.id, s.name
    `

	rows, err := l.db.Query(query, userID, time.Now().Add(-duration))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []SonStats
	for rows.Next() {
		var stat SonStats
		var sonID string
		err := rows.Scan(&sonID, &stat.Name, &stat.Executions, &stat.Success, &stat.Failure)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	// Cache the result
	statsJSON, err := json.Marshal(stats)
	if err == nil {
		l.redis.Set(ctx, cacheKey, statsJSON, 5*time.Minute)
	}

	return stats, nil
}
