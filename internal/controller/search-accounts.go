package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type searchAccounts struct {
	config config.Config
	db     *sql.DB
}

func SearchAccounts(config config.Config, db *sql.DB) http.Handler {
	return &searchAccounts{
		config: config,
		db:     db,
	}
}

func (h *searchAccounts) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("search accounts")
	accounts, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		resMap := make(map[string][]data.Account)
		resMap["accounts"] = accounts
		_ = writeJsonResponse(writer, resMap)
	}
}

func (h *searchAccounts) process(r *http.Request) ([]data.Account, error) {
	var res []data.Account

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	companySearch, err := readAccountsSearch(r)
	if err != nil {
		return res, err
	}

	search := companySearch.Search
	companyId := companySearch.CompanyId

	userId := user.UserId

	res, err = serv.SearchAccounts(search, companyId, userId)

	return res, err
}
