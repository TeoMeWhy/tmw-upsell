package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

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

func FormatQueryFilters(query string, filters map[string]string) (string, []interface{}) {

	where := "WHERE"
	values := []string{}
	for k, v := range filters {
		values = append(values, v)
		where += fmt.Sprintf(` %s = ? AND`, k)
	}

	where = strings.TrimSuffix(where, "AND")

	valuesInterface := []interface{}{}
	for _, v := range values {
		valuesInterface = append(valuesInterface, v)
	}

	return query + where, valuesInterface
}
