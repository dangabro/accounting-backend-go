package service

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/util"
)

func (d *servc) UpdateTransaction(adding bool, transaction data.Transaction, userId string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = checkTransaction(transaction)
	if err != nil {
		return err
	}

	// check company belongs to user and if not adding that the company exists
	err = checkTransactionBelongsToUser(tx, adding, transaction, userId)
	if err != nil {
		return err
	}

	err = checkAccountsBelongToCompanyId(tx, transaction)
	if err != nil {
		return err
	}

	err = checkGeneralTransactionLayout(transaction)
	if err != nil {
		return err
	}

	// format the transaction date
	formattedDate, err := util.ParseStringDate(transaction.TransactionDate)
	if err != nil {
		return err
	}

	transaction.TransactionDate = formattedDate

	// delete existent transaction
	if !adding {
		err = dao.DeletePositionsForTransactions(tx, []string{transaction.TransactionId})
		if err != nil {
			return err
		}

		err = dao.DeleteTransactions(tx, []string{transaction.TransactionId})
		if err != nil {
			return err
		}
	}

	// add new transaction
	err = dao.InsertTransaction(tx, transaction)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func checkGeneralTransactionLayout(transaction data.Transaction) error {
	// at least two positions
	nrPos := len(transaction.Positions)
	if nrPos < 2 {
		return data.CreateIdError(false, "transaction must have at least two positions")
	}
	// debit equal with credit
	debit := 0.0
	credit := 0.0

	pos := transaction.Positions
	for _, val := range pos {
		debit += val.Debit
		credit += val.Credit
	}

	if debit != credit {
		return data.CreateIdError(false, "total debit and credit are different")
	}

	// not the same account twice in the transaction
	acctMap := make(map[string]string)
	for _, val := range pos {
		accountId := val.AccountId

		_, ok := acctMap[accountId]
		if !ok {
			acctMap[accountId] = ""
		} else {
			return data.CreateIdError(false, "at least one account present for two different transaction positions")
		}
	}

	// no position with both debit and credit not null
	for _, val := range pos {
		debit := val.Debit
		credit := val.Credit

		if debit != 0.0 && credit != 0.0 {
			return data.CreateIdError(false, "at least one position has both debit and credit not zero which is not allowed; only either debit or credit can be not zero")
		}
	}

	// date should be filled out and parseable to something ok

	return nil
}

func checkAccountsBelongToCompanyId(tx *sql.Tx, transaction data.Transaction) error {
	var accountIds []string
	pos := transaction.Positions
	for _, val := range pos {
		accountId := val.AccountId

		accountIds = append(accountIds, accountId)
	}

	cids, err := dao.GetCompanyIdsForAccountIds(tx, accountIds)
	if err != nil {
		return err
	}

	if len(cids) != 1 || cids[0] != transaction.CompanyId {
		return data.CreateIdError(false, "somehow not all the accounts in the transaction belong to the posted company")
	}

	return nil
}

func checkTransactionBelongsToUser(tx *sql.Tx, adding bool, transaction data.Transaction, userId string) error {
	companyId := transaction.CompanyId
	if !adding {
		providedCompanyId, err := dao.GetTransactionCompanyId(tx, transaction.TransactionId)
		if err != nil {
			return err
		}

		if providedCompanyId != companyId {
			return data.CreateIdError(false, "the company id provided by transaction and the one in the database are different")
		}
	}

	userIds, err := dao.GetUsersForCompanyIds(tx, []string{companyId})
	if err != nil {
		return err
	}

	if len(userIds) != 1 || userIds[0] != userId {
		return data.CreateIdError(false, "the current company does not belong to the logged in user")
	}

	return nil
}
