package sheet

import (
	"fmt"
	"github.com/dgb9/db-account-server/internal/data"
	"github.com/dgb9/db-account-server/internal/sheet/style"
	"github.com/xuri/excelize/v2"
)

func createCompanySheet(f *excelize.File, st style.Style, company data.CompanyData) error {
	sheet := "Company"
	_, err := f.NewSheet(sheet)
	if err != nil {
		return err
	}

	name := company.Name
	month := company.Month
	day := company.Day

	closing := fmt.Sprintf("Closing month: %d and day: %d", month, day)

	styleTitle := st.GetStyle(style.Title)

	// we don't except issues with the current cell information
	_ = f.SetCellValue(sheet, "A1", "Company Information")
	_ = f.SetCellValue(sheet, "A2", name)
	_ = f.SetCellValue(sheet, "A3", closing)
	_ = f.SetColStyle(sheet, "A", styleTitle)
	_ = f.SetColWidth(sheet, "A", "A", 100)
	_ = f.SetRowHeight(sheet, 1, 50)
	_ = f.SetRowHeight(sheet, 2, 50)
	_ = f.SetRowHeight(sheet, 3, 50)

	return nil
}
