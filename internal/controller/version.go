package controller

import (
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/util"
	"net/http"
)

type version struct {
}

func Version() http.Handler {
	return &version{}
}

func (h *version) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	res := data.VersionData{Version: util.Version}

	_ = writeJsonResponse(writer, res)
}
