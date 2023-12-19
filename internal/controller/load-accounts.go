package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

type loadAccounts struct {
	config config.Config
	db     *sql.DB
}

func LoadAccounts(config config.Config, db *sql.DB) http.Handler {
	return &loadAccounts{
		config: config,
		db:     db,
	}
}

func (h *loadAccounts) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("load accounts")
	res, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		acctMap := make(map[string][]data.Account)
		acctMap["accounts"] = res

		_ = writeJsonResponse(writer, acctMap)
	}
}

func (h *loadAccounts) process(r *http.Request) ([]data.Account, error) {
	res := make([]data.Account, 0)

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	params := mux.Vars(r)
	companyId := params["companyId"]

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	userId := user.UserId

	res, err = serv.GetAllAccounts(companyId, userId)

	return res, err
}
