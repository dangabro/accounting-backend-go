package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"net/http"
)

type companies struct {
	config config.Config
	db     *sql.DB
}

func Companies(config config.Config, db *sql.DB) http.Handler {
	return &companies{
		config: config,
		db:     db,
	}
}

func (h *companies) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	success, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, success)
	}
}

func (h *companies) process(r *http.Request) ([]data.CompanyData, error) {
	var res []data.CompanyData

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	userId := user.UserId

	res, err = serv.GetAllCompanies(userId)

	return res, err
}
