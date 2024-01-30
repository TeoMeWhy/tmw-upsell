package customer

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Customer struct {
	UUID   string
	Name   string
	Email  string
	CPF    string
	Points int
}

func GetCustomers(tx *sql.Tx) ([]Customer, error) {
	query := "SELECT * FROM tb_customers"
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	data := []Customer{}

	for rows.Next() {
		var id, name, email, cpf string
		var points int
		rows.Scan(&id, &name, &email, &cpf, &points)
		data = append(data, Customer{id, name, email, cpf, points})
	}
	return data, nil
}

func GetCustomerByField(idValue, idField string, tx *sql.Tx) (Customer, error) {

	query := fmt.Sprintf("SELECT * FROM tb_customers WHERE %s = ?", idField)
	statement, err := tx.Prepare(query)
	if err != nil {
		return Customer{}, err
	}

	rows, err := statement.Query(idValue)
	if err != nil {
		return Customer{}, err
	}

	data := Customer{}
	for rows.Next() {
		rows.Scan(&data.UUID, &data.Name, &data.Email, &data.CPF, &data.Points)
	}

	return data, nil
}

func CreateCustomer(newCustomer Customer, tx *sql.Tx) (Customer, error) {

	id := uuid.New().String()
	newCustomer.UUID = id

	statement, err := tx.Prepare("INSERT INTO tb_customers VALUES (?,?,?,?,?);")
	if err != nil {
		return Customer{}, err
	}

	_, err = statement.Exec(
		newCustomer.UUID,
		newCustomer.Name,
		newCustomer.Email,
		newCustomer.CPF,
		newCustomer.Points)
	return newCustomer, err
}

func UpdateCustomerPoints(points int, idCustomer string, tx *sql.Tx) error {

	statement, err := tx.Prepare("UPDATE tb_customers SET Points = ? WHERE UUID = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(points, idCustomer)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCustomerEmail(newEmail string, idCustomer string, tx *sql.Tx) error {

	statement, err := tx.Prepare("UPDATE tb_customers SET Email = ? WHERE UUID = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(newEmail, idCustomer)
	if err != nil {
		return err
	}

	return nil
}
