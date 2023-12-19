package service

import (
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/util"
)

func (d *servc) SearchTransactions(search string, companyId string, userId string) ([]data.Transaction, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, nil
	}
	defer tx.Rollback()

	userIds, err := dao.GetUsersForCompanyIds(tx, []string{companyId})
	if len(userIds) != 1 || userIds[0] != userId {
		return nil, data.CreateIdError(false, "the provided company does not belong to the logged in user")
	}

	dt, err := util.ParseSearchString(search)
	if err != nil {
		return nil, err
	}

	transactions, err := dao.SearchTransaction(tx, dt, companyId)
	if err != nil {
		return nil, err
	}

	// for returning txns, format them with the display format
	tx.Commit()

	return transactions, nil
}
