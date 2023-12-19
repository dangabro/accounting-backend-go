package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

type excelReport struct {
	config config.Config
	db     *sql.DB
}

func ExcelReport(config config.Config, db *sql.DB) http.Handler {
	return &excelReport{
		config: config,
		db:     db,
	}
}

func (h *excelReport) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	log.Info().Msg("trial balance report")
	res, err := h.process(r)

	if err != nil {
		_ = writeJsonResponse(writer, err)
	} else {
		fileName := res.FileName
		disposition := fmt.Sprintf("attachment; filename=%s", fileName)
		writer.Header().Set("Content-Disposition", disposition)
		_, _ = res.DataBuffer.WriteTo(writer)
	}
}

func (h *excelReport) process(r *http.Request) (data.ExcelResponse, error) {
	res := data.ExcelResponse{}
	token, err := getToken(r)
	if err != nil {
		return res, err
	}

	serv := service.New(h.db, h.config)
	user, err := serv.ValidateToken(token)
	if err != nil {
		return res, err
	}

	excelReport, err := readExcelReport(r)
	if err != nil {
		return res, err
	}

	userId := user.UserId
	companyId := excelReport.CompanyId
	end := excelReport.EndDate
	start := excelReport.StartDate
	excelResponse, err := serv.ExcelReport(start, end, companyId, userId)

	if err != nil {
		return excelResponse, err
	}

	return excelResponse, nil
}

func readExcelReport(r *http.Request) (data.ExcelReportRequest, error) {
	var res data.ExcelReportRequest
	err := json.NewDecoder(r.Body).Decode(&res)

	return res, err
}
