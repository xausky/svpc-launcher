package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func SendConfirmLogin(uuid string, values url.Values) bool {

	uri := "https://service.mkey.163.com/mpay/api/qrcode/confirm_login"
	method := "POST"

	requestBodyParams, err := url.ParseQuery("uuid=xxxxx&device_id=amawkzyaasy32xwt-d&token=xxxxx&jf_game_id=ma68&pay_channel=netease&is_remember=0&game_id=aecfrugltuaaaajo-g-ma68&gv=120&gvn=4.3.20&cv=a3.29.0&sv=33&app_type=games&app_mode=2&app_channel=netease&mcount_app_key=EEkEEXLymcNjM42yLY3Bn6AO15aGy4yq&mcount_transaction_id=4514532b-8ec3-11ee-9a36-4b696ee8d233-1")
	if err != nil {
		panic(err)
	}
	for k, _ := range requestBodyParams {
		if values.Has(k) && requestBodyParams.Get(k) != values.Get(k) {
			fmt.Println("replace", k, requestBodyParams.Get(k), values.Get(k))
			requestBodyParams.Set(k, values.Get(k))
		}
	}
	requestBodyParams.Set("uuid", uuid)

	payload := strings.NewReader(requestBodyParams.Encode())

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept-Charset", "UTF-8")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("POST", uri)
	fmt.Println(requestBodyParams.Encode())
	fmt.Println(string(body))
	if res.StatusCode == 200 && gjson.Get(string(body), "code").Int() == 0 {
		return true
	} else {
		return false
	}
}

func ScanRequestSend(uuid string, values url.Values) bool {
	uri := "https://service.mkey.163.com/mpay/api/qrcode/scan?uuid=xxxxx&cv=a3.29.0&game_id=aecfrugltuaaaajo-g-ma68&gv=120&gvn=4.3.20&sv=33&app_type=games&app_mode=2&app_channel=netease&mcount_app_key=EEkEEXLymcNjM42yLY3Bn6AO15aGy4yq&mcount_transaction_id=4514532b-8ec3-11ee-9a36-4b696ee8d233-1"

	reqUrl, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	queryParams := reqUrl.Query()

	for k, _ := range queryParams {
		if values.Has(k) && queryParams.Get(k) != values.Get(k) {
			fmt.Println("replace", k, queryParams.Get(k), values.Get(k))
			queryParams.Set(k, values.Get(k))
		}
	}

	queryParams.Set("uuid", uuid)
	reqUrl.RawQuery = queryParams.Encode()

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, reqUrl.String(), nil)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept-Charset", "UTF-8")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("GET", reqUrl.String())
	fmt.Println(string(body))
	return true
}
