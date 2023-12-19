package service

import (
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
)

func (d *servc) GetAccountTypeCds() ([]data.AccountTypeData, error) {
	res := make([]data.AccountTypeData, 0)
	tx, err := d.db.Begin()
	if err != nil {
		return res, nil
	}
	defer tx.Rollback()

	res, err = dao.GetAccountTypeCds(tx)
	if err != nil {
		return res, err
	}

	_ = tx.Commit()

	return res, err
}
