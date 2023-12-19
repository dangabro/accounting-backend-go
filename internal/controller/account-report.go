package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type accountReport struct {
	config config.Config
	db     *sql.DB
}

func AccountReport(config config.Config, db *sql.DB) http.Handler {
	return &accountReport{
		config: config,
		db:     db,
	}
}

func (h *accountReport) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("accountReport")
	success, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, success)
	}
}

func (h *accountReport) process(r *http.Request) (data.AccountReportResult, error) {
	var res data.AccountReportResult

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	accountReportRequest, err := readAccountReportRequest(r)
	if err != nil {
		return res, err
	}

	userId := user.UserId
	start := accountReportRequest.StartDate
	end := accountReportRequest.EndDate
	accountId := accountReportRequest.AccountId

	res, err = serv.AccountReport(start, end, accountId, userId)
	if err != nil {
		return res, err
	}

	return res, err
}

func readAccountReportRequest(r *http.Request) (data.AccountReportRequest, error) {
	var res data.AccountReportRequest
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}
