package controller

import (
	"encoding/json"
	"github.com/dgb9/db-account-server/internal/data"
	"net/http"
)

func writeJsonResponse(writer http.ResponseWriter, payload any) error {
	status := http.StatusOK
	processPayload := payload

	currentError, ok := payload.(data.IdError)
	if ok {
		processPayload = data.PayloadError{Message: currentError.Error()}

		if currentError.IsSystem() {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
	} else {
		resError, ok := payload.(error)
		if ok {
			processPayload = data.PayloadError{Message: resError.Error()}
			status = http.StatusInternalServerError
		}
	}

	header := writer.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Cache-Control", "no-cache")

	writer.WriteHeader(status)

	// write the actual object here
	bts, err := json.Marshal(processPayload)
	if err != nil {
		return err
	}

	_, err = writer.Write(bts)
	if err != nil {
		return err
	}

	return nil
}

func readIdsCollection(r *http.Request) (data.IdsCollection, error) {
	var res data.IdsCollection
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}

func readAccountsSearch(r *http.Request) (data.AccountSearch, error) {
	var res data.AccountSearch
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}

func readUpdateAccountData(r *http.Request) (data.UpdateAccountRequest, error) {
	var res data.UpdateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}

func readSearchTransactions(r *http.Request) (data.SearchTransactionsRequest, error) {
	var res data.SearchTransactionsRequest
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}

func readUpdateTransaction(r *http.Request) (data.UpdateTransactionRequest, error) {
	var res data.UpdateTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}

func readTransactionReportRequest(r *http.Request) (data.TransactionReportRequest, error) {
	var res data.TransactionReportRequest
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}
