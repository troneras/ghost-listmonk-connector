package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/troneras/ghost-listmonk-connector/database"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

var (
	ErrSonAlreadyExists = errors.New("son with this ID already exists")
	ErrSonNotFound      = errors.New("son not found")
)

type SonStorage struct {
	db                    *sql.DB
	recentActivityService *RecentActivityService
}

// Update the NewSonStorage function
func NewSonStorage(recentActivityService *RecentActivityService) *SonStorage {
	return &SonStorage{
		db:                    database.GetDB(),
		recentActivityService: recentActivityService,
	}
}

func (s *SonStorage) Create(son *models.Son) error {
	actionsJSON, err := json.Marshal(son.Actions)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to marshal actions: %v", err)
		return err
	}

	_, err = s.db.Exec(
		"INSERT INTO sons (id, user_id, name, trigger_event, delay, actions, enabled, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
		son.ID, son.UserID, son.Name, son.Trigger, son.Delay, actionsJSON, son.Enabled,
	)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to create Son: %v", err)
		return err
	}

	if err := s.recentActivityService.LogActivity(son.UserID, "son_created", fmt.Sprintf("Created Son: %s", son.Name)); err != nil {
		utils.ErrorLogger.Printf("Failed to log activity: %v", err)
	}

	utils.InfoLogger.Infof("Created new Son with ID: %s", son.ID)
	return nil
}

func (s *SonStorage) Get(id string) (models.Son, error) {
	var son models.Son
	var actionsJSON []byte

	err := s.db.QueryRow(
		"SELECT id, user_id, name, trigger_event, delay, actions, enabled, created_at, updated_at FROM sons WHERE id = ?",
		id,
	).Scan(&son.ID, &son.UserID, &son.Name, &son.Trigger, &son.Delay, &actionsJSON, &son.Enabled, &son.CreatedAt, &son.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorLogger.Errorf("Failed to get Son: %v", ErrSonNotFound)
			return models.Son{}, ErrSonNotFound
		}
		utils.ErrorLogger.Errorf("Failed to get Son: %v", err)
		return models.Son{}, err
	}

	err = json.Unmarshal(actionsJSON, &son.Actions)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to unmarshal actions: %v", err)
		return models.Son{}, err
	}

	utils.InfoLogger.Infof("Retrieved Son with ID: %s", id)
	return son, nil
}

func (s *SonStorage) Update(son models.Son) error {
	actionsJSON, err := json.Marshal(son.Actions)
	if err != nil {
		return err
	}

	result, err := s.db.Exec(
		"UPDATE sons SET name = ?, trigger_event = ?, delay = ?, actions = ?, enabled = ?, updated_at = NOW() WHERE id = ? AND user_id = ?",
		son.Name, son.Trigger, son.Delay, actionsJSON, son.Enabled, son.ID, son.UserID,
	)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to update Son: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		utils.ErrorLogger.Errorf("Failed to update Son: %v", ErrSonNotFound)
		return ErrSonNotFound
	}

	if err := s.recentActivityService.LogActivity(son.UserID, "son_updated", fmt.Sprintf("Updated Son: %s", son.Name)); err != nil {
		utils.ErrorLogger.Printf("Failed to log activity: %v", err)
	}

	utils.InfoLogger.Infof("Updated Son with ID: %s", son.ID)
	return nil
}

func (s *SonStorage) Delete(id string, userID string) error {
	result, err := s.db.Exec("DELETE FROM sons WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to delete Son: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		utils.ErrorLogger.Errorf("Failed to delete Son: %v", ErrSonNotFound)
		return ErrSonNotFound
	}

	if err := s.recentActivityService.LogActivity(userID, "son_deleted", fmt.Sprintf("Deleted Son: %s", id)); err != nil {
		utils.ErrorLogger.Printf("Failed to log activity: %v", err)
	}

	utils.InfoLogger.Infof("Deleted Son with ID: %s", id)
	return nil
}

func (s *SonStorage) List(userID string) ([]models.Son, error) {
	rows, err := s.db.Query("SELECT id, user_id, name, trigger_event, delay, actions, enabled, created_at, updated_at FROM sons WHERE user_id = ?", userID)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to list Sons: %v", err)
		return nil, err
	}
	defer rows.Close()

	var sons []models.Son
	for rows.Next() {
		var son models.Son
		var actionsJSON []byte

		err := rows.Scan(&son.ID, &son.UserID, &son.Name, &son.Trigger, &son.Delay, &actionsJSON, &son.Enabled, &son.CreatedAt, &son.UpdatedAt)
		if err != nil {
			utils.ErrorLogger.Errorf("Failed to scan Son: %v", err)
			continue
		}

		err = json.Unmarshal(actionsJSON, &son.Actions)
		if err != nil {
			utils.ErrorLogger.Errorf("Failed to unmarshal actions: %v", err)
			continue
		}

		sons = append(sons, son)
	}

	if sons == nil {
		sons = []models.Son{}
	}

	utils.InfoLogger.Infof("Retrieved list of %d Sons for user %s", len(sons), userID)
	return sons, nil
}
