package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type searchTransactions struct {
	config config.Config
	db     *sql.DB
}

func SearchTransactions(config config.Config, db *sql.DB) http.Handler {
	return &searchTransactions{
		config: config,
		db:     db,
	}
}

func (h *searchTransactions) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("search transactions")
	transactions, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		// wrap this in a javascript object under the name transactions
		txnMap := make(map[string][]data.Transaction)
		txnMap["transactions"] = transactions

		_ = writeJsonResponse(writer, txnMap)
	}
}

func (h *searchTransactions) process(r *http.Request) ([]data.Transaction, error) {
	var res []data.Transaction

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	transactionSearch, err := readSearchTransactions(r)
	if err != nil {
		return res, err
	}

	search := transactionSearch.Search
	companyId := transactionSearch.CompanyId
	userId := user.UserId

	res, err = serv.SearchTransactions(search, companyId, userId)

	return res, err
}
