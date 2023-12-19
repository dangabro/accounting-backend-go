package service

import (
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
)

func (s *servc) DeleteTransactions(ids []string, userId string) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	// see the company belongs to our user
	users, err := dao.GetUsersForTransactionIds(tx, ids)
	if err != nil {
		return err
	}

	if len(users) != 1 || users[0] != userId {
		return data.CreateIdError(false, "provided transactions do not all belong to the logged in user")
	}

	// we are good now
	err = dao.DeletePositionsForTransactions(tx, ids)
	if err != nil {
		return err
	}
	err = dao.DeleteTransactions(tx, ids)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
