package sheet

import (
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/sheet/style"
	"github.com/xuri/excelize/v2"
	"strings"
)

func createTrialBalance(f *excelize.File, st style.Style, balance data.BalanceResult, typeMap AcctTypeMap) error {
	sheet := "Trial Balance"
	_, err := f.NewSheet(sheet)
	if err != nil {
		return err
	}

	_ = f.MergeCell(sheet, "A1", "I1")
	_ = f.MergeCell(sheet, "A2", "I2")
	_ = f.MergeCell(sheet, "A3", "I3")
	_ = f.SetRowHeight(sheet, 1, 30)
	_ = f.SetRowHeight(sheet, 2, 20)
	_ = f.SetRowHeight(sheet, 3, 20)
	_ = f.SetRowStyle(sheet, 2, 3, st.GetStyle(style.VerticalCenter))

	_ = f.SetCellValue(sheet, "A1", "Trial Balance")
	_ = f.SetCellStyle(sheet, "A1", "A1", st.GetStyle(style.Title))
	_ = f.SetCellValue(sheet, "A2", "Starting with: "+balance.StartDate)
	endDate := balance.EndDate
	trimmedEndDate := strings.TrimSpace(endDate)
	lengthEndDate := len(trimmedEndDate)

	valEnd := ""
	if lengthEndDate > 0 {
		valEnd = "Ending with: " + balance.EndDate
	}
	_ = f.SetCellValue(sheet, "A3", valEnd)

	// titles
	_ = f.MergeCell(sheet, "A5", "C5")
	_ = f.MergeCell(sheet, "D5", "E5")
	_ = f.MergeCell(sheet, "F5", "G5")
	_ = f.MergeCell(sheet, "H5", "I5")

	_ = f.SetCellValue(sheet, "A5", "Account")
	_ = f.SetCellValue(sheet, "D5", "Start")
	_ = f.SetCellValue(sheet, "F5", "Runs")
	_ = f.SetCellValue(sheet, "H5", "End")
	_ = f.SetCellValue(sheet, "A6", "Code")
	_ = f.SetCellValue(sheet, "B6", "Name")
	_ = f.SetCellValue(sheet, "C6", "Type")
	_ = f.SetCellValue(sheet, "D6", "Debit")
	_ = f.SetCellValue(sheet, "E6", "Credit")
	_ = f.SetCellValue(sheet, "F6", "Debit")
	_ = f.SetCellValue(sheet, "G6", "Credit")
	_ = f.SetCellValue(sheet, "H6", "Debit")
	_ = f.SetCellValue(sheet, "I6", "Credit")
	_ = f.SetRowHeight(sheet, 6, 20)

	_ = f.SetCellStyle(sheet, "A5", "H5", st.GetStyle(style.CenterBold))
	_ = f.SetCellStyle(sheet, "A6", "C6", st.GetStyle(style.Bold))
	_ = f.SetCellStyle(sheet, "D6", "I6", st.GetStyle(style.RightBold))
	_ = f.SetRowHeight(sheet, 5, 20)

	starterRow := 7 // this is the first row
	for index, item := range balance.Values {
		rw := starterRow + index
		accountTypeName := typeMap.GetAccountTypeName(item.AccountTypeCd)

		_ = f.SetCellValue(sheet, getCell("A", rw), item.Code)
		_ = f.SetCellValue(sheet, getCell("B", rw), item.Name)
		_ = f.SetCellValue(sheet, getCell("C", rw), accountTypeName)
		_ = f.SetCellValue(sheet, getCell("D", rw), item.StartBalance.Debit)
		_ = f.SetCellValue(sheet, getCell("E", rw), item.StartBalance.Credit)
		_ = f.SetCellValue(sheet, getCell("F", rw), item.Runs.Debit)
		_ = f.SetCellValue(sheet, getCell("G", rw), item.Runs.Credit)
		_ = f.SetCellValue(sheet, getCell("H", rw), item.EndBalance.Debit)
		_ = f.SetCellValue(sheet, getCell("I", rw), item.EndBalance.Credit)

		startCell := getCell("D", rw)
		endCell := getCell("I", rw)
		_ = f.SetCellStyle(sheet, startCell, endCell, st.GetStyle(style.NumberFormat))
	}

	// the totals contain the word Total right aligned and then the values
	totalRow := starterRow + len(balance.Values)
	_ = f.MergeCell(sheet, getCell("A", totalRow), getCell("C", totalRow))
	_ = f.SetCellValue(sheet, getCell("A", totalRow), "Total:")

	totals := balance.Totals

	_ = f.SetCellValue(sheet, getCell("D", totalRow), totals.StartBalance.Debit)
	_ = f.SetCellValue(sheet, getCell("E", totalRow), totals.StartBalance.Credit)
	_ = f.SetCellValue(sheet, getCell("F", totalRow), totals.Runs.Debit)
	_ = f.SetCellValue(sheet, getCell("G", totalRow), totals.Runs.Credit)
	_ = f.SetCellValue(sheet, getCell("H", totalRow), totals.EndBalance.Debit)
	_ = f.SetCellValue(sheet, getCell("I", totalRow), totals.EndBalance.Credit)
	_ = f.SetRowHeight(sheet, totalRow, 25)
	_ = f.SetRowStyle(sheet, totalRow, totalRow, st.GetStyle(style.NumberFormatBold))
	_ = f.SetCellStyle(sheet, getCell("A", totalRow), getCell("A", totalRow), st.GetStyle(style.RightBold))

	return nil
}
