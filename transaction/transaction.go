package transaction

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	UUID          string
	IdCustomer    string
	DtTransaction string
	Points        int
}

func MakeTransaction(points int, IdCustomer string, tx *sql.Tx) (string, error) {

	idTransaction := uuid.New().String()
	dtTransaction := time.Now().UTC().String()

	statement, err := tx.Prepare("INSERT INTO tb_transactions VALUES(?,?,?,?);")
	if err != nil {
		return "", err
	}

	if _, err := statement.Exec(idTransaction, IdCustomer, dtTransaction, points); err != nil {
		return "", err
	}
	return idTransaction, nil
}
