package controller

import (
	"github.com/dgb9/db-account-server/internal/data"
	"net/http"
)

func getToken(r *http.Request) (string, error) {
	res := ""

	// get the token from the request
	token := r.Header.Get("Authorization")
	if token == "" {
		return res, data.CreateIdError(false, "cannot find token in the payload")
	}

	return token, nil
}
