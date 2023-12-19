package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type transactionReport struct {
	config config.Config
	db     *sql.DB
}

func TransactionReport(config config.Config, db *sql.DB) http.Handler {
	return &transactionReport{
		config: config,
		db:     db,
	}
}

func (h *transactionReport) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("transaction report")
	txnReportResult, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		// wrap this in a javascript object under the name transactions

		_ = writeJsonResponse(writer, txnReportResult)
	}
}

func (h *transactionReport) process(r *http.Request) (data.TransactionReportResult, error) {
	var res data.TransactionReportResult

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	transactionReportRequest, err := readTransactionReportRequest(r)
	if err != nil {
		return res, err
	}

	start := transactionReportRequest.Start
	end := transactionReportRequest.End
	companyId := transactionReportRequest.CompanyId

	res, err = serv.TransactionReport(start, end, companyId, user.UserId)

	return res, err
}
