package main

import (
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"net/url"
	"os/exec"
	"time"

	"github.com/gonutz/w32/v2"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func main() {
	render := qrcode.NewQRCodeReader()
	tokenUrlQuery := LoadToken(false)
	err := exec.Command("./影之诗.exe").Start()
	if err != nil {
		fmt.Println(err)
	}
	for {
		time.Sleep(time.Second)
		hWnd := w32.FindWindow("MPAY_LOGIN", "登录")
		if hWnd == 0 {
			fmt.Println("等待登录窗口出现。")
			continue
		}
		shot := TakeScreenshot(hWnd)
		bitmap, err := gozxing.NewBinaryBitmapFromImage(shot)
		if err != nil {
			fmt.Println("二维码解码错误，等待重试。", err)
			return
		}
		result, err := render.Decode(bitmap, nil)
		if err != nil {
			fmt.Println("二维码解码错误，等待重试。", err)
			continue
		}
		scanUrl, err := url.Parse(result.GetText())
		if err != nil {
			fmt.Println("二维码解码错误，等待重试。", err)
			continue
		}
		scanQueryParams, err := url.ParseQuery(scanUrl.RawQuery)
		if err != nil {
			fmt.Println("登录URL解析错误，等待重试。", err)
			return
		}
		uuid := scanQueryParams.Get("uuid")
		fmt.Println("二维码扫码成功，UUID: ", uuid)
		ScanRequestSend(uuid, tokenUrlQuery)
		if SendConfirmLogin(uuid, tokenUrlQuery) {
			break
		} else {
			tokenUrlQuery = LoadToken(true)
		}
	}
}
