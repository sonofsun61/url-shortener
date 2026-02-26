package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite" // init sqlite driver
)

type Storage struct {
	db *sql.DB
}

func (storage *Storage) NewStorage(storagePath string) error {
	const op = "storage.sqlite.New"
	var err error
	storage.db, err = sql.Open("sqlite", storagePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	query := (`CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	stmt, err := storage.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (storage *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storsge.sqlite.SaveURL"
	query := "INSERT INTO  url(url, alias) VALUES(?, ?)"
	stmt, err := storage.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	result, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		//TODO: Add database duplicate check
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil
}

func (storage *Storage) GetURL(alias string) (string, error) {
	const op = "storsge.sqlite.GetURL"
	query := "SELECT url FROM url WHERE alias = ?"
	stmt, err := storage.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	var res string
	err = stmt.QueryRow(alias).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: no matched urls found: %w", op, err)
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return res, nil
}

func (storage *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"
	query := "DELETE FROM url WHERE alias = ?"
	stmt, err := storage.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: failed to delete URL: %w", op, err)
	}
	return nil
}
