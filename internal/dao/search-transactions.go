package dao

import (
	"database/sql"
	"fmt"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/util"
)

func SearchTransaction(tx *sql.Tx, sd util.TransactionSearchData, companyId string) ([]data.Transaction, error) {
	ids, err := findIds(tx, sd, companyId)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	transactions, err := loadTransactionsByIds(tx, ids)
	return transactions, err
}

func processFirst(p ProcessIn, first bool, connect string) bool {
	if first {
		p.AddSql(" where ")
	} else {
		p.AddSql(connect)
	}

	return false
}

func findIds(tx *sql.Tx, sd util.TransactionSearchData, companyId string) ([]string, error) {
	sql := `select t.transaction_id from transaction t
			left join transaction_position tp on t.transaction_id = tp.transaction_id
			left join account a on tp.account_id = a.account_id`

	pin := NewProcessIn(sql)
	first := true

	first = processFirst(pin, first, "and")
	pin.AddSqlAndParam(" t.company_id = ? ", companyId)

	if len(sd.Start) > 0 {
		first = processFirst(pin, first, "and")
		pin.AddSqlAndParam("t.transaction_date >= ?", sd.Start)
	}

	if len(sd.End) > 0 {
		first = processFirst(pin, first, "and")
		pin.AddSqlAndParam("t.transaction_date <= ?", sd.End)
	}

	if len(sd.Accounts) > 0 {
		pin.AddSql(" and (")
		start := true

		for _, account := range sd.Accounts {
			padded := fmt.Sprintf("%%%s%%", account)

			if start {
				start = false
			} else {
				pin.AddSql(" or ")
			}

			pin.AddSqlAndParam(" a.code like ? ", padded)
		}
		pin.AddSql(" ) ")
	}

	if len(sd.Comments) > 0 {
		pin.AddSql(" and (")
		start := true

		for _, comment := range sd.Comments {
			padded := fmt.Sprintf("%%%s%%", comment)

			if start {
				start = false
			} else {
				pin.AddSql(" or ")
			}

			pin.AddSqlAndParam(" t.comments like ? ", padded)
		}

		pin.AddSql(" ) ")
	}

	query, pars := pin.GetFinalSql()

	rows, err := tx.Query(query, pars...)
	if err != nil {
		return nil, err
	}

	defer closeRows(rows)

	var res []string
	var id string
	for rows.Next() {
		err = rows.Scan(&id)

		if err != nil {
			return nil, err
		}

		res = append(res, id)
	}

	return res, nil
}

func loadTransactionsByIds(tx *sql.Tx, ids []string) ([]data.Transaction, error) {
	res := make([]data.Transaction, 0)
	// if no ids, then just return an empty array
	if len(ids) == 0 {
		return res, nil
	}

	startQuery := `select t.transaction_id,
			   t.company_id,
			   t.transaction_date,
			   t.sequence,
			   t.comments,
			   tp.transaction_position_id,
			   tp.account_id,
			   tp.sequence tp_sequence,
			   tp.amount_debit,
			   tp.amount_credit,
			   tp.comments tp_comments
		from transaction t
			left outer join transaction_position tp on t.transaction_id = tp.transaction_id
			where t.transaction_id in (`

	endQuery := ` ) order by t.transaction_date, t.sequence, tp.sequence`

	procIn := NewProcessIn("")
	var paramsAny []any
	for _, id := range ids {
		paramsAny = append(paramsAny, id)
	}

	procIn.AddIn(startQuery, endQuery, paramsAny)
	query, pars := procIn.GetFinalSql()

	rows, err := tx.Query(query, pars...)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	txnArray := make([]*data.Transaction, 0)
	txnMap := make(map[string]*data.Transaction)

	for rows.Next() {
		var transactionId string
		var loadedCompanyId string
		var transactionDate string
		var sequence int32
		var comments string

		var tpId sql.NullString
		var accountId sql.NullString
		var tpSequence sql.NullInt32
		var debit sql.NullFloat64
		var credit sql.NullFloat64
		var tpComments sql.NullString

		err = rows.Scan(&transactionId, &loadedCompanyId, &transactionDate, &sequence, &comments, &tpId, &accountId, &tpSequence, &debit, &credit, &tpComments)
		if err != nil {
			return nil, err
		}

		var processedDate string
		processedDate, err = util.FormatDisplayDate(transactionDate)
		if err != nil {
			return nil, err
		}

		var txn *data.Transaction
		txn, ok := txnMap[transactionId]
		if !ok {
			txn = &data.Transaction{
				TransactionId:   transactionId,
				CompanyId:       loadedCompanyId,
				TransactionDate: processedDate,
				Sequence:        sequence,
				Comments:        comments,
				Positions:       make([]data.TransactionPosition, 0),
			}

			txnArray = append(txnArray, txn)
			txnMap[transactionId] = txn
		}

		// now we got the txn, we try to load the transaction positions and then
		// attach them to txn
		if tpId.Valid {
			position := data.TransactionPosition{
				TransactionPositionId: tpId.String,
				TransactionId:         transactionId,
				AccountId:             accountId.String,
				Sequence:              tpSequence.Int32,
				Debit:                 debit.Float64,
				Credit:                credit.Float64,
				Comments:              tpComments.String,
			}

			txn.Positions = append(txn.Positions, position)
		}
	}

	for _, val := range txnArray {
		res = append(res, *val)
	}

	return res, nil
}

func findIdsByAccountId(tx *sql.Tx, start string, end string, accountId string) ([]string, error) {
	sql := `select t.transaction_id from transaction t inner join transaction_position p on t.transaction_id = p.transaction_id`

	pin := NewProcessIn(sql)
	first := true

	first = processFirst(pin, first, "and")
	pin.AddSqlAndParam(" p.account_id = ? ", accountId)

	if len(start) > 0 {
		first = processFirst(pin, first, "and")
		pin.AddSqlAndParam("t.transaction_date >= ?", start)
	}

	if len(end) > 0 {
		first = processFirst(pin, first, "and")
		pin.AddSqlAndParam("t.transaction_date <= ?", end)
	}

	query, pars := pin.GetFinalSql()

	rows, err := tx.Query(query, pars...)
	if err != nil {
		return nil, err
	}

	defer closeRows(rows)

	var res []string
	var id string
	for rows.Next() {
		err = rows.Scan(&id)

		if err != nil {
			return nil, err
		}

		res = append(res, id)
	}

	return res, nil
}

func SearchTransactionsForAccount(tx *sql.Tx, start string, end string, accountId string) ([]data.Transaction, error) {
	ids, err := findIdsByAccountId(tx, start, end, accountId)
	if err != nil {
		return nil, err
	}

	txns, err := loadTransactionsByIds(tx, ids)
	if err != nil {
		return nil, err
	}

	return txns, nil
}
