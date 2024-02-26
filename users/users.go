package users

import (
	"database/sql"
	"log"
	"points_mgmt/auth"

	"github.com/google/uuid"
)

type User struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	IdOrg string `json:"idorg"`
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (u *User) CreateUser(tx *sql.Tx) error {

	query := `
	INSERT INTO
	tb_users (UUID,Name,Email,idOrg,Token,Role)
	VALUES (?,?,?,?,?,?);
	`

	statement, err := tx.Prepare(query)
	if err != nil {
		log.Println("Erro no statement de novo usuário")
		return err
	}

	if _, err := statement.Exec(
		u.UUID,
		u.Name,
		u.Email,
		u.IdOrg,
		u.Token,
		u.Role,
	); err != nil {
		log.Println("Erro ao executar a transação de salvar usuário")
	}

	return nil
}

func (u *User) DeleteUser(tx *sql.Tx) error {
	query := `
	DELETE
	FROM tb_users
	WHERE UUID = ?;
	`

	statement, err := tx.Prepare(query)
	if err != nil {
		log.Println("Erro no statement de novo usuário")
		return err
	}

	if _, err := statement.Exec(u.UUID); err != nil {
		log.Println("Erro ao executar a transação de salvar usuário")
	}

	return nil
}

func (u *User) UpdateUser(con *sql.DB) error {

	tx, err := con.Begin()
	if err != nil {
		log.Println("Erro ao abrir a transaction do usuário", err)
		return err
	}

	if err := u.DeleteUser(tx); err != nil {
		tx.Rollback()
		log.Println("Erro na deleção do usuário", err)
		return err
	}

	if err := u.CreateUser(tx); err != nil {
		tx.Rollback()
		log.Println("Erro na deleção do usuário", err)
		return err
	}

	return nil
}

func NewUser(name, email, idorg, role string) (*User, error) {

	id_user := uuid.New().String()
	token, err := auth.GenerateToken()
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

func GetUser(email string, tx *sql.Tx) (*User, error) {

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
	`

	statement, err := tx.Prepare(query)
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
