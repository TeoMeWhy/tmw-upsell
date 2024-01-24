package transaction

import (
	"database/sql"

	"github.com/google/uuid"
)

type TransactionCart struct {
	UUID          string
	IDTransaction string
	Products      map[string]int
}

func MakeTransactionCart(idTransaction string, products map[string]int, tx *sql.Tx) error {

	statement, err := tx.Prepare("INSERT INTO tb_transactions_cart VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	for product, qtde := range products {
		_, err := statement.Exec(uuid.New().String(), idTransaction, product, qtde)
		if err != nil {
			return err
		}
	}

	return nil
}
