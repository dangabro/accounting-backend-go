package service

import "github.com/dgb9/db-account-server/internal/dao"

func (h *servc) Logout(token string) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer rollbackTx(tx)

	// see if the session exists and load it
	session, err := dao.GetSessionByToken(tx, token)
	if err != nil {
		return err
	}

	// if not already canceled, cancel it
	if !session.Expired {
		err = dao.ExpireSession(tx, session.SessionId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}
