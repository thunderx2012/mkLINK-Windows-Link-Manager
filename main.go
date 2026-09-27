//go:build windows

package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

const (
	SRCCOPY    = 0x00CC0020
	appTitle   = "mkLINK"
	appVersion = "1.5.22"
	creatorURL = "https://www.threads.com/@thunderx2012"

	WS_OVERLAPPEDWINDOW = 0x00CF0000
	WS_CHILD            = 0x40000000
	WS_VISIBLE          = 0x10000000
	WS_TABSTOP          = 0x00010000
	WS_VSCROLL          = 0x00200000
	SBS_VERT            = 0x0001
	SIF_RANGE           = 0x0001
	SIF_PAGE            = 0x0002
	SIF_POS             = 0x0004
	SIF_TRACKPOS        = 0x0010
	SB_CTL              = 0x0002
	SB_LINEUP           = 0
	SB_LINEDOWN         = 1
	SB_PAGEUP           = 2
	SB_PAGEDOWN         = 3
	SB_LINELEFT         = SB_LINEUP
	SB_LINERIGHT        = SB_LINEDOWN
	SB_PAGELEFT         = SB_PAGEUP
	SB_PAGERIGHT        = SB_PAGEDOWN
	SB_THUMBPOSITION    = 4
	SB_THUMBTRACK       = 5
	SB_TOP              = 6
	SB_BOTTOM           = 7
	MF_STRING           = 0x0000
	MF_GRAYED           = 0x0001
	MF_CHECKED          = 0x0008
	MF_SEPARATOR        = 0x0800
	TPM_RIGHTBUTTON     = 0x0002
	TPM_RETURNCMD       = 0x0100
	IDM_PASTE_PATH      = 4101
	IDM_REMOVE_SELECTED = 4102
	IDM_LANG_ZHTW       = 4201
	IDM_LANG_ENUS       = 4202
	IDM_LANG_JAJP       = 4203
	IDM_LANG_ZHCN       = 4204
	IDM_LANG_CUSTOM     = 4205
	ES_AUTOHSCROLL      = 0x0080
	ES_LEFT             = 0x0000
	ES_MULTILINE        = 0x0004
	ES_AUTOVSCROLL      = 0x0040
	SBS_HORZ            = 0x0000
	SB_HORZ             = 0x0000
	IDC_SIZEWE          = 32644

	WM_CREATE        = 0x0001
	WM_DESTROY       = 0x0002
	WM_SIZE          = 0x0005
	WM_SETFOCUS      = 0x0007
	WM_KILLFOCUS     = 0x0008
	WM_CLOSE         = 0x0010
	WM_KEYDOWN       = 0x0100
	WM_CHAR          = 0x0102
	WM_MOUSEWHEEL    = 0x020A
	WM_RBUTTONUP     = 0x0205
	WM_HSCROLL       = 0x0114
	WM_VSCROLL       = 0x0115
	WM_LBUTTONDOWN   = 0x0201
	WM_LBUTTONUP     = 0x0202
	WM_LBUTTONDBLCLK = 0x0203
	WM_MOUSEMOVE     = 0x0200
	WM_DROPFILES     = 0x0233
	WM_PAINT         = 0x000F
	WM_ERASEBKGND    = 0x0014
	WM_GETMINMAXINFO = 0x0024
	WM_DPICHANGED    = 0x02E0
	WM_CTLCOLOREDIT  = 0x0133
	WM_COMMAND       = 0x0111
	WM_SETFONT       = 0x0030
	WM_GETTEXT       = 0x000D
	WM_GETTEXTLENGTH = 0x000E
	WM_SETTEXT       = 0x000C
	WM_SETCUEBANNER  = 0x1501
	EM_SETSEL        = 0x00B1
	EM_REPLACESEL    = 0x00C2
	BM_GETCHECK      = 0x00F0
	BN_CLICKED       = 0
	EN_CHANGE        = 0x0300
	WM_SETCURSOR     = 0x0020
	HTCLIENT         = 1
	IDC_HAND         = 32649

	VK_CONTROL = 0x11
	VK_SHIFT   = 0x10
	VK_A       = 0x41
	VK_V       = 0x56
	VK_DELETE  = 0x2E
	VK_ESCAPE  = 0x1B
	VK_RETURN  = 0x0D
	VK_UP      = 0x26
	VK_DOWN    = 0x28
	VK_HOME    = 0x24
	VK_END     = 0x23

	MK_LBUTTON = 0x0001
	MK_SHIFT   = 0x0004
	MK_CONTROL = 0x0008

	SW_SHOW       = 5
	SW_HIDE       = 0
	SW_SHOWNORMAL = 1
	CW_USEDEFAULT = 0x80000000
	IDC_ARROW     = 32512
	DI_NORMAL     = 0x0003

	DT_LEFT         = 0x00000000
	DT_CENTER       = 0x00000001
	DT_RIGHT        = 0x00000002
	DT_VCENTER      = 0x00000004
	DT_SINGLELINE   = 0x00000020
	DT_END_ELLIPSIS = 0x00008000
	DT_NOPREFIX     = 0x00000800
	DT_WORDBREAK    = 0x00000010
	DT_CALCRECT     = 0x00000400

	PS_SOLID    = 0
	TRANSPARENT = 1
	NULL_BRUSH  = 5
	FW_NORMAL   = 400
	FW_SEMIBOLD = 600
	FW_BOLD     = 700
	TRUE        = 1

	OFN_EXPLORER         = 0x00080000
	OFN_FILEMUSTEXIST    = 0x00001000
	OFN_ALLOWMULTISELECT = 0x00000200
	OFN_PATHMUSTEXIST    = 0x00000800
	OFN_HIDEREADONLY     = 0x00000004
	FNERR_BUFFERTOOSMALL = 0x3003

	BIF_RETURNONLYFSDIRS = 0x00000001
	BIF_NEWDIALOGSTYLE   = 0x00000040
	BIF_EDITBOX          = 0x00000010
	BIF_VALIDATE         = 0x00000020

	CF_UNICODETEXT = 13

	GENERIC_WRITE                          = 0x40000000
	FILE_SHARE_READ                        = 0x00000001
	FILE_SHARE_WRITE                       = 0x00000002
	FILE_SHARE_DELETE                      = 0x00000004
	OPEN_EXISTING                          = 3
	FILE_FLAG_OPEN_REPARSE_POINT           = 0x00200000
	FILE_FLAG_BACKUP_SEMANTICS             = 0x02000000
	INVALID_HANDLE_VALUE                   = ^uintptr(0)
	FSCTL_SET_REPARSE_POINT                = 0x000900A4
	IO_REPARSE_TAG_MOUNT_POINT             = 0xA0000003
	SYMBOLIC_LINK_FLAG_DIRECTORY           = 0x1
	SYMBOLIC_LINK_FLAG_UNPRIVILEGED_CREATE = 0x2
	FORMAT_MESSAGE_FROM_SYSTEM             = 0x00001000
	FORMAT_MESSAGE_IGNORE_INSERTS          = 0x00000200
	TOKEN_QUERY                            = 0x0008
	TokenElevation                         = 20
	MB_YESNO                               = 0x00000004
	MB_ICONWARNING                         = 0x00000030
	IDYES                                  = 6

	GWLP_WNDPROC = ^uintptr(3)
)

func rgb(r, g, b byte) uint32 { return uint32(r) | uint32(g)<<8 | uint32(b)<<16 }

type HWND uintptr
type HINSTANCE uintptr
type HICON uintptr
type HBRUSH uintptr
type HFONT uintptr
type HDC uintptr

type RECT struct{ Left, Top, Right, Bottom int32 }
type POINT struct{ X, Y int32 }
type TOKEN_ELEVATION struct{ TokenIsElevated uint32 }
type MSG struct {
	Hwnd           HWND
	Message        uint32
	WParam, LParam uintptr
	Time           uint32
	PtX, PtY       int32
	LPrivate       uint32
}
type PAINTSTRUCT struct {
	Hdc         HDC
	FErase      int32
	RcPaint     RECT
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}
type MINMAXINFO struct {
	PtReserved     POINT
	PtMaxSize      POINT
	PtMaxPosition  POINT
	PtMinTrackSize POINT
	PtMaxTrackSize POINT
}
type SCROLLINFO struct {
	CbSize    uint32
	FMask     uint32
	NMin      int32
	NMax      int32
	NPage     uint32
	NPos      int32
	NTrackPos int32
}
type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         HICON
	HCursor       uintptr
	HbrBackground HBRUSH
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       HICON
}
type OPENFILENAMEW struct {
	LStructSize       uint32
	HwndOwner         HWND
	HInstance         HINSTANCE
	LpstrFilter       *uint16
	LpstrCustomFilter *uint16
	NMaxCustFilter    uint32
	NFilterIndex      uint32
	LpstrFile         *uint16
	NMaxFile          uint32
	LpstrFileTitle    *uint16
	NMaxFileTitle     uint32
	LpstrInitialDir   *uint16
	LpstrTitle        *uint16
	Flags             uint32
	NFileOffset       uint16
	NFileExtension    uint16
	LpstrDefExt       *uint16
	LCustData         uintptr
	LpfnHook          uintptr
	LpTemplateName    *uint16
	PvReserved        uintptr
	DwReserved        uint32
	FlagsEx           uint32
}
type BROWSEINFOW struct {
	HwndOwner      HWND
	PidlRoot       uintptr
	PszDisplayName *uint16
	LpszTitle      *uint16
	UlFlags        uint32
	Lpfn           uintptr
	LParam         uintptr
	IImage         int32
}

type SourceItem struct {
	Path, Kind string
	Selected   bool
}

