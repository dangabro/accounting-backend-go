package dao

import (
	"database/sql"
	"fmt"
	"github.com/dgb9/db-account-server/internal/data"
	"strings"
)

func BalanceReport(tx *sql.Tx, start string, end string, accounts []data.Account) ([]*data.TrialBalance, error) {
	var res []*data.TrialBalance
	aMap := make(map[string]*data.TrialBalance)
	var ids []any

	// you can run it with no accounts but then must return immediately
	if len(accounts) == 0 {
		return res, nil
	}

	// fill the map
	for _, account := range accounts {
		accountId := account.AccountId
		code := account.Code
		name := account.Name
		accountTypeCd := account.AccountTypeCd

		ids = append(ids, accountId)

		container := data.TrialBalance{
			AccountId:     accountId,
			Code:          code,
			Name:          name,
			AccountTypeCd: accountTypeCd,
			StartBalance: data.AccountValues{
				Debit:  0,
				Credit: 0,
			},
			Runs: data.AccountValues{
				Debit:  0,
				Credit: 0,
			},
		}

		aMap[accountId] = &container
		res = append(res, &container)
	}

	// put the query
	first := true
	proc := NewProcessIn(starterSqlTrialBalance)

	strEnd := strings.TrimSpace(end)
	if len(strEnd) > 0 {
		first = proc.ProcessFirst(first, " and ", " where ")
		proc.AddSqlAndParam(" t.transaction_date <= ? ", strEnd)
	}

	first = proc.ProcessFirst(first, " and ", " where ")
	proc.AddIn(" tp.account_id in (", ")", ids)

	finalSql, params := proc.GetFinalSql()
	rows, err := tx.Query(finalSql, params...)
	if err != nil {
		return nil, err
	}

	defer closeRows(rows)

	for rows.Next() {
		var txnId string
		var dt string
		var acctId string
		var debit float64
		var credit float64

		err = rows.Scan(&txnId, &dt, &acctId, &debit, &credit)
		if err != nil {
			return nil, err
		}

		// get the account id
		trb, ok := aMap[acctId]
		if !ok {
			message := fmt.Sprintf("programmer error cannot find entry for trial balance for accoutId %s", acctId)
			return nil, data.CreateIdError(false, message)
		}

		if dt < start {
			trb.StartBalance.Debit += debit
			trb.StartBalance.Credit += credit
		} else {
			trb.Runs.Debit += debit
			trb.Runs.Credit += credit
		}
	}

	return res, nil
}
