package sheet

import "github.com/dgb9/db-account-server/internal/data"

type AcctTypeMap interface {
	GetAccountTypeName(code string) string
}

type acctTypeStruct struct {
	accountTypeMap map[string]data.AccountTypeData
}

func NewAcctTypeMap(types []data.AccountTypeData) AcctTypeMap {
	info := acctTypeStruct{}
	info.accountTypeMap = make(map[string]data.AccountTypeData)

	for _, dt := range types {
		info.accountTypeMap[dt.AccountTypeCd] = dt
	}

	return &info
}

func (c *acctTypeStruct) GetAccountTypeName(code string) string {
	val, ok := c.accountTypeMap[code]
	res := "unknown"

	if ok {
		res = val.Name
	}

	return res
}
