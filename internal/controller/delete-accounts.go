package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"net/http"
)

type deleteAccounts struct {
	config config.Config
	db     *sql.DB
}

func DeleteAccounts(config config.Config, db *sql.DB) http.Handler {
	return &deleteAccounts{
		config: config,
		db:     db,
	}
}

func (h *deleteAccounts) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	success, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, success)
	}
}

func (h *deleteAccounts) process(r *http.Request) (data.SuccessData, error) {
	res := data.GetSuccessData()

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	idsRequest, err := readIdsCollection(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	userId := user.UserId

	err = serv.DeleteAccounts(idsRequest.Ids, userId)

	return res, err
}
