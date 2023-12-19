package dao

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/data"
)

func GetTransactionCompanyId(tx *sql.Tx, txnId string) (string, error) {
	res := ""
	sql := "select company_id from transaction where transaction_id = ?"
	rows, err := tx.Query(sql, txnId)
	if err != nil {
		return res, err
	}
	defer closeRows(rows)

	if rows.Next() {
		err = rows.Scan(&res)

		if err != nil {
			return res, err
		}
	} else {
		return res, data.CreateIdError(false, "cannot find transaction with the mentioned id")
	}

	return res, nil
}

func InsertTransaction(tx *sql.Tx, transaction data.Transaction) error {
	err := insertRootTransaction(tx, transaction)
	if err != nil {
		return err
	}

	pos := transaction.Positions
	for _, currentPos := range pos {
		err = insertPosition(tx, transaction.TransactionId, currentPos)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertRootTransaction(tx *sql.Tx, transaction data.Transaction) error {
	sql := "insert into transaction (transaction_id, company_id, transaction_date, sequence, comments) values (?,?,?,?,?)"
	tid := transaction.TransactionId
	cid := transaction.CompanyId
	dat := transaction.TransactionDate
	seq := transaction.Sequence
	comm := transaction.Comments

	_, err := tx.Exec(sql, tid, cid, dat, seq, comm)
	return err
}

func insertPosition(tx *sql.Tx, txnId string, pos data.TransactionPosition) error {
	sql := `insert into transaction_position 
        (transaction_position_id, transaction_id, account_id, sequence, 
         amount_debit, amount_credit, comments) 
    values 
    (?, ?, ?, ?, ?, ?, ?)`

	id := pos.TransactionPositionId
	accountId := pos.AccountId
	seq := pos.Sequence
	debit := pos.Debit
	credit := pos.Credit
	comments := pos.Comments

	_, err := tx.Exec(sql, id, txnId, accountId, seq, debit, credit, comments)

	return err
}
