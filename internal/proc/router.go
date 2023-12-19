package proc

import (
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/dgb9/db-account-server/internal/controller"
	"github.com/gorilla/mux"
)

func Router(config config.Config, db *sql.DB) *mux.Router {
	context := config.Context()

	r := mux.NewRouter()
	r.Handle(context+"/", controller.Root(config, db)).Methods("GET")
	r.Handle(context+"/login", controller.Login(config, db)).Methods("POST")
	r.Handle(context+"/logout", controller.Logout(config, db)).Methods("POST")
	r.Handle(context+"/company", controller.Company(config, db)).Methods("POST")
	r.Handle(context+"/companies", controller.Companies(config, db)).Methods("GET")
	r.Handle(context+"/accountType", controller.AccountType(config, db)).Methods("GET")
	r.Handle(context+"/deleteAccounts", controller.DeleteAccounts(config, db)).Methods("POST")
	r.Handle(context+"/loadAccounts/{companyId}", controller.LoadAccounts(config, db)).Methods("GET")
	r.Handle(context+"/deleteCompanies", controller.DeleteCompanies(config, db)).Methods("POST")
	r.Handle(context+"/deleteTransactions", controller.DeleteTransactions(config, db)).Methods("POST")
	r.Handle(context+"/searchAccounts", controller.SearchAccounts(config, db)).Methods("POST")
	r.Handle(context+"/updateAccount", controller.UpdateAccount(config, db)).Methods("POST")
	r.Handle(context+"/searchTransactions", controller.SearchTransactions(config, db)).Methods("POST")
	r.Handle(context+"/updateTransaction", controller.UpdateTransaction(config, db)).Methods("POST")
	r.Handle(context+"/trialBalanceReport", controller.TrialBalanceReport(config, db)).Methods("POST")
	r.Handle(context+"/accountReport", controller.AccountReport(config, db)).Methods("POST")
	r.Handle(context+"/transactionReport", controller.TransactionReport(config, db)).Methods("POST")
	r.Handle(context+"/excelReport", controller.ExcelReport(config, db)).Methods("POST")
	r.Handle(context+"/version", controller.Version()).Methods("GET")

	return r
}