type App struct {
	hwnd                                                                                                            HWND
	destination                                                                                                     HWND
	dpi                                                                                                             int32
	uiScale                                                                                                         float64
	icon                                                                                                            HICON
	font, titleFont, sectionFont, bodyFont, smallFont, monoFont, buttonFont                                         HFONT
	items                                                                                                           []SourceItem
	sourceScrollbar                                                                                                 HWND
	sourceHScrollbar                                                                                                HWND
	languageRect                                                                                                    RECT
	sourceFocus                                                                                                     bool
	selectedType                                                                                                    int
	scroll                                                                                                          int
	hoverID                                                                                                         int
	pressedID                                                                                                       int
	anchorIndex                                                                                                     int
	languageLocale                                                                                                  string
	languageFile                                                                                                    string
	lang, fallbackLang                                                                                              LanguagePack
	columnOrder                                                                                                     [3]int
	columnWidths                                                                                                    [3]int32
	hScroll                                                                                                         int32
	columnDragMode                                                                                                  int
	columnDragCol                                                                                                   int
	columnDragStartX                                                                                                int32
	columnDragOrigWidth                                                                                             int32
	columnDragFrom, columnDragTo                                                                                    int
	creatorHover                                                                                                    bool
	compatOK                                                                                                        bool
	status                                                                                                          string
	sourceRect, sourceListArea, modeRect, destRect, compatRect, previewRect, optionsRect, destEditRect, creatorRect RECT
	toolbarRects                                                                                                    []RECT
	modeRects                                                                                                       [3]RECT
	browseRect, createRect                                                                                          RECT
	pfOK                                                                                                            bool
	pfIssues, pfCmds                                                                                                []string
	pfItems                                                                                                         []SourceItem
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	advapi32 = syscall.NewLazyDLL("advapi32.dll")

	pOpenProcessToken    = advapi32.NewProc("OpenProcessToken")
	pGetTokenInformation = advapi32.NewProc("GetTokenInformation")

	pRegisterClassExW              = user32.NewProc("RegisterClassExW")
	pCreateWindowExW               = user32.NewProc("CreateWindowExW")
	pDefWindowProcW                = user32.NewProc("DefWindowProcW")
	pDestroyWindow                 = user32.NewProc("DestroyWindow")
	pShowWindow                    = user32.NewProc("ShowWindow")
	pUpdateWindow                  = user32.NewProc("UpdateWindow")
	pGetMessageW                   = user32.NewProc("GetMessageW")
	pTranslateMessage              = user32.NewProc("TranslateMessage")
	pDispatchMessageW              = user32.NewProc("DispatchMessageW")
	pPostQuitMessage               = user32.NewProc("PostQuitMessage")
	pMoveWindow                    = user32.NewProc("MoveWindow")
	pSetScrollInfo                 = user32.NewProc("SetScrollInfo")
	pGetScrollInfo                 = user32.NewProc("GetScrollInfo")
	pCreatePopupMenu               = user32.NewProc("CreatePopupMenu")
	pAppendMenuW                   = user32.NewProc("AppendMenuW")
	pTrackPopupMenu                = user32.NewProc("TrackPopupMenu")
	pDestroyMenu                   = user32.NewProc("DestroyMenu")
	pGetClientRect                 = user32.NewProc("GetClientRect")
	pGetDpiForWindow               = user32.NewProc("GetDpiForWindow")
	pGetAsyncKeyState              = user32.NewProc("GetAsyncKeyState")
	pSetFocusProc                  = user32.NewProc("SetFocus")
	pSetCapture                    = user32.NewProc("SetCapture")
	pReleaseCapture                = user32.NewProc("ReleaseCapture")
	pMessageBoxW                   = user32.NewProc("MessageBoxW")
	pLoadCursorW                   = user32.NewProc("LoadCursorW")
	pLoadIconW                     = user32.NewProc("LoadIconW")
	pSendMessageW                  = user32.NewProc("SendMessageW")
	pInvalidateRect                = user32.NewProc("InvalidateRect")
	pBeginPaint                    = user32.NewProc("BeginPaint")
	pEndPaint                      = user32.NewProc("EndPaint")
	pFillRect                      = user32.NewProc("FillRect")
	pSetTextColor                  = gdi32.NewProc("SetTextColor")
	pSetBkColor                    = gdi32.NewProc("SetBkColor")
	pSetBkMode                     = gdi32.NewProc("SetBkMode")
	pCreateSolidBrush              = gdi32.NewProc("CreateSolidBrush")
	pCreatePen                     = gdi32.NewProc("CreatePen")
	pSelectObject                  = gdi32.NewProc("SelectObject")
	pCreateCompatibleDC            = gdi32.NewProc("CreateCompatibleDC")
	pCreateCompatibleBitmap        = gdi32.NewProc("CreateCompatibleBitmap")
	pBitBlt                        = gdi32.NewProc("BitBlt")
	pDeleteDC                      = gdi32.NewProc("DeleteDC")
	pDeleteObject                  = gdi32.NewProc("DeleteObject")
	pRoundRect                     = gdi32.NewProc("RoundRect")
	pRectangle                     = gdi32.NewProc("Rectangle")
	pEllipse                       = gdi32.NewProc("Ellipse")
	pMoveToEx                      = gdi32.NewProc("MoveToEx")
	pLineTo                        = gdi32.NewProc("LineTo")
	pPolygon                       = gdi32.NewProc("Polygon")
	pCreateFontW                   = gdi32.NewProc("CreateFontW")
	pTextOutW                      = gdi32.NewProc("TextOutW")
	pGetStockObject                = gdi32.NewProc("GetStockObject")
	pDrawTextW                     = user32.NewProc("DrawTextW")
	pDrawIconEx                    = user32.NewProc("DrawIconEx")
	pShellExecuteW                 = shell32.NewProc("ShellExecuteW")
	pSetProcessDpiAwarenessContext = user32.NewProc("SetProcessDpiAwarenessContext")
	dwmapi                         = syscall.NewLazyDLL("dwmapi.dll")
	pDwmSetWindowAttribute         = dwmapi.NewProc("DwmSetWindowAttribute")
	pGetClipboardData              = user32.NewProc("GetClipboardData")
	pOpenClipboard                 = user32.NewProc("OpenClipboard")
	pCloseClipboard                = user32.NewProc("CloseClipboard")
	pGlobalLock                    = kernel32.NewProc("GlobalLock")
	pGlobalUnlock                  = kernel32.NewProc("GlobalUnlock")
	pGlobalSize                    = kernel32.NewProc("GlobalSize")
	pRtlMoveMemory                 = kernel32.NewProc("RtlMoveMemory")
	pFormatMessageW                = kernel32.NewProc("FormatMessageW")
	pCreateSymbolicLinkW           = kernel32.NewProc("CreateSymbolicLinkW")
	pCreateHardLinkW               = kernel32.NewProc("CreateHardLinkW")
	pCreateFileW                   = kernel32.NewProc("CreateFileW")
	pCloseHandle                   = kernel32.NewProc("CloseHandle")
	pCreateDirectoryW              = kernel32.NewProc("CreateDirectoryW")
	pRemoveDirectoryW              = kernel32.NewProc("RemoveDirectoryW")
	pDeviceIoControl               = kernel32.NewProc("DeviceIoControl")
	pGetOpenFileNameW              = comdlg32.NewProc("GetOpenFileNameW")
	pCommDlgExtendedError          = comdlg32.NewProc("CommDlgExtendedError")
	pSHBrowseForFolderW            = shell32.NewProc("SHBrowseForFolderW")
	pSHGetPathFromIDListEx         = shell32.NewProc("SHGetPathFromIDListEx")
	pCoTaskMemFree                 = ole32.NewProc("CoTaskMemFree")
	pDragAcceptFiles               = shell32.NewProc("DragAcceptFiles")
	pDragQueryFileW                = shell32.NewProc("DragQueryFileW")
	pDragFinish                    = shell32.NewProc("DragFinish")
	pGetCursorPos                  = user32.NewProc("GetCursorPos")
	pScreenToClient                = user32.NewProc("ScreenToClient")
	pSetCursor                     = user32.NewProc("SetCursor")
)

var hHandCursor uintptr
var hSizeWECursor uintptr

var globalApp *App
var wndProcPtr uintptr
var brInput HBRUSH
var toolbarFont HFONT

// mkLINK v1.5.5 design system: the supplied reference is the visual source of truth.
var (
	cBG       = rgb(8, 18, 31)
	cHeader   = rgb(13, 27, 44)
	cToolbar  = rgb(15, 31, 49)
	cPanel    = rgb(17, 34, 53)
	cSurface  = rgb(13, 27, 43)
	cSurface2 = rgb(21, 41, 62)
	cBorder   = rgb(40, 67, 94)
	cBorder2  = rgb(57, 88, 119)
	cText     = rgb(235, 243, 252)
	cText2    = rgb(181, 199, 220)
	cMuted    = rgb(126, 151, 178)
	cBlue     = rgb(56, 169, 255)
	cBlue2    = rgb(19, 143, 236)
	cGreen    = rgb(54, 211, 133)
	cGold     = rgb(239, 184, 72)
	cRed      = rgb(255, 104, 110)
	cHover    = rgb(25, 51, 77)
	cSelected = rgb(23, 63, 95)
)

func u16(s string) *uint16 { p, _ := syscall.UTF16PtrFromString(s); return p }
func utf16NoNul(s string) []uint16 {
	r := []rune(s)
	out := make([]uint16, 0, len(r))
	for _, ch := range r {
		if ch == 0 {
			continue
		}
		if ch <= 0xffff {
			out = append(out, uint16(ch))
		} else {
			v := ch - 0x10000
			out = append(out, uint16(0xd800|(v>>10)), uint16(0xdc00|(v&0x3ff)))
		}
	}
	return out
}
func utf16MultiString(parts ...string) []uint16 {
	out := []uint16{}
	for _, p := range parts {
		out = append(out, utf16NoNul(p)...)
		out = append(out, 0)
	}
	out = append(out, 0)
	return out
}
func loword(v uintptr) uint16 { return uint16(v & 0xffff) }
func hiword(v uintptr) uint16 { return uint16((v >> 16) & 0xffff) }
func pointFromLParam(v uintptr) (int32, int32) {
	return int32(int16(v & 0xffff)), int32(int16((v >> 16) & 0xffff))
}
func trimText(s string) string { return strings.TrimSpace(strings.Trim(s, "\"")) }
func expandWindowsEnv(s string) string {
	s = os.ExpandEnv(s)
	for pass := 0; pass < 32; pass++ {
		changed := false
		var b strings.Builder
		for i := 0; i < len(s); {
			if s[i] != '%' {
				b.WriteByte(s[i])
				i++
				continue
			}
			j := strings.IndexByte(s[i+1:], '%')
			if j < 0 {
				b.WriteString(s[i:])
				break
			}
			j += i + 1
			name := s[i+1 : j]
			if name != "" {
				if v, ok := os.LookupEnv(name); ok {
					b.WriteString(v)
					changed = true
				} else {
					b.WriteString(s[i : j+1])
				}
			} else {
				b.WriteString("%%")
			}
			i = j + 1
		}
		ns := b.String()
		if !changed || ns == s {
			return ns
		}
		s = ns
	}
	return s
}
func absPath(s string) string {
	s = trimText(expandWindowsEnv(s))
	if s == "" {
		return ""
	}
	p, err := filepath.Abs(filepath.Clean(s))
	if err != nil {
		return ""
	}
	return p
}
func norm(s string) string { p := absPath(s); return strings.ToLower(filepath.Clean(p)) }
func kindOf(p string) string {
	st, err := os.Stat(p)
	if err != nil {
		return ""
	}
	if st.IsDir() {
		return "folder"
	}
	return "file"
}
func isDir(p string) bool       { return kindOf(p) == "folder" }
func volumeKey(p string) string { return strings.ToLower(filepath.VolumeName(absPath(p))) }
func displayName(p string) string {
	clean := filepath.Clean(p)
	vol := filepath.VolumeName(clean)
	t := strings.TrimRight(clean, `\\/`)
	if vol != "" && strings.EqualFold(strings.TrimRight(vol, `\\/`), strings.TrimRight(t, `\\/`)) {
		v := strings.TrimRight(vol, `:\\\/`)
		if v != "" {
			return v + "_root"
		}
	}
	b := filepath.Base(t)
	if b == "." || b == "" || b == string(filepath.Separator) {
		return "link"
	}
	return b
}
func isSubpath(p, parent string) bool {
	pp, _ := filepath.Abs(p)
	pr, _ := filepath.Abs(parent)
	rel, err := filepath.Rel(pr, pp)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)))
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer func() {
		if r := recover(); r != nil {
			showFatal(fmt.Sprintf(trText("dialog.fatal_startup"), r))
		}
	}()
	pSetProcessDpiAwarenessContext.Call(^uintptr(3)) // PER_MONITOR_AWARE_V2 (-4)
	wndProcPtr = syscall.NewCallback(wndProc)
	inst, _, _ := kernel32.NewProc("GetModuleHandleW").Call(0)
	cls := u16("mkLINKMainWindow")
	icon := loadIcon(inst)
	hHandCursor, _, _ = pLoadCursorW.Call(0, uintptr(IDC_HAND))
	hSizeWECursor, _, _ = pLoadCursorW.Call(0, uintptr(IDC_SIZEWE))
	wc := WNDCLASSEXW{CbSize: uint32(unsafe.Sizeof(WNDCLASSEXW{})), LpfnWndProc: wndProcPtr, HInstance: HINSTANCE(inst), HIcon: HICON(icon), HCursor: loadCursor(), HbrBackground: 0, LpszClassName: cls, HIconSm: HICON(icon)}
	if r, _, e := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 && e != syscall.Errno(1410) {
		showFatal(fmt.Sprintf(trText("dialog.window_class_fail"), e))
		return
	}
	title := u16(appTitle)
	a := &App{dpi: 96, uiScale: 1, icon: HICON(icon), selectedType: 0, anchorIndex: -1, columnDragCol: -1, columnDragFrom: -1, columnDragTo: -1}
	a.loadSettings()
	a.initLanguage(a.languageLocale, a.languageFile)
	globalApp = a
	h, _, e := pCreateWindowExW.Call(0, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(title)), WS_OVERLAPPEDWINDOW, CW_USEDEFAULT, CW_USEDEFAULT, 1105, 805, 0, 0, inst, 0)
	if h == 0 {
		globalApp = nil
		showFatal(fmt.Sprintf(trText("dialog.window_create_fail"), e))
		return
	}
	a.hwnd = HWND(h)
	a.initFonts()
	toolbarFont = a.buttonFont
	setDarkTitleBar(a.hwnd)
	a.createDestinationEdit()
	a.createSourceScrollbar()
	a.createSourceHScrollbar()
	pDragAcceptFiles.Call(uintptr(a.hwnd), 1)
	a.layout()
	a.refresh()
	pShowWindow.Call(h, SW_SHOW)
	pUpdateWindow.Call(h)
	var msg MSG
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func showFatal(s string) {
	pMessageBoxW.Call(0, uintptr(unsafe.Pointer(u16(s))), uintptr(unsafe.Pointer(u16(appTitle))), 0x10)
}
func loadCursor() uintptr           { r, _, _ := pLoadCursorW.Call(0, uintptr(IDC_ARROW)); return r }
func loadIcon(inst uintptr) uintptr { r, _, _ := pLoadIconW.Call(inst, 1); return r }
func (a *App) sp(v int32) int32 {
	s := a.uiScale
	if s <= 0 {
		s = 1
	}
	return int32(float64(v*a.dpi) / 96.0 * s)
}
func (a *App) releaseFonts() {
	for _, f := range []HFONT{a.font, a.titleFont, a.sectionFont, a.bodyFont, a.smallFont, a.monoFont, a.buttonFont} {
		if f != 0 {
			pDeleteObject.Call(uintptr(f))
		}
	}
	a.font, a.titleFont, a.sectionFont, a.bodyFont, a.smallFont, a.monoFont, a.buttonFont = 0, 0, 0, 0, 0, 0, 0
}
func setDarkTitleBar(hwnd HWND) {
	dark := uint32(1)
	// DWMWA_USE_IMMERSIVE_DARK_MODE: 20 on current Windows 10/11, 19 on older builds.
	if pDwmSetWindowAttribute != nil {
		pDwmSetWindowAttribute.Call(uintptr(hwnd), 20, uintptr(unsafe.Pointer(&dark)), unsafe.Sizeof(dark))
		pDwmSetWindowAttribute.Call(uintptr(hwnd), 19, uintptr(unsafe.Pointer(&dark)), unsafe.Sizeof(dark))
	}
}
func (a *App) fp(v int32) int32 {
	s := a.uiScale
	if s < 0.92 {
		s = 0.92
	}
	return int32(float64(v*a.dpi) / 96.0 * s)
}
func (a *App) initFonts() {
	if a.dpi <= 0 {
		a.dpi = 96
	}
	if a.uiScale <= 0 {
		a.uiScale = 1
	}
	// Keep typography optically stable when the window is reduced. Geometry may
	// scale, but text must not collapse into tiny, low-legibility glyphs.
	face := "Microsoft JhengHei UI"
	a.font = makeFont(a.fp(-16), FW_NORMAL, face)
	a.bodyFont = makeFont(a.fp(-16), FW_NORMAL, face)
	a.smallFont = makeFont(a.fp(-14), FW_NORMAL, face)
	a.sectionFont = makeFont(a.fp(-19), FW_SEMIBOLD, face)
	a.titleFont = makeFont(a.fp(-24), FW_BOLD, face)
	a.buttonFont = makeFont(a.fp(-16), FW_NORMAL, face)
	a.monoFont = makeFont(a.fp(-14), FW_NORMAL, "Consolas")
}

