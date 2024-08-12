package services

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &Database{db: db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *Database) createTables() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS sons (
			id TEXT PRIMARY KEY,
			name TEXT,
			trigger TEXT,
			delay INTEGER,
			created_at DATETIME,
			updated_at DATETIME
		);
		CREATE TABLE IF NOT EXISTS actions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			son_id TEXT,
			type TEXT,
			parameters TEXT,
			FOREIGN KEY (son_id) REFERENCES sons(id)
		);
	`)
	return err
}

func (d *Database) CreateSon(son *models.Son) error {
	son.ID = utils.GenerateUUID()
	son.CreatedAt = time.Now()
	son.UpdatedAt = time.Now()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO sons (id, name, trigger, delay, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		son.ID, son.Name, son.Trigger, int64(son.Delay), son.CreatedAt, son.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, action := range son.Actions {
		parametersJSON, err := json.Marshal(action.Parameters)
		if err != nil {
			return err
		}
		_, err = tx.Exec(
			"INSERT INTO actions (son_id, type, parameters) VALUES (?, ?, ?)",
			son.ID, action.Type, string(parametersJSON),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) GetSon(id string) (*models.Son, error) {
	son := &models.Son{}
	err := d.db.QueryRow(
		"SELECT id, name, trigger, delay, created_at, updated_at FROM sons WHERE id = ?",
		id,
	).Scan(&son.ID, &son.Name, &son.Trigger, &son.Delay, &son.CreatedAt, &son.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSonNotFound
		}
		return nil, err
	}

	rows, err := d.db.Query("SELECT type, parameters FROM actions WHERE son_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var action models.Action
		var parametersJSON string
		err := rows.Scan(&action.Type, &parametersJSON)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(parametersJSON), &action.Parameters)
		if err != nil {
			return nil, err
		}
		son.Actions = append(son.Actions, action)
	}

	return son, nil
}

func (d *Database) UpdateSon(son *models.Son) error {
	son.UpdatedAt = time.Now()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE sons SET name = ?, trigger = ?, delay = ?, updated_at = ? WHERE id = ?",
		son.Name, son.Trigger, int64(son.Delay), son.UpdatedAt, son.ID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM actions WHERE son_id = ?", son.ID)
	if err != nil {
		return err
	}

	for _, action := range son.Actions {
		parametersJSON, err := json.Marshal(action.Parameters)
		if err != nil {
			return err
		}
		_, err = tx.Exec(
			"INSERT INTO actions (son_id, type, parameters) VALUES (?, ?, ?)",
			son.ID, action.Type, string(parametersJSON),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) DeleteSon(id string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM actions WHERE son_id = ?", id)
	if err != nil {
		return err
	}

	result, err := tx.Exec("DELETE FROM sons WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrSonNotFound
	}

	return tx.Commit()
}

func (d *Database) ListSons() ([]models.Son, error) {
	rows, err := d.db.Query("SELECT id, name, trigger, delay, created_at, updated_at FROM sons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sons []models.Son
	for rows.Next() {
		var son models.Son
		err := rows.Scan(&son.ID, &son.Name, &son.Trigger, &son.Delay, &son.CreatedAt, &son.UpdatedAt)
		if err != nil {
			return nil, err
		}
		sons = append(sons, son)
	}

	for i, son := range sons {
		actionRows, err := d.db.Query("SELECT type, parameters FROM actions WHERE son_id = ?", son.ID)
		if err != nil {
			return nil, err
		}
		defer actionRows.Close()

		for actionRows.Next() {
			var action models.Action
			var parametersJSON string
			err := actionRows.Scan(&action.Type, &parametersJSON)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal([]byte(parametersJSON), &action.Parameters)
			if err != nil {
				return nil, err
			}
			sons[i].Actions = append(sons[i].Actions, action)
		}
	}

	return sons, nil
}
