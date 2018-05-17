package db

type Account struct {
	ID      int    `json:"id" db:"id"`
	Address string `json:"address" db:"address"` // 账号
	Private string `json:"private" db:"private"` // 私钥
}

type TxOrder struct {
	ID       int64   `json:"id" db:"id"`
	TXID     string  `json:"txid" db:"txid"`           // 交易hash
	Type     int     `json:"type" db:"type"`           // 交易类型 1、NEO 2、GAS
	From     string  `json:"from" db:"from"`           // 输出账号
	To       string  `json:"to" db:"to"`               // 输出账号
	Amount   float64 `json:"amount" db:"amount"`       // 输出金额
	GasPrice float64 `json:"gas_price" db:"gas_price"` // 交易费
	SysFee   float64 `json:"sys_fee" db:"sys_fee"`     // 系统费
	NetFee   float64 `json:"net_fee" db:"net_fee"`     // 网络费
}
