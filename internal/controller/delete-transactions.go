package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type deleteTransactions struct {
	config config.Config
	db     *sql.DB
}

func DeleteTransactions(config config.Config, db *sql.DB) http.Handler {
	return &deleteTransactions{
		config: config,
		db:     db,
	}
}

func (h *deleteTransactions) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("delete transactions")
	success, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, success)
	}
}

func (h *deleteTransactions) process(r *http.Request) ([]data.SuccessData, error) {
	var res []data.SuccessData

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	idsRequest, err := readIdsCollection(r)
	if err != nil {
		return res, err
	}

	userId := user.UserId

	err = serv.DeleteTransactions(idsRequest.Ids, userId)

	return res, err
}
