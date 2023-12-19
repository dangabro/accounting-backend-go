package service

import (
	"fmt"
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
)

func (d *servc) UpdateCompany(adding bool, company data.CompanyData, userId string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer rollbackTx(tx)

	// check company filled out etc
	err = checkCompany(company)
	if err != nil {
		return err
	}

	if adding {
		companyUserId := company.UserId
		if companyUserId != userId {
			return data.CreateIdError(false, "the provided company and the token user id are different")
		}

		err = dao.InsertCompany(tx, company)
	} else {
		comp, err1 := dao.GetCompany(tx, company.CompanyId)
		if err1 != nil {
			return err1
		}

		// the company exists; then verify the user is the same
		companyUserId := comp.UserId
		if companyUserId != userId {
			return data.CreateIdError(false, "the company you try to update does not belong to you")
		}

		err1 = dao.UpdateCompany(tx, company)
		if err1 != nil {
			return err1
		}
	}

	tx.Commit()

	return nil
}

func (d *servc) GetAllCompanies(userId string) ([]data.CompanyData, error) {
	var res []data.CompanyData

	tx, err := d.db.Begin()
	if err != nil {
		return res, err
	}
	defer tx.Rollback()

	res, err = dao.GetAllCompanies(tx, userId)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return res, err
}

func (d *servc) DeleteCompanies(ids []string, userId string) error {
	tx, err := d.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	users, err := dao.GetUsersForCompanyIds(tx, ids)
	if err != nil {
		return err
	}

	if len(users) != 1 || users[0] != userId {
		return data.CreateIdError(false, "the companies do not belong to the logged in user")
	}

	// see there aren't accounts linked to the company
	count, err := dao.GetCountAccountForCompanyIds(tx, ids)
	if err != nil {
		return err
	}

	if count > 0 {
		message := fmt.Sprintf("there are %d accounts depending on the companies you want to delete", count)
		return data.CreateIdError(false, message)
	}

	// all right, proceed with the deletion
	err = dao.DeleteCompanies(tx, ids)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
