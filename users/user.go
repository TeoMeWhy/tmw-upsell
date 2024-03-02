package users

import (
	"database/sql"
	"log"
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
