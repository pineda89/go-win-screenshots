package main

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	modgdi32               = syscall.NewLazyDLL("gdi32.dll")
	moduser32              = syscall.NewLazyDLL("user32.dll")
	modkernel32            = syscall.NewLazyDLL("kernel32.dll")
	procGetDC              = moduser32.NewProc("GetDC")
	procReleaseDC          = moduser32.NewProc("ReleaseDC")
	procDeleteDC           = modgdi32.NewProc("DeleteDC")
	procBitBlt             = modgdi32.NewProc("BitBlt")
	procDeleteObject       = modgdi32.NewProc("DeleteObject")
	procSelectObject       = modgdi32.NewProc("SelectObject")
	procCreateDIBSection   = modgdi32.NewProc("CreateDIBSection")
	procCreateCompatibleDC = modgdi32.NewProc("CreateCompatibleDC")
	procGetDeviceCaps      = modgdi32.NewProc("GetDeviceCaps")
	procGetLastError       = modkernel32.NewProc("GetLastError")
)

const (
	BI_RGB           = 0
	HORZRES          = 8
	VERTRES          = 10
	DIB_RGB_COLORS   = 0
	InvalidParameter = 2
	SRCCOPY          = 0x00CC0020
)

type (
	HANDLE  uintptr
	HWND    HANDLE
	HGDIOBJ HANDLE
	HDC     HANDLE
	HBITMAP HANDLE
)

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors *RGBQUAD
}

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

func GetDC(hwnd HWND) HDC {
	ret, _, _ := procGetDC.Call(uintptr(hwnd))

	return HDC(ret)
}

func GetDeviceCaps(hdc HDC, index int) int {
	ret, _, _ := procGetDeviceCaps.Call(uintptr(hdc), uintptr(index))

	return int(ret)
}

func ReleaseDC(hwnd HWND, hDC HDC) bool {
	ret, _, _ := procReleaseDC.Call(uintptr(hwnd), uintptr(hDC))

	return ret != 0
}

func DeleteDC(hdc HDC) bool {
	ret, _, _ := procDeleteDC.Call(uintptr(hdc))
	return ret != 0
}

func CreateCompatibleDC(hdc HDC) HDC {
	ret, _, _ := procCreateCompatibleDC.Call(
		uintptr(hdc))

	return HDC(ret)
}

func SelectObject(hdc HDC, hgdiobj HGDIOBJ) HGDIOBJ {
	ret, _, _ := procSelectObject.Call(uintptr(hdc), uintptr(hgdiobj))
	return HGDIOBJ(ret)
}

func BitBlt(hdcDest HDC, nXDest, nYDest, nWidth, nHeight int, hdcSrc HDC, nXSrc, nYSrc int, dwRop uint) error {
	ret, _, _ := procBitBlt.Call(uintptr(hdcDest), uintptr(nXDest), uintptr(nYDest), uintptr(nWidth), uintptr(nHeight), uintptr(hdcSrc), uintptr(nXSrc), uintptr(nYSrc), uintptr(dwRop))

	if ret == 0 {
		return errors.New("BitBlt error")
	}

	return nil
}

func GetLastError() uint32 {
	ret, _, _ := procGetLastError.Call()
	return uint32(ret)
}

func DeleteObject(hObject HGDIOBJ) bool {
	ret, _, _ := procDeleteObject.Call(uintptr(hObject))

	return ret != 0
}

func CreateDIBSection(hdc HDC, pbmi *BITMAPINFO, iUsage uint, ppvBits *unsafe.Pointer, hSection HANDLE, dwOffset uint) HBITMAP {
	ret, _, _ := procCreateDIBSection.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(pbmi)),
		uintptr(iUsage),
		uintptr(unsafe.Pointer(ppvBits)),
		uintptr(hSection),
		uintptr(dwOffset))

	return HBITMAP(ret)
}
