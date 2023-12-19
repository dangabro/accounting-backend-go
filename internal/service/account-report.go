package service

import (
	"fmt"
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/util"
	"strings"
)

func (d *servc) AccountReport(start string, end string, accountId string, userId string) (data.AccountReportResult, error) {
	res := data.AccountReportResult{}

	tx, err := d.db.Begin()
	if err != nil {
		return res, err
	}
	defer rollbackTx(tx)

	// check accountId is filled out
	err = checkRange(accountId, 1, 64, "accountId is mandatory please provide account id")
	if err != nil {
		return res, err
	}

	// check accountId belongs to the userId
	userIds, err := dao.GetUserIdsForAccountIds(tx, []string{accountId})
	if err != nil {
		return res, err
	}

	if len(userIds) != 1 || userIds[0] != userId {
		return res, data.CreateIdError(false, "the provided account id does not belong to the logged in user")
	}

	// check the start is filled out
	strStart := strings.TrimSpace(start)
	strEnd := strings.TrimSpace(end)

	err = checkRange(strStart, 1, 64, "start date must be filled out")
	if err != nil {
		return res, err
	}

	// search ids
	parseStart, err := util.ParseStringDate(strStart)
	if err != nil {
		return res, err
	}

	parseEnd := ""
	if len(strEnd) > 0 {
		parseEnd, err = util.ParseStringDate(strEnd)
		if err != nil {
			return res, err
		}

	}

	startBalance, err := dao.GetAccountRuns(tx, parseStart, accountId)
	if err != nil {
		return res, err
	}

	// load the account
	types, err := dao.GetAccountTypeCds(tx)
	if err != nil {
		return res, err
	}

	typeMap := make(map[string]string)
	for _, tp := range types {
		typeMap[tp.AccountTypeCd] = tp.BalanceTypeCd
	}

	// load the account by accountId
	account, err := dao.GetAccountById(tx, accountId)
	if err != nil {
		return res, err
	}

	accountTypeCd := account.AccountTypeCd
	balanceTypeCd, _ := typeMap[accountTypeCd]

	if balanceTypeCd == "credit" {
		startBalance.Credit = startBalance.Credit - startBalance.Debit
		startBalance.Debit = 0
	} else if balanceTypeCd == "debit" {
		startBalance.Debit = startBalance.Debit - startBalance.Credit
		startBalance.Credit = 0
	}

	// load all transactions that are between the two dates and then contain
	transactions, err := dao.SearchTransactionsForAccount(tx, parseStart, parseEnd, accountId)
	if err != nil {
		return res, err
	}

	resStart, err := util.FormatDisplayDate(parseStart)
	if err != nil {
		return res, err
	}

	resEnd := ""
	if len(parseEnd) > 0 {
		resEnd, err = util.FormatDisplayDate(parseEnd)

		if err != nil {
			return res, err
		}
	}

	// perform calculations
	res.AccountId = accountId
	res.Code = account.Code
	res.Name = account.Name
	res.AccountTypeCd = accountTypeCd
	res.Start = resStart
	res.End = resEnd
	res.StartBalance = startBalance

	var details []data.AccountReportDetail

	companyId := account.CompanyId
	accounts, err := dao.GetAccountsByCompanyId(tx, companyId)
	if err != nil {
		return res, err
	}

	accountMap := make(map[string]data.Account)
	for _, acct := range accounts {
		accountMap[acct.AccountId] = acct
	}

	// for each transaction now we fill out the values
	balance := startBalance
	for _, txn := range transactions {
		detail, newBal, err := calculateAccountReportDetail(accountId, txn, balance, accountMap, balanceTypeCd)
		if err != nil {
			return res, nil
		}

		balance.Debit = newBal.Debit
		balance.Credit = newBal.Credit

		details = append(details, detail)
	}

	res.Details = details

	// calculate totals
	totals := data.AccountValues{}
	for _, det := range details {
		totals.Debit += det.Amount.Debit
		totals.Credit += det.Amount.Credit
	}

	res.Totals = totals

	// and now the final balance
	calculateFinalBalance(&res, balanceTypeCd)

	_ = tx.Commit()
	return res, nil
}

func calculateFinalBalance(d *data.AccountReportResult, balanceTypeCd string) {
	startDebit := d.StartBalance.Debit
	startCredit := d.StartBalance.Credit

	totalDebit := d.Totals.Debit
	totalCredit := d.Totals.Credit

	var balance data.AccountValues
	if balanceTypeCd == "debit" {
		balanceDebit := startDebit + totalDebit - totalCredit - startCredit
		balanceCredit := 0.0

		balance = data.AccountValues{
			Debit:  balanceDebit,
			Credit: balanceCredit,
		}
	} else {
		balanceCredit := startCredit + totalCredit - totalDebit - startDebit
		balanceDebit := 0.0

		balance = data.AccountValues{
			Debit:  balanceDebit,
			Credit: balanceCredit,
		}
	}

	d.FinalBalance = balance
}

func calculateAccountReportDetail(id string, txn data.Transaction, balance data.AccountValues, mp map[string]data.Account, cd string) (data.AccountReportDetail, data.AccountValues, error) {
	res := data.AccountReportDetail{}
	newBal := data.AccountValues{}

	// positions
	pos := txn.Positions
	amount := 0.0

	comments := []string{txn.Comments}
	var debitCode []string // comma delimited codes
	var creditCode []string

	for _, p := range pos {
		amount += p.Debit // total debit equal total credit, so we just consider one of them
		currentId := p.AccountId
		account, ok := mp[currentId]
		if !ok {
			message := fmt.Sprintf("cannot find account with id %s", currentId)
			return res, newBal, data.CreateIdError(false, message)
		}

		if currentId == id {
			res.TransactionPositionId = p.TransactionPositionId
			res.Amount.Debit = p.Debit
			res.Amount.Credit = p.Credit

			// add to comments
			positionComments := strings.TrimSpace(p.Comments)
			if len(positionComments) > 0 {
				comments = append(comments, positionComments)
			}
		} else {
			if p.Debit != 0.0 {
				debitCode = append(debitCode, account.Code)
			} else {
				creditCode = append(creditCode, account.Code)
			}
		}
	} // for

	res.Comments = strings.Join(comments, ", ")
	res.Date = txn.TransactionDate
	res.TransactionAmount = amount
	res.CreditCodes = strings.Join(creditCode, ",")
	res.DebitCodes = strings.Join(debitCode, ",")

	// calculate the balance
	newBal.Debit = balance.Debit + res.Amount.Debit
	newBal.Credit = balance.Credit + res.Amount.Credit

	if cd == "debit" {
		newBal.Debit = newBal.Debit - newBal.Credit
		newBal.Credit = 0
		res.CurrentBalance = newBal.Debit
	} else if cd == "credit" {
		newBal.Credit = newBal.Credit - newBal.Debit
		newBal.Debit = 0
		res.CurrentBalance = newBal.Credit
	}

	return res, newBal, nil
}
