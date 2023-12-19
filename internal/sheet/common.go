package sheet

import "fmt"

func getCell(col string, row int) string {
	return fmt.Sprintf("%s%d", col, row)
}
