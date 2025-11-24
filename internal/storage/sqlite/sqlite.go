package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Alex322322/short-url/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string, timestamp time.Time) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias, created_at) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(urlToSave, alias, timestamp)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUrlExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get id of last insert element %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrUrlNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}


func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.sqlite.DeleteUrl"
	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}