package sheet

import (
	"fmt"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/sheet/style"
	"github.com/xuri/excelize/v2"
)

func createTransactionReport(f *excelize.File, styles style.Style, start string, end string, accounts []data.Account, transactions []data.Transaction) error {
	sheet := "Transactions"
	_, _ = f.NewSheet(sheet)

	_ = f.SetCellValue(sheet, "A1", "Transactions")
	_ = f.SetCellStyle(sheet, "A1", "A1", styles.GetStyle(style.Title))
	_ = f.SetRowHeight(sheet, 1, 30)

	startDate := fmt.Sprintf("Starting with: %s", start)
	_ = f.SetCellValue(sheet, "A2", startDate)

	if len(end) > 0 {
		endDate := fmt.Sprintf("Ending with: %s", end)
		_ = f.SetCellValue(sheet, "A3", endDate)
	}

	_ = f.SetRowHeight(sheet, 2, 20)
	_ = f.SetRowHeight(sheet, 3, 20)

	_ = f.MergeCell(sheet, "A1", "F1")
	_ = f.MergeCell(sheet, "A2", "F2")
	_ = f.MergeCell(sheet, "A3", "F3")

	_ = f.SetRowStyle(sheet, 2, 3, styles.GetStyle(style.VerticalCenter))

	// the title here
	_ = f.SetCellValue(sheet, "A4", "Date")
	_ = f.SetCellValue(sheet, "B4", "Comments")
	_ = f.SetCellValue(sheet, "C4", "Amount")
	_ = f.SetCellValue(sheet, "D4", "Account")
	_ = f.SetCellStyle(sheet, "A4", "B4", styles.GetStyle(style.Bold))
	_ = f.SetCellStyle(sheet, "C4", "C4", styles.GetStyle(style.RightBold))
	_ = f.SetCellStyle(sheet, "D4", "D4", styles.GetStyle(style.Bold))

	_ = f.SetCellValue(sheet, "E4", "Debit")
	_ = f.SetCellValue(sheet, "F4", "Credit")
	_ = f.SetCellStyle(sheet, "E4", "F4", styles.GetStyle(style.RightBold))
	_ = f.SetRowHeight(sheet, 4, 20)

	// building account map
	acctMap := make(map[string]string)
	for _, acct := range accounts {
		acctMap[acct.AccountId] = fmt.Sprintf("%s - %s", acct.Code, acct.Name)
	}

	first := true
	newRow := 5

	for _, transaction := range transactions {
		first, newRow = processTransaction(f, sheet, styles, transaction, acctMap, first, newRow)
	}

	return nil
}

func processTransaction(f *excelize.File, sheet string, styles style.Style, transaction data.Transaction, acctMap map[string]string, first bool, row int) (bool, int) {

	positions := transaction.Positions
	rowCount := len(positions)

	for index, pos := range positions {
		currentRow := row + index
		if index == 0 {
			_ = f.SetCellValue(sheet, getCell("A", currentRow), transaction.TransactionDate)
			_ = f.SetCellValue(sheet, getCell("B", currentRow), transaction.Comments)

			amountCell := getCell("C", currentRow)
			_ = f.SetCellValue(sheet, amountCell, getAmount(positions))

			// except the first transaction line, each first line in transaction is distanced
			if !first {
				_ = f.SetRowHeight(sheet, currentRow, 25)
				_ = f.SetRowStyle(sheet, currentRow, currentRow, styles.GetStyle(style.VerticalBottom))
			}

			_ = f.SetCellStyle(sheet, amountCell, amountCell, styles.GetStyle(style.NumberFormat))
		}

		// account
		strAccount, ok := acctMap[pos.AccountId]
		if !ok {
			strAccount = "not found"
		}

		_ = f.SetCellValue(sheet, getCell("D", currentRow), strAccount)

		// debit and credit
		if pos.Debit != 0 {
			debitCell := getCell("E", currentRow)
			_ = f.SetCellValue(sheet, debitCell, pos.Debit)
			_ = f.SetCellStyle(sheet, debitCell, debitCell, styles.GetStyle(style.NumberFormat))
		}

		if pos.Credit != 0 {
			creditCell := getCell("F", currentRow)
			_ = f.SetCellValue(sheet, creditCell, pos.Credit)
			_ = f.SetCellStyle(sheet, creditCell, creditCell, styles.GetStyle(style.NumberFormat))
		}
	}

	return false, row + rowCount
}

func getAmount(positions []data.TransactionPosition) float64 {
	amount := 0.0

	for _, pos := range positions {
		amount += pos.Debit
	}

	return amount
}
