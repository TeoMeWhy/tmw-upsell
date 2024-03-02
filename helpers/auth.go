package helpers

import (
	"log"
)

func IsValidToken(token string) (bool, error) {

	if token == "" {
		return false, nil
	}

	query := `
	SELECT count(*)
	FROM tb_users
	WHERE Token = ?
	`

	statement, err := con.Prepare(query)
	if err != nil {
		log.Println("Erro ao preparar a query de token:", err)
		return false, err
	}

	rows, err := statement.Query(token)
	if err != nil {
		log.Println("Erro ao executar a query de token:", err)
		return false, err
	}

	count := 0
	for rows.Next() {
		rows.Scan(&count)
	}

	return count > 0, nil
}
