package services

import (
	"database/sql"
	"time"

	"github.com/troneras/ghost-listmonk-connector/database"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type WebhookService struct {
	db *sql.DB
}

func NewWebhookService() *WebhookService {
	return &WebhookService{db: database.GetDB()}
}

func (s *WebhookService) CreateWebhook(userID string) (*models.Webhook, error) {
	id := utils.GenerateUUID()
	endpoint := id
	secret := utils.GenerateSecret()
	now := time.Now()

	_, err := s.db.Exec("INSERT INTO webhooks (id, user_id, endpoint, secret, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		id, userID, endpoint, secret, now, now)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create webhook for user %s: %v", userID, err)
		return nil, err
	}

	utils.InfoLogger.Printf("Created webhook for user %s", userID)
	return &models.Webhook{
		ID:        id,
		UserID:    userID,
		Endpoint:  endpoint,
		Secret:    secret,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *WebhookService) GetWebhooksByUserID(userID string) ([]models.Webhook, error) {
	rows, err := s.db.Query("SELECT id, user_id, endpoint, secret, created_at, updated_at FROM webhooks WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []models.Webhook
	for rows.Next() {
		var webhook models.Webhook
		err := rows.Scan(&webhook.ID, &webhook.UserID, &webhook.Endpoint, &webhook.Secret, &webhook.CreatedAt, &webhook.UpdatedAt)
		if err != nil {
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

func (s *WebhookService) GetWebhookByEndpoint(endpoint string) (*models.Webhook, error) {
	var webhook models.Webhook
	err := s.db.QueryRow("SELECT id, user_id, endpoint, secret, created_at, updated_at FROM webhooks WHERE endpoint = ?", endpoint).
		Scan(&webhook.ID, &webhook.UserID, &webhook.Endpoint, &webhook.Secret, &webhook.CreatedAt, &webhook.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}
