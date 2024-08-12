// services/recent_activity_service.go
package services

import (
	"database/sql"

	"github.com/troneras/ghost-listmonk-connector/database"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type RecentActivityService struct {
	db *sql.DB
}

func NewRecentActivityService() *RecentActivityService {
	return &RecentActivityService{db: database.GetDB()}
}

func (s *RecentActivityService) LogActivity(userID, actionType, description string) error {
	id := utils.GenerateUUID()
	_, err := s.db.Exec(`
        INSERT INTO recent_activity (id, user_id, action_type, description)
        VALUES (?, ?, ?, ?)
    `, id, userID, actionType, description)
	return err
}

func (s *RecentActivityService) GetRecentActivity(userID string, limit int) ([]models.RecentActivity, error) {
	rows, err := s.db.Query(`
        SELECT id, user_id, action_type, description, timestamp
        FROM recent_activity
        WHERE user_id = ?
        ORDER BY timestamp DESC
        LIMIT ?
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.RecentActivity
	for rows.Next() {
		var activity models.RecentActivity
		err := rows.Scan(&activity.ID, &activity.UserID, &activity.ActionType, &activity.Description, &activity.Timestamp)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}
