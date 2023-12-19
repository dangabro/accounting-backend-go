package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"net/http"
)

type company struct {
	config config.Config
	db     *sql.DB
}

func Company(config config.Config, db *sql.DB) http.Handler {
	return &company{
		config: config,
		db:     db,
	}
}

func (h *company) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	success, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, success)
	}
}

func (h *company) process(r *http.Request) (data.SuccessData, error) {
	res := data.GetSuccessData()

	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	companyRequest, err := readCompanyData(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	userId := user.UserId

	adding := companyRequest.Adding
	company := companyRequest.Company

	err = serv.UpdateCompany(adding, company, userId)

	return res, err
}

func readCompanyData(r *http.Request) (data.UpdateCompanyRequest, error) {
	var res data.UpdateCompanyRequest
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}
