package sheet

import (
	"bytes"
	"fmt"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/sheet/style"
	"github.com/xuri/excelize/v2"
	"regexp"
	"strings"
	"time"
)

func Proc(dt data.ExcelReport) (*bytes.Buffer, string, error) {
	f := excelize.NewFile()

	var err error

	styles := style.NewStyle(f)

	err = createCompanySheet(f, styles, dt.Company)
	if err != nil {
		return nil, "", err
	}

	err = deleteDefaultSheet(f)
	if err != nil {
		return nil, "", err
	}

	typesMap := NewAcctTypeMap(dt.AccountTypes)

	err = createAccountPlan(f, styles, dt.Accounts, typesMap)
	if err != nil {
		return nil, "", err
	}

	err = createTrialBalance(f, styles, dt.Balance, typesMap)
	if err != nil {
		return nil, "", err
	}

	err = createTransactionReport(f, styles, dt.Start, dt.End, dt.Accounts, dt.Transactions)
	if err != nil {
		return nil, "", err
	}

	err = createAccountReports(f, styles, dt.Start, dt.End, typesMap, dt.AccountsReport)
	if err != nil {
		return nil, "", err
	}

	fileName := getFileName(dt.Start, dt.End, dt.Company)
	buffer, err := f.WriteToBuffer()
	return buffer, fileName, err
}

func getFileName(start string, end string, company data.CompanyData) string {
	// create file name based on the entries
	companyName := company.Name
	currentDate := time.Now().Format("2006-01-02 150405")
	strEnd := strings.TrimSpace(end)
	if len(strEnd) == 0 {
		strEnd = "onwards"
	}

	fileName := fmt.Sprintf("%s_from_%s_to_%s_on_%s.xlsx", companyName, start, strEnd, currentDate)

	re := regexp.MustCompile("\\s+")
	fileName = re.ReplaceAllString(fileName, "_")

	return fileName
}

func deleteDefaultSheet(f *excelize.File) error {
	err := f.DeleteSheet("Sheet1")
	return err
}
