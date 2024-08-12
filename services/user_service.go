package services

import (
	"database/sql"
	"time"

	"github.com/troneras/ghost-listmonk-connector/database"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type UserService struct {
	db             *sql.DB
	webhookService *WebhookService
}

func NewUserService(webhookService *WebhookService) *UserService {
	return &UserService{
		db:             database.GetDB(),
		webhookService: webhookService,
	}
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, email, role, subscription_level, created_at, updated_at FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.Role, &user.SubscriptionLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, email, role, subscription_level, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Email, &user.Role, &user.SubscriptionLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(email string) (*models.User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to start transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	id := utils.GenerateUUID()
	now := time.Now()

	_, err = tx.Exec("INSERT INTO users (id, email, role, subscription_level, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		id, email, models.RoleUser, models.SubscriptionFree, now, now)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create user: %v", err)
		return nil, err
	}

	// Commit the transaction to create the user
	if err = tx.Commit(); err != nil {
		utils.ErrorLogger.Printf("Failed to commit transaction: %v", err)
		return nil, err
	}

	user := &models.User{
		ID:                id,
		Email:             email,
		Role:              models.RoleUser,
		SubscriptionLevel: models.SubscriptionFree,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Create default webhook in a separate operation
	_, err = s.webhookService.CreateWebhook(id)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create default webhook: %v", err)
		// Note: We don't return here because the user has been created successfully
		// You may want to implement a cleanup or retry mechanism for the webhook creation
	}

	return user, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	_, err := s.db.Exec("UPDATE users SET email = ?, role = ?, subscription_level = ?, updated_at = ? WHERE id = ?",
		user.Email, user.Role, user.SubscriptionLevel, time.Now(), user.ID)
	return err
}
