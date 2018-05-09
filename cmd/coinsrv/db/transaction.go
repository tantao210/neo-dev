package db

import (
	"github.com/jmoiron/sqlx"
)

// CreateTransaction 创建交易信息
func CreateTransaction(db *sqlx.DB, order *TxOrder) (int64, error) {
	sql := "INSERT INTO tx_order(txid, `type`, `from`, `to`, amount, gas_price, sys_fee, net_fee) values(?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(sql, order.TXID, order.Type, order.From, order.To, order.Amount, order.GasPrice, order.SysFee, order.NetFee)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	order.ID = id
	return id, err
}

func GetTransactionList(db *sqlx.DB) ([]*TxOrder, error) {
	sql := "SELECT txid, `type`, `from`, `to`, amount, gas_price, sys_fee, net_fee FROM tx_order WHERE MARKS=0 AND IFNULL(txid, '') <> '';"
	orderList := make([]*TxOrder, 0)
	err := db.Select(&orderList, sql)
	return orderList, err
}