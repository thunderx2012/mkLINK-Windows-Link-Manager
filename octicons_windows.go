//go:build windows

package main

import (
	"bytes"
	_ "embed"
	"image/png"
	"sync"
	"syscall"
	"unsafe"
)

const (
	DIB_RGB_COLORS = 0
	BI_RGB         = 0
	AC_SRC_OVER    = 0
	AC_SRC_ALPHA   = 1
)

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
	Blue     byte
	Green    byte
	Red      byte
	Reserved byte
}

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors [1]RGBQUAD
}

type BLENDFUNCTION struct {
	BlendOp             byte
	BlendFlags          byte
	SourceConstantAlpha byte
	AlphaFormat         byte
}

type octiconKey struct {
	name  string
	color uint32
}

type octiconBitmap struct {
	bmp    uintptr
	width  int32
	height int32
}

var (
	msimg32       = syscall.NewLazyDLL("msimg32.dll")
	pAlphaBlend   = msimg32.NewProc("AlphaBlend")
	pCreateDIBSec = gdi32.NewProc("CreateDIBSection")

	octiconMu    sync.Mutex
	octiconCache = map[octiconKey]octiconBitmap{}

	// The PNG assets are used as antialiased alpha masks for the Win32/GDI renderer.
	//go:embed assets/octicons/check-circle-24.png
	octCheckCircle []byte
	//go:embed assets/octicons/eye-24.png
	octEye []byte
	//go:embed assets/octicons/file-added-24.png
	octFileAdded []byte
	//go:embed assets/octicons/file-directory-24.png
	octFolder []byte
	//go:embed assets/octicons/file-directory-symlink-24.png
	octFolderSymlink []byte
	//go:embed assets/octicons/file-symlink-file-24.png
	octFileSymlink []byte
	//go:embed assets/octicons/gear-24.png
	octGear []byte
	//go:embed assets/octicons/link.png
	octLink16 []byte
	//go:embed assets/octicons/paste-24.png
	octPaste []byte
	//go:embed assets/octicons/sync-24.png
	octSync []byte
	//go:embed assets/octicons/trash-24.png
	octTrash []byte
	//go:embed assets/octicons/x-24.png
	octX []byte
	//go:embed assets/octicons/shield-lock-24.png
	octShieldLock []byte
)

func initOcticonIndex() map[string][]byte {
	return map[string][]byte{
		"check-circle":           octCheckCircle,
		"eye":                    octEye,
		"file-added":             octFileAdded,
		"file-directory":         octFolder,
		"file-directory-symlink": octFolderSymlink,
		"file-symlink-file":      octFileSymlink,
		"gear":                   octGear,
		"link":                   octLink16,
		"paste":                  octPaste,
		"sync":                   octSync,
		"trash":                  octTrash,
		"x":                      octX,
		"shield-lock":            octShieldLock,
	}
}

var octiconSources = initOcticonIndex()

func octiconSize(name string) (int32, int32) {
	if name == "link" {
		return 16, 16
	}
	return 24, 24
}

func decodeOcticonMask(data []byte, width, height int32, c uint32, hdc uintptr) (octiconBitmap, bool) {
	im, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return octiconBitmap{}, false
	}
	bounds := im.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return octiconBitmap{}, false
	}
	var bi BITMAPINFO
	bi.BmiHeader.BiSize = uint32(unsafe.Sizeof(BITMAPINFOHEADER{}))
	bi.BmiHeader.BiWidth = int32(w)
	bi.BmiHeader.BiHeight = -int32(h) // top-down DIB: matches image coordinates
	bi.BmiHeader.BiPlanes = 1
	bi.BmiHeader.BiBitCount = 32
	bi.BmiHeader.BiCompression = BI_RGB
	var bits uintptr
	bmp, _, _ := pCreateDIBSec.Call(
		hdc,
		uintptr(unsafe.Pointer(&bi)),
		uintptr(DIB_RGB_COLORS),
		uintptr(unsafe.Pointer(&bits)),
		0,
		0,
	)
	if bmp == 0 || bits == 0 {
		return octiconBitmap{}, false
	}

	// DIB memory is BGRA. AlphaBlend with AC_SRC_ALPHA expects premultiplied
	// RGB, so premultiply the requested foreground color by the mask alpha.
	r := uint32(c & 0xff)
	g := uint32((c >> 8) & 0xff)
	b := uint32((c >> 16) & 0xff)
	buf := unsafe.Slice((*byte)(unsafe.Pointer(bits)), w*h*4)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			rr, gg, bb, aa := im.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			_ = rr
			_ = gg
			_ = bb
			alpha := uint32(aa >> 8)
			i := (y*w + x) * 4
			buf[i+0] = byte((b*alpha + 127) / 255)
			buf[i+1] = byte((g*alpha + 127) / 255)
			buf[i+2] = byte((r*alpha + 127) / 255)
			buf[i+3] = byte(alpha)
		}
	}
	return octiconBitmap{bmp: bmp, width: int32(w), height: int32(h)}, true
}