func newBrush(c uint32) HBRUSH { r, _, _ := pCreateSolidBrush.Call(uintptr(c)); return HBRUSH(r) }
func makeFont(h, w int32, face string) HFONT {
	p := u16(face)
	r, _, _ := pCreateFontW.Call(uintptr(h), 0, 0, 0, uintptr(w), 0, 0, 0, 1, 0, 0, 5, 0, uintptr(unsafe.Pointer(p)))
	return HFONT(r)
}
func (a *App) createDestinationEdit() {
	class, empty := u16("EDIT"), u16("")
	h, _, _ := pCreateWindowExW.Call(0, uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(empty)), WS_CHILD|WS_VISIBLE|ES_AUTOHSCROLL|WS_TABSTOP, 0, 0, 0, 0, uintptr(a.hwnd), 2100, 0, 0)
	a.destination = HWND(h)
	brInput = HBRUSH(newBrush(cSurface))
	pSendMessageW.Call(uintptr(a.destination), WM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(u16(a.tr("destination.placeholder")))))
	pSendMessageW.Call(uintptr(a.destination), WM_SETFONT, uintptr(a.buttonFont), TRUE)
}
func (a *App) createSourceScrollbar() {
	class, empty := u16("SCROLLBAR"), u16("")
	h, _, _ := pCreateWindowExW.Call(0, uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(empty)), WS_CHILD|SBS_VERT, 0, 0, 0, 0, uintptr(a.hwnd), 2101, 0, 0)
	a.sourceScrollbar = HWND(h)
}
func (a *App) createSourceHScrollbar() {
	class, empty := u16("SCROLLBAR"), u16("")
	h, _, _ := pCreateWindowExW.Call(0, uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(empty)), WS_CHILD|SBS_HORZ, 0, 0, 0, 0, uintptr(a.hwnd), 2102, 0, 0)
	a.sourceHScrollbar = HWND(h)
}

func (a *App) sourceMetrics() (rowH int32, visible, maxScroll int) {
	headerH := a.sp(46)
	rowH = a.sp(54)
	contentH := a.sourceListArea.Bottom - a.sourceListArea.Top - headerH
	if contentH < rowH {
		visible = 1
	} else {
		visible = int(contentH / rowH)
	}
	if visible < 1 {
		visible = 1
	}
	maxScroll = len(a.items) - visible
	if maxScroll < 0 {
		maxScroll = 0
	}
	return
}

func (a *App) clampSourceScroll() {
	_, _, maxScroll := a.sourceMetrics()
	if a.scroll < 0 {
		a.scroll = 0
	}
	if a.scroll > maxScroll {
		a.scroll = maxScroll
	}
}

func (a *App) layout() {
	var rc RECT
	pGetClientRect.Call(uintptr(a.hwnd), uintptr(unsafe.Pointer(&rc)))
	cw, ch := rc.Right, rc.Bottom
	// Geometry follows the supplied reference composition, but every rectangle
	// is derived from the current client area. Nothing is allowed to extend
	// beyond the actual window.
	ws := float64(cw) / 1536.0
	hs := float64(ch) / 864.0
	s := ws
	if hs < s {
		s = hs
	}
	if s > 1 {
		s = 1
	}
	if s < 0.72 {
		s = 0.72
	}
	a.uiScale = s
	if a.dpi <= 0 {
		a.dpi = 96
	}
	a.releaseFonts()
	a.initFonts()
	toolbarFont = a.buttonFont

	pad := a.sp(22)
	gap := a.sp(18)
	headerH := a.sp(92)
	toolbarH := a.sp(92)
	footerH := a.sp(88)
	mainTop := headerH + toolbarH + a.sp(18)
	mainBottom := ch - footerH - a.sp(18)
	if mainBottom < mainTop+a.sp(320) {
		mainBottom = mainTop + a.sp(320)
	}
	// Four compact primary actions. Paste Path and Remove Selected are now
	// available directly in the source-area context menu and via Ctrl+V/Delete.
	a.toolbarRects = nil
	toolbarGap := a.sp(10)
	buttonCount := 4
	toolbarAvail := cw - 2*pad - int32(buttonCount-1)*toolbarGap
	buttonW := toolbarAvail / int32(buttonCount)
	if buttonW < a.sp(150) {
		buttonW = a.sp(150)
	}
	x := pad
	for i := 0; i < buttonCount; i++ {
		w := buttonW
		if i == buttonCount-1 {
			w = cw - pad - x
		}
		a.toolbarRects = append(a.toolbarRects, RECT{x, headerH + a.sp(12), x + w, headerH + a.sp(80)})
		x += w + toolbarGap
	}

	usableW := cw - 2*pad - gap
	leftW := int32(float64(usableW) * 0.56)
	minLeft := a.sp(500)
	minRight := a.sp(370)
	if leftW < minLeft {
		leftW = minLeft
	}
	if leftW > usableW-minRight {
		leftW = usableW - minRight
	}
	if leftW < a.sp(420) {
		leftW = a.sp(420)
	}
	rightX := pad + leftW + gap
	if rightX >= cw-pad {
		rightX = pad + (usableW*55)/100
	}
	a.sourceRect = RECT{pad, mainTop, pad + leftW, mainBottom}

	// The source list now fills the panel. The drag/drop hint is rendered inside
	// the same area only while it is empty. A native vertical scrollbar is used
	// for rows and a native horizontal scrollbar appears only when the user
	// configured columns wider than the available viewport.
	// The actual source-table geometry is finalized by layoutSourceTable().
	a.sourceListArea = RECT{a.sourceRect.Left + a.sp(12), a.sourceRect.Top + a.sp(70), a.sourceRect.Right - a.sp(29), a.sourceRect.Bottom - a.sp(46)}
	a.layoutSourceTable()

	right := RECT{rightX, mainTop, cw - pad, mainBottom}
	inner := a.sp(16)
	available := right.Bottom - right.Top
	baseGap := a.sp(12)
	// Reference proportions; clamp as a group so panels never overlap.
	modeH := a.sp(214)
	destH := a.sp(122)
	compatH := a.sp(100)
	previewMin := a.sp(118)
	needed := modeH + destH + compatH + previewMin + baseGap*3
	if available < needed {
		extra := needed - available
		for extra > 0 && modeH > a.sp(170) {
			modeH -= a.sp(4)
			extra -= a.sp(4)
		}
		for extra > 0 && destH > a.sp(92) {
			destH -= a.sp(3)
			extra -= a.sp(3)
		}
		for extra > 0 && compatH > a.sp(78) {
			compatH -= a.sp(2)
			extra -= a.sp(2)
		}
	}
	previewH := available - modeH - destH - compatH - baseGap*3
	if previewH < previewMin {
		previewH = previewMin
	}
	a.modeRect = RECT{right.Left, right.Top, right.Right, right.Top + modeH}
	a.destRect = RECT{right.Left, a.modeRect.Bottom + baseGap, right.Right, a.modeRect.Bottom + baseGap + destH}
	a.compatRect = RECT{right.Left, a.destRect.Bottom + baseGap, right.Right, a.destRect.Bottom + baseGap + compatH}
	a.previewRect = RECT{right.Left, a.compatRect.Bottom + baseGap, right.Right, right.Bottom}
	// Final safety clamp: each rectangle stays inside the client area.
	if a.previewRect.Bottom > mainBottom {
		a.previewRect.Bottom = mainBottom
	}
	if a.previewRect.Top > a.previewRect.Bottom-a.sp(40) {
		a.previewRect.Top = a.previewRect.Bottom - a.sp(40)
	}

	browseW := a.sp(104)
	editH := a.sp(44)
	destY := a.destRect.Top + a.sp(58)
	if destY+editH > a.destRect.Bottom-a.sp(10) {
		destY = a.destRect.Bottom - a.sp(10) - editH
	}
	if destY < a.destRect.Top+a.sp(48) {
		destY = a.destRect.Top + a.sp(48)
	}
	a.destEditRect = RECT{a.destRect.Left + inner, destY, a.destRect.Right - inner - browseW - a.sp(10), destY + editH}
	if a.destEditRect.Right < a.destEditRect.Left+a.sp(120) {
		a.destEditRect.Right = a.destEditRect.Left + a.sp(120)
	}
	a.browseRect = RECT{a.destRect.Right - inner - browseW, destY, a.destRect.Right - inner, destY + editH}
	pMoveWindow.Call(uintptr(a.destination), uintptr(a.destEditRect.Left), uintptr(a.destEditRect.Top), uintptr(a.destEditRect.Right-a.destEditRect.Left), uintptr(a.destEditRect.Bottom-a.destEditRect.Top), 1)
	pSendMessageW.Call(uintptr(a.destination), WM_SETFONT, uintptr(a.buttonFont), TRUE)

	cardGap := a.sp(10)
	cardW := (a.modeRect.Right - a.modeRect.Left - 2*inner - 2*cardGap) / 3
	for i := 0; i < 3; i++ {
		l := a.modeRect.Left + inner + int32(i)*(cardW+cardGap)
		top := a.modeRect.Top + a.sp(72)
		bottom := a.modeRect.Bottom - a.sp(14)
		if bottom < top+a.sp(64) {
			top = a.modeRect.Top + a.sp(60)
			bottom = a.modeRect.Bottom - a.sp(8)
		}
		a.modeRects[i] = RECT{l, top, l + cardW, bottom}
	}

	languageW := a.sp(176)
	languageGap := a.sp(12)
	versionW := a.sp(160)
	a.languageRect = RECT{cw - pad - versionW - languageGap - languageW, a.sp(18), cw - pad - versionW - languageGap, a.sp(58)}
	createW := a.sp(228)
	a.createRect = RECT{cw - createW - pad, ch - a.sp(70), cw - pad, ch - a.sp(18)}
	creatorW := a.sp(280)
	a.creatorRect = RECT{a.createRect.Left - creatorW - a.sp(18), ch - a.sp(68), a.createRect.Left - a.sp(18), ch - a.sp(22)}
	a.optionsRect = RECT{}
	a.updateSourceScrollbars()
	pInvalidateRect.Call(uintptr(a.hwnd), 0, 1)
}

func fillRect(hdc uintptr, r RECT, c uint32) {
	b, _, _ := pCreateSolidBrush.Call(uintptr(c))
	if b == 0 {
		return
	}
	pFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), b)
	pDeleteObject.Call(b)
}
func frameRect(hdc uintptr, r RECT, c uint32) {
	pen, _, _ := pCreatePen.Call(PS_SOLID, 1, uintptr(c))
	if pen == 0 {
		return
	}
	defer pDeleteObject.Call(pen)
	null, _, _ := pGetStockObject.Call(NULL_BRUSH)
	op, _, _ := pSelectObject.Call(hdc, pen)
	ob, _, _ := pSelectObject.Call(hdc, null)
	pRectangle.Call(hdc, uintptr(r.Left), uintptr(r.Top), uintptr(r.Right), uintptr(r.Bottom))
	pSelectObject.Call(hdc, ob)
	pSelectObject.Call(hdc, op)
}
func textRect(hdc uintptr, r RECT, s string, c uint32, f HFONT, flags uintptr) {
	if f != 0 {
		old, _, _ := pSelectObject.Call(hdc, uintptr(f))
		defer pSelectObject.Call(hdc, old)
	}
	pSetBkMode.Call(hdc, TRANSPARENT)
	pSetTextColor.Call(hdc, uintptr(c))
	b := utf16NoNul(s)
	if len(b) == 0 {
		return
	}
	rr := r
	pDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(unsafe.Pointer(&rr)), flags)
}

