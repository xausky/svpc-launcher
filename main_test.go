package main

import (
	"github.com/gonutz/w32/v2"
	"image/png"
	"os"
	"testing"
)

func TestScanQR(t *testing.T) {
	hWnd := w32.FindWindow("MPAY_LOGIN", "登录")
	if hWnd == 0 {
		t.Fatal("window not found")
	}
	shot := TakeScreenshot(hWnd)
	open, _ := os.OpenFile("test.png", os.O_CREATE|os.O_TRUNC, 0644)
	png.Encode(open, shot)
	params, err := ScanQR(shot)
	if err != nil || !params.Has("uuid") {
		t.Fatal(err)
	}
	t.Log(params)
}
