package service

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/data"
)

type servc struct {
	db     *sql.DB
	config config.Config
}

type Service interface {
	Logout(token string) error
	Login(loginData data.LoginData) (data.LoginResult, error)
	UpdateCompany(adding bool, company data.CompanyData, userId string) error
	ValidateToken(token string) (data.User, error)
	GetAllCompanies(userId string) ([]data.CompanyData, error)
	GetAccountTypeCds() ([]data.AccountTypeData, error)
	DeleteAccounts(ids []string, userId string) error
	GetAllAccounts(companyId string, userId string) ([]data.Account, error)
	DeleteCompanies(ids []string, userId string) error
	DeleteTransactions(ids []string, userId string) error
	SearchAccounts(search string, companyId string, userId string) ([]data.Account, error)
	UpdateAccount(adding bool, account data.Account, userId string) error
	SearchTransactions(search string, companyId string, userId string) ([]data.Transaction, error)
	UpdateTransaction(adding bool, transaction data.Transaction, userId string) error
	TrialBalance(start string, end string, companyId string, userId string) (data.BalanceResult, error)
	AccountReport(start string, end string, accountId string, userId string) (data.AccountReportResult, error)
	TransactionReport(start string, end string, companyId string, userId string) (data.TransactionReportResult, error)
	ExcelReport(start string, end string, companyId string, userId string) (data.ExcelResponse, error)
}

func New(db *sql.DB, config config.Config) Service {
	return &servc{
		db:     db,
		config: config,
	}
}
