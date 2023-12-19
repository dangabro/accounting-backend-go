package dao

import (
	"database/sql"
	"fmt"
	"github.com/dgb9/db-account-server/internal/data"
)

func GetAccountTypeCds(tx *sql.Tx) ([]data.AccountTypeData, error) {
	var res []data.AccountTypeData
	sql := "select account_type_cd, name, balance_type_cd from account_type order by account_type_cd"

	rows, err := tx.Query(sql)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)

	for rows.Next() {
		current := data.AccountTypeData{}
		err = rows.Scan(&current.AccountTypeCd, &current.Name, &current.BalanceTypeCd)

		if err != nil {
			return res, err
		}

		// all good
		res = append(res, current)
	}

	return res, nil
}

func GetUserIdsForAccountIds(tx *sql.Tx, ids []string) ([]string, error) {
	res := make([]string, 0)
	if len(ids) == 0 {
		return res, nil
	}

	strSql := `select distinct u.user_id 
		from user u 
			inner join company c on u.user_id = c.user_id
			inner join account a on c.company_id = a.company_id
		where 
			a.account_id in (`

	first := true
	var idAny []any

	for _, val := range ids {
		idAny = append(idAny, val)

		if first {
			first = false
		} else {
			strSql += ", "
		}

		strSql += "?"
	}

	strSql += ")"

	rows, err := tx.Query(strSql, idAny...)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)
	for rows.Next() {
		var strId string

		err = rows.Scan(&strId)
		if err != nil {
			return res, err
		}

		res = append(res, strId)
	}

	return res, nil
}

func HasTransactionsForAccountIds(tx *sql.Tx, ids []string) (bool, error) {
	res := false
	if len(ids) == 0 {
		return res, nil
	}

	var params []any
	strSql := `select count(*) from transaction_position where account_id in (`
	first := true

	for _, id := range ids {
		if first {
			first = false
		} else {
			strSql += ", "
		}

		strSql += "?"
		params = append(params, id)
	}

	strSql += ")"

	rows, err := tx.Query(strSql, params...)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)
	var count int
	rows.Next() // count * always returns something
	err = rows.Scan(&count)
	if err != nil {
		return res, err
	}

	res = count > 0
	return res, nil
}

func DeleteAccounts(tx *sql.Tx, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	strSql := `delete from account where account_id in (`
	first := true
	var idAny []any

	for _, id := range ids {
		if first {
			first = false
		} else {
			strSql += ","
		}

		strSql += "?"
		idAny = append(idAny, id)
	}

	strSql = strSql + ")"

	_, err := tx.Exec(strSql, idAny...)

	return err
}

func GetAccountsByCompanyId(tx *sql.Tx, companyId string) ([]data.Account, error) {
	res := make([]data.Account, 0)
	sql := `select
			a.account_id, a.code, a.name, a.account_type_cd, a.company_id
		from
			account a
			inner join account_type t on a.account_type_cd = t.account_type_cd
		where company_id = ?
		order by t.sequence_no, a.code`

	rows, err := tx.Query(sql, companyId)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)

	for rows.Next() {
		account := data.Account{}
		err = rows.Scan(&account.AccountId, &account.Code, &account.Name, &account.AccountTypeCd, &account.CompanyId)

		if err != nil {
			return res, err
		}

		res = append(res, account)
	}

	return res, nil
}

func GetCountAccountForCompanyIds(tx *sql.Tx, companyIds []string) (int, error) {
	res := 0

	sql := `select 
			count(*) current 
		from account a 
			inner join company c on a.company_id = c.company_id 
		where 
			c.company_id in (`

	first := true
	var idsAny []any

	for _, id := range companyIds {
		if first {
			first = false
		} else {
			sql += ","
		}

		sql += "?"
		idsAny = append(idsAny, id)
	}

	sql += ")"

	rows, err := tx.Query(sql, idsAny...)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)
	rows.Next()
	err = rows.Scan(&res)

	return res, err
}

