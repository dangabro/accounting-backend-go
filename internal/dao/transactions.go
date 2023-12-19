package dao

import (
	"database/sql"
	"fmt"
)

func GetUsersForTransactionIds(tx *sql.Tx, ids []string) ([]string, error) {
	res := make([]string, 0)

	sqlStart := `select distinct user_id 
		from 
			company c
		inner join transaction t on c.company_id = t.company_id
		where 
			t.transaction_id in (`

	var newIds []any
	for _, id := range ids {
		newIds = append(newIds, id)
	}

	procIn := NewProcessIn("")
	procIn.AddIn(sqlStart, ")", newIds)
	sql, params := procIn.GetFinalSql()

	rows, err := tx.Query(sql, params...)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)
	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)

		if err != nil {
			return res, err
		}

		res = append(res, userId)
	}

	return res, nil
}

func DeleteTransactions(tx *sql.Tx, ids []string) error {
	return DeleteTableIds(tx, ids, "transaction", "transaction_id")
}

func DeletePositionsForTransactions(tx *sql.Tx, ids []string) error {
	return DeleteTableIds(tx, ids, "transaction_position", "transaction_id")
}

func DeleteTableIds(tx *sql.Tx, ids []string, tableName string, field string) error {
	sql := fmt.Sprintf("delete from %s where %s in (", tableName, field)
	procIn := NewProcessIn("")

	var newIds []any
	for _, id := range ids {
		newIds = append(newIds, id)
	}

	procIn.AddIn(sql, ")", newIds)
	sql, params := procIn.GetFinalSql()

	_, err := tx.Exec(sql, params...)

	return err
}
