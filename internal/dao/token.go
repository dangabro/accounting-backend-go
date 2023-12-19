package dao

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/data"
	"strings"
	"time"
)

func GetTokenEntry(tx *sql.Tx, token string) (data.Session, error) {
	res := data.Session{}
	sql := "select session_id, user_id, token, expired_ind, expiry_dt from session where token = ?"
	rows, err := tx.Query(sql, token)
	if err != nil {
		return res, err
	}
	defer closeRows(rows)

	if !rows.Next() {
		return res, data.CreateIdError(false, "cannot find the token with the value you provided")
	}

	expiredInd := "N"
	err = rows.Scan(&res.SessionId, &res.UserId, &res.Token, &expiredInd, &res.ExpiryDt)
	if err != nil {
		return res, err
	}

	res.Expired = expiredInd == "Y"

	return res, nil
}

func ExpireSession(tx *sql.Tx, sessionId string) error {
	sql := "update session set expired_ind = 'Y', expiry_dt = ? where session_id = ?"

	newDate := time.Now()
	_, err := tx.Exec(sql, newDate, sessionId)

	return err
}

func GetSessionByToken(tx *sql.Tx, token string) (data.Session, error) {
	res := data.Session{}

	sql := "select session_id, user_id, token, expired_ind, expiry_dt from session where token = ?"
	rows, err := tx.Query(sql, token)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)

	if rows.Next() {
		var sessionId string
		var userId string
		var loadedToken string
		var expiredInd string
		var expiryDate time.Time

		err = rows.Scan(&sessionId, &userId, &loadedToken, &expiredInd, &expiryDate)
		if err != nil {
			return res, err
		}

		// transform expired
		expiredInd = strings.TrimSpace(expiredInd)
		expiredInd = strings.ToUpper(expiredInd)

		expired := expiredInd == "Y"

		res = data.Session{
			SessionId: sessionId,
			UserId:    userId,
			Token:     loadedToken,
			Expired:   expired,
			ExpiryDt:  expiryDate,
		}
	} else {
		return res, data.CreateIdError(false, "cannot find session with provided token")
	}

	return res, nil
}
