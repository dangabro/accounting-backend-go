package service

import (
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/util"
	"strings"
)

func (d *servc) TrialBalance(start string, end string, companyId string, userId string) (data.BalanceResult, error) {
	var res data.BalanceResult
	tx, err := d.db.Begin()
	if err != nil {
		return res, err
	}
	defer rollbackTx(tx)

	// check company id is good
	company, err := dao.GetCompany(tx, companyId)
	if err != nil {
		return res, err
	}

	// check company id belongs to the user
	if company.UserId != userId {
		return res, data.CreateIdError(false, "The company does not belong to the logged in user")
	}

	// get accounts for company
	accounts, err := dao.GetAccountsByCompanyId(tx, companyId)
	if err != nil {
		return res, err
	}

	// map of account id and account type cd
	atcMap := make(map[string]string)
	for _, account := range accounts {
		atcMap[account.AccountId] = account.AccountTypeCd
	}

	strStart := strings.TrimSpace(start)
	strEnd := strings.TrimSpace(end)

	// first date must be filled out, second does not have, however if it is it will be parsed
	if len(strStart) == 0 {
		return res, data.CreateIdError(false, "start date is mandatory for trial balance")
	}

	strStart, err = util.ParseStringDate(strStart)
	if err != nil {
		return res, err
	}

	if len(strEnd) > 0 {
		strEnd, err = util.ParseStringDate(strEnd)

		if err != nil {
			return res, err
		}
	}

	// run the query
	trialBalances, err := dao.BalanceReport(tx, strStart, strEnd, accounts)
	if err != nil {
		return res, err
	}

	res.CompanyId = companyId
	formattedStart, err := util.FormatDisplayDate(strStart)
	if err != nil {
		return res, err
	}
	res.StartDate = formattedStart

	var formattedEnd string
	if len(strEnd) > 0 {
		formattedEnd, err = util.FormatDisplayDate(strEnd)
		if err != nil {
			return res, err
		}

		res.EndDate = formattedEnd
	}

	// get balance type cd
	types, err := dao.GetAccountTypeCds(tx)
	if err != nil {
		return res, err
	}
	typeMap := make(map[string]string)
	for _, accountType := range types {
		accountTypeCd := accountType.AccountTypeCd
		balanceTypeCd := accountType.BalanceTypeCd

		typeMap[accountTypeCd] = balanceTypeCd
	}

	// calculate starter balances based on balance type cd
	for _, trb := range trialBalances {
		debit := trb.StartBalance.Debit
		credit := trb.StartBalance.Credit

		accountTypeCd, ok := atcMap[trb.AccountId]
		if !ok {
			return res, data.CreateIdError(false, "cannot find account in the map")
		}

		balanceType, ok := typeMap[accountTypeCd]

		if balanceType == "debit" {
			trb.StartBalance.Debit = debit - credit
			trb.StartBalance.Credit = 0

			endDebit := trb.StartBalance.Debit + trb.Runs.Debit - trb.Runs.Credit
			trb.EndBalance.Debit = endDebit
			trb.EndBalance.Credit = 0
		} else if balanceType == "credit" {
			trb.StartBalance.Debit = 0
			trb.StartBalance.Credit = credit - debit

			endCredit := trb.StartBalance.Credit + trb.Runs.Credit - trb.Runs.Debit
			trb.EndBalance.Debit = 0
			trb.EndBalance.Credit = endCredit
		} else {
			return res, data.CreateIdError(false, "cannot find the balance type: "+balanceType)
		}
	}

	var finalValues []data.TrialBalance
	for _, point := range trialBalances {
		finalValues = append(finalValues, *point)
	}

	res.Values = finalValues

	// calculate totals now
	totals := data.TrialBalance{}
	for _, v := range finalValues {
		totals.StartBalance.Debit += v.StartBalance.Debit
		totals.StartBalance.Credit += v.StartBalance.Credit

		totals.Runs.Debit += v.Runs.Debit
		totals.Runs.Credit += v.Runs.Credit

		totals.EndBalance.Debit += v.EndBalance.Debit
		totals.EndBalance.Credit += v.EndBalance.Credit
	}

	res.Totals = totals

	err = tx.Commit()
	if res.Values == nil {
		res.Values = []data.TrialBalance{}
	}

	return res, err
}
