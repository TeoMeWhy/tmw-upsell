package customers

import (
	"database/sql"
	"log"
	"points_mgmt/db"

	"github.com/google/uuid"
)

type Customer struct {
	UUID           string
	IdOrg          string
	Name           string
	Email          string
	CPF            string
	Points         int
	TelResidencial string
	TelComercial   string
	Instagram      string
}

func GetCustomers(tx *sql.Tx) ([]Customer, error) {
	query := `
	SELECT
		COALESCE(UUID, '') AS UUID,
		COALESCE(IdOrg, '') AS IdOrg,
		COALESCE(Name, '') AS Name,
		COALESCE(Email, '') AS Email,
		COALESCE(CPF, '') AS CPF,
		COALESCE(Points, 0) AS Points,
		COALESCE(TelResidencial, '') AS TelResidencial,
		COALESCE(TelComercial, '') AS TelComercial,
		COALESCE(Instagram, '') AS Instagram
	FROM tb_customers`

	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	data := []Customer{}

	for rows.Next() {
		user := Customer{}
		rows.Scan(
			&user.UUID,
			&user.IdOrg,
			&user.Name,
			&user.Email,
			&user.CPF,
			&user.Points,
			&user.TelResidencial,
			&user.TelComercial,
			&user.Instagram,
		)
		data = append(data, user)
	}
	return data, nil
}

func GetCustomerByField(keyValue map[string]string, tx *sql.Tx) (Customer, error) {

	query := `
	SELECT
		COALESCE(UUID, '') AS UUID,
		COALESCE(IdOrg, '') AS IdOrg,
		COALESCE(Name, '') AS Name,
		COALESCE(Email, '') AS Email,
		COALESCE(CPF, '') AS CPF,
		COALESCE(Points, 0) AS Points,
		COALESCE(TelResidencial, '') AS TelResidencial,
		COALESCE(TelComercial, '') AS TelComercial,
		COALESCE(Instagram, '') AS Instagram
	FROM tb_customers
	`

	query, values := db.FormatQueryFilters(query, keyValue)

	log.Println(query)

	statement, err := tx.Prepare(query)
	if err != nil {
		return Customer{}, err
	}

	log.Println(values)

	rows, err := statement.Query(values...)
	if err != nil {
		return Customer{}, err
	}

	data := Customer{}
	for rows.Next() {
		rows.Scan(
			&data.UUID,
			&data.IdOrg,
			&data.Name,
			&data.Email,
			&data.CPF,
			&data.Points,
			&data.TelResidencial,
			&data.TelComercial,
			&data.Instagram,
		)
	}
	return data, nil
}

func CreateCustomer(newCustomer Customer, tx *sql.Tx) (Customer, error) {

	newCustomer.UUID = uuid.New().String()

	statement, err := tx.Prepare("INSERT INTO tb_customers VALUES (?,?,?,?,?,?,?,?,?);")
	if err != nil {
		return Customer{}, err
	}

	_, err = statement.Exec(
		newCustomer.UUID,
		newCustomer.IdOrg,
		newCustomer.Name,
		newCustomer.Email,
		newCustomer.CPF,
		newCustomer.Points,
		newCustomer.TelResidencial,
		newCustomer.TelComercial,
		newCustomer.Instagram)
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

func UpdateCustomerEmail(customer Customer, tx *sql.Tx) error {

	query := `
	UPDATE tb_customers SET Email = ?
	WHERE UUID = ?
	AND IdOrg = ?`

	statement, err := tx.Prepare(query)
	if err != nil {
		log.Println("Erro ao preparar o update de email")
		return err
	}

	_, err = statement.Exec(customer.Email, customer.UUID, customer.IdOrg)
	if err != nil {
		log.Println("Erro ao executar o update de email")
		return err
	}

	return nil
}

func UpdateCustomer(customer Customer, tx *sql.Tx) error {

	query := `
	UPDATE tb_customers 
	
	SET
	Name=?,
	Email=?,
	CPF=?,
	Points=?,
	TelResidencial=?,
	TelComercial=?,
	Instagram=?
	
	WHERE UUID = ?
	AND IdOrg = ?;`

	statement, err := tx.Prepare(query)
	if err != nil {
		log.Println("Erro ao preparar o statement de UpdateCustomer")
		return err
	}

	if _, err := statement.Exec(
		customer.Name,
		customer.Email,
		customer.CPF,
		customer.Points,
		customer.TelResidencial,
		customer.TelComercial,
		customer.Instagram,
		customer.UUID,
		customer.IdOrg,
	); err != nil {
		log.Println("Erro ao realizar a execução do statemente de UpdateCustomer")
	}

	return nil
}
