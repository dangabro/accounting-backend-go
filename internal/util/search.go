package util

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type TransactionSearchData struct {
	Start    string
	End      string
	Accounts []string
	Comments []string
}

func ParseSearchString(search string) (TransactionSearchData, error) {
	res := TransactionSearchData{}

	reminder, err := proceedParse(search, &res, "date:", processDateParse)
	if err != nil {
		return res, err
	}

	reminder, err = proceedParse(reminder, &res, "accounts:", processAccountsParse)
	if err != nil {
		return res, err
	}

	err = parseComments(reminder, &res)

	return res, nil
}

func parseComments(reminder string, t *TransactionSearchData) error {
	trimmed := strings.TrimSpace(reminder)
	items := strings.Split(trimmed, " ")

	var comments []string
	for _, item := range items {
		trm := strings.TrimSpace(item)
		if len(trm) > 0 {
			comments = append(comments, trm)
		}
	}

	t.Comments = comments

	return nil
}

func processAccountsParse(val string, t *TransactionSearchData) error {
	r, err := regexp.Compile("(.*\\()(.*)(\\))")
	if err != nil {
		return err
	}

	interior := r.ReplaceAllString(val, "$2")

	parts := strings.Split(interior, ",")
	var resParts []string
	for _, str := range parts {
		resParts = append(resParts, strings.TrimSpace(str))
	}

	t.Accounts = resParts

	return nil
}

func processDateParse(val string, t *TransactionSearchData) error {
	r, err := regexp.Compile("(.*\\()(.*)(\\))")
	if err != nil {
		return err
	}

	interior := r.ReplaceAllString(val, "$2")

	parts := strings.Split(interior, "-")
	lnParts := len(parts)
	if lnParts == 1 {
		t.Start = strings.TrimSpace(parts[0])
		t.End = ""
	} else if lnParts == 2 {
		t.Start = strings.TrimSpace(parts[0])
		t.End = strings.TrimSpace(parts[1])

	} else {
		return errors.New("the date should have maximum two parts separated by dash")
	}

	// now we try to parse the date and then replace the values
	var parsedDate string
	start := t.Start
	if len(start) > 0 {
		parsedDate, err = ParseStringDate(start)
		if err != nil {
			return err
		}

		t.Start = parsedDate
	}

	end := t.End
	if len(end) > 0 {
		parsedDate, err = ParseStringDate(end)
		if err != nil {
			return err
		}

		t.End = parsedDate
	}

	return nil
}

func proceedParse(reminder string, t *TransactionSearchData, prefix string, further func(string, *TransactionSearchData) error) (string, error) {
	res := reminder

	format := fmt.Sprintf("(.*)(%s\\s*\\([^()]*\\))(.*)", prefix)
	reg, err := regexp.Compile(format)
	if err != nil {
		return "", err
	}

	if reg.MatchString(reminder) {
		before := reg.ReplaceAllString(reminder, "$1")
		val := reg.ReplaceAllString(reminder, "$2")
		after := reg.ReplaceAllString(reminder, "$3")

		err = further(val, t)
		if err != nil {
			return "", err
		}

		trimmedBefore := strings.TrimSpace(before)
		trimmedAfter := strings.TrimSpace(after)
		res = fmt.Sprintf("%s %s", trimmedBefore, trimmedAfter)
	}

	return res, nil
}
