package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type updateTransaction struct {
	config config.Config
	db     *sql.DB
}

func UpdateTransaction(config config.Config, db *sql.DB) http.Handler {
	return &updateTransaction{
		config: config,
		db:     db,
	}
}

func (h *updateTransaction) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("update transaction")
	success, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, success)
	}
}

func (h *updateTransaction) process(r *http.Request) (data.SuccessData, error) {
	res := data.SuccessData{Success: false}

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	updateTransactionRequest, err := readUpdateTransaction(r)
	if err != nil {
		return res, err
	}

	userId := user.UserId
	adding := updateTransactionRequest.Adding
	transaction := updateTransactionRequest.Transaction

	err = serv.UpdateTransaction(adding, transaction, userId)
	if err != nil {
		return res, err
	}

	res.Success = true
	return res, err
}
