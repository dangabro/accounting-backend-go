package service

import "database/sql"

func rollbackTx(tx *sql.Tx) {
	_ = tx.Rollback()
}
