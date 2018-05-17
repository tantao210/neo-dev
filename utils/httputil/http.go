package httputil

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"neo-dev/utils/errorutil"
	"net/http"
	"strings"
	// "fmt"
)

var (
	PostMethod = "POST"
	GetMethod  = "GET"
)

func BuildParams(params map[string]string) string {
	var args string
	for key, value := range params {
		args += key + "=" + value + "&"
	}
	if len(args) > 0 {
		return args[0 : len(args)-1]
	}
	return args
}

func HTTPRequest(url string, method string, head map[string]string, data interface{}) ([]byte, error) {
	client := &http.Client{}
	var (
		req     *http.Request
		reqData []byte
		err     error
	)
	if method == PostMethod {
		reqData, err = json.Marshal(data)
		if err != nil {
			return nil, errors.New(errorutil.System_Error)
		}
		req, err = http.NewRequest(method, url, strings.NewReader(string(reqData)))
	} else {
		params := data.(map[string]string)
		url = url + "?" + BuildParams(params)
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	for key, value := range head {
		req.Header.Add(key, value)
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(reqData))
	return resp, nil
}
