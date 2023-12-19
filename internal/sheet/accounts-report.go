package sheet

import (
	"fmt"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/sheet/style"
	"github.com/xuri/excelize/v2"
)

func createAccountReports(f *excelize.File, styles style.Style, start string, end string, typesMap AcctTypeMap, items []data.AccountReportResult) error {
	sheet := "Account Details"
	_, _ = f.NewSheet(sheet)

	_ = f.SetCellValue(sheet, "A1", "Account Details")
	_ = f.SetRowHeight(sheet, 1, 30)
	_ = f.SetRowStyle(sheet, 1, 1, styles.GetStyle(style.Title))

	_ = f.SetCellValue(sheet, "A2", fmt.Sprintf("Starting from: %s", start))
	_ = f.SetCellStyle(sheet, "A2", "A2", styles.GetStyle(style.Bold))

	if len(end) > 0 {
		_ = f.SetCellValue(sheet, "A3", fmt.Sprintf("Ending with: %s", end))
		_ = f.SetCellStyle(sheet, "A3", "A3", styles.GetStyle(style.Bold))
	}

	_ = f.MergeCell(sheet, "A1", "H1")
	_ = f.MergeCell(sheet, "A2", "H2")
	_ = f.MergeCell(sheet, "A3", "H3")

	startRow := 4
	first := true
	for _, acctData := range items {
		startRow, first = processAccountData(f, sheet, styles, typesMap, acctData, startRow, first)
	}

	return nil
}

func processAccountData(f *excelize.File, sheet string, styles style.Style, typesMap AcctTypeMap, acctData data.AccountReportResult, row int, first bool) (int, bool) {
	details := acctData.Details

	// about 5 lines plus the length of the details, this makes info, dates, start balance, title and end balance plus one line per details
	code := acctData.Code
	name := acctData.Name
	accountType := typesMap.GetAccountTypeName(acctData.AccountTypeCd)
	info := fmt.Sprintf("Account: %s - %s (%s)", code, name, accountType)

	// the info row
	infoCell := getCell("A", row)
	_ = f.SetCellValue(sheet, infoCell, info)
	_ = f.SetCellStyle(sheet, infoCell, infoCell, styles.GetStyle(style.Bold))
	_ = f.MergeCell(sheet, infoCell, getCell("H", row))

	if !first {
		_ = f.SetCellStyle(sheet, infoCell, infoCell, styles.GetStyle(style.VerticalBottomBold))
		_ = f.SetRowHeight(sheet, row, 30)
	}

	// dates
	rowDates := row + 1
	start := acctData.Start
	end := acctData.End
	dates := fmt.Sprintf("Dates %s - %s", start, end)
	datesCell := getCell("A", rowDates)
	endCell := getCell("H", rowDates)

	_ = f.SetCellValue(sheet, datesCell, dates)
	_ = f.MergeCell(sheet, datesCell, endCell)

	// title
	rowTitle := row + 2
	_ = f.SetCellValue(sheet, getCell("A", rowTitle), "Date")
	_ = f.SetCellStyle(sheet, getCell("A", rowTitle), getCell("A", rowTitle), styles.GetStyle(style.VerticalCenterRightBold))
	_ = f.SetCellValue(sheet, getCell("B", rowTitle), "Amount")
	_ = f.SetCellStyle(sheet, getCell("B", rowTitle), getCell("B", rowTitle), styles.GetStyle(style.CenterBold))

	_ = f.SetCellValue(sheet, getCell("D", rowTitle), "Total Amount")
	_ = f.SetCellStyle(sheet, getCell("D", rowTitle), getCell("D", rowTitle), styles.GetStyle(style.VerticalCenterRightBold))

	_ = f.SetCellValue(sheet, getCell("E", rowTitle), "Comments")
	_ = f.SetCellStyle(sheet, getCell("E", rowTitle), getCell("E", rowTitle), styles.GetStyle(style.CenterBold))

	_ = f.SetCellValue(sheet, getCell("F", rowTitle), "Codes")
	_ = f.SetCellStyle(sheet, getCell("F", rowTitle), getCell("F", rowTitle), styles.GetStyle(style.CenterBold))

	_ = f.SetCellValue(sheet, getCell("H", rowTitle), "Balance")
	_ = f.SetCellStyle(sheet, getCell("H", rowTitle), getCell("H", rowTitle), styles.GetStyle(style.VerticalCenterRightBold))

	rowTitleSecond := row + 3

	_ = f.SetCellValue(sheet, getCell("B", rowTitleSecond), "Debit")
	_ = f.SetCellValue(sheet, getCell("C", rowTitleSecond), "Credit")
	_ = f.SetCellValue(sheet, getCell("F", rowTitleSecond), "Debit")
	_ = f.SetCellValue(sheet, getCell("G", rowTitleSecond), "Credit")
	_ = f.SetCellStyle(sheet, getCell("B", rowTitleSecond), getCell("C", rowTitleSecond), styles.GetStyle(style.RightBold))
	_ = f.SetCellStyle(sheet, getCell("F", rowTitleSecond), getCell("G", rowTitleSecond), styles.GetStyle(style.RightBold))

	// merge the cells
	_ = f.MergeCell(sheet, getCell("A", rowTitle), getCell("A", rowTitleSecond)) // date
	_ = f.MergeCell(sheet, getCell("B", rowTitle), getCell("C", rowTitle))       // amount
	_ = f.MergeCell(sheet, getCell("D", rowTitle), getCell("D", rowTitleSecond)) // total amount
	_ = f.MergeCell(sheet, getCell("E", rowTitle), getCell("E", rowTitleSecond)) // comments
	_ = f.MergeCell(sheet, getCell("F", rowTitle), getCell("G", rowTitle))       // codes
	_ = f.MergeCell(sheet, getCell("H", rowTitle), getCell("H", rowTitleSecond)) // balance

	// now produce all
	startRow := row + 4
	_ = f.SetCellValue(sheet, getCell("A", startRow), "Start:")
	_ = f.SetCellStyle(sheet, getCell("A", startRow), getCell("A", startRow), styles.GetStyle(style.RightBold))

	_ = f.SetCellValue(sheet, getCell("B", startRow), acctData.StartBalance.Debit)
	_ = f.SetCellValue(sheet, getCell("C", startRow), acctData.StartBalance.Credit)
	_ = f.SetCellStyle(sheet, getCell("B", startRow), getCell("C", startRow), styles.GetStyle(style.NumberFormatBold))
	_ = f.MergeCell(sheet, getCell("D", startRow), getCell("H", startRow))

	for index, detail := range details {
		currentRow := index + row + 5
		_ = f.SetCellValue(sheet, getCell("A", currentRow), detail.Date)

		_ = f.SetCellValue(sheet, getCell("B", currentRow), detail.Amount.Debit)
		_ = f.SetCellValue(sheet, getCell("C", currentRow), detail.Amount.Credit)
		_ = f.SetCellStyle(sheet, getCell("B", currentRow), getCell("C", currentRow), styles.GetStyle(style.NumberFormat))

		_ = f.SetCellValue(sheet, getCell("D", currentRow), detail.TransactionAmount)
		_ = f.SetCellStyle(sheet, getCell("D", currentRow), getCell("D", currentRow), styles.GetStyle(style.NumberFormat))

		_ = f.SetCellValue(sheet, getCell("E", currentRow), detail.Comments)
		_ = f.SetCellValue(sheet, getCell("F", currentRow), detail.DebitCodes)
		_ = f.SetCellValue(sheet, getCell("G", currentRow), detail.CreditCodes)
		_ = f.SetCellValue(sheet, getCell("H", currentRow), detail.CurrentBalance)
		_ = f.SetCellStyle(sheet, getCell("H", currentRow), getCell("H", currentRow), styles.GetStyle(style.NumberFormat))
	}

	endRow := row + 5 + len(details)

	_ = f.SetCellValue(sheet, getCell("A", endRow), "End:")
	_ = f.SetCellStyle(sheet, getCell("A", endRow), getCell("A", endRow), styles.GetStyle(style.RightBold))

	_ = f.SetCellValue(sheet, getCell("B", endRow), acctData.FinalBalance.Debit)
	_ = f.SetCellValue(sheet, getCell("C", endRow), acctData.FinalBalance.Credit)
	_ = f.SetCellStyle(sheet, getCell("B", endRow), getCell("C", endRow), styles.GetStyle(style.NumberFormatBold))
	_ = f.MergeCell(sheet, getCell("D", endRow), getCell("H", endRow))

	return endRow + 1, false
}
