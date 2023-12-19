package controller

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/rs/zerolog/log"
	"net/http"
)

type root struct {
	config config.Config
	db     *sql.DB
}

func Root(config config.Config, db *sql.DB) http.Handler {
	return &root{
		config: config,
		db:     db,
	}
}

func (h *root) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("ROOT")
	rootHealthData := data.GetRootHealthData(h.db)
	_ = writeJsonResponse(writer, rootHealthData)
}
