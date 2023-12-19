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

type trialBalanceReport struct {
	config config.Config
	db     *sql.DB
}

func TrialBalanceReport(config config.Config, db *sql.DB) http.Handler {
	return &trialBalanceReport{
		config: config,
		db:     db,
	}
}

func (h *trialBalanceReport) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("trial balance report")
	success, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, success)
	}
}

func (h *trialBalanceReport) process(r *http.Request) (data.BalanceResult, error) {
	var res data.BalanceResult

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	trialBalanceRequest, err := readTrialBalanceRequest(r)
	if err != nil {
		return res, err
	}

	userId := user.UserId
	start := trialBalanceRequest.StartDate
	end := trialBalanceRequest.EndDate
	companyId := trialBalanceRequest.CompanyId

	res, err = serv.TrialBalance(start, end, companyId, userId)
	if err != nil {
		return res, err
	}

	return res, err
}

func readTrialBalanceRequest(r *http.Request) (data.TrialBalanceRequest, error) {
	var res data.TrialBalanceRequest
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}
