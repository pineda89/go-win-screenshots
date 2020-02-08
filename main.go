package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

func main() {
	// GetDC(0) is equals than GetDC(NULL). This will take screenshot of entire screen
	screen_hdc = GetDC(0)
	screen_horzres = GetDeviceCaps(screen_hdc, HORZRES)
	screen_vertres = GetDeviceCaps(screen_hdc, VERTRES)
	screen_rectangle = image.Rect(0, 0, screen_horzres, screen_vertres)

	img, _ := CaptureScreen()
	myImg := image.Image(img)
	outputfilename := strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg"
	imgToFile(&myImg, outputfilename)
	checkPixelColor(&myImg, 230, 480)
}

func checkPixelColor(i *image.Image, x int, y int) {
	color := (*i).At(x, y)
	r, g, b, _ := color.RGBA()
	nr := r * 255 / 65535
	ng := g * 255 / 65535
	nb := b * 255 / 65535
	fmt.Println(x, y, "color is", nr, ng, nb)
}

var (
	screen_hdc       HDC
	compatibleDC     HDC
	screen_horzres   int
	screen_vertres   int
	screen_rectangle image.Rectangle
)

func CaptureScreen() (*image.RGBA, error) {
	x, y := screen_rectangle.Dx(), screen_rectangle.Dy()

	bt := BITMAPINFO{}
	bt.BmiHeader.BiSize = uint32(reflect.TypeOf(bt.BmiHeader).Size())
	bt.BmiHeader.BiWidth = int32(x)
	bt.BmiHeader.BiHeight = int32(-y)
	bt.BmiHeader.BiPlanes = 1
	bt.BmiHeader.BiBitCount = 32
	bt.BmiHeader.BiCompression = BI_RGB

	ptr := unsafe.Pointer(uintptr(0))

	compatibleDC = CreateCompatibleDC(screen_hdc)

	DIBSection := CreateDIBSection(compatibleDC, &bt, DIB_RGB_COLORS, &ptr, 0, 0)
	if DIBSection == 0 {
		return nil, fmt.Errorf("Could not Create DIB Section err:%d.\n", GetLastError())
	}
	if DIBSection == InvalidParameter {
		return nil, fmt.Errorf("CreateDIBSection invalid params\n")
	}
	defer DeleteObject(HGDIOBJ(DIBSection))

	obj := SelectObject(compatibleDC, HGDIOBJ(DIBSection))
	if obj == 0 {
		return nil, fmt.Errorf("SelectObject err:%d.\n", GetLastError())
	}
	if obj == 0xffffffff { //GDI_ERROR
		return nil, fmt.Errorf("GDI_ERROR err:%d.\n", GetLastError())
	}
	defer DeleteObject(obj)

	BitBlt(compatibleDC, 0, 0, x, y, screen_hdc, screen_rectangle.Min.X, screen_rectangle.Min.Y, SRCCOPY)

	DeleteDC(compatibleDC)

	var slice []byte
	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	hdrp.Data = uintptr(ptr)
	hdrp.Len = x * y * 4
	hdrp.Cap = x * y * 4

	imageBytes := make([]byte, len(slice))

	for i := 0; i < len(imageBytes); i += 4 {
		imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
	}

	img := &image.RGBA{imageBytes, 4 * x, image.Rect(0, 0, x, y)}
	return img, nil
}

func imgToFile(img *image.Image, filename string) {
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, *img, nil)
	send_s3 := buf.Bytes()
	ioutil.WriteFile(filename, send_s3, 0777)
}
