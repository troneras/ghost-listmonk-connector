package services

import (
	"database/sql"
	"time"

	"github.com/troneras/ghost-listmonk-connector/database"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type MagicLinkService struct {
	db *sql.DB
}

func NewMagicLinkService() *MagicLinkService {
	return &MagicLinkService{db: database.GetDB()}
}

func (s *MagicLinkService) CreateToken(userID string) (string, error) {
	token := utils.GenerateUUID()
	expiresAt := time.Now().Add(15 * time.Minute)

	_, err := s.db.Exec("INSERT INTO magic_links (token, user_id, expires_at) VALUES (?, ?, ?)", token, userID, expiresAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *MagicLinkService) VerifyToken(token string) (string, error) {
	var userID string
	var expiresAt time.Time

	err := s.db.QueryRow("SELECT user_id, expires_at FROM magic_links WHERE token = ?", token).Scan(&userID, &expiresAt)
	if err != nil {
		return "", err
	}

	if time.Now().After(expiresAt) {
		return "", utils.NewError("TokenExpired", "Token has expired")
	}

	// Delete the used token
	_, err = s.db.Exec("DELETE FROM magic_links WHERE token = ?", token)
	if err != nil {
		return "", err
	}

	return userID, nil
}
