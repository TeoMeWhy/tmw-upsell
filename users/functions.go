package users

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"

	"github.com/google/uuid"
)

func NewUser(name, email, idorg, role string) (*User, error) {

	id_user := uuid.New().String()
	token, err := GenerateToken()
	if err != nil {
		log.Println("Erro ao gerar o token do usuário novo")
		return nil, err
	}

	user := &User{
		UUID:  id_user,
		Name:  name,
		Email: email,
		IdOrg: idorg,
		Token: token,
		Role:  role,
	}

	return user, nil
}

func GetUser(email, org string, con *sql.DB) (*User, error) {

	query := `
	SELECT
		COALESCE( UUID, '') AS UUID,
		COALESCE( Name, '') AS Name,
		COALESCE( Email, '') AS Email,
		COALESCE( idOrg, '') AS IdOrg,
		COALESCE( Token, '') AS Token,
		COALESCE( Role, '') AS Role
	FROM tb_users
	WHERE Email = ?
	AND idOrg = ?
	`

	statement, err := con.Prepare(query)
	if err != nil {
		log.Println("Erro ao preparar o GetUser", err)
		return nil, err
	}

	rows, err := statement.Query(email)
	if err != nil {
		log.Println("Erro ao executar a query de GetUser", err)
		return nil, err
	}

	u := &User{}
	for rows.Next() {
		rows.Scan(
			&u.UUID,
			&u.Name,
			&u.Email,
			&u.IdOrg,
			&u.Token,
			&u.Role,
		)
	}

	return u, nil

}

func GetUserByToken(token string, con *sql.DB) (*User, error) {

	query := `
	SELECT
		COALESCE( UUID, '') AS UUID,
		COALESCE( Name, '') AS Name,
		COALESCE( Email, '') AS Email,
		COALESCE( idOrg, '') AS IdOrg,
		COALESCE( Token, '') AS Token,
		COALESCE( Role, '') AS Role
	FROM tb_users
	WHERE Token = ?
	`

	statement, err := con.Prepare(query)
	if err != nil {
		log.Println("Erro ao preparar o GetUser", err)
		return nil, err
	}

	rows, err := statement.Query(token)
	if err != nil {
		log.Println("Erro ao executar a query de GetUser", err)
		return nil, err
	}

	u := &User{}
	for rows.Next() {
		rows.Scan(
			&u.UUID,
			&u.Name,
			&u.Email,
			&u.IdOrg,
			&u.Token,
			&u.Role,
		)
	}

	return u, nil

}

func UserExists(email, idorg string, tx *sql.Tx) (bool, error) {

	query := `
	SELECT COUNT(*)
	FROM tb_users
	WHERE Email = ?
	AND idOrg = ?
	`

	statement, err := tx.Prepare(query)
	if err != nil {
		log.Println("Erro ao preparar a query para buscar usuário")
		return false, err
	}

	rows, err := statement.Query(email, idorg)
	if err != nil {
		log.Println("Erro ao executar a query de busca de usuario")
		return false, err
	}

	count := 0
	for rows.Next() {
		rows.Scan(&count)
	}

	return count > 0, nil

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
