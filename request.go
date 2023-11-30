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

	url := "https://service.mkey.163.com/mpay/api/qrcode/confirm_login"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("uuid=%v&device_id=amawkzyaasy32xwt-d&token=%v&jf_game_id=ma68&pay_channel=netease&is_remember=0&game_id=aecfrugltuaaaajo-g-ma68&gv=120&gvn=4.3.20&cv=a3.29.0&sv=33&app_type=games&app_mode=2&app_channel=netease&mcount_app_key=EEkEEXLymcNjM42yLY3Bn6AO15aGy4yq&mcount_transaction_id=4514532b-8ec3-11ee-9a36-4b696ee8d233-1", uuid, values.Get("token")))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept-Charset", "UTF-8")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
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
	fmt.Println(string(body))
	if res.StatusCode == 200 && gjson.Get(string(body), "code").Int() == 0 {
		return true
	} else {
		return false
	}
}

func ScanRequestSend(uuid string, values url.Values) bool {
	url := fmt.Sprintf("https://service.mkey.163.com/mpay/api/qrcode/scan?uuid=%v&cv=a3.29.0&game_id=aecfrugltuaaaajo-g-ma68&gv=120&gvn=4.3.20&sv=33&app_type=games&app_mode=2&app_channel=netease&mcount_app_key=EEkEEXLymcNjM42yLY3Bn6AO15aGy4yq&mcount_transaction_id=4514532b-8ec3-11ee-9a36-4b696ee8d233-1", uuid)

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

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
	fmt.Println(string(body))
	return true
}
