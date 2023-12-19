package dao

import (
	"database/sql"
	"strings"
)

func closeRows(rows *sql.Rows) {
	_ = rows.Close()
}

type ProcessIn interface {
	AddSql(sql string)
	AddSqlAndParam(sql string, param any)
	GetFinalSql() (string, []any)
	AddIn(start string, end string, params []any)
	ProcessFirst(first bool, strIfFalse string, strIfTrue string) bool
}

type processIn struct {
	sql    []string
	params []any
}

func NewProcessIn(start string) ProcessIn {
	return &processIn{
		sql: []string{start},
	}
}

func (d *processIn) AddSqlAndParam(sql string, param any) {
	d.sql = append(d.sql, sql)
	d.params = append(d.params, param)
}

func (d *processIn) AddSql(sql string) {
	d.sql = append(d.sql, sql)
}

func (d *processIn) GetFinalSql() (string, []any) {
	finalSql := strings.Join(d.sql, " ")

	return finalSql, d.params
}

func (d *processIn) AddIn(start string, end string, params []any) {
	d.sql = append(d.sql, start)

	first := true
	for _, id := range params {
		if first {
			first = false
		} else {
			d.AddSql(",")
		}

		d.AddSqlAndParam("?", id)
	}

	d.sql = append(d.sql, end)
}

func (d *processIn) ProcessFirst(first bool, strIfFalse string, strIfTrue string) bool {
	if first {
		d.AddSql(strIfTrue)
	} else {
		d.AddSql(strIfFalse)
	}

	return false
}
