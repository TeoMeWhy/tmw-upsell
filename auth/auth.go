package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"points_mgmt/db"
)

var con, _ = db.Connect()

func IsValidToken(token, email string) (bool, error) {

	if token == "" || email == "" {
		return false, nil
	}

	query := `
	SELECT Token
	FROM tb_users
	WHERE Email = ?
	`

	statement, err := con.Prepare(query)
	if err != nil {
		log.Println("Erro ao preparar a query de token:", err)
		return false, err
	}

	rows, err := statement.Query(email)
	if err != nil {
		log.Println("Erro ao executar a query de token:", err)
		return false, err
	}

	var user_token string
	for rows.Next() {
		rows.Scan(&user_token)
	}

	return user_token == token, nil
}

func GenerateToken() (string, error) {

	bytes := make([]byte, 100/4*3)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Println("Erro ao ler os bytes")
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	return token, nil
}
