package service

import (
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/util"
	"strings"
)

func (d *servc) TransactionReport(start string, end string, companyId string, userId string) (data.TransactionReportResult, error) {
	res := data.TransactionReportResult{}
	tx, err := d.db.Begin()
	if err != nil {
		return res, err
	}
	defer rollbackTx(tx)

	userIds, err := dao.GetUsersForCompanyIds(tx, []string{companyId})
	if len(userIds) != 1 || userIds[0] != userId {
		return res, data.CreateIdError(false, "the provided company does not belong to the logged in user")
	}

	dates, err := parseDates([]string{start, end})
	if err != nil {
		return res, err
	}

	parsedStart := dates[0]
	parsedEnd := dates[1]

	txns, err := dao.SearchTransaction(tx, util.TransactionSearchData{Start: parsedStart, End: parsedEnd}, companyId)
	if err != nil {
		return res, err
	}

	res.Start, err = util.FormatDisplayDate(parsedStart)
	if err != nil {
		return res, err
	}

	res.End, err = util.FormatDisplayDate(parsedEnd)
	if err != nil {
		return res, err
	}

	res.CompanyId = companyId
	res.Transactions = txns
	if res.Transactions == nil {
		res.Transactions = []data.Transaction{}
	}

	_ = tx.Commit() // it is just a search after all
	return res, nil
}

func parseDates(dates []string) ([]string, error) {
	var res []string

	for _, date := range dates {
		strDate := strings.TrimSpace(date)
		parsedDate := ""
		var err error

		if len(strDate) > 0 {
			parsedDate, err = util.ParseStringDate(date)
			if err != nil {
				return nil, err
			}
		}

		res = append(res, parsedDate)
	}

	return res, nil
}
