package data

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/util"
)

type RootHealthData struct {
	Version   string `json:"version"`
	Success   bool   `json:"success"`
	DbSuccess bool   `json:"dbSuccess"`
	DbError   string `json:"dbError"`
}

func GetRootHealthData(db *sql.DB) RootHealthData {
	dbSuccess, dbError := getDbSuccess(db)

	return RootHealthData{
		Version:   util.Version,
		Success:   true,
		DbSuccess: dbSuccess,
		DbError:   dbError,
	}
}

func getDbSuccess(db *sql.DB) (bool, string) {
	res := true
	message := ""

	// get a connection
	err := db.Ping()
	if err != nil {
		return false, err.Error()
	}

	// then close the connection
	return res, message
}
