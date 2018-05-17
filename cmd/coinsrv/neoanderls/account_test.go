package neoanderls

import (
	"encoding/json"
	"neo-dev/cmd/coinsrv/db"
	"neo-dev/cmd/coinsrv/global"
	"neo-dev/utils/httputil"
	"strconv"
	"testing"
	"time"
)

/**
 * cmmand
 * go test neo-dev\cmd\coinsrv\neoanderls -run TestNeo -timeout 20000s
 *
**/

func TestCreate(t *testing.T) {

}

// TestNeo 测试创建账号 发起交易 查询余额 查询交易功能
func TestNeo(t *testing.T) {
	var (
		// err error
		host        string
		createUrl   string
		accountNum  int     // 创建账号的数量
		neoAmount   int     // 转 neo 的数量
		gasAmount   float64 // 转 Gas 的数量
		sendTranUrl string
		balanceUrl  string
		getTxUrl    string
		header      map[string]string
	)
	// t.Errorf("调试过程的固定Error")
	// 初始化基本参数
	host = `http://127.0.0.1:6000`
	accountNum = 1
	neoAmount = 1
	gasAmount = 0.5
	// 创建账号
	createUrl = host + `/neo/account/create`
	sendTranUrl = host + `/neo/trade/send`
	balanceUrl = host + `/neo/account/query`
	getTxUrl = host + `/neo/trade/query`
	header = make(map[string]string)
	header["Content-Type"] = "application/json"
	for i := 0; i < accountNum; i++ {
		data, err := httputil.HTTPRequest(createUrl, httputil.PostMethod, header, nil)
		if err != nil {
			t.Errorf("请求创建账号服务失败~%s", err)
		}
		address := new(Address)
		err = json.Unmarshal(data, address)
		if err != nil {
			failResp := new(failResponse)
			err = json.Unmarshal(data, failResp)
			if err != nil {
				t.Errorf("解析创建账号返回结果失败 ~ %s", err)
			}
			if !ResponseIsSuccess(failResp.Code) {
				t.Errorf("创建账号失败 ~ %s", failResp.Msg)
			}
		} else {
			if address.Address == "" {
				t.Errorf("创建账号不成功，账号信息为空~")
			} else {
				t.Logf("创建账号: %d/%d", i+1, accountNum)
			}
		}
	}
	t.Logf("已完成 [%d] 个账号的创建", accountNum)

	// 初始华数据库操作对象
	global.NewWebApp()

	// 发起交易
	mainAccountList, err := db.GetAccountByMain(global.DB, 1)
	if err != nil {
		t.Errorf("主账号不存在 ~ %s", err)
	}
	mainAccount := mainAccountList[0]

	accountList, err := db.GetAccountByMain(global.DB, 0)
	if err != nil {
		t.Errorf("账号不存在 ~ %s", err)
	}

	// fmt.Println("开始交易 ", len(accountList))

	for i, account := range accountList {
		t.Logf("发起交易: %d/%d", i+1, len(accountList))
		for tradeType := 1; tradeType < 3; tradeType++ {
			tranReq := new(Transaction)
			tranReq.Type = tradeType
			tranReq.From = mainAccount.Address
			tranReq.To = account.Address
			if tradeType == 1 {
				tranReq.Amount = float64(neoAmount)
			} else {
				tranReq.Amount = gasAmount
			}
			tranReq.Private = mainAccount.Private
			data, err := httputil.HTTPRequest(sendTranUrl, httputil.PostMethod, header, tranReq)
			if err != nil {
				t.Errorf("请求发起交易服务失败~%s", err)
			}
			failResp := new(failResponse)
			err = json.Unmarshal(data, failResp)
			if err != nil {
				t.Errorf("解析发起交易返回结果失败 ~ %s", err)
			} else {
				if failResp.Code != "" && !ResponseIsSuccess(failResp.Code) {
					t.Errorf("发起交易失败 ~ %s", failResp.Msg)
				} else {
					// fmt.Println(string(data))
					// t.Logf(string(data))
					neoRespones := new(NeoResponse)
					transactionResult := new(TransactionResult)
					neoRespones.Result = transactionResult
					err = json.Unmarshal(data, neoRespones)
					// t.Logf(string(data))
					if err != nil {
						t.Errorf("解析发起交易返回结果失败 ~ %s", err)
					}
					if neoRespones.Error.Code != 0 {
						t.Errorf("交易失败 ~ %s", neoRespones.Error.Message)
					}
					if transactionResult.TXID != "" {
						t.Logf("交易成功 ~ %s", transactionResult.TXID)
					}
				}
			}
			time.Sleep(time.Second * 15)
		}
	}

	// fmt.Println("开始查询余额 ", len(accountList))
	// 查询余额
	for i, account := range accountList {
		t.Logf("查询余额: %d / %d", i+1, len(accountList))
		balanceRequest := new(BalanceRequest)
		balanceRequest.Address = account.Address
		data, err := httputil.HTTPRequest(balanceUrl, httputil.PostMethod, header, balanceRequest)
		if err != nil {
			t.Errorf("请求查询余额服务失败~%s", err)
		}
		// neoRespones := new(NeoResponse)
		result := new(BalanceResult)
		// neoRespones.Result = result
		err = json.Unmarshal(data, result)
		if err != nil {
			failResp := new(failResponse)
			err = json.Unmarshal(data, failResp)
			if err != nil {
				t.Errorf("解析查询余额返回结果失败 ~ %s", err)
			}
			if !ResponseIsSuccess(failResp.Code) {
				t.Errorf("查询余额失败 ~ %s", failResp.Msg)
			}
		} else {
			t.Logf("查询余额成功 ~ Balance: %s\t Confirmd: %s\t GasBalance: %s\t GasConfirmed: %s", result.Balance, result.Confirmed, result.GasBalance, result.GasConfirmed)
		}
	}
	// 查询交易
	// fmt.Println("开始获取交易记录 ", len(accountList))
	orderList, err := db.GetTransactionList(global.DB)
	if err != nil {
		t.Errorf("获取交易记录失败 ~ %s", err)
	}
	for i, order := range orderList {
		t.Logf("查询交易: %d / %d", i+1, len(orderList))
		request := new(QueryTransactionRequest)
		request.TXID = order.TXID
		data, err := httputil.HTTPRequest(getTxUrl, httputil.PostMethod, header, request)
		if err != nil {
			t.Errorf("请求查询交易服务失败~%s", err)
		}
		response := new(QueryTransactionResponse)
		err = json.Unmarshal(data, response)
		if err != nil {
			failResp := new(failResponse)
			err = json.Unmarshal(data, failResp)
			if err != nil {
				t.Errorf("解析查询交易返回结果失败 ~ %s", err)
			}
			if !ResponseIsSuccess(failResp.Code) {
				t.Errorf("查询交易失败 ~ %s", failResp.Msg)
			}
		} else {
			// t.Logf("Neo\tBalance\tConfirmed\t\tGAS\tBalance\tConfirmed\nNEO\t%d\t%d\t\tGAS\t%f\t%f", response)
			t.Logf("TXID: %s\nFrom: %s\nTo: %s\nAmount: %v\nType: %d", response.TXID, response.From, response.To, response.Amount, response.Type)
		}
	}
}

