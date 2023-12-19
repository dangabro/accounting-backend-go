package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func GetConnectionPool(config Config) (*sql.DB, error) {
	user, password, machine, port, database := config.Database()

	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, machine, port, database)
	db, err := sql.Open("mysql", url)

	if err != nil {
		return nil, err
	}

	return db, nil
}
