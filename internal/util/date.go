package util

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func ParseStringDate(date string) (string, error) {
	strShortFormat := "Jan 2"
	strLongFormat := "Jan 2, 2006"

	_, err := time.Parse(strShortFormat, date)
	currentDate := date
	if err == nil {
		trimmed := strings.TrimSpace(date)
		currentYear := time.Now().Year()
		currentDate = fmt.Sprintf("%s, %d", trimmed, currentYear)
	}

	// now with the current date we attempt to format
	parsed, err := time.Parse(strLongFormat, currentDate)
	if err != nil {
		return "", err
	}

	strParsed := fmt.Sprintf("%04d%02d%02d", parsed.Year(), parsed.Month(), parsed.Day())

	return strParsed, nil
}

func FormatDisplayDate(date string) (string, error) {
	res := ""

	strDate := strings.TrimSpace(date)
	if len(strDate) > 0 {
		if len(strDate) != 8 {
			return res, errors.New("cannot format the date which has to have eight characters precisely")
		}

		format := "20060102"
		currentTime, err := time.Parse(format, strDate)
		if err != nil {
			return res, err
		}

		res = currentTime.Format("Jan 2, 2006")
	}

	return res, nil
}
