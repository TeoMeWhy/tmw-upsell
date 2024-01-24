package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sql.DB, error) {
	con, err := sql.Open("sqlite3", "data/database.db")
	if err != nil {
		return nil, err
	}
	return con, nil
}

func ImportQuery(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	query := string(content)
	return query, nil
}