func drawModeTitle(hdc uintptr, r RECT, s string, c uint32, f HFONT) {
	if f == 0 {
		return
	}
	b := utf16NoNul(s)
	if len(b) == 0 {
		return
	}
	old, _, _ := pSelectObject.Call(hdc, uintptr(f))
	defer pSelectObject.Call(hdc, old)
	pSetBkMode.Call(hdc, TRANSPARENT)
	pSetTextColor.Call(hdc, uintptr(c))

	// First measure the wrapped text with the exact width available in the
	// card. DrawTextW will break at word boundaries (and at CJK boundaries)
	// as needed.
	width := r.Right - r.Left
	if width <= 0 || r.Bottom <= r.Top {
		return
	}
	measure := RECT{Left: 0, Top: 0, Right: width, Bottom: 0}
	flags := uintptr(DT_CENTER | DT_WORDBREAK | DT_NOPREFIX | DT_CALCRECT)
	pDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(unsafe.Pointer(&measure)), flags)
	textH := measure.Bottom - measure.Top
	if textH <= 0 {
		return
	}
	boxH := r.Bottom - r.Top
	if textH > boxH {
		textH = boxH
	}
	top := r.Top + (boxH-textH)/2
	rr := RECT{Left: r.Left, Top: top, Right: r.Right, Bottom: top + textH}
	pDrawTextW.Call(hdc, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(unsafe.Pointer(&rr)), uintptr(DT_CENTER|DT_WORDBREAK|DT_NOPREFIX))
}
func (a *App) paint() {
	var ps PAINTSTRUCT
	hdc, _, _ := pBeginPaint.Call(uintptr(a.hwnd), uintptr(unsafe.Pointer(&ps)))
	if hdc == 0 {
		return
	}
	defer pEndPaint.Call(uintptr(a.hwnd), uintptr(unsafe.Pointer(&ps)))
	var rc RECT
	pGetClientRect.Call(uintptr(a.hwnd), uintptr(unsafe.Pointer(&rc)))
	w := rc.Right - rc.Left
	h := rc.Bottom - rc.Top
	if w <= 0 || h <= 0 {
		return
	}

	// Render the complete frame into an off-screen bitmap first. This prevents
	// visible GDI tearing/flicker when WM_MOUSEMOVE changes hover state.
	memDC, _, _ := pCreateCompatibleDC.Call(hdc)
	if memDC == 0 {
		fillRect(hdc, rc, cBG)
		a.drawHeader(hdc, rc)
		a.drawToolbar(hdc)
		a.drawMain(hdc, rc)
		a.drawFooter(hdc, rc)
		return
	}
	defer pDeleteDC.Call(memDC)
	bmp, _, _ := pCreateCompatibleBitmap.Call(hdc, uintptr(w), uintptr(h))
	if bmp == 0 {
		fillRect(hdc, rc, cBG)
		a.drawHeader(hdc, rc)
		a.drawToolbar(hdc)
		a.drawMain(hdc, rc)
		a.drawFooter(hdc, rc)
		return
	}
	defer pDeleteObject.Call(bmp)
	oldBmp, _, _ := pSelectObject.Call(memDC, bmp)
	defer pSelectObject.Call(memDC, oldBmp)

	fillRect(memDC, rc, cBG)
	a.drawHeader(memDC, rc)
	a.drawToolbar(memDC)
	a.drawMain(memDC, rc)
	a.drawFooter(memDC, rc)

	pBitBlt.Call(hdc, 0, 0, uintptr(w), uintptr(h), memDC, 0, 0, uintptr(SRCCOPY))
}
func drawRoundedPanel(hdc uintptr, r RECT, fill, border uint32, rad int32) {
	b, _, _ := pCreateSolidBrush.Call(uintptr(fill))
	if b == 0 {
		return
	}
	defer pDeleteObject.Call(b)
	pen, _, _ := pCreatePen.Call(PS_SOLID, 1, uintptr(border))
	if pen == 0 {
		return
	}
	defer pDeleteObject.Call(pen)

	ob, _, _ := pSelectObject.Call(hdc, b)
	op, _, _ := pSelectObject.Call(hdc, pen)
	pRoundRect.Call(hdc, uintptr(r.Left), uintptr(r.Top), uintptr(r.Right), uintptr(r.Bottom), uintptr(rad), uintptr(rad))
	pSelectObject.Call(hdc, ob)
	pSelectObject.Call(hdc, op)
}
func (a *App) drawHeader(hdc uintptr, rc RECT) {
	h := a.sp(92)
	fillRect(hdc, RECT{0, 0, rc.Right, h}, cHeader)
	fillRect(hdc, RECT{0, h - 1, rc.Right, h}, cBorder)
	if a.icon != 0 {
		pDrawIconEx.Call(hdc, uintptr(a.sp(18)), uintptr(a.sp(17)), uintptr(a.icon), uintptr(a.sp(42)), uintptr(a.sp(42)), 0, 0, DI_NORMAL)
	}
	textRect(hdc, RECT{a.sp(64), a.sp(12), a.sp(182), a.sp(46)}, appTitle, cText, a.titleFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	textRect(hdc, RECT{a.sp(190), a.sp(16), a.sp(440), a.sp(44)}, a.tr("app.subtitle"), cText2, a.bodyFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	textRect(hdc, RECT{a.sp(64), a.sp(49), a.sp(430), a.sp(76)}, a.tr("app.tagline"), cText2, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	langBG, langBorder := cSurface2, cBorder2
	if a.hoverID == 40 {
		langBG, langBorder = cHover, cBlue2
	}
	drawRoundedPanel(hdc, a.languageRect, langBG, langBorder, a.sp(6))
	textRect(hdc, RECT{a.languageRect.Left + a.sp(10), a.languageRect.Top, a.languageRect.Right - a.sp(30), a.languageRect.Bottom}, a.languageDisplayName(), cText2, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
	textRect(hdc, RECT{a.languageRect.Right - a.sp(28), a.languageRect.Top, a.languageRect.Right - a.sp(8), a.languageRect.Bottom}, a.tr("create.dropdown"), cText2, a.smallFont, DT_CENTER|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	textRect(hdc, RECT{rc.Right - a.sp(160), a.sp(12), rc.Right - a.sp(24), a.sp(38)}, "v"+appVersion, cText, a.bodyFont, DT_RIGHT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
}
func (a *App) drawToolbar(hdc uintptr) {
	for i, r := range a.toolbarRects {
		a.drawReferenceToolbar(hdc, r, i, a.hoverID == i+1, a.pressedID == i+1)
	}
}
func (a *App) drawReferenceToolbar(hdc uintptr, r RECT, idx int, hover, pressed bool) {
	bg := cToolbar
	if hover {
		bg = cHover
	}
	if pressed {
		bg = cSelected
	}
	drawRoundedPanel(hdc, r, bg, cBorder, a.sp(8))
	iconR := RECT{r.Left + a.sp(18), r.Top + 14, r.Left + a.sp(64), r.Bottom - 14}
	drawCommandIcon(hdc, iconR, idx, cBlue)
	labels := []string{a.tr("toolbar.add_file"), a.tr("toolbar.add_folder"), a.tr("toolbar.clear"), a.tr("toolbar.admin")}
	subs := []string{a.tr("toolbar.add_file_sub"), a.tr("toolbar.add_folder_sub"), a.tr("toolbar.clear_sub"), a.tr("toolbar.admin_sub")}
	textRect(hdc, RECT{r.Left + a.sp(76), r.Top + a.sp(11), r.Right - a.sp(10), r.Top + a.sp(39)}, labels[idx], cText, a.bodyFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	textRect(hdc, RECT{r.Left + a.sp(76), r.Top + a.sp(40), r.Right - a.sp(10), r.Bottom - a.sp(8)}, subs[idx], cMuted, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
}
func (a *App) drawMain(hdc uintptr, rc RECT) {
	// Compute preflight (filesystem stat/lstat checks) exactly once per paint
	// and share it with drawCompatModern/drawPreviewModern below. Previously
	// each of those called a.preflight() independently, so every repaint —
	// including ones triggered purely by a hover-state change on mouse move —
	// re-ran the same disk I/O for every source item twice. With many items,
	// or a source/destination on a slow or disconnected network path, that
	// doubled I/O on every hover crossing was a real stutter/hang risk.
	a.pfOK, a.pfIssues, a.pfCmds, a.pfItems, _ = a.preflight()
	a.drawSourceModern(hdc)
	a.drawModeModern(hdc)
	a.drawDestModern(hdc)
	a.drawCompatModern(hdc)
	a.drawPreviewModern(hdc)
}
func (a *App) panelTitle(hdc uintptr, r RECT, title string, iconIdx int, count string, f HFONT) {
	textRect(hdc, RECT{r.Left + a.sp(62), r.Top + a.sp(12), r.Right - a.sp(12), r.Top + a.sp(45)}, title, cText, f, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	if count != "" {
		textRect(hdc, RECT{r.Right - a.sp(150), r.Top + 15, r.Right - a.sp(22), r.Top + a.sp(43)}, count, cText2, a.smallFont, DT_RIGHT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	}
	_ = iconIdx
}
func (a *App) drawSourceModern(hdc uintptr) {
	r := a.sourceRect
	drawRoundedPanel(hdc, r, cPanel, cBorder, a.sp(10))
	drawFolderIcon(hdc, RECT{r.Left + a.sp(20), r.Top + a.sp(17), r.Left + a.sp(52), r.Top + a.sp(49)}, cBlue)
	a.panelTitle(hdc, r, a.tr("source.title"), 2, fmt.Sprintf(a.tr("source.count"), len(a.items)), a.sectionFont)
	textRect(hdc, RECT{r.Left + a.sp(20), r.Top + a.sp(50), r.Right - a.sp(20), r.Top + a.sp(70)}, a.tr("source.description"), cMuted, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
	list := a.sourceListArea
	if len(a.items) == 0 {
		fillRect(hdc, list, cSurface)
		frameRect(hdc, list, cBorder2)
		cx := (list.Left + list.Right) / 2
		cy := (list.Top + list.Bottom) / 2
		drawDropIcon(hdc, RECT{cx - a.sp(28), cy - a.sp(36), cx + a.sp(28), cy + a.sp(20)}, cBlue)
		textRect(hdc, RECT{list.Left + a.sp(24), cy + a.sp(28), list.Right - a.sp(24), cy + a.sp(58)}, a.tr("source.empty.title"), cText, a.sectionFont, DT_CENTER|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
		textRect(hdc, RECT{list.Left + a.sp(24), cy + a.sp(59), list.Right - a.sp(24), cy + a.sp(84)}, a.tr("source.empty.subtitle"), cText2, a.bodyFont, DT_CENTER|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
	} else {
		a.drawSourceTable(hdc, list)
	}
	textRect(hdc, RECT{r.Left + a.sp(20), r.Bottom - a.sp(36), r.Right - a.sp(20), r.Bottom - a.sp(12)}, fmt.Sprintf(a.tr("source.selected.count"), len(a.selectedItems())), cText2, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
}
func (a *App) drawCheckbox(hdc uintptr, r RECT, on bool) {
	drawRoundedPanel(hdc, r, func() uint32 {
		if on {
			return cBlue
		}
		return cSurface2
	}(), func() uint32 {
		if on {
			return cBlue
		}
		return cBorder2
	}(), checkboxCornerRadius)
	if on {
		drawCheckMark(hdc, r, cText)
	}
}

const checkboxCornerRadius int32 = 5

func drawCheckMark(hdc uintptr, r RECT, c uint32) {
	pen := iconPen(c)
	if pen == 0 {
		return
	}
	defer pDeleteObject.Call(pen)
	old, _, _ := pSelectObject.Call(hdc, pen)
	defer pSelectObject.Call(hdc, old)
	w := r.Right - r.Left
	h := r.Bottom - r.Top
	if w <= 0 || h <= 0 {
		return
	}
	// Draw proportionally within the checkbox so the mark remains balanced
	// across DPI scales and any future checkbox size adjustments.
	x0 := r.Left + w*20/100
	y0 := r.Top + h*50/100
	xm := r.Left + w*40/100
	ym := r.Top + h*72/100
	x1 := r.Right - w*15/100
	y1 := r.Top + h*22/100
	pMoveToEx.Call(hdc, uintptr(x0), uintptr(y0), 0)
	pLineTo.Call(hdc, uintptr(xm), uintptr(ym))
	pLineTo.Call(hdc, uintptr(x1), uintptr(y1))
}
func drawItemIcon(hdc uintptr, r RECT, kind string) {
	if kind == "folder" {
		drawFolderIcon(hdc, r, cGold)
	} else {
		drawFileIcon(hdc, r, cText2)
	}
}
func (a *App) drawModeModern(hdc uintptr) {
	r := a.modeRect
	drawRoundedPanel(hdc, r, cPanel, cBorder, a.sp(10))
	drawChainIcon(hdc, RECT{r.Left + a.sp(20), r.Top + a.sp(18), r.Left + a.sp(52), r.Top + a.sp(50)}, cBlue)
	textRect(hdc, RECT{r.Left + a.sp(62), r.Top + a.sp(12), r.Right, r.Top + a.sp(44)}, a.tr("mode.title"), cText, a.sectionFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	textRect(hdc, RECT{r.Left + a.sp(62), r.Top + a.sp(43), r.Right - a.sp(20), r.Top + a.sp(66)}, a.tr("mode.subtitle"), cMuted, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
	titles := []string{a.tr("mode.symbolic"), a.tr("mode.hard"), a.tr("mode.junction")}
	desc1 := []string{a.tr("mode.symbolic_desc1"), a.tr("mode.hard_desc1"), a.tr("mode.junction_desc1")}
	desc2 := []string{a.tr("mode.symbolic_desc2"), a.tr("mode.hard_desc2"), a.tr("mode.junction_desc2")}
	for i, cr := range a.modeRects {
		selected := a.selectedType == i
		bg := cSurface
		border := cBorder2
		switch {
		case selected:
			bg = cSelected
			border = cBlue
		case a.hoverID == 20+i:
			bg = cHover
			border = cBorder2
		}
		drawRoundedPanel(hdc, cr, bg, border, a.sp(8))
		drawModeIcon(hdc, RECT{cr.Left + a.sp(16), cr.Top + a.sp(16), cr.Left + a.sp(54), cr.Top + a.sp(54)}, i, []uint32{cBlue, cText2, cGold}[i])
		modeTitleFont := a.bodyFont
		modeTextLeft := a.sp(66)
		if cr.Right-cr.Left < a.sp(150) {
			modeTitleFont = a.smallFont
			modeTextLeft = a.sp(56)
		}
		// The title is allowed to wrap to two lines so English names such as
		// "Symbolic Link" and "Directory Junction" remain fully visible in
		// the existing fixed-width cards. The measured text block is vertically
		// centered against the icon instead of shrinking the entire UI.
		drawModeTitle(hdc, RECT{cr.Left + modeTextLeft, cr.Top + a.sp(10), cr.Right - a.sp(8), cr.Top + a.sp(62)}, titles[i], cText, modeTitleFont)
		textRect(hdc, RECT{cr.Left + a.sp(16), cr.Top + a.sp(72), cr.Right - a.sp(12), cr.Top + a.sp(96)}, desc1[i], cText2, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
		textRect(hdc, RECT{cr.Left + a.sp(16), cr.Top + a.sp(98), cr.Right - a.sp(12), cr.Top + a.sp(122)}, desc2[i], cText2, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
	}
}
func (a *App) drawDestModern(hdc uintptr) {
	r := a.destRect
	drawRoundedPanel(hdc, r, cPanel, cBorder, a.sp(10))
	drawFolderIcon(hdc, RECT{r.Left + a.sp(20), r.Top + a.sp(18), r.Left + a.sp(52), r.Top + a.sp(50)}, cBlue)
	textRect(hdc, RECT{r.Left + a.sp(62), r.Top + a.sp(12), r.Right, r.Top + a.sp(44)}, a.tr("destination.title"), cText, a.sectionFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	bBg, bBorder := cSurface2, cBorder2
	switch {
	case a.pressedID == 30:
		bBg, bBorder = cSelected, cBlue
	case a.hoverID == 30:
		bBg, bBorder = cHover, cBlue2
	}
	drawRoundedPanel(hdc, a.browseRect, bBg, bBorder, a.sp(6))
	drawFolderIcon(hdc, RECT{a.browseRect.Left + a.sp(12), a.browseRect.Top + a.sp(8), a.browseRect.Left + a.sp(38), a.browseRect.Bottom - a.sp(8)}, cBlue)
	textRect(hdc, RECT{a.browseRect.Left + a.sp(40), a.browseRect.Top, a.browseRect.Right - a.sp(8), a.browseRect.Bottom}, a.tr("destination.browse"), cText, a.smallFont, DT_CENTER|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
}
func (a *App) drawCompatModern(hdc uintptr) {
	r := a.compatRect
	drawRoundedPanel(hdc, r, cPanel, cBorder, a.sp(10))
	drawGearIcon(hdc, RECT{r.Left + a.sp(20), r.Top + a.sp(18), r.Left + a.sp(52), r.Top + a.sp(50)}, cBlue)
	textRect(hdc, RECT{r.Left + a.sp(62), r.Top + a.sp(12), r.Right, r.Top + a.sp(44)}, a.tr("preflight.title"), cText, a.sectionFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	if a.pfOK {
		drawStatusDot(hdc, RECT{r.Left + a.sp(22), r.Top + a.sp(62), r.Left + a.sp(46), r.Top + a.sp(86)}, cGreen)
		textRect(hdc, RECT{r.Left + a.sp(58), r.Top + a.sp(58), r.Right - a.sp(16), r.Top + a.sp(90)}, a.tr("preflight.ok"), cGreen, a.bodyFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
	} else {
		drawStatusDot(hdc, RECT{r.Left + a.sp(22), r.Top + a.sp(62), r.Left + a.sp(46), r.Top + a.sp(86)}, cRed)
		msg := a.tr("preflight.fail")
		if len(a.pfIssues) > 0 {
			msg += a.pfIssues[0]
		}
		textRect(hdc, RECT{r.Left + a.sp(58), r.Top + a.sp(58), r.Right - a.sp(16), r.Top + a.sp(90)}, msg, cRed, a.bodyFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
	}
}
func (a *App) drawPreviewModern(hdc uintptr) {
	r := a.previewRect
	drawRoundedPanel(hdc, r, cPanel, cBorder, a.sp(10))
	drawEyeIcon(hdc, RECT{r.Left + a.sp(20), r.Top + a.sp(18), r.Left + a.sp(52), r.Top + a.sp(50)}, cBlue)
	textRect(hdc, RECT{r.Left + a.sp(62), r.Top + a.sp(12), r.Right, r.Top + a.sp(44)}, a.tr("preview.title"), cText, a.sectionFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	box := RECT{r.Left + a.sp(18), r.Top + a.sp(58), r.Right - a.sp(18), r.Bottom - a.sp(18)}
	drawRoundedPanel(hdc, box, cSurface, cBorder2, a.sp(6))
	textRect(hdc, RECT{box.Left + a.sp(12), box.Top + a.sp(10), box.Right - a.sp(12), box.Bottom - a.sp(10)}, a.previewText(), cText2, a.smallFont, DT_LEFT|DT_WORDBREAK|DT_NOPREFIX)
}
func (a *App) drawFooter(hdc uintptr, rc RECT) {
	y := rc.Bottom - a.sp(88)
	fillRect(hdc, RECT{0, y, rc.Right, rc.Bottom}, cHeader)
	fillRect(hdc, RECT{0, y, rc.Right, y + 1}, cBorder)
	textRect(hdc, RECT{a.sp(36), y + a.sp(22), a.sp(150), rc.Bottom - a.sp(20)}, a.tr("footer.ready"), cText2, a.bodyFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	footerStatus := a.status
	if hover := a.modeHoverText(); hover != "" {
		footerStatus = hover
	}
	textRect(hdc, RECT{a.sp(170), y + a.sp(22), a.creatorRect.Left - a.sp(18), rc.Bottom - a.sp(20)}, footerStatus, cText2, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
	creatorColor := cBlue
	if a.hoverID == 32 {
		creatorColor = cText
	}
	textRect(hdc, a.creatorRect, a.tr("footer.creator"), creatorColor, a.smallFont, DT_CENTER|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
	bg := cBlue2
	switch {
	case !a.pfOK:
		bg = cSurface2
	case a.pressedID == 31:
		bg = rgb(13, 108, 168)
	case a.hoverID == 31:
		bg = cBlue
	}
	drawRoundedPanel(hdc, a.createRect, bg, cBlue, a.sp(7))
	drawChainIcon(hdc, RECT{a.createRect.Left + a.sp(20), a.createRect.Top + 15, a.createRect.Left + a.sp(54), a.createRect.Bottom - 15}, cText)
	textRect(hdc, RECT{a.createRect.Left + a.sp(62), a.createRect.Top, a.createRect.Right - a.sp(42), a.createRect.Bottom}, a.tr("create.title"), cText, a.buttonFont, DT_CENTER|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
	textRect(hdc, RECT{a.createRect.Right - a.sp(34), a.createRect.Top, a.createRect.Right - a.sp(10), a.createRect.Bottom}, a.tr("create.dropdown"), cText, a.bodyFont, DT_CENTER|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)
}
func (a *App) modeHoverText() string {
	switch a.hoverID {
	case 20:
		return a.tr("mode.hover.symbolic")
	case 21:
		return a.tr("mode.hover.hard")
	case 22:
		return a.tr("mode.hover.junction")
	default:
		return ""
	}
}

func (a *App) previewText() string {
	if !a.pfOK {
		return a.tr("preview.issues") + "\n" + strings.Join(a.pfIssues, "\n")
	}
	if len(a.pfItems) == 0 {
		return a.tr("preview.none")
	}
	return a.tr("preview.plan") + "\n" + strings.Join(a.pfCmds, "\n")
}

func drawStatusDot(hdc uintptr, r RECT, c uint32) {
	// Keep the existing call site, but use the supplied Octicons for both states.
	if c == cRed {
		drawOcticon(hdc, r, "x", c)
		return
	}
	drawStatusIcon(hdc, r, c)
}
func iconPen(c uint32) uintptr {
	p, _, _ := pCreatePen.Call(PS_SOLID, 2, uintptr(c))
	return p
}

func (a *App) hitToolbarIndex(x, y int32) int {
	for i, r := range a.toolbarRects {
		if pointIn(r, x, y) {
			return i
		}
	}
	return -1
}
func (a *App) hitTest(x, y int32) int {
	if i := a.hitToolbarIndex(x, y); i >= 0 {
		return i + 1
	}
	for i, r := range a.modeRects {
		if pointIn(r, x, y) {
			return 20 + i
		}
	}
	if pointIn(a.languageRect, x, y) {
		return 40
	}
	if pointIn(a.browseRect, x, y) {
		return 30
	}
	if pointIn(a.createRect, x, y) {
		return 31
	}
	if pointIn(a.creatorRect, x, y) {
		return 32
	}
	if pointIn(a.sourceListArea, x, y) {
		return 60
	}
	return 0
}
func pointIn(r RECT, x, y int32) bool {
	return x >= r.Left && x <= r.Right && y >= r.Top && y <= r.Bottom
}
func (a *App) wndCommand(hit int32, x, y int32) {
	switch {
	case hit >= 1 && hit <= 4:
		switch hit {
		case 1:
			a.addFiles()
		case 2:
			a.addFolder()
		case 3:
			a.clearItems()
		case 4:
			a.runAsAdministrator()
		}
	case hit >= 20 && hit <= 22:
		a.selectedType = int(hit - 20)
		a.refresh()
	case hit == 30:
		a.chooseDestination()
	case hit == 31:
		a.createLinks()
	case hit == 60:
		a.selectSourceAt(x, y)
	case hit == 32:
		a.openCreatorLink()
	case hit == 40:
		a.showLanguageMenu()
	}
}
func (a *App) sourceListAreaForHit() RECT { return a.sourceListArea }
func (a *App) sourceIndexAt(x, y int32) int {
	r := a.sourceListAreaForHit()
	if !pointIn(r, x, y) || len(a.items) == 0 {
		return -1
	}
	headerH := a.sp(46)
	if y < r.Top+headerH {
		return -1
	}
	rowH, _, _ := a.sourceMetrics()
	idx := a.scroll + int((y-r.Top-headerH)/rowH)
	if idx < 0 || idx >= len(a.items) {
		return -1
	}
	return idx
}
func (a *App) selectSourceAt(x, y int32) {
	a.sourceFocus = true
	r := a.sourceListAreaForHit()
	if !pointIn(r, x, y) {
		return
	}
	idx := a.sourceIndexAt(x, y)
	if idx < 0 {
		return
	}
	// The checkbox is an explicit hit target. A plain left-click on the box
	// toggles that item's state directly, matching familiar list-selection UI.
	// A checkbox click also becomes the selection anchor for subsequent Shift-clicks.
	rowH, _, _ := a.sourceMetrics()
	headerH := a.sp(46)
	rowTop := r.Top + headerH + int32(idx-a.scroll)*rowH
	checkbox := RECT{r.Left + a.sp(15), rowTop + a.sp(15), r.Left + a.sp(39), rowTop + a.sp(39)}
	if pointIn(checkbox, x, y) {
		a.items[idx].Selected = !a.items[idx].Selected
		a.anchorIndex = idx
		a.status = fmt.Sprintf(a.tr("source.status.selection"), len(a.selectedItems()))
		a.refresh()
		setFocus(a.hwnd)
		return
	}
	ctrl := false
	shift := false
	if r, _, _ := pGetAsyncKeyState.Call(uintptr(VK_CONTROL)); int16(r) < 0 {
		ctrl = true
	}
	if r, _, _ := pGetAsyncKeyState.Call(uintptr(VK_SHIFT)); int16(r) < 0 {
		shift = true
	}
	if !ctrl && !shift {
		for i := range a.items {
			a.items[i].Selected = false
		}
		a.items[idx].Selected = true
		a.anchorIndex = idx
	} else if ctrl {
		a.items[idx].Selected = !a.items[idx].Selected
		a.anchorIndex = idx
	} else {
		startIdx := a.anchorIndex
		if startIdx < 0 || startIdx >= len(a.items) {
			startIdx = idx
		}
		lo, hi := startIdx, idx
		if lo > hi {
			lo, hi = hi, lo
		}
		for i := range a.items {
			a.items[i].Selected = i >= lo && i <= hi
		}
	}
	a.status = fmt.Sprintf(a.tr("source.status.selection"), len(a.selectedItems()))
	a.refresh()
	setFocus(a.hwnd)
}

func (a *App) ensureSelectionState() {
	if len(a.items) == 0 {
		a.status = a.tr("source.status.empty")
	} else {
		a.status = fmt.Sprintf(a.tr("source.status.selected"), len(a.items), len(a.selectedItems()))
	}
}
func (a *App) selectedItems() []SourceItem {
	out := []SourceItem{}
	for _, it := range a.items {
		if it.Selected {
			out = append(out, it)
		}
	}
	if len(out) == 0 {
		return append([]SourceItem(nil), a.items...)
	}
	return out
}
func (a *App) hasExplicitSelection() bool {
	for _, it := range a.items {
		if it.Selected {
			return true
		}
	}
	return false
}

func (a *App) addItem(p string) bool {
	p = absPath(p)
	k := kindOf(p)
	if k == "" {
		return false
	}
	for _, it := range a.items {
		if norm(it.Path) == norm(p) {
			return false
		}
	}
	a.items = append(a.items, SourceItem{Path: p, Kind: k})
	return true
}
func (a *App) addFiles() {
	size := 65536
	for size <= 4*1024*1024 {
		buf := make([]uint16, size)
		filter := utf16MultiString(a.tr("dialog.all_files"), "*.*")
		ofn := OPENFILENAMEW{LStructSize: uint32(unsafe.Sizeof(OPENFILENAMEW{})), HwndOwner: a.hwnd, LpstrFilter: &filter[0], LpstrFile: &buf[0], NMaxFile: uint32(len(buf)), Flags: OFN_EXPLORER | OFN_FILEMUSTEXIST | OFN_PATHMUSTEXIST | OFN_ALLOWMULTISELECT | OFN_HIDEREADONLY, NFilterIndex: 1}
		r, _, _ := pGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
		if r != 0 {
			vals := splitMulti(buf)
			if len(vals) == 1 {
				a.addItem(vals[0])
			} else if len(vals) > 1 {
				dir := vals[0]
				for _, name := range vals[1:] {
					a.addItem(filepath.Join(dir, name))
				}
			}
			a.ensureSelectionState()
			a.refresh()
			return
		}
		err, _, _ := pCommDlgExtendedError.Call()
		if err == 0 {
			return
		}
		if err == FNERR_BUFFERTOOSMALL {
			size *= 2
			continue
		}
		message(a.hwnd, appTitle, fmt.Sprintf(a.tr("dialog.addfile_error"), err), 0x10)
		return
	}
	message(a.hwnd, appTitle, a.tr("dialog.addfile_too_large"), 0x30)
}
func splitMulti(buf []uint16) []string {
	vals := []string{}
	start := 0
	for i, v := range buf {
		if v == 0 {
			if i > start {
				vals = append(vals, syscall.UTF16ToString(buf[start:i]))
			} else {
				break
			}
			start = i + 1
		}
	}
	return vals
}
func (a *App) addFolder() {
	if p := pickFolder(a.hwnd, a.tr("dialog.addfolder_title")); p != "" {
		if a.addItem(p) {
			a.ensureSelectionState()
			a.refresh()
		}
	}
}
func pickFolder(owner HWND, title string) string {
	disp := make([]uint16, 260)
	bi := BROWSEINFOW{HwndOwner: owner, PszDisplayName: &disp[0], LpszTitle: u16(title), UlFlags: BIF_RETURNONLYFSDIRS | BIF_NEWDIALOGSTYLE | BIF_EDITBOX | BIF_VALIDATE}
	pidl, _, _ := pSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl == 0 {
		return ""
	}
	defer pCoTaskMemFree.Call(pidl)
	out := make([]uint16, 32768)
	ok, _, _ := pSHGetPathFromIDListEx.Call(pidl, uintptr(unsafe.Pointer(&out[0])), uintptr(len(out)), 0)
	if ok == 0 {
		return ""
	}
	return syscall.UTF16ToString(out)
}
func getClipboardText() string {
	if r, _, _ := pOpenClipboard.Call(0); r == 0 {
		return ""
	}
	defer pCloseClipboard.Call()
	h, _, _ := pGetClipboardData.Call(CF_UNICODETEXT)
	if h == 0 {
		return ""
	}
	sz, _, _ := pGlobalSize.Call(h)
	if sz < 2 || sz > 64*1024*1024 {
		return ""
	}
	p, _, _ := pGlobalLock.Call(h)
	if p == 0 {
		return ""
	}
	defer pGlobalUnlock.Call(h)
	units := int(sz / 2)
	buf := make([]uint16, units)
	pRtlMoveMemory.Call(uintptr(unsafe.Pointer(&buf[0])), p, uintptr(units*2))
	return syscall.UTF16ToString(buf)
}
func (a *App) pastePaths() {
	raw := getClipboardText()
	if raw == "" {
		message(a.hwnd, appTitle, a.tr("dialog.paste_empty"), 0x40)
		return
	}
	n := 0
	for _, line := range strings.Split(strings.ReplaceAll(raw, "\r\n", "\n"), "\n") {
		s := trimText(line)
		if s != "" && a.addItem(s) {
			n++
		}
	}
	a.ensureSelectionState()
	a.refresh()
	message(a.hwnd, appTitle, fmt.Sprintf(a.tr("dialog.paste_result"), n), 0x40)
}
func (a *App) chooseDestination() {
	if p := pickFolder(a.hwnd, a.tr("destination.pick_title")); p != "" {
		setText(a.destination, p)
		a.refresh()
	}
}
func getText(h HWND) string {
	n, _, _ := pSendMessageW.Call(uintptr(h), WM_GETTEXTLENGTH, 0, 0)
	if n == 0 {
		return ""
	}
	buf := make([]uint16, int(n)+1)
	pSendMessageW.Call(uintptr(h), WM_GETTEXT, uintptr(len(buf)), uintptr(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}
func setText(h HWND, s string) {
	p := u16(s)
	pSendMessageW.Call(uintptr(h), WM_SETTEXT, 0, uintptr(unsafe.Pointer(p)))
}
func (a *App) linkType() string {
	switch a.selectedType {
	case 1:
		return "hard"
	case 2:
		return "junction"
	}
	return "symbolic"
}

func (a *App) preflight() (bool, []string, []string, []SourceItem, string) {
	items := a.selectedItems()
	dest := absPath(getText(a.destination))
	lt := a.linkType()
	issues := []string{}
	preview := []string{}
	if len(items) == 0 {
		issues = append(issues, a.tr("preflight.no_source"))
	}
	if dest == "" {
		issues = append(issues, a.tr("preflight.no_destination"))
	} else if !isDir(dest) {
		issues = append(issues, a.tr("preflight.destination_not_dir"))
	}
	seen := map[string]bool{}
	for _, it := range items {
		cur := kindOf(it.Path)
		if cur == "" {
			issues = append(issues, fmt.Sprintf(a.tr("preflight.source_unavailable"), it.Path))
			continue
		}
		if cur != it.Kind {
			issues = append(issues, fmt.Sprintf(a.tr("preflight.source_kind_changed"), it.Path))
		}
		if lt == "hard" && it.Kind == "folder" {
			issues = append(issues, fmt.Sprintf(a.tr("preflight.hard_file_only"), it.Path))
		}
		if lt == "junction" && it.Kind == "file" {
			issues = append(issues, fmt.Sprintf(a.tr("preflight.junction_folder_only"), it.Path))
		}
	}
	if isDir(dest) {
		for _, it := range items {
			name := displayName(it.Path)
			link := filepath.Join(dest, name)
			nk := norm(link)
			if seen[nk] {
				issues = append(issues, fmt.Sprintf(a.tr("preflight.duplicate_target"), link))
			}
			seen[nk] = true
			if norm(link) == norm(it.Path) {
				issues = append(issues, fmt.Sprintf(a.tr("preflight.self_target"), link))
			}
			if _, err := os.Lstat(link); err == nil {
				issues = append(issues, fmt.Sprintf(a.tr("preflight.target_exists"), link))
			}
			if lt == "hard" && volumeKey(it.Path) != volumeKey(dest) {
				issues = append(issues, fmt.Sprintf(a.tr("preflight.hard_cross_volume"), it.Path))
			}
			if it.Kind == "folder" && (lt == "junction" || lt == "symbolic") && isSubpath(dest, it.Path) {
				issues = append(issues, fmt.Sprintf(a.tr("preflight.recursive"), it.Path, dest))
			}
		}
	}
	if len(issues) == 0 {
		for _, it := range items {
			link := filepath.Join(dest, displayName(it.Path))
			switch lt {
			case "symbolic":
				preview = append(preview, fmt.Sprintf(`%s  →  %s`, link, it.Path))
			case "hard":
				preview = append(preview, fmt.Sprintf(`%s  →  %s`, link, it.Path))
			case "junction":
				preview = append(preview, fmt.Sprintf(`%s  →  %s`, link, it.Path))
			}
		}
	}
	return len(issues) == 0, issues, preview, items, dest
}
func (a *App) refresh() {
	ok, issues, _, items, _ := a.preflight()
	a.compatOK = ok
	a.updateSourceScrollbars()
	if len(items) == 0 {
		a.status = a.tr("source.status.empty")
	} else if ok {
		a.status = fmt.Sprintf(a.tr("source.status.ready"), len(a.items), len(items))
	} else {
		a.status = fmt.Sprintf(a.tr("source.status.invalid"), len(a.items))
	}
	_ = issues
	pInvalidateRect.Call(uintptr(a.hwnd), 0, 1)
}

func formatSystemMessage(code uint32) string {
	if code == 0 {
		return ""
	}
	buf := make([]uint16, 1024)
	n, _, _ := pFormatMessageW.Call(FORMAT_MESSAGE_FROM_SYSTEM|FORMAT_MESSAGE_IGNORE_INSERTS, 0, uintptr(code), 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), 0)
	if n == 0 {
		return ""
	}
	return strings.TrimSpace(syscall.UTF16ToString(buf[:n]))
}
func errnoCode(e error) uint32 {
	if e == nil {
		return 0
	}
	if n, ok := e.(syscall.Errno); ok {
		return uint32(n)
	}
	return 0
}
func formatWinError(prefix string, e error) string {
	code := errnoCode(e)
	sys := formatSystemMessage(code)
	if code != 0 && sys != "" {
		return fmt.Sprintf(trText("winerror.with_code"), prefix, sys, code)
	}
	return fmt.Sprintf(trText("winerror.with_error"), prefix, e)
}
func isProcessElevated() bool {
	var token uintptr
	r, _, _ := pOpenProcessToken.Call(^uintptr(0), TOKEN_QUERY, uintptr(unsafe.Pointer(&token)))
	if r == 0 || token == 0 {
		return false
	}
	defer pCloseHandle.Call(token)
	var elevation TOKEN_ELEVATION
	var retLen uint32
	r, _, _ = pGetTokenInformation.Call(token, TokenElevation, uintptr(unsafe.Pointer(&elevation)), uintptr(unsafe.Sizeof(elevation)), uintptr(unsafe.Pointer(&retLen)))
	return r != 0 && elevation.TokenIsElevated != 0
}

func destinationWriteProbe(dest string) (bool, error) {
	if dest == "" || !isDir(dest) {
		return false, fmt.Errorf("%s", trText("winerror.probe_invalid"))
	}
	f, err := os.CreateTemp(dest, ".mkLINK-access-test-*")
	if err != nil {
		return false, err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if err := f.Close(); err != nil {
		return false, err
	}
	return true, nil
}

func confirmCreatePrivilege(owner HWND, dest, lt string) bool {
	elevated := isProcessElevated()
	writeOK, writeErr := destinationWriteProbe(dest)

	if lt == "symbolic" && !elevated {
		msg := trText("permission.symbolic_unprivileged")
		if !writeOK && writeErr != nil {
			msg += fmt.Sprintf(trText("permission.symbolic_write_fail_extra"), writeErr.Error())
		}
		r, _, _ := pMessageBoxW.Call(uintptr(owner), uintptr(unsafe.Pointer(u16(msg))), uintptr(unsafe.Pointer(u16(trText("permission.title")))), MB_YESNO|MB_ICONWARNING)
		return r == IDYES
	}

	if !writeOK {
		if writeErr == nil {
			writeErr = fmt.Errorf("%s", trText("winerror.probe_invalid"))
		}
		if !elevated {
			msg := fmt.Sprintf(trText("permission.write_fail_unprivileged"), writeErr.Error())
			r, _, _ := pMessageBoxW.Call(uintptr(owner), uintptr(unsafe.Pointer(u16(msg))), uintptr(unsafe.Pointer(u16(trText("permission.title")))), MB_YESNO|MB_ICONWARNING)
			return r == IDYES
		}
		msg := fmt.Sprintf(trText("permission.write_fail_elevated"), writeErr.Error())
		r, _, _ := pMessageBoxW.Call(uintptr(owner), uintptr(unsafe.Pointer(u16(msg))), uintptr(unsafe.Pointer(u16(trText("permission.title")))), MB_YESNO|MB_ICONWARNING)
		return r == IDYES
	}

	return true
}

func createSymbolicLink(link, source, kind string) error {
	flags := uint32(SYMBOLIC_LINK_FLAG_UNPRIVILEGED_CREATE)
	if kind == "folder" {
		flags |= SYMBOLIC_LINK_FLAG_DIRECTORY
	}
	r, _, e := pCreateSymbolicLinkW.Call(uintptr(unsafe.Pointer(u16(link))), uintptr(unsafe.Pointer(u16(source))), uintptr(flags))
	if r == 0 {
		return fmt.Errorf("%s", formatWinError(trText("winerror.symbolic"), e))
	}
	return nil
}
func createHardLink(link, source string) error {
	r, _, e := pCreateHardLinkW.Call(uintptr(unsafe.Pointer(u16(link))), uintptr(unsafe.Pointer(u16(source))), 0)
	if r == 0 {
		return fmt.Errorf("%s", formatWinError(trText("winerror.hard"), e))
	}
	return nil
}
func junctionSubstitutePath(target string) string {
	p := filepath.Clean(target)
	if strings.HasPrefix(p, `\\?\UNC\`) {
		return `\??\UNC\` + p[len(`\\?\`):]
	}
	if strings.HasPrefix(p, `\\?\`) {
		return `\??\` + p[len(`\\?\`):]
	}
	if strings.HasPrefix(p, `\\`) {
		return `\??\UNC\` + strings.TrimPrefix(p, `\\`)
	}
	return `\??\` + p
}
func buildMountPointReparseBuffer(target string) []byte {
	sub := utf16NoNul(junctionSubstitutePath(target))
	pr := utf16NoNul(filepath.Clean(target))
	path := append([]uint16{}, sub...)
	path = append(path, 0)
	printOffset := uint16(len(path) * 2)
	path = append(path, pr...)
	path = append(path, 0)
	dataLen := uint16(8 + len(path)*2)
	buf := make([]byte, 16+len(path)*2)
	binary.LittleEndian.PutUint32(buf[0:4], IO_REPARSE_TAG_MOUNT_POINT)
	binary.LittleEndian.PutUint16(buf[4:6], dataLen)
	binary.LittleEndian.PutUint16(buf[6:8], 0)
	binary.LittleEndian.PutUint16(buf[8:10], 0)
	binary.LittleEndian.PutUint16(buf[10:12], uint16(len(sub)*2))
	binary.LittleEndian.PutUint16(buf[12:14], printOffset)
	binary.LittleEndian.PutUint16(buf[14:16], uint16(len(pr)*2))
	for i, v := range path {
		binary.LittleEndian.PutUint16(buf[16+i*2:], v)
	}
	return buf
}
func createJunction(link, source string) error {
	r, _, e := pCreateDirectoryW.Call(uintptr(unsafe.Pointer(u16(link))), 0)
	if r == 0 {
		return fmt.Errorf("%s", formatWinError(trText("winerror.mkdir"), e))
	}
	ok := false
	defer func() {
		if !ok {
			pRemoveDirectoryW.Call(uintptr(unsafe.Pointer(u16(link))))
		}
	}()
	h, _, e := pCreateFileW.Call(uintptr(unsafe.Pointer(u16(link))), GENERIC_WRITE, FILE_SHARE_READ|FILE_SHARE_WRITE|FILE_SHARE_DELETE, 0, OPEN_EXISTING, FILE_FLAG_OPEN_REPARSE_POINT|FILE_FLAG_BACKUP_SEMANTICS, 0)
	if h == INVALID_HANDLE_VALUE {
		return fmt.Errorf("%s", formatWinError(trText("winerror.open"), e))
	}
	defer pCloseHandle.Call(h)
	buf := buildMountPointReparseBuffer(source)
	if len(buf) > 16384 {
		return fmt.Errorf("%s", trText("winerror.reparse_size"))
	}
	var out uint32
	r, _, e = pDeviceIoControl.Call(h, FSCTL_SET_REPARSE_POINT, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), 0, 0, uintptr(unsafe.Pointer(&out)), 0)
	if r == 0 {
		return fmt.Errorf("%s", formatWinError(trText("winerror.reparse"), e))
	}
	ok = true
	return nil
}
func runLink(link, source, lt, kind string) error {
	switch lt {
	case "symbolic":
		return createSymbolicLink(link, source, kind)
	case "hard":
		return createHardLink(link, source)
	case "junction":
		return createJunction(link, source)
	}
	return fmt.Errorf("%s", trText("winerror.unknown"))
}

func (a *App) createLinks() {
	ok, issues, _, items, dest := a.preflight()
	if !ok {
		message(a.hwnd, a.tr("create.error_title"), strings.Join(issues, "\n\n"), 0x10)
		return
	}
	if !confirmCreatePrivilege(a.hwnd, dest, a.linkType()) {
		a.status = a.tr("permission.cancelled")
		pInvalidateRect.Call(uintptr(a.hwnd), 0, 0)
		return
	}
	success, fail := 0, 0
	lines := []string{}
	succeeded := map[string]bool{}
	for _, it := range items {
		link := filepath.Join(dest, displayName(it.Path))
		if err := runLink(link, it.Path, a.linkType(), it.Kind); err != nil {
			fail++
			lines = append(lines, fmt.Sprintf(a.tr("create.fail_line"), link, it.Path, err))
			continue
		}
		if _, err := os.Lstat(link); err != nil {
			fail++
			lines = append(lines, fmt.Sprintf(a.tr("create.verify_fail"), link, err))
			continue
		}
		success++
		succeeded[norm(it.Path)] = true
		lines = append(lines, fmt.Sprintf(a.tr("create.success_line"), link, it.Path))
	}
	if success > 0 {
		keep := make([]SourceItem, 0, len(a.items))
		for _, it := range a.items {
			if !succeeded[norm(it.Path)] {
				keep = append(keep, it)
			}
		}
		a.items = keep
	}
	for i := range a.items {
		a.items[i].Selected = false
	}
	a.refresh()
	if fail == 0 {
		message(a.hwnd, a.tr("create.complete_title"), fmt.Sprintf(a.tr("create.complete"), success, strings.Join(lines, "\r\n\r\n")), 0x40)
	} else {
		message(a.hwnd, a.tr("create.result_title"), fmt.Sprintf(a.tr("create.result"), success, fail, strings.Join(lines, "\r\n\r\n")), 0x10)
	}
}

func (a *App) removeSelected() {
	if len(a.items) == 0 {
		return
	}
	explicit := a.hasExplicitSelection()
	if !explicit {
		message(a.hwnd, appTitle, a.tr("dialog.remove_none"), 0x40)
		return
	}
	anchorPath := ""
	if a.anchorIndex >= 0 && a.anchorIndex < len(a.items) {
		anchorPath = norm(a.items[a.anchorIndex].Path)
	}
	keep := make([]SourceItem, 0, len(a.items))
	for _, it := range a.items {
		if !it.Selected {
			keep = append(keep, it)
		}
	}
	a.items = keep
	a.clampSourceScroll()
	a.anchorIndex = -1
	if anchorPath != "" {
		for i := range a.items {
			if norm(a.items[i].Path) == anchorPath {
				a.anchorIndex = i
				break
			}
		}
	}
	a.refresh()
}
func (a *App) clearItems() {
	a.items = nil
	a.scroll = 0
	a.hScroll = 0
	a.anchorIndex = -1
	a.refresh()
}

func (a *App) runAsAdministrator() {
	if isProcessElevated() {
		message(a.hwnd, a.tr("admin.confirm_title"), a.tr("admin.already"), 0x40)
		return
	}
	r, _, _ := pMessageBoxW.Call(uintptr(a.hwnd), uintptr(unsafe.Pointer(u16(a.tr("admin.confirm")))), uintptr(unsafe.Pointer(u16(a.tr("admin.confirm_title")))), MB_YESNO|MB_ICONWARNING)
	if r != IDYES {
		return
	}
	exe, err := os.Executable()
	if err != nil || exe == "" {
		message(a.hwnd, a.tr("admin.confirm_title"), a.tr("admin.launch_fail"), 0x10)
		return
	}
	verb := u16("runas")
	file := u16(exe)
	result, _, _ := pShellExecuteW.Call(uintptr(a.hwnd), uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), 0, 0, SW_SHOWNORMAL)
	if result <= 32 {
		message(a.hwnd, a.tr("admin.confirm_title"), a.tr("admin.launch_fail"), 0x10)
		return
	}
	pDestroyWindow.Call(uintptr(a.hwnd))
}

func (a *App) showLanguageMenu() {
	menu, _, _ := pCreatePopupMenu.Call()
	if menu == 0 {
		return
	}
	defer pDestroyMenu.Call(menu)
	items := []struct {
		id     uint32
		text   string
		locale string
	}{
		{IDM_LANG_ZHTW, a.tr("language.zhTW"), "zh-TW"},
		{IDM_LANG_ENUS, a.tr("language.enUS"), "en-US"},
		{IDM_LANG_JAJP, a.tr("language.jaJP"), "ja-JP"},
		{IDM_LANG_ZHCN, a.tr("language.zhCN"), "zh-CN"},
	}
	for _, item := range items {
		flags := uint32(MF_STRING)
		if a.languageFile == "" && a.languageLocale == item.locale {
			flags |= MF_CHECKED
		}
		pAppendMenuW.Call(menu, uintptr(flags), uintptr(item.id), uintptr(unsafe.Pointer(u16(item.text))))
	}
	pAppendMenuW.Call(menu, uintptr(MF_SEPARATOR), 0, 0)
	pAppendMenuW.Call(menu, uintptr(MF_STRING), uintptr(IDM_LANG_CUSTOM), uintptr(unsafe.Pointer(u16(a.tr("language.custom")))))
	var pt POINT
	pGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	cmd, _, _ := pTrackPopupMenu.Call(menu, TPM_RETURNCMD|TPM_RIGHTBUTTON, uintptr(pt.X), uintptr(pt.Y), 0, uintptr(a.hwnd), 0)
	switch cmd {
	case IDM_LANG_ZHTW:
		a.selectLanguageLocale("zh-TW")
	case IDM_LANG_ENUS:
		a.selectLanguageLocale("en-US")
	case IDM_LANG_JAJP:
		a.selectLanguageLocale("ja-JP")
	case IDM_LANG_ZHCN:
		a.selectLanguageLocale("zh-CN")
	case IDM_LANG_CUSTOM:
		a.chooseExternalLanguageFile()
	}
}

func (a *App) openCreatorLink() {
	r, _, _ := pShellExecuteW.Call(uintptr(a.hwnd), uintptr(unsafe.Pointer(u16("open"))), uintptr(unsafe.Pointer(u16(creatorURL))), 0, 0, SW_SHOWNORMAL)
	if r <= 32 {
		message(a.hwnd, appTitle, a.tr("dialog.creator_fail"), 0x10)
	}
}
func (a *App) showSourceContextMenu(x, y int32) {
	if a == nil {
		return
	}
	// Like Explorer, right-clicking an unselected row makes that row the active
	// selection; right-clicking an already-selected row preserves multi-select.
	if idx := a.sourceIndexAt(x, y); idx >= 0 && !a.items[idx].Selected {
		for i := range a.items {
			a.items[i].Selected = false
		}
		a.items[idx].Selected = true
		a.refresh()
	}
	a.sourceFocus = true
	setFocus(a.hwnd)
	menu, _, _ := pCreatePopupMenu.Call()
	if menu == 0 {
		return
	}
	defer pDestroyMenu.Call(menu)
	removeFlags := uint32(MF_STRING)
	if !a.hasExplicitSelection() {
		removeFlags |= MF_GRAYED
	}
	pAppendMenuW.Call(menu, uintptr(MF_STRING), IDM_PASTE_PATH, uintptr(unsafe.Pointer(u16(a.tr("context.paste")))))
	pAppendMenuW.Call(menu, uintptr(MF_SEPARATOR), 0, 0)
	pAppendMenuW.Call(menu, uintptr(removeFlags), IDM_REMOVE_SELECTED, uintptr(unsafe.Pointer(u16(a.tr("context.remove")))))
	var pt POINT
	pGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	cmd, _, _ := pTrackPopupMenu.Call(menu, TPM_RIGHTBUTTON|TPM_RETURNCMD, uintptr(pt.X), uintptr(pt.Y), 0, uintptr(a.hwnd), 0)
	switch cmd {
	case IDM_PASTE_PATH:
		a.pastePaths()
	case IDM_REMOVE_SELECTED:
		a.removeSelected()
	}
}

func (a *App) handleDrop(hDrop uintptr) {
	defer pDragFinish.Call(hDrop)
	count, _, _ := pDragQueryFileW.Call(hDrop, ^uintptr(0), 0, 0)
	added := 0
	for i := uintptr(0); i < count; i++ {
		need, _, _ := pDragQueryFileW.Call(hDrop, i, 0, 0)
		if need == 0 {
			continue
		}
		buf := make([]uint16, int(need)+1)
		n, _, _ := pDragQueryFileW.Call(hDrop, i, uintptr(unsafe.Pointer(&buf[0])), need+1)
		if n > 0 && a.addItem(syscall.UTF16ToString(buf[:n])) {
			added++
		}
	}
	a.ensureSelectionState()
	if added == 0 {
		a.status = a.tr("source.status.drop.none")
	} else {
		a.status = fmt.Sprintf(a.tr("source.status.drop"), added)
	}
	a.refresh()
}

func wndProc(hwnd HWND, msg uint32, wParam, lParam uintptr) (ret uintptr) {
	defer func() {
		if r := recover(); r != nil {
			message(hwnd, appTitle, fmt.Sprintf(trText("dialog.gui_error"), r), 0x10)
			ret = 0
		}
	}()
	a := globalApp
	switch msg {
	case WM_ERASEBKGND:
		return 1
	case WM_GETMINMAXINFO:
		if lParam != 0 {
			// WM_GETMINMAXINFO is a same-process message: lParam already points at
			// the real MINMAXINFO the system pre-filled with sane defaults
			// (ptMaxSize/ptMaxPosition/ptMaxTrackSize based on the monitor work
			// area). We only want to raise ptMinTrackSize, so edit the struct in
			// place instead of overwriting it with a zero-valued copy — the old
			// pRtlMoveMemory round-trip zeroed ptMaxSize/ptMaxTrackSize, which made
			// double-clicking the title bar (or Win+Up) "maximize" the window to
			// 0x0 instead of actually maximizing it.
			mi := (*MINMAXINFO)(unsafe.Pointer(lParam))
			dpi := int32(96)
			if a != nil && a.dpi > 0 {
				dpi = a.dpi
			}
			mi.PtMinTrackSize = POINT{X: 980 * dpi / 96, Y: 720 * dpi / 96}
			mi.PtMaxSize = POINT{X: 1105 * dpi / 96, Y: 805 * dpi / 96}
			mi.PtMaxTrackSize = POINT{X: 1105 * dpi / 96, Y: 805 * dpi / 96}
		}
		return 0
	case WM_PAINT:
		if a != nil {
			a.paint()
			return 0
		}
	case WM_SIZE, WM_DPICHANGED:
		if a != nil {
			if msg == WM_DPICHANGED {
				if r, _, _ := pGetDpiForWindow.Call(uintptr(hwnd)); r != 0 {
					a.dpi = int32(r)
				}
				a.releaseFonts()
				a.initFonts()
				toolbarFont = a.buttonFont
			}
			a.layout()
			return 0
		}
	case WM_MOUSEMOVE:
		if a != nil {
			x, y := pointFromLParam(lParam)
			if a.columnDragMode != columnDragNone {
				a.updateColumnDrag(x)
				return 0
			}
			old := a.hoverID
			a.hoverID = a.hitTest(x, y)
			if old != a.hoverID {
				pInvalidateRect.Call(uintptr(hwnd), 0, 0)
			}
			return 0
		}
	case WM_MOUSEWHEEL:
		if a != nil {
			var pt POINT
			pt.X = int32(int16(lParam & 0xffff))
			pt.Y = int32(int16((lParam >> 16) & 0xffff))
			pScreenToClient.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pt)))
			if !pointIn(a.sourceListArea, pt.X, pt.Y) {
				break
			}
			d := int16(uint16((wParam >> 16) & 0xffff))
			_, visible, maxScroll := a.sourceMetrics()
			step := 2
			if visible < 8 {
				step = 1
			}
			if d > 0 {
				a.scroll -= step
			} else if d < 0 {
				a.scroll += step
			}
			if a.scroll < 0 {
				a.scroll = 0
			}
			if a.scroll > maxScroll {
				a.scroll = maxScroll
			}
			a.updateSourceScrollbars()
			pInvalidateRect.Call(uintptr(hwnd), 0, 0)
			return 0
		}
	case WM_LBUTTONDOWN:
		if a != nil {
			x, y := pointFromLParam(lParam)
			if a.beginColumnDrag(x, y) {
				return 0
			}
			a.pressedID = a.hitTest(x, y)
			a.sourceFocus = (a.pressedID == 60)
			return 0
		}
	case WM_RBUTTONUP:
		if a != nil {
			x, y := pointFromLParam(lParam)
			if pointIn(a.sourceRect, x, y) {
				a.showSourceContextMenu(x, y)
				return 0
			}
		}
	case WM_HSCROLL:
		if a != nil && a.sourceHScrollbar != 0 && HWND(lParam) == a.sourceHScrollbar {
			code := int(loword(wParam))
			maxScroll := a.sourceHorizontalMaxScroll()
			switch code {
			case SB_LINELEFT:
				a.hScroll--
			case SB_LINERIGHT:
				a.hScroll++
			case SB_PAGELEFT:
				a.hScroll -= maxInt(1, a.sourceDataViewportWidth()/2)
			case SB_PAGERIGHT:
				a.hScroll += maxInt(1, a.sourceDataViewportWidth()/2)
			case SB_TOP:
				a.hScroll = 0
			case SB_BOTTOM:
				a.hScroll = maxScroll
			case SB_THUMBPOSITION, SB_THUMBTRACK:
				var info SCROLLINFO
				info.CbSize = uint32(unsafe.Sizeof(SCROLLINFO{}))
				info.FMask = SIF_TRACKPOS | SIF_POS
				if r, _, _ := pGetScrollInfo.Call(uintptr(a.sourceHScrollbar), SB_CTL, uintptr(unsafe.Pointer(&info))); r != 0 {
					a.hScroll = info.NTrackPos
				}
			}
			a.clampSourceHScroll()
			a.updateSourceScrollbars()
			pInvalidateRect.Call(uintptr(hwnd), 0, 0)
			return 0
		}
	case WM_VSCROLL:
		if a != nil && a.sourceScrollbar != 0 && HWND(lParam) == a.sourceScrollbar {
			code := int(loword(wParam))
			_, visible, maxScroll := a.sourceMetrics()
			switch code {
			case SB_LINEUP:
				a.scroll--
			case SB_LINEDOWN:
				a.scroll++
			case SB_PAGEUP:
				a.scroll -= visible
			case SB_PAGEDOWN:
				a.scroll += visible
			case SB_TOP:
				a.scroll = 0
			case SB_BOTTOM:
				a.scroll = maxScroll
			case SB_THUMBPOSITION, SB_THUMBTRACK:
				var info SCROLLINFO
				info.CbSize = uint32(unsafe.Sizeof(SCROLLINFO{}))
				info.FMask = SIF_TRACKPOS | SIF_POS
				if r, _, _ := pGetScrollInfo.Call(uintptr(a.sourceScrollbar), SB_CTL, uintptr(unsafe.Pointer(&info))); r != 0 {
					a.scroll = int(info.NTrackPos)
				}
			}
			a.clampSourceScroll()
			a.updateSourceScrollbars()
			pInvalidateRect.Call(uintptr(hwnd), 0, 0)
			return 0
		}
	case WM_LBUTTONUP:
		if a != nil {
			x, y := pointFromLParam(lParam)
			if a.columnDragMode != columnDragNone {
				a.finishColumnDrag()
				a.pressedID = 0
				pInvalidateRect.Call(uintptr(hwnd), 0, 0)
				return 0
			}
			h := a.hitTest(x, y)
			if h != 0 && h == a.pressedID {
				a.wndCommand(int32(h), x, y)
			}
			a.pressedID = 0
			pInvalidateRect.Call(uintptr(hwnd), 0, 0)
			return 0
		}
	case WM_KEYDOWN:
		if a != nil {
			switch wParam {
			case VK_A:
				if ctrlKeyDown() {
					a.selectAll()
				}
			case VK_V:
				if ctrlKeyDown() && a.sourceFocus {
					a.pastePaths()
				}
			case VK_DELETE:
				a.removeSelected()
			case VK_RETURN:
				a.createLinks()
			case VK_ESCAPE:
				for i := range a.items {
					a.items[i].Selected = false
				}
				a.refresh()
			}
			return 0
		}
	case WM_DROPFILES:
		if a != nil {
			a.handleDrop(wParam)
			return 0
		}
	case WM_COMMAND:
		if a != nil && HWND(lParam) == a.destination && hiword(wParam) == EN_CHANGE {
			// The destination EDIT control is native, so typing or pasting a path
			// directly into it (as opposed to using "瀏覽…") never invalidated the
			// window before: the 預檢結果／建立預覽 panels kept showing stale
			// results until some other action happened to call refresh().
			a.refresh()
			return 0
		}
	case WM_SETCURSOR:
		if a != nil && loword(lParam) == HTCLIENT {
			var pt POINT
			pGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
			pScreenToClient.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pt)))
			if a.columnDragMode == columnDragResize {
				if hSizeWECursor != 0 {
					pSetCursor.Call(hSizeWECursor)
				}
				return 1
			}
			if a.columnDragMode == columnDragReorder {
				if hHandCursor != 0 {
					pSetCursor.Call(hHandCursor)
				}
				return 1
			}
			if a.isSourceHeader(pt.X, pt.Y) {
				if a.sourceHeaderResizeAt(pt.X, pt.Y) {
					if hSizeWECursor != 0 {
						pSetCursor.Call(hSizeWECursor)
					}
				} else if hHandCursor != 0 {
					pSetCursor.Call(hHandCursor)
				}
				return 1
			}
			if a.hitTest(pt.X, pt.Y) != 0 {
				if hHandCursor != 0 {
					pSetCursor.Call(hHandCursor)
				}
				return 1
			}
		}
	case WM_CTLCOLOREDIT:
		if a != nil {
			pSetTextColor.Call(wParam, uintptr(cText))
			pSetBkColor.Call(wParam, uintptr(cSurface))
			if brInput != 0 {
				return uintptr(brInput)
			}
		}
	case WM_DESTROY:
		if a != nil {
			a.saveSettings()
		}
		if brInput != 0 {
			pDeleteObject.Call(uintptr(brInput))
			brInput = 0
		}
		if a != nil {
			a.releaseFonts()
		}
		releaseOcticons()
		pPostQuitMessage.Call(0)
		return 0
	case WM_CLOSE:
		pDestroyWindow.Call(uintptr(hwnd))
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return r
}
func ctrlKeyDown() bool { r, _, _ := pGetAsyncKeyState.Call(uintptr(VK_CONTROL)); return int16(r) < 0 }
func maxInt(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func (a *App) selectAll() {
	for i := range a.items {
		a.items[i].Selected = true
	}
	a.refresh()
}
func message(owner HWND, title, text string, flags uintptr) {
	pMessageBoxW.Call(uintptr(owner), uintptr(unsafe.Pointer(u16(text))), uintptr(unsafe.Pointer(u16(title))), flags)
}
func setFocus(h HWND) { pSetFocusProc.Call(uintptr(h)) }

var _ = os.ErrNotExist
