package controller

import (
	"database/sql"
	"fmt"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type logout struct {
	config config.Config
	db     *sql.DB
}

func Logout(config config.Config, db *sql.DB) http.Handler {
	return &logout{
		config: config,
		db:     db,
	}
}

func (h *logout) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("logout")
	loginRes, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		_ = writeJsonResponse(writer, loginRes)
	}
}

func (h *logout) process(r *http.Request) (data.SuccessData, error) {
	res := data.GetSuccessData()

	// we close the session if it is valid
	token, err := getToken(r)
	if err == nil {
		serv := service.New(h.db, h.config)
		err = serv.Logout(token)
		if err != nil {
			errMessage := err.Error()
			message := fmt.Sprintf("error closing the session: %s", errMessage)

			log.Error().Msg(message)
		}
	} else {
		errMessage := err.Error()
		message := fmt.Sprintf("error retrieving the token: %s", errMessage)

		log.Error().Msg(message)
	}

	return res, nil
}
