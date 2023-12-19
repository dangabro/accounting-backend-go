package service

import (
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/sheet"
)

func (d *servc) ExcelReport(start string, end string, companyId string, userId string) (data.ExcelResponse, error) {
	res := data.ExcelResponse{}
	tx, err := d.db.Begin()
	if err != nil {
		return res, err
	}
	defer rollbackTx(tx)

	err = checkRange(companyId, 1, 64, "company id is mandatory, please provide company id")
	if err != nil {
		return res, err
	}

	ids, err := dao.GetUsersForCompanyIds(tx, []string{companyId})
	if err != nil {
		return res, err
	}

	if len(ids) != 1 || ids[0] != userId {
		err = data.CreateIdError(false, "the provided company does not belong to the logged in user")
		return res, err
	}

	// ok now we process the dates
	parsedDates, err := parseDates([]string{start, end})
	if err != nil {
		return res, err
	}

	if len(parsedDates[0]) == 0 {
		return res, data.CreateIdError(false, "at least the start date must be filled out and it is not")
	}

	strStart := parsedDates[0]
	strEnd := parsedDates[1]
	company, err := dao.GetCompany(tx, companyId)
	if err != nil {
		return res, err
	}

	accounts, err := dao.GetAccountsByCompanyId(tx, companyId)
	if err != nil {
		return res, err
	}

	accountTypes, err := dao.GetAccountTypeCds(tx)
	if err != nil {
		return res, err
	}

	transactReport, err := d.TransactionReport(start, end, companyId, userId)
	if err != nil {
		return res, err
	}

	balance, err := d.TrialBalance(start, end, companyId, userId)
	if err != nil {
		return res, err
	}

	var accountResults []data.AccountReportResult = nil
	var result data.AccountReportResult

	for _, account := range accounts {
		id := account.AccountId
		result, err = d.AccountReport(start, end, id, userId)
		if err != nil {
			return res, err
		}

		accountResults = append(accountResults, result)
	}

	excelReport := data.ExcelReport{
		Start:          strStart,
		End:            strEnd,
		Company:        company,
		Balance:        balance,
		Accounts:       accounts,
		AccountTypes:   accountTypes,
		Transactions:   transactReport.Transactions,
		AccountsReport: accountResults,
	}

	buffer, fileName, err := sheet.Proc(excelReport)

	err = tx.Commit()
	if err != nil {
		return res, err
	}

	res.DataBuffer = buffer
	res.FileName = fileName

	return res, err
}
