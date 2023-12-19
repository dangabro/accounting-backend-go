package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgb9/db-account-server/internal/dao"
	"github.com/dgb9/db-account-server/internal/data"
	"io/ioutil"
	"net/http"
	"time"
)

func (d *servc) Login(lg data.LoginData) (data.LoginResult, error) {
	// call over internet to get logged in
	res := data.LoginResult{}
	userSSO, err := d.loginSSO(lg)

	// if not logged in, return error
	if err != nil {
		return res, err
	}

	// if logged in, see if the right is there, if not, error
	rights := userSSO.Rights
	login := userSSO.Login
	_, accessRight := d.config.Access()
	found := false
	for _, right := range rights {
		if right == accessRight {
			found = true
			break
		}
	}
	if !found {
		message := fmt.Sprintf("user %s does not have the require right %s", login, accessRight)

		return res, data.CreateIdError(false, message)
	}

	// see if the user exists, if not, create it
	tx, err := d.db.Begin()
	if err != nil {
		return res, err
	}

	defer tx.Rollback()

	exists, err := dao.UserExists(tx, login)
	if err != nil {
		return res, err
	}

	var user data.User
	if !exists {
		user = data.User{
			UserId:     GetNewUUID(),
			ProvidedId: userSSO.Id,
			Login:      login,
			Name:       userSSO.Name,
		}

		err = dao.CreateUser(tx, user)
		if err != nil {
			return res, nil
		}
	} else {
		user, err = dao.GetUserByLogin(tx, login)
		if err != nil {
			return res, err
		}
	}

	userId := user.UserId

	// now create the token
	token := GetNewUUID()
	sessionId := GetNewUUID()

	tokenValidSeconds := d.config.SecondsToken()
	timeValid := time.Now().Add(time.Second * time.Duration(tokenValidSeconds))

	session := data.Session{
		SessionId: sessionId,
		UserId:    userId,
		Token:     token,
		Expired:   false,
		ExpiryDt:  timeValid,
	}

	err = dao.AddSession(tx, session)
	if err != nil {
		return res, err
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return res, err
	}

	// prepare result
	res = data.LoginResult{
		Login: login,
		Token: token,
		Id:    userId,
		Name:  userSSO.Name,
	}

	return res, nil
}

func (d *servc) loginSSO(login data.LoginData) (data.UserSSO, error) {
	res := data.UserSSO{}
	config := d.config
	url, _ := config.Access()

	// getting the bytes from the url
	bt, err := json.Marshal(login)
	if err != nil {
		return res, err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(bt))
	if err != nil {
		return res, err
	}

	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	result, err := client.Do(r)
	if err != nil {
		return res, err
	}

	defer result.Body.Close() //

	// parse the result to data.UserSSO
	if result.StatusCode == http.StatusOK {
		err = json.NewDecoder(result.Body).Decode(&res)
	} else {
		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			return res, err
		}

		err = data.CreateIdError(false, string(body))
		return res, err
	}

	return res, err
}
