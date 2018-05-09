package neoanderls

import (
	"github.com/labstack/echo"
	"neo-dev/utils/errorutil"
	"neo-dev/utils/httputil"
	"fmt"
	"net/http"
	"neo-dev/configure"
	"encoding/json"
	"neo-dev/cmd/coinsrv/db"
	"neo-dev/cmd/coinsrv/global"
	"strconv"
)

// Balance 查询余额
func Balance(ctx echo.Context) error {
	request := new(BalanceRequest)
	if err := ctx.Bind(request); err != nil {
		err = fmt.Errorf("解析JSON失败")
		data, _ := FailResponse(errorutil.Network_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	neo := NewNeoRequest()
	neo.Method = "balance"
	neo.Params = make([]interface{}, 0)
	// neo.Params =  append(neo.Params, "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b") // NEO 资产
	// neo.Params =  append(neo.Params, "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7") // GAS 资产
	neo.Params =  append(neo.Params, request.Address)
	respData, err := httputil.HTTPRequest("http://" + configure.RPCAddr, httputil.PostMethod, nil, neo)
	if err != nil {
		// fmt.Println(err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	// rpc 返回结果说明 
	/*
		{
			"jsonrpc": "2.0",
			"id": 1,
			"result": {
				"balance": "neo 真实余额",
				"confirmed": "neo 已确认金额，只有已确认金额可以用来转账, balance 和 confirmed 二者可能会不相等，仅在从钱包中转出一笔钱，而且有找零未确认时时，confirmed 值会小于balance。当这笔交易确定后，二者会变得相等",
				"gasbalance": "gas 真实余额",
				"gasconfirmed": "gas 已确认金额，只有已确认金额可以用来转账"
			}
		}
	*/
	// 返回结果数据转换
	neoRespones := new(NeoResponse)
	result := new(BalanceResult)
	neoRespones.Result = result
	err = json.Unmarshal(respData, neoRespones)
	if err != nil {
		err = fmt.Errorf("解析返回结果失败~%s", err)
		fmt.Println(err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	if neoRespones.Error.Code != 0 {
		err = fmt.Errorf("查询余额失败~%s", neoRespones.Error.Message)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	data, _ := json.Marshal(result)
	return ctx.JSONBlob(http.StatusOK, data)
}

// Create 创建账号
func Create(ctx echo.Context) error {
	neo := NewNeoRequest()
	neo.Method = "newaddress"
	neo.Params = make([]interface{}, 0)
	respData, err := httputil.HTTPRequest("http://" + configure.RPCAddr, httputil.PostMethod, nil, neo)
	if err != nil {
		// fmt.Println(err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	neoRespones := new(NeoResponse)
	addressResponse := new(AddressResponse)
	neoRespones.Result = addressResponse
	err = json.Unmarshal(respData, neoRespones)
	if err != nil {
		err = fmt.Errorf("解析返回结果失败~%s", err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	if addressResponse.Code != 0 {
		err = fmt.Errorf("创建账号失败~%s", addressResponse.Message)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	// 保存账号信息
	_, err = db.CreateAccount(global.DB, addressResponse.Address.Address, addressResponse.Address.Private)
	if err != nil {
		err = fmt.Errorf("保存账号信息失败~%s", err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	data, _ := json.Marshal(addressResponse.Address)
	return ctx.JSONBlob(http.StatusOK, data)
}

// TransactionInfoByHash 使用hash查询交易详情
func TransactionInfoByHash(ctx echo.Context) error {
	var (
		
		err error
		response QueryTransactionResponse
	)
	request := new(QueryTransactionRequest)
	if err = ctx.Bind(request); err != nil {
		// fmt.Println(err)
		err = fmt.Errorf("解析JSON失败")
		data, _ := FailResponse(errorutil.Network_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	neo := NewNeoRequest()
	neo.Method = "getrawtransaction"
	neo.Params = make([]interface{}, 0)
	neo.Params = append(neo.Params, request.TXID, 1) 
	respData, err := httputil.HTTPRequest("http://" + configure.RPCAddr, httputil.PostMethod, nil, neo)
	if err != nil {
		// fmt.Println(err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	// return ctx.JSONBlob(http.StatusOK, respData)
	
	// 保存交易信息
	neoRespones := new(NeoResponse)
	transactionResult := new(TransactionResult)
	neoRespones.Result = transactionResult
	err = json.Unmarshal(respData, neoRespones)
	if err != nil {
		err = fmt.Errorf("解析返回结果失败~%s", err)
		fmt.Println(err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	if neoRespones.Error.Code != 0 {
		err = fmt.Errorf("查询交易失败~%s", neoRespones.Error.Message)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	response.TXID = transactionResult.TXID
	response.NetFee, _ = strconv.ParseFloat(transactionResult.NetFee, 64)
	response.SysFee, _ = strconv.ParseFloat(transactionResult.SysFee, 64)
	// response.From = transactionResult.Vin[0].
	response.To = transactionResult.Vout[0].Address
	response.Amount, _ = strconv.ParseFloat(transactionResult.Vout[0].Value, 64)
	response.SetType(transactionResult.Vout[0].Asset)
	// response.GasPrice = transaction.GasLimit
	data, _ := json.Marshal(response)
	return ctx.JSONBlob(http.StatusOK, data)
	
}

// SendTransaction 发起交易
func SendTransaction(ctx echo.Context) error{
	var (
		err error
	)
	transaction := new(Transaction)
	if err = ctx.Bind(transaction); err != nil {
		err = fmt.Errorf("解析JSON失败 ~ %s", err)
		data, _ := FailResponse(errorutil.Network_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	if transaction.GetType() == "" {
		err = fmt.Errorf("参数错误")
		data, _ := FailResponse(errorutil.Network_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data) 
	}
	neo := NewNeoRequest()
	neo.Method = "sendfromto"
	neo.Params = make([]interface{}, 0)
	neo.Params = append(neo.Params, transaction.GetType(), transaction.From, transaction.To, transaction.Amount, transaction.Private)
	respData, err := httputil.HTTPRequest("http://" + configure.RPCAddr, httputil.PostMethod, nil, neo)
	if err != nil {
		// fmt.Println(err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	// 保存交易信息
	neoRespones := new(NeoResponse)
	transactionResult := new(TransactionResult)
	neoRespones.Result = transactionResult
	err = json.Unmarshal(respData, neoRespones)
	if err != nil {
		err = fmt.Errorf("解析返回结果失败~%s", err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	// fmt.Println("Code: ", neoRespones.Error.Code)
	if neoRespones.Error.Code != 0 {
		err = fmt.Errorf("发起交易失败~%s", neoRespones.Error.Message)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	order := new(db.TxOrder)
	order.TXID = transactionResult.TXID
	order.From = transaction.From
	order.To = transaction.To
	order.Amount = transaction.Amount
	order.Type = transaction.Type
	order.GasPrice = transaction.GasLimit
	order.NetFee, _ =  strconv.ParseFloat(transactionResult.NetFee, 64)
	order.SysFee, _ =  strconv.ParseFloat(transactionResult.SysFee, 64)
	_, err = db.CreateTransaction(global.DB, order)
	if err != nil {
		err = fmt.Errorf("保存交易信息失败~%s", err)
		data, _ := FailResponse(errorutil.System_Error, err.Error())
		return ctx.JSONBlob(http.StatusOK, data)
	}
	// data, _ := json.Marshal(order)
	return ctx.JSONBlob(http.StatusOK, respData)
}