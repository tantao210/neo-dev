package db

import (
	"github.com/jmoiron/sqlx"
)

// CreateAccount 创建一个账号
func CreateAccount(db *sqlx.DB, address, private string) (int64, error) {
	sql := `INSERT INTO account(address, private) VALUES(?, ?);`
	result, err := db.Exec(sql, address, private)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAccountByMain 查询钱包账号信息
func GetAccountByMain(db *sqlx.DB, ismain int) ([]*Account, error) {
	sql := `SELECT address, private FROM account WHERE marks=0 AND ismain=?;`
	var accountList []*Account
	accountList = make([]*Account, 0)
	err := db.Select(&accountList, sql, ismain)
	return accountList, err
}