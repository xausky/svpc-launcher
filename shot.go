package main

import (
	"fmt"
	"github.com/gonutz/w32/v2"
	"image"
	"image/color"
	"unsafe"
)

func TakeScreenshot(hWnd w32.HWND) image.Image {
	s := float64(w32.GetDeviceCaps(w32.GetDC(0), w32.DESKTOPHORZRES)) / float64(w32.GetSystemMetrics(0))
	rect := w32.GetWindowRect(hWnd)
	width := int32(float64(rect.Right-rect.Left) * s)
	height := int32(float64(rect.Bottom-rect.Top) * s)
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
