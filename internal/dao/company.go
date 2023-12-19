package dao

import (
	"database/sql"
	"fmt"
	"github.com/dgb9/db-account-server/internal/data"
)

func InsertCompany(tx *sql.Tx, dt data.CompanyData) error {
	sql := "insert into company (company_id, user_id, name, month_end, day_end) values (?, ?, ?, ?, ?)"
	_, err := tx.Exec(sql, dt.CompanyId, dt.UserId, dt.Name, dt.Month, dt.Day)

	return err
}

func UpdateCompany(tx *sql.Tx, dt data.CompanyData) error {
	// you cannot update the user id
	sql := "update company set name = ?, month_end = ?, day_end = ? where company_id = ?"
	_, err := tx.Exec(sql, dt.Name, dt.Month, dt.Day, dt.CompanyId)

	return err
}

func GetCompany(tx *sql.Tx, companyId string) (data.CompanyData, error) {
	res := data.CompanyData{}
	sql := "select company_id, user_id, name, month_end, day_end from company where company_id = ?"
	rows, err := tx.Query(sql, companyId)
	if err != nil {
		return res, err
	}
	defer closeRows(rows)

	if !rows.Next() {
		message := fmt.Sprintf("cannot find company with the id: %s", companyId)
		return res, data.CreateIdError(false, message)
	}

	err = rows.Scan(&res.CompanyId, &res.UserId, &res.Name, &res.Month, &res.Day)
	return res, err
}

func GetAllCompanies(tx *sql.Tx, userId string) ([]data.CompanyData, error) {
	res := make([]data.CompanyData, 0)
	sql := "select company_id, name, month_end, day_end from company where user_id = ? order by name"
	rows, err := tx.Query(sql, userId)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)

	for rows.Next() {
		var companyId string
		var name string
		var month int
		var day int

		err = rows.Scan(&companyId, &name, &month, &day)
		if err != nil {
			return res, err
		}

		currentData := data.CompanyData{
			CompanyId: companyId,
			UserId:    userId,
			Name:      name,
			Month:     month,
			Day:       day,
		}

		res = append(res, currentData)
	}

	return res, nil
}

func GetUsersForCompanyIds(tx *sql.Tx, ids []string) ([]string, error) {
	res := make([]string, 0)
	if len(ids) == 0 {
		return res, nil
	}

	sql := "select distinct user_id from company where company_id in ("
	first := true
	var idsAny []any

	for _, id := range ids {
		if first {
			first = false
		} else {
			sql += ", "
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

func DeleteCompanies(tx *sql.Tx, ids []string) error {
	if len(ids) == 0 {
		return nil // nothing to delete
	}

	sql := "delete from company where company_id in ("
	first := true
	var idsAny []any

	for _, id := range ids {
		if first {
			first = false
		} else {
			sql += ","
		}

		sql += "?"
		idsAny = append(idsAny, id)
	}

	sql += ")"

	_, err := tx.Exec(sql, idsAny...)

	return err
}
