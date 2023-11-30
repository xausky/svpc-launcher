package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/makiuchi-d/gozxing"

	"github.com/gonutz/w32/v2"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func TakeScreenshot(hWnd w32.HWND) image.Image {
	rect := w32.GetWindowRect(hWnd)
	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top
	fmt.Println(width, height)

	hdc := w32.GetDC(hWnd)
	defer w32.ReleaseDC(hWnd, hdc)

	dcDest := w32.CreateCompatibleDC(hdc)
	defer w32.DeleteDC(dcDest)

	hBitmap := w32.CreateCompatibleBitmap(hdc, int(width), int(height))
	defer w32.DeleteObject(w32.HGDIOBJ(hBitmap))
	hOld := w32.SelectObject(dcDest, w32.HGDIOBJ(hBitmap))
	w32.BitBlt(dcDest, 0, 0, int(width), int(height), hdc, 0, 0, w32.SRCCOPY)
	w32.SelectObject(dcDest, hOld)

	bitmapInfo := w32.BITMAPINFO{
		BmiHeader: w32.BITMAPINFOHEADER{
			BiSize:        uint32(unsafe.Sizeof(w32.BITMAPINFOHEADER{})),
			BiPlanes:      1,
			BiBitCount:    32,
			BiWidth:       width,
			BiHeight:      -height,
			BiCompression: w32.BI_RGB,
			BiSizeImage:   uint32(width * height * 4),
		},
	}

	pixels := make([]byte, width*height*4)
	w32.GetDIBits(hdc, hBitmap, 0, uint(height), unsafe.Pointer(&pixels[0]), &bitmapInfo, w32.DIB_RGB_COLORS)

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			offset := (y*int(width) + x) * 4
			img.SetRGBA(x, y, color.RGBA{R: pixels[offset+2], G: pixels[offset+1], B: pixels[offset], A: 255})
		}
	}

	return img
}

func ParserTokenUrl(tokenUrl string) (url.Values, error) {
	parse, err := url.Parse(tokenUrl)
	if err != nil {
		return nil, err
	}
	query, err := url.ParseQuery(parse.RawQuery)
	if err != nil {
		return nil, err
	}
	if !query.Has("token") {
		return nil, errors.New("must has token params")
	}
	return query, nil
}

func InputTokenUrl() string {
	var tokenUrl string
	for {
		fmt.Println("登录信息无效，请重新输入授权地址，获取教程：https://blog.xausky.cn")
		reader := bufio.NewReaderSize(os.Stdin, 65536)
		tokenUrlBytes, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("授权地址有误请严格按照教程执行。", err)
			continue
		}
		tokenUrl = string(tokenUrlBytes)
		_, err = ParserTokenUrl(tokenUrl)
		if err != nil {
			fmt.Println("授权地址有误请严格按照教程执行。", err)
			continue
		} else {
			break
		}
	}
	return tokenUrl
}

func LoadToken() url.Values {
	tokenUrlBytes, err := os.ReadFile("svpc-launcher-token.txt")
	var tokenUrl string
	if os.IsNotExist(err) {
		tokenUrl = InputTokenUrl()
		err = os.WriteFile("svpc-launcher-token.txt", []byte(tokenUrl), 0644)
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}
	tokenUrl = string(tokenUrlBytes)
	tokenUrlParams, err := ParserTokenUrl(tokenUrl)
	if err != nil {
		tokenUrl = InputTokenUrl()
		tokenUrlParams, _ = ParserTokenUrl(tokenUrl)
		err = os.WriteFile("svpc-launcher-token.txt", []byte(tokenUrl), 0644)
		if err != nil {
			panic(err)
		}
	}
	return tokenUrlParams
}

func main() {
	render := qrcode.NewQRCodeReader()
	tokenUrlQuery := LoadToken()
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
		ScanRequestSend(uuid, tokenUrlQuery)
		SendConfirmLogin(uuid, tokenUrlQuery)
	}
}

func SendConfirmLogin(uuid string, values url.Values) {

	url := "https://service.mkey.163.com/mpay/api/qrcode/confirm_login"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("uuid=%v&device_id=amawkzyaasy32xwt-d&token=%v&jf_game_id=ma68&pay_channel=netease&is_remember=0&game_id=aecfrugltuaaaajo-g-ma68&gv=120&gvn=4.3.20&cv=a3.29.0&sv=33&app_type=games&app_mode=2&app_channel=netease&mcount_app_key=EEkEEXLymcNjM42yLY3Bn6AO15aGy4yq&mcount_transaction_id=4514532b-8ec3-11ee-9a36-4b696ee8d233-1", uuid, values.Get("token")))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept-Charset", "UTF-8")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func ScanRequestSend(uuid string, values url.Values) {
	url := fmt.Sprintf("https://service.mkey.163.com/mpay/api/qrcode/scan?uuid=%v&cv=a3.29.0&game_id=aecfrugltuaaaajo-g-ma68&gv=120&gvn=4.3.20&sv=33&app_type=games&app_mode=2&app_channel=netease&mcount_app_key=EEkEEXLymcNjM42yLY3Bn6AO15aGy4yq&mcount_transaction_id=4514532b-8ec3-11ee-9a36-4b696ee8d233-1", uuid)

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept-Charset", "UTF-8")
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
