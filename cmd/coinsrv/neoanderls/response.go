package neoanderls

import (
	"encoding/json"
	"reflect"
)

type failResponse struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
}

type successResponse struct {
	Code string `json:"code"`
}

func FailResponse(code string, msg string) (data []byte, err error) {
	res := new(failResponse)
	res.Code, res.Msg = code, msg
	return json.Marshal(res)
}

func SuccessResponse(res interface{}) (data []byte, err error) {
	if reflect.TypeOf(res) == nil {
		res := new(successResponse)
		res.Code = "S001"
		return json.Marshal(res)
	}
	return json.Marshal(res)
}

func ResponseIsSuccess(code string) bool {
	return code == "S001"
}