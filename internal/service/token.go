package service

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
	"time"
)

func (d *servc) ValidateToken(token string) (data.User, error) {
	res := data.User{}
	tx, err := d.db.Begin()
	if err != nil {
		return res, err
	}

	defer tx.Rollback()

	// find the token in the database
	// if not found, error
	session, err := dao.GetTokenEntry(tx, token)
	if err != nil {
		return res, err
	}

	// if found, see if it is expired
	if session.Expired {
		return res, data.CreateIdError(false, "the token you provided is marked expired")
	}

	// it is not marked expired, but lets' check the time
	now := time.Now()
	expired := session.ExpiryDt
	if now.After(expired) {
		return res, data.CreateIdError(false, "current timestamp past the token timestamp")
	}

	// if not expired, find the user to whom it belongs
	userId := session.UserId

	// load the user and return it
	res, err = dao.GetUserById(tx, userId)
	if err != nil {
		return res, err
	}

	// and now move the session updating the expiry date
	secondsToken := d.config.SecondsToken()
	id := session.SessionId

	// calculate new valid token expiry dt
	timeValid := time.Now().Add(time.Second * time.Duration(secondsToken))
	err = moveSession(tx, id, timeValid)
	if err != nil {
		return res, err
	}

	err = tx.Commit()

	return res, err
}

func moveSession(tx *sql.Tx, id string, newValidTime time.Time) error {
	sql := "update session set expiry_dt = ? where session_id = ?"

	_, err := tx.Exec(sql, newValidTime, id)
	return err
}
