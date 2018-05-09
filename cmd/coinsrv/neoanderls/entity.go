package neoanderls

import (
	"neo-dev/configure"
)

// BaseNeo neo 基础
type BaseNeo struct {
	ID int `json:"id"`
	JsonRPC string `json:"jsonrpc"`
}

// NeoFailResponse 失败响应
type NeoFailResponse struct {
	Code int64 `json:"code"`
	Message string `json:"message"`
}

// NeoResponse NEO Response
type NeoResponse struct {
	BaseNeo
	Result interface{} `json:"result"`
	Error NeoFailResponse `json:"error"`
}

// BalanceRequest 余额查询请求
type BalanceRequest struct {
	Address string `json:"address"`
}

// BalanceResult 余额查询结果
type BalanceResult struct {
	Balance string `json:"balance"` // NEO 真实余额 
	Confirmed string `json:"confirmed"` // NEO 已确认金额，只有已确认金额可以用来转账, balance 和 confirmed 二者可能会不相等，仅在从钱包中转出一笔钱，而且有找零未确认时时，confirmed 值会小于balance。当这笔交易确定后，二者会变得相等
	GasBalance string `json:"gasbalance"`// GAS 真实余额
	GasConfirmed string `json:"gasconfirmed"` // GAS 已确认金额，只有已确认金额可以用来转账, balance 和 confirmed 二者可能会不相等，仅在从钱包中转出一笔钱，而且有找零未确认时时，confirmed 值会小于balance。当这笔交易确定后，二者会变得相等
}

// NeoRequest neo rpc request
type NeoRequest struct {
	BaseNeo
	Method string `json:"method"`
	Params []interface{} `json:"params"`
	
}

// NewNeoRequest
func NewNeoRequest() *NeoRequest {
	request := new(NeoRequest)
	request.JsonRPC = "2.0"
	request.ID = 1
	return request
}

// AddressResponse 创建账号返回结果
type AddressResponse struct {
	NeoFailResponse // 失败才有的返回结果
	Address // 成功返回
}

// Address 账号信息
type Address struct {
	Address string `json:"address"`
	Private string `json:"private"`
}

// Transaction 发起交易请求
type Transaction struct {
	Type int `json:"type"` // 类型
	From string `json:"from"`  // output account
	To string `json:"to"` // input account
	Amount float64 `json:"amount"` // amount
	Private string `json:"private"` // private key
	GasLimit float64 `json:"gaslimit"` // gas limit  暂时不用此参数
	ChangeAddress string `json:"change_address"` // 找零地址 默认 input account  暂时不用此参数
}

// GetType 交易类型
func (t *Transaction) GetType() string {
	if t.Type == 1 { // NEO
		return configure.NeoAsset
	} else if t.Type == 2 {// GAS
		return configure.GasAsset
	}
	return ""
}

// TransactionResult 交易返回结果
type TransactionResult struct {
	TXID string `json:"txid"` // 交易Hash
	Size int `json:"size"` // 大小
	Version int `json:"version"` // 版本
	Attributes []*TransactionAttribute `json:"attributes"`// 该交易所具备的额外特性
	Vin []*CoinReference `json:"vin"` // 输入列表
	Vout []*TransactionOutput `json:"vout"`// 输出列表
	SysFee string `json:"sys_fee"` // 系统费
	NetFee string `json:"net_fee"` // 网络费
	Scripts []*Witness `json:"scripts"` // 用于验证交易的脚本列表
}

// Witness 合约脚本
type Witness struct {
	Invocation string `json:"invocation"` // 调用
	Verification string `json:"verification"` // 核查
}

// TransactionOutput 交易输出
type TransactionOutput struct {
	N int `json:"n"` // 该交易输出在交易中的索引
	Asset string `json:"asset"` // 资产编号
	Value string `json:"value"` // 金额
	Address string `json:"address"` // 账号的ScriptHash
}

// CoinReference 交易输入
type CoinReference struct {
	TXID string `json:"txid"` // 引用交易的散列值
	Vout int `json:"vout"` // 引用交易输出的索引 
}

// TransactionAttribute 外部数据
type TransactionAttribute struct {
	Usage string `json:"usage"` // 用途
	Data string `json:"data"`// 特定于用途的外部数据
}

// QueryTransactionRequest 查询交易请求
type QueryTransactionRequest struct {
	TXID string `json:"txid"`
}

// QueryTransactionResponse 查询交易响应
type QueryTransactionResponse struct {
	TXID string `json:"txid"` // 交易hash
	Type int `json:"type"` // 交易类型 1、NEO 2、GAS
	From string `json:"from"` // 输出账号
	To string `json:"to"` // 输出账号
	Amount float64 `json:"amount"` // 输出金额
	GasPrice float64 `json:"gas_price"` // 交易费
	SysFee float64 `json:"sys_fee"` // 系统费
	NetFee float64 `json:"net_fee"` // 网络费
}

// SetType 设置交易类型
func (q *QueryTransactionResponse) SetType(asset string) {
	if asset[2:] == configure.NeoAsset { // NEO
		q.Type = 1
	} else if asset[2:] == configure.GasAsset { // GAS
		q.Type = 2
	}
}