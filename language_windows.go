//go:build windows

package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

type LanguagePack struct {
	SchemaVersion int               `json:"schema_version"`
	Locale        string            `json:"locale"`
	LanguageName  string            `json:"language_name"`
	Strings       map[string]string `json:"strings"`
}

// Four supported built-in language packs are also embedded as a safe fallback.
// The release ships the same JSON files externally under ./languages/ so users
// can replace them without modifying the executable.
//
//go:embed languages/*.json
var embeddedLanguageFS embed.FS

var languageLocales = []string{"zh-TW", "en-US", "ja-JP", "zh-CN"}

var requiredLanguageKeys = []string{
	"app.subtitle", "app.tagline", "header.language_hint",
	"toolbar.add_file", "toolbar.add_file_sub", "toolbar.add_folder", "toolbar.add_folder_sub", "toolbar.clear", "toolbar.clear_sub", "toolbar.admin", "toolbar.admin_sub",
	"source.title", "source.count", "source.description", "source.empty.title", "source.empty.subtitle", "source.col.name", "source.col.type", "source.col.path", "source.kind.file", "source.kind.folder", "source.selected.count",
	"source.status.empty", "source.status.selected", "source.status.ready", "source.status.invalid", "source.status.drop", "source.status.drop.none", "source.status.paste", "source.status.selection",
	"mode.title", "mode.subtitle", "mode.symbolic", "mode.hard", "mode.junction", "mode.symbolic_desc1", "mode.hard_desc1", "mode.junction_desc1", "mode.symbolic_desc2", "mode.hard_desc2", "mode.junction_desc2", "mode.hover.symbolic", "mode.hover.hard", "mode.hover.junction",
	"destination.title", "destination.placeholder", "destination.browse", "destination.pick_title",
	"preflight.title", "preflight.ok", "preflight.fail", "preview.title", "preview.none", "preview.issues", "preview.plan",
	"footer.ready", "footer.creator", "create.title", "create.dropdown",
	"dialog.addfile_error", "dialog.addfile_too_large", "dialog.all_files", "dialog.addfolder_title", "dialog.paste_empty", "dialog.paste_result", "dialog.remove_none", "dialog.creator_fail", "dialog.fatal_startup", "dialog.window_class_fail", "dialog.window_create_fail", "dialog.gui_error",
	"preflight.no_source", "preflight.no_destination", "preflight.destination_not_dir", "preflight.source_unavailable", "preflight.source_kind_changed", "preflight.hard_file_only", "preflight.junction_folder_only", "preflight.duplicate_target", "preflight.self_target", "preflight.target_exists", "preflight.hard_cross_volume", "preflight.recursive",
	"winerror.unknown", "winerror.symbolic", "winerror.hard", "winerror.mkdir", "winerror.open", "winerror.reparse_size", "winerror.reparse", "winerror.probe_invalid", "winerror.with_code", "winerror.with_error",
	"permission.title", "permission.symbolic_unprivileged", "permission.symbolic_write_fail_extra", "permission.write_fail_unprivileged", "permission.write_fail_elevated", "permission.cancelled", "permission.already_elevated",
	"admin.confirm", "admin.confirm_title", "admin.launch_fail", "admin.already",
	"create.error_title", "create.complete_title", "create.complete", "create.result_title", "create.result", "create.success_line", "create.fail_line", "create.verify_fail", "create.cancelled",
	"context.paste", "context.remove",
	"language.menu_title", "language.zhTW", "language.enUS", "language.jaJP", "language.zhCN", "language.custom", "language.file_filter", "language.file_error", "language.file_invalid", "language.loaded", "language.load_failed", "language.default",
	"admin.sub", "admin.cursor", "source.drag.cursor", "source.reorder.cursor",
}

func parseLanguageData(data []byte) (LanguagePack, error) {
	var p LanguagePack
	if err := json.Unmarshal(data, &p); err != nil {
		return LanguagePack{}, err
	}
	if p.SchemaVersion != 1 || p.LanguageName == "" || p.Locale == "" || p.Strings == nil {
		return LanguagePack{}, fmt.Errorf("invalid language pack schema")
	}
	for _, key := range requiredLanguageKeys {
		if _, ok := p.Strings[key]; !ok {
			return LanguagePack{}, fmt.Errorf("missing key: %s", key)
		}
	}
	return p, nil
}

func embeddedLanguage(locale string) (LanguagePack, error) {
	data, err := embeddedLanguageFS.ReadFile(filepath.ToSlash(filepath.Join("languages", locale+".json")))
	if err != nil {
		return LanguagePack{}, err
	}
	return parseLanguageData(data)
}