func SearchAccounts(tx *sql.Tx, search string, companyId string) ([]data.Account, error) {
	res := make([]data.Account, 0)

	sql := `select a.account_id, a.code, a.name, a.account_type_cd, a.company_id
		from account a
				 inner join account_type t on a.account_type_cd = t.account_type_cd
		where a.company_id = ?
		  and (a.code like ? or a.name like ?)
		order by t.sequence_no, a.code`

	src := fmt.Sprintf("%%%s%%", search)

	rows, err := tx.Query(sql, companyId, src, src)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)
	for rows.Next() {
		var dt data.Account

		err = rows.Scan(&dt.AccountId, &dt.Code, &dt.Name, &dt.AccountTypeCd, &dt.CompanyId)
		if err != nil {
			return res, err
		}

		res = append(res, dt)
	}

	return res, nil
}

func EnsureAccountExists(tx *sql.Tx, id string) (bool, error) {
	res := false
	sql := "select account_id from account where account_id = ?"
	rows, err := tx.Query(sql, id)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)
	if rows.Next() {
		res = true
	}

	return res, nil
}

func AddAccount(tx *sql.Tx, account data.Account) error {
	sql := "insert into account (account_id, code, name, account_type_cd, company_id) values (?, ?, ?, ?, ?)"
	_, err := tx.Exec(sql, account.AccountId, account.Code, account.Name, account.AccountTypeCd, account.CompanyId)

	return err
}

// attention, the company id is not being updated
func UpdateAccount(tx *sql.Tx, account data.Account) error {
	sql := "update account set code = ?, name = ?, account_type_cd = ? where account_id = ?"
	_, err := tx.Exec(sql, account.Code, account.Name, account.AccountTypeCd, account.AccountId)

	return err
}

func CheckCompanyExists(tx *sql.Tx, companyId string) (bool, error) {
	res := false
	sql := "select company_id from company where company_id = ?"
	rows, err := tx.Query(sql, companyId)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)

	if rows.Next() {
		res = true
	}

	return res, nil
}

func GetCompanyIdsForAccountIds(tx *sql.Tx, accountIds []string) ([]string, error) {
	res := make([]string, 0)

	sql := "select distinct company_id from account where account_id in ("
	processIn := NewProcessIn("")
	var params []any
	for _, val := range accountIds {
		params = append(params, val)
	}

	processIn.AddIn(sql, ")", params)
	strSql, pars := processIn.GetFinalSql()
	rows, err := tx.Query(strSql, pars...)
	if err != nil {
		return res, err
	}
	defer closeRows(rows)

	for rows.Next() {
		var cid string
		err = rows.Scan(&cid)
		if err != nil {
			return res, err
		}

		res = append(res, cid)
	}

	return res, nil
}

func GetAccountRuns(tx *sql.Tx, beforeDate string, accountId string) (data.AccountValues, error) {
	res := data.AccountValues{}

	sql := `select 
			p.amount_debit, p.amount_credit 
		from 
			transaction t 
			inner join transaction_position p on t.transaction_id = p.transaction_id 
		where t.transaction_date < ? and p.account_id = ?`

	rows, err := tx.Query(sql, beforeDate, accountId)
	if err != nil {
		return res, err
	}
	defer closeRows(rows)

	for rows.Next() {
		var debit float64
		var credit float64

		err = rows.Scan(&debit, &credit)
		if err != nil {
			return res, err
		}

		res.Debit += debit
		res.Credit += credit
	}

	return res, nil
}

func GetAccountById(tx *sql.Tx, companyId string) (data.Account, error) {
	res := data.Account{}

	sql := `select
			a.account_id, a.code, a.name, a.account_type_cd, a.company_id
		from
			account a
		where account_id = ?`

	rows, err := tx.Query(sql, companyId)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)

	if rows.Next() {
		err = rows.Scan(&res.AccountId, &res.Code, &res.Name, &res.AccountTypeCd, &res.CompanyId)

		if err != nil {
			return res, err
		}
	} else {
		return res, data.CreateIdError(false, "cannot find account with the id you provided")
	}

	return res, nil
}
