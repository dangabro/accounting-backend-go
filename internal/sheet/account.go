package sheet

import (
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/sheet/style"
	"github.com/xuri/excelize/v2"
)

func createAccountPlan(f *excelize.File, styles style.Style, accounts []data.Account, types AcctTypeMap) error {
	sheet := "Accounts"
	_, _ = f.NewSheet(sheet)

	_ = f.SetCellValue(sheet, "A1", "Accounts")
	_ = f.MergeCell(sheet, "A1", "C1")
	_ = f.SetCellStyle(sheet, "A1", "A1", styles.GetStyle(style.Title))

	_ = f.SetCellValue(sheet, "A2", "Code")
	_ = f.SetCellValue(sheet, "B2", "Name")
	_ = f.SetCellValue(sheet, "C2", "Account Type")
	_ = f.SetRowHeight(sheet, 1, 30)
	_ = f.SetRowHeight(sheet, 2, 30)
	_ = f.SetRowStyle(sheet, 2, 2, styles.GetStyle(style.Bold))

	for index, account := range accounts {
		row := 2 + index + 1

		code := account.Code
		name := account.Name
		accountType := types.GetAccountTypeName(account.AccountTypeCd)

		_ = f.SetCellValue(sheet, getCell("A", row), code)
		_ = f.SetCellValue(sheet, getCell("B", row), name)
		_ = f.SetCellValue(sheet, getCell("C", row), accountType)
	}

	return nil
}
