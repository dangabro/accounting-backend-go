package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type updateAccount struct {
	config config.Config
	db     *sql.DB
}

func UpdateAccount(config config.Config, db *sql.DB) http.Handler {
	return &updateAccount{
		config: config,
		db:     db,
	}
}

func (h *updateAccount) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("update account")
	success, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, success)
	}
}

func (h *updateAccount) process(r *http.Request) (data.SuccessData, error) {
	res := data.GetSuccessData()

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	updateAccountRequest, err := readUpdateAccountData(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	userId := user.UserId
	adding := updateAccountRequest.Adding
	account := updateAccountRequest.Account

	err = serv.UpdateAccount(adding, account, userId)

	return res, err
}
