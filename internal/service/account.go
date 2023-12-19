package service

import (
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
)

func (s *servc) DeleteAccounts(ids []string, userId string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// get the userIds associated with the accountIds
	users, err := dao.GetUserIdsForAccountIds(tx, ids)
	if err != nil {
		return err
	}

	// check if userId is the only user that owns those account ids
	if len(users) != 1 {
		return data.CreateIdError(false, "Accounts do not belong (only) to the logged in user")
	} else {
		onlyUser := users[0]

		if onlyUser != userId {
			return data.CreateIdError(false, "Accounts do not belong to the logged in user")
		}
	}

	// search for transactions that have accounts
	exist, err := dao.HasTransactionsForAccountIds(tx, ids)
	if err != nil {
		return err
	}

	if exist {
		return data.CreateIdError(false, "There are still transactions that use those accounts")
	}

	// if not found, delete accounts
	err = dao.DeleteAccounts(tx, ids)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (s *servc) GetAllAccounts(companyId string, userId string) ([]data.Account, error) {
	res := make([]data.Account, 0)
	tx, err := s.db.Begin()

	if err != nil {
		return res, err
	}
	defer tx.Rollback()

	// see the company belongs to our user
	ids := []string{companyId}
	users, err := dao.GetUsersForCompanyIds(tx, ids)
	if err != nil {
		return res, err
	}

	if len(users) != 1 || users[0] != userId {
		return res, data.CreateIdError(false, "provided company does not belong to the logged in user")
	}

	// we are good now
	res, err = dao.GetAccountsByCompanyId(tx, companyId)
	if err != nil {
		return res, err
	}

	tx.Commit()

	return res, nil
}

func (s *servc) SearchAccounts(search string, companyId string, userId string) ([]data.Account, error) {
	res := make([]data.Account, 0)
	tx, err := s.db.Begin()
	if err != nil {
		return res, err
	}

	defer tx.Rollback()

	users, err := dao.GetUsersForCompanyIds(tx, []string{companyId})
	if err != nil {
		return res, err
	}

	if len(users) != 1 || users[0] != userId {
		return res, data.CreateIdError(false, "the company you search for does not belong to the logged in user")
	}

	res, err = dao.SearchAccounts(tx, search, companyId)

	tx.Commit()

	return res, nil
}

func (s *servc) UpdateAccount(adding bool, account data.Account, userId string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// check the acount
	err = checkAccount(account)
	if err != nil {
		return err
	}

	users, err := dao.GetUsersForCompanyIds(tx, []string{account.CompanyId})
	if err != nil {
		return err
	}

	if len(users) != 1 || users[0] != userId {
		return data.CreateIdError(false, "The provided account does not belong to the logged in user")
	}

	// get the user id associated with the account
	if adding {
		var exists bool
		exists, err = dao.CheckCompanyExists(tx, account.CompanyId)
		if err != nil {
			return err
		}

		if !exists {
			return data.CreateIdError(false, "cannot find the company with this id")
		}
		err = dao.AddAccount(tx, account)
	} else {
		var exists bool
		exists, err = dao.EnsureAccountExists(tx, account.AccountId)
		if err != nil {
			return err
		}

		if exists {
			err = dao.UpdateAccount(tx, account)
		} else {
			return data.CreateIdError(false, "cannot find the account with this id")
		}
	}

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
