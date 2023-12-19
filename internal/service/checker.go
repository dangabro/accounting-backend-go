package service

import "github.com/dgb9/db-account-server/internal/data"

func checkRange(str string, minLength int, maxLength int, msg string) error {
	ln := len(str)
	if ln < minLength || ln > maxLength {
		return data.CreateIdError(false, msg)
	}

	return nil
}

func checkAccount(account data.Account) error {
	accountId := account.AccountId
	err := checkRange(accountId, 1, 64, "account id is mandatory and must be between 1 and 64 characters long")
	if err != nil {
		return err
	}

	code := account.Code
	err = checkRange(code, 1, 64, "code is mandatory and must be between 1 and 64 characters long")
	if err != nil {
		return err
	}

	name := account.Name
	err = checkRange(name, 1, 255, "code is mandatory and must be between 1 and 255 characters long")
	if err != nil {
		return err
	}

	accountType := account.AccountTypeCd
	err = checkRange(accountType, 1, 64, "accountType is mandatory and must be between 1 and 64 characters long")
	if err != nil {
		return err
	}

	companyId := account.CompanyId
	err = checkRange(companyId, 1, 64, "company ID is mandatory and must be between 1 and 64 characters long")
	if err != nil {
		return err
	}

	return nil
}

func checkTxnRecord(txn data.Transaction) error {
	id := txn.TransactionId
	err := checkRange(id, 1, 64, "txn id mandatory [1-64] characters")
	if err != nil {
		return err
	}

	cid := txn.CompanyId
	err = checkRange(cid, 1, 64, "company id mandatory [1-64] characters")
	if err != nil {
		return err
	}

	date := txn.TransactionDate
	err = checkRange(date, 1, 32, "date is mandatory: 1 to 32 characters")
	if err != nil {
		return err
	}

	return nil
}

func checkTransactionPosition(p data.TransactionPosition) error {
	// for now all the checks are done by other modules
	return nil
}

func checkTransaction(txn data.Transaction) error {
	err := checkTxnRecord(txn)
	if err != nil {
		return err
	}

	pos := txn.Positions
	for _, p := range pos {
		err = checkTransactionPosition(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkCompany(p data.CompanyData) error {
	id := p.CompanyId
	userId := p.UserId
	companyName := p.Name

	err := checkRange(id, 1, 64, "company id is mandatory [1-64] length")
	if err != nil {
		return err
	}

	err = checkRange(userId, 1, 64, "user id is mandatory [1-64] length")
	if err != nil {
		return err
	}

	err = checkRange(companyName, 1, 255, "company name is mandatory [1-255] length")
	if err != nil {
		return err
	}

	return nil
}