func loadLanguageFile(path string) (LanguagePack, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return LanguagePack{}, err
	}
	return parseLanguageData(data)
}

func (a *App) languageFilePath(locale string) string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(exe), "languages", locale+".json")
}

func (a *App) initLanguage(locale, customPath string) {
	fallback, err := embeddedLanguage("zh-TW")
	if err != nil {
		panic(fmt.Errorf("embedded zh-TW language pack unavailable: %w", err))
	}
	a.fallbackLang = fallback
	a.languageLocale = "zh-TW"
	a.languageFile = ""
	a.lang = fallback

	if customPath != "" {
		if p, err := loadLanguageFile(customPath); err == nil {
			a.lang = p
			a.languageLocale = p.Locale
			a.languageFile = customPath
			return
		}
	}
	for _, loc := range languageLocales {
		if loc != locale {
			continue
		}
		path := a.languageFilePath(loc)
		if path != "" {
			if p, err := loadLanguageFile(path); err == nil {
				a.lang = p
				a.languageLocale = p.Locale
				return
			}
		}
		if p, err := embeddedLanguage(loc); err == nil {
			a.lang = p
			a.languageLocale = p.Locale
			return
		}
	}
	// zh-TW is already installed as fallback.
}

func trText(key string) string {
	if globalApp != nil {
		return globalApp.tr(key)
	}
	if p, err := embeddedLanguage("zh-TW"); err == nil {
		if v, ok := p.Strings[key]; ok && v != "" {
			return v
		}
	}
	return key
}

func (a *App) tr(key string) string {
	if a != nil {
		if v, ok := a.lang.Strings[key]; ok && v != "" {
			return v
		}
		if v, ok := a.fallbackLang.Strings[key]; ok {
			return v
		}
	}
	return key
}

func (a *App) languageDisplayName() string {
	if a == nil {
		if p, err := embeddedLanguage("zh-TW"); err == nil && p.LanguageName != "" {
			return p.LanguageName
		}
		return "mkLINK"
	}
	if a.lang.LanguageName == "" {
		return a.tr("language.default")
	}
	return a.lang.LanguageName
}

func (a *App) selectLanguageLocale(locale string) {
	old := a.lang
	path := a.languageFilePath(locale)
	var p LanguagePack
	var err error
	if path != "" {
		p, err = loadLanguageFile(path)
	}
	if err != nil {
		p, err = embeddedLanguage(locale)
	}
	if err != nil {
		a.lang = old
		message(a.hwnd, old.LanguageName, fmt.Sprintf(a.tr("language.file_error"), err), 0x10)
		return
	}
	a.lang = p
	a.languageLocale = p.Locale
	a.languageFile = ""
	a.saveSettings()
	a.layout()
	a.status = fmt.Sprintf(a.tr("language.loaded"), p.LanguageName)
	pInvalidateRect.Call(uintptr(a.hwnd), 0, 1)
}

func (a *App) chooseExternalLanguageFile() {
	path := a.pickLanguageFile()
	if path == "" {
		return
	}
	p, err := loadLanguageFile(path)
	if err != nil {
		message(a.hwnd, a.tr("language.menu_title"), fmt.Sprintf(a.tr("language.file_invalid"), err), 0x10)
		return
	}
	a.lang = p
	a.languageLocale = p.Locale
	a.languageFile = path
	a.saveSettings()
	a.layout()
	a.status = fmt.Sprintf(a.tr("language.loaded"), p.LanguageName)
	pInvalidateRect.Call(uintptr(a.hwnd), 0, 1)
}

func (a *App) pickLanguageFile() string {
	buf := make([]uint16, 32768)
	raw := a.tr("language.file_filter")
	filter := make([]uint16, 0, len(raw)+1)
	for _, r := range raw {
		filter = append(filter, uint16(r))
	}
	filter = append(filter, 0)
	ofn := OPENFILENAMEW{
		LStructSize:  uint32(unsafe.Sizeof(OPENFILENAMEW{})),
		HwndOwner:    a.hwnd,
		LpstrFilter:  &filter[0],
		LpstrFile:    &buf[0],
		NMaxFile:     uint32(len(buf)),
		Flags:        OFN_EXPLORER | OFN_FILEMUSTEXIST | OFN_PATHMUSTEXIST | OFN_HIDEREADONLY,
		NFilterIndex: 1,
	}
	r, _, _ := pGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	if r == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}