// TestCollect 归集所有账号的余额到主账号
func TestCollect(t *testing.T) {
	var (
		// err error
		host        string
		sendTranUrl string
		balanceUrl  string
		header      map[string]string
	)
	t.Errorf("调试的固定Error")
	// 初始化基本参数
	host = `http://127.0.0.1:6000`
	sendTranUrl = host + `/neo/trade/send`
	balanceUrl = host + `/neo/account/query`
	header = make(map[string]string)
	header["Content-Type"] = "application/json"
	global.NewWebApp()
	mainAccountList, err := db.GetAccountByMain(global.DB, 1)
	if err != nil {
		t.Errorf("主账号不存在 ~ %s", err)
	}
	mainAccount := mainAccountList[0]

	accountList, err := db.GetAccountByMain(global.DB, 0)
	if err != nil {
		t.Errorf("账号不存在 ~ %s", err)
	}
	// 查询余额
	for i, account := range accountList {
		t.Logf("归集查询余额: %d / %d", i+1, len(accountList))
		balanceRequest := new(BalanceRequest)
		balanceRequest.Address = account.Address
		data, err := httputil.HTTPRequest(balanceUrl, httputil.PostMethod, header, balanceRequest)
		if err != nil {
			t.Errorf("请求查询余额服务失败~%s", err)
		}
		// neoRespones := new(NeoResponse)
		result := new(BalanceResult)
		// neoRespones.Result = result
		err = json.Unmarshal(data, result)
		if err != nil {
			failResp := new(failResponse)
			err = json.Unmarshal(data, failResp)
			if err != nil {
				t.Errorf("解析查询余额返回结果失败 ~ %s", err)
			}
			if !ResponseIsSuccess(failResp.Code) {
				t.Errorf("查询余额失败 ~ %s", failResp.Msg)
			}
		} else {
			t.Logf("查询余额成功 ~ Balance: %s\t Confirmd: %s\t GasBalance: %s\t GasConfirmed: %s", result.Balance, result.Confirmed, result.GasBalance, result.GasConfirmed)

			// 归集所有余额
			// t.Logf("开始归集所有余额")
			var data []byte
			if confirmed, err := strconv.Atoi(result.Confirmed); err == nil && confirmed > 0 {
				t.Logf("开始归集账号 [%s] 的 NEO [%d]", account.Address, confirmed)
				tranReq := new(Transaction)
				tranReq.Type = 1
				tranReq.From = account.Address
				tranReq.To = mainAccount.Address
				tranReq.Amount = float64(confirmed)
				tranReq.Private = account.Private
				data, err = httputil.HTTPRequest(sendTranUrl, httputil.PostMethod, header, tranReq)
				if err != nil {
					t.Errorf("归集所有余额请求发起NEO交易服务失败~%s", err)
				}
			}
			if gasconfirmed, err := strconv.ParseFloat(result.GasConfirmed, 64); err == nil && gasconfirmed > 0 {
				t.Logf("开始归集账号 [%s] 的 GAS [%f]", account.Address, gasconfirmed)
				tranReq := new(Transaction)
				tranReq.Type = 2
				tranReq.From = account.Address
				tranReq.To = mainAccount.Address
				tranReq.Amount = gasconfirmed
				tranReq.Private = account.Private
				data, err = httputil.HTTPRequest(sendTranUrl, httputil.PostMethod, header, tranReq)
				if err != nil {
					t.Errorf("归集所有余额请求发起GAS交易服务失败~%s", err)
				}
			}
			if len(data) > 0 {
				failResp := new(failResponse)
				err = json.Unmarshal(data, failResp)
				if err != nil {
					t.Errorf("解析发起交易返回结果失败 ~ %s ~ %s", err, string(data))
				} else {
					if failResp.Code != "" && !ResponseIsSuccess(failResp.Code) {
						t.Errorf("发起交易失败 ~ %s", failResp.Msg)
					} else {
						// fmt.Println(string(data))
						// t.Logf(string(data))
						neoRespones := new(NeoResponse)
						transactionResult := new(TransactionResult)
						neoRespones.Result = transactionResult
						err = json.Unmarshal(data, neoRespones)
						// t.Logf(string(data))
						if err != nil {
							t.Errorf("解析发起交易返回结果失败 ~ %s", err)
						}
						if neoRespones.Error.Code != 0 {
							t.Errorf("交易失败 ~ %s", neoRespones.Error.Message)
						}
						if transactionResult.TXID != "" {
							t.Logf("交易成功 ~ %s", transactionResult.TXID)
						}
					}
				}
			}
		}
	}
	t.Logf("归集完成")
}
