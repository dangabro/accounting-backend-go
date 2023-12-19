package dao

import (
	"database/sql"
	"fmt"
	"github.com/dgb9/db-account-server/internal/data"
)

func UserExists(tx *sql.Tx, login string) (bool, error) {
	sql := "select login from user where login = ?"
	rows, err := tx.Query(sql, login)
	if err != nil {
		return false, err
	}

	defer closeRows(rows)
	res := rows.Next()

	return res, nil
}

func CreateUser(tx *sql.Tx, user data.User) error {
	sql := "insert into user (user_id, provided_id, login, name) values (?, ?, ?, ?)"
	userId := user.UserId
	providedId := user.ProvidedId
	login := user.Login
	name := user.Name

	_, err := tx.Exec(sql, userId, providedId, login, name)

	return err
}

func GetUserById(tx *sql.Tx, userId string) (user data.User, err error) {
	res := data.User{}
	sql := "select user_id, provided_id, login, name from user where user_id = ?"
	rows, err := tx.Query(sql, userId)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)
	if !rows.Next() {
		message := fmt.Sprintf("can't find the user with the user id %s", userId)
		return res, data.CreateIdError(false, message)
	}

	err = rows.Scan(&res.UserId, &res.ProvidedId, &res.Login, &res.Name)
	return res, err
}

func GetUserByLogin(tx *sql.Tx, login string) (user data.User, err error) {
	res := data.User{}
	sql := "select user_id, provided_id, login, name from user where login = ?"
	rows, err := tx.Query(sql, login)
	if err != nil {
		return res, err
	}

	defer closeRows(rows)
	if rows.Next() {
		err = rows.Scan(&res.UserId, &res.ProvidedId, &res.Login, &res.Name)

		return res, err
	} else {
		message := fmt.Sprintf("cannot find user with login: %s", login)
		return res, data.CreateIdError(false, message)
	}

	return res, nil
}

func AddSession(tx *sql.Tx, session data.Session) error {
	sql := "insert into session (session_id, user_id, token, expired_ind, expiry_dt) values (?, ?, ?, ?, ?)"
	expired := "N"

	_, err := tx.Exec(sql, session.SessionId, session.UserId, session.Token, expired, session.ExpiryDt)
	return err
}
