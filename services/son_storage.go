package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

var (
	ErrSonAlreadyExists = errors.New("son with this ID already exists")
	ErrSonNotFound      = errors.New("son not found")
)

type SonStorage struct {
	db *sql.DB
}

func NewSonStorage(dbPath string) (*SonStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	storage := &SonStorage{db: db}
	if err := storage.createTables(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *SonStorage) createTables() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS sons (
			id TEXT PRIMARY KEY,
			name TEXT,
			trigger TEXT,
			delay INTEGER,
			actions TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)
	`)
	return err
}

func (s *SonStorage) Create(son *models.Son) error {
	son.ID = utils.GenerateUUID()
	son.CreatedAt = time.Now()
	son.UpdatedAt = time.Now()

	actionsJSON, err := json.Marshal(son.Actions)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		"INSERT INTO sons (id, name, trigger, delay, actions, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		son.ID, son.Name, son.Trigger, int64(son.Delay), actionsJSON, son.CreatedAt, son.UpdatedAt,
	)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to create Son: %v", err)
		return err
	}

	utils.InfoLogger.Infof("Created new Son with ID: %s", son.ID)
	return nil
}

func (s *SonStorage) Get(id string) (models.Son, error) {
	var son models.Son
	var actionsJSON []byte
	var delayInt int64

	err := s.db.QueryRow(
		"SELECT id, name, trigger, delay, actions, created_at, updated_at FROM sons WHERE id = ?",
		id,
	).Scan(&son.ID, &son.Name, &son.Trigger, &delayInt, &actionsJSON, &son.CreatedAt, &son.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorLogger.Errorf("Failed to get Son: %v", ErrSonNotFound)
			return models.Son{}, ErrSonNotFound
		}
		utils.ErrorLogger.Errorf("Failed to get Son: %v", err)
		return models.Son{}, err
	}

	son.Delay = models.Duration(time.Duration(delayInt))

	err = json.Unmarshal(actionsJSON, &son.Actions)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to unmarshal actions: %v", err)
		return models.Son{}, err
	}

	utils.InfoLogger.Infof("Retrieved Son with ID: %s", id)
	return son, nil
}

func (s *SonStorage) Update(son models.Son) error {
	son.UpdatedAt = time.Now()

	actionsJSON, err := json.Marshal(son.Actions)
	if err != nil {
		return err
	}

	result, err := s.db.Exec(
		"UPDATE sons SET name = ?, trigger = ?, delay = ?, actions = ?, updated_at = ? WHERE id = ?",
		son.Name, son.Trigger, int64(son.Delay), actionsJSON, son.UpdatedAt, son.ID,
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

	utils.InfoLogger.Infof("Updated Son with ID: %s", son.ID)
	return nil
}

func (s *SonStorage) Delete(id string) error {
	result, err := s.db.Exec("DELETE FROM sons WHERE id = ?", id)
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

	utils.InfoLogger.Infof("Deleted Son with ID: %s", id)
	return nil
}

func (s *SonStorage) List() ([]models.Son, error) {
	rows, err := s.db.Query("SELECT id, name, trigger, delay, actions, created_at, updated_at FROM sons")
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to list Sons: %v", err)
		return nil, err
	}
	defer rows.Close()

	var sons []models.Son
	for rows.Next() {
		var son models.Son
		var actionsJSON []byte
		var delayInt int64

		err := rows.Scan(&son.ID, &son.Name, &son.Trigger, &delayInt, &actionsJSON, &son.CreatedAt, &son.UpdatedAt)
		if err != nil {
			utils.ErrorLogger.Errorf("Failed to scan Son: %v", err)
			continue
		}

		son.Delay = models.Duration(time.Duration(delayInt))

		err = json.Unmarshal(actionsJSON, &son.Actions)
		if err != nil {
			utils.ErrorLogger.Errorf("Failed to unmarshal actions: %v", err)
			continue
		}

		sons = append(sons, son)
	}

	if sons == nil {
		sons = []models.Son{} // Ensure we always return an array, even if empty
	}

	utils.InfoLogger.Infof("Retrieved list of %d Sons", len(sons))
	return sons, nil
}
