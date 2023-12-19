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

type login struct {
	config config.Config
	db     *sql.DB
}

func Login(config config.Config, db *sql.DB) http.Handler {
	return &login{
		config: config,
		db:     db,
	}
}

func (h *login) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("login")
	loginRes, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, loginRes)
	}
}

func (h *login) process(r *http.Request) (data.LoginResult, error) {
	res := data.LoginResult{}
	loginData, err := readLoginData(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	res, err = serv.Login(loginData)
	return res, err
}

func readLoginData(r *http.Request) (data.LoginData, error) {
	var res data.LoginData
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}
