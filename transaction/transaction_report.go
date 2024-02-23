package transaction

import (
	"database/sql"
	"log"
)

type TransactionLineReport struct {
	CPF           string `json:"cpf"`
	Name          string `json:"Name"`
	DtTransaction string `json:"dtTransaction"`
	Points        int    `json:"points"`
	Product       string `json:"product"`
	QtdeProduct   int    `json:"qtdeProduct"`
}

type TransactionReport []TransactionLineReport

func GetCustomerTransactions(IdCustomer string, con *sql.DB) (TransactionReport, error) {

	query := `
    SELECT 
        t1.CPF,
        t1.Name,
		date(substr(t2.DtTransaction, 0,11)) AS DtTransaction,
        t2.Points,
		t3.Product,
        t3.Quantity AS Qtdeproduct

    FROM tb_customers AS t1

    LEFT JOIN tb_transactions AS t2
    ON t1.UUID = t2.idCustomer

    LEFT JOIN tb_transactions_cart As t3
    ON t2.UUID = t3.idTransaction

    WHERE t1.UUID = ?`

	report := TransactionReport{}

	prepare, err := con.Prepare(query)
	if err != nil {
		log.Println("Erro na preparação da query")
		return report, err
	}

	rows, err := prepare.Query(IdCustomer)
	if err != nil {
		log.Println("Erro na execução da query")
		return report, err
	}

	for rows.Next() {

		line := TransactionLineReport{}
		rows.Scan(
			&line.CPF,
			&line.Name,
			&line.DtTransaction,
			&line.Points,
			&line.Product,
			&line.QtdeProduct,
		)

		log.Println(line)

		report = append(report, line)
	}

	return report, nil

}