func getOcticon(name string, c uint32, hdc uintptr) (octiconBitmap, bool) {
	data, ok := octiconSources[name]
	if !ok || len(data) == 0 {
		return octiconBitmap{}, false
	}
	key := octiconKey{name: name, color: c}
	octiconMu.Lock()
	defer octiconMu.Unlock()
	if bm, ok := octiconCache[key]; ok && bm.bmp != 0 {
		return bm, true
	}
	w, h := octiconSize(name)
	bm, ok := decodeOcticonMask(data, w, h, c, hdc)
	if !ok {
		return octiconBitmap{}, false
	}
	octiconCache[key] = bm
	return bm, true
}

func drawOcticon(hdc uintptr, box RECT, name string, c uint32) {
	bm, ok := getOcticon(name, c, hdc)
	if !ok || bm.bmp == 0 || bm.width <= 0 || bm.height <= 0 {
		return
	}
	memDC, _, _ := pCreateCompatibleDC.Call(hdc)
	if memDC == 0 {
		return
	}
	defer pDeleteDC.Call(memDC)
	old, _, _ := pSelectObject.Call(memDC, bm.bmp)
	defer pSelectObject.Call(memDC, old)

	x := (box.Left + box.Right - bm.width) / 2
	y := (box.Top + box.Bottom - bm.height) / 2
	bf := BLENDFUNCTION{BlendOp: AC_SRC_OVER, BlendFlags: 0, SourceConstantAlpha: 255, AlphaFormat: AC_SRC_ALPHA}
	pAlphaBlend.Call(
		hdc,
		uintptr(x), uintptr(y), uintptr(bm.width), uintptr(bm.height),
		memDC,
		0, 0, uintptr(bm.width), uintptr(bm.height),
		uintptr(*(*uint32)(unsafe.Pointer(&bf))),
	)
}

func releaseOcticons() {
	octiconMu.Lock()
	defer octiconMu.Unlock()
	for k, bm := range octiconCache {
		if bm.bmp != 0 {
			pDeleteObject.Call(bm.bmp)
		}
		delete(octiconCache, k)
	}
}

func drawFolderIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "file-directory", c)
}

func drawFileIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "file-added", c)
}

func drawClipboardIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "paste", c)
}

func drawRefreshIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "sync", c)
}

func drawTrashIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "trash", c)
}

func drawClearIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "x", c)
}

func drawGearIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "gear", c)
}

func drawEyeIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "eye", c)
}

func drawStatusIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "check-circle", c)
}

func drawCommandIcon(hdc uintptr, r RECT, idx int, c uint32) {
	var name string
	switch idx {
	case 0:
		name = "file-added"
	case 1:
		name = "file-directory"
	case 2:
		name = "x"
	case 3:
		name = "shield-lock"
	default:
		return
	}
	drawOcticon(hdc, r, name, c)
}

func drawModeIcon(hdc uintptr, r RECT, idx int, c uint32) {
	switch idx {
	case 0:
		drawOcticon(hdc, r, "file-symlink-file", c)
	case 1:
		drawOcticon(hdc, r, "link", c)
	case 2:
		drawOcticon(hdc, r, "file-directory-symlink", c)
	}
}

func drawChainIcon(hdc uintptr, r RECT, c uint32) {
	drawOcticon(hdc, r, "link", c)
}

func drawDropIcon(hdc uintptr, r RECT, c uint32) {
	cx := (r.Left + r.Right) / 2
	drawOcticon(hdc, RECT{cx - 22, r.Top + 3, cx + 2, r.Top + 27}, "file-added", c)
	drawOcticon(hdc, RECT{cx - 1, r.Top + 20, cx + 23, r.Top + 44}, "file-directory", c)
}
