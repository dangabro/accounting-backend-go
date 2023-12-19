package data

import (
	"bytes"
	"time"
)

type Session struct {
	SessionId string
	UserId    string
	Token     string
	Expired   bool
	ExpiryDt  time.Time
}

type User struct {
	UserId     string
	ProvidedId string
	Login      string
	Name       string
}

type UserSSO struct {
	Id     string   `json:"id"`
	Login  string   `json:"login"`
	Name   string   `json:"name"`
	Rights []string `json:"rights"`
}

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResult struct {
	Login string `json:"login"`
	Token string `json:"token"`
	Id    string `json:"id"`
	Name  string `json:"name"`
}

type CompanyData struct {
	CompanyId string `json:"id"`
	UserId    string `json:"userId"`
	Name      string `json:"name"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
}

type UpdateCompanyRequest struct {
	Adding  bool        `json:"adding"`
	Company CompanyData `json:"company"`
}

type AccountTypeData struct {
	AccountTypeCd string `json:"accountTypeCd"`
	Name          string `json:"name"`
	BalanceTypeCd string `json:"balanceTypeCd"`
}

type IdsCollection struct {
	Ids []string `json:"ids"`
}

type Account struct {
	AccountId     string `json:"accountId"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	AccountTypeCd string `json:"accountTypeCd"`
	CompanyId     string `json:"companyId"`
}

type AccountSearch struct {
	CompanyId string `json:"companyId"`
	Search    string `json:"search"`
}

type UpdateAccountRequest struct {
	Adding  bool    `json:"adding"`
	Account Account `'json:"account"'`
}

type TransactionPosition struct {
	TransactionPositionId string  `json:"transactionPositionId"`
	TransactionId         string  `json:"transactionId"`
	AccountId             string  `json:"accountId"`
	Sequence              int32   `json:"sequence"`
	Debit                 float64 `json:"debit"`
	Credit                float64 `json:"credit"`
	Comments              string  `json:"comments"`
}

type Transaction struct {
	TransactionId   string                `json:"transactionId"`
	CompanyId       string                `json:"companyId"`
	TransactionDate string                `json:"transactionDate"`
	Sequence        int32                 `json:"sequence"`
	Comments        string                `json:"comments"`
	Positions       []TransactionPosition `json:"positions"`
}

type SearchTransactionsRequest struct {
	CompanyId string `json:"companyId"`
	Search    string `json:"search"`
}

type UpdateTransactionRequest struct {
	Transaction Transaction `json:"transaction"`
	Adding      bool        `json:"adding"`
}

type AccountValues struct {
	Debit  float64 `json:"debit"`
	Credit float64 `json:"credit"`
}

type TrialBalance struct {
	AccountId     string        `json:"accountId"`
	Code          string        `json:"code"`
	Name          string        `json:"name"`
	AccountTypeCd string        `json:"accountTypeCd"`
	StartBalance  AccountValues `json:"startBalance"`
	Runs          AccountValues `json:"runs"`
	EndBalance    AccountValues `json:"endBalance"`
}

type BalanceResult struct {
	StartDate string         `json:"start"`
	EndDate   string         `json:"end"`
	CompanyId string         `json:"companyId"`
	Values    []TrialBalance `json:"items"`
	Totals    TrialBalance   `json:"totals"`
}

type TrialBalanceRequest struct {
	StartDate string `json:"start"`
	EndDate   string `json:"end"`
	CompanyId string `json:"companyId"`
}

type ExcelReportRequest struct {
	StartDate string `json:"start"`
	EndDate   string `json:"end"`
	CompanyId string `json:"companyId"`
}

type AccountReportRequest struct {
	StartDate string `json:"start"`
	EndDate   string `json:"end"`
	AccountId string `json:"accountId"`
}

type AccountReportDetail struct {
	TransactionPositionId string        `json:"transactionPositionId"`
	Date                  string        `json:"date"`
	Amount                AccountValues `json:"amount"`
	TransactionAmount     float64       `json:"transactionAmount"`
	Comments              string        `json:"comments"`
	DebitCodes            string        `json:"debitCodes"`
	CreditCodes           string        `json:"creditCodes"`
	CurrentBalance        float64       `json:"currentBalance"`
}

type AccountReportResult struct {
	AccountId     string                `json:"accountId"`
	Code          string                `json:"code"`
	Name          string                `json:"name"`
	AccountTypeCd string                `json:"accountTypeCd"`
	Start         string                `json:"start"`
	End           string                `json:"end"`
	StartBalance  AccountValues         `json:"startBalance"`
	Details       []AccountReportDetail `json:"details"`
	Totals        AccountValues         `json:"totals"`
	FinalBalance  AccountValues         `json:"finalBalance"`
}

type TransactionReportRequest struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	CompanyId string `json:"companyId"`
}

type TransactionReportResult struct {
	Start        string        `json:"start"`
	End          string        `json:"end"`
	CompanyId    string        `json:"companyId"`
	Transactions []Transaction `json:"transactions"`
}

type ExcelReport struct {
	Start          string
	End            string
	Company        CompanyData
	Balance        BalanceResult
	Accounts       []Account
	AccountTypes   []AccountTypeData
	Transactions   []Transaction
	AccountsReport []AccountReportResult
}

type ExcelResponse struct {
	DataBuffer *bytes.Buffer
	FileName   string
}

type VersionData struct {
	Version string `json:"version"`
}
