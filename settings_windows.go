//go:build windows

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type mkLINKSettings struct {
	Version            int      `json:"version"`
	LanguageLocale     string   `json:"language_locale"`
	CustomLanguageFile string   `json:"custom_language_file"`
	ColumnOrder        []string `json:"column_order"`
	ColumnWidths       []int32  `json:"column_widths"`
}

func defaultColumnOrder() [3]int {
	return [3]int{0, 1, 2}
}

func defaultColumnWidths() [3]int32 {
	// Logical 96-DPI units; tuned to match the existing table composition
	// without forcing the horizontal scrollbar on first launch.
	return [3]int32{130, 84, 245}
}

func settingsDir() string {
	dir, err := os.UserConfigDir()
	if err != nil || dir == "" {
		return ""
	}
	return filepath.Join(dir, "mkLINK")
}

func settingsPath() string {
	dir := settingsDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, "settings.json")
}

func (a *App) loadSettings() {
	a.columnOrder = defaultColumnOrder()
	a.columnWidths = defaultColumnWidths()
	a.languageLocale = "zh-TW"
	a.languageFile = ""

	path := settingsPath()
	if path == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var s mkLINKSettings
	if json.Unmarshal(data, &s) != nil || s.Version != 1 {
		return
	}
	if len(s.ColumnOrder) == 3 && validColumnOrder(s.ColumnOrder) {
		for i, v := range s.ColumnOrder {
			a.columnOrder[i] = columnIndexFromKey(v)
		}
	}
	if len(s.ColumnWidths) == 3 {
		for i, v := range s.ColumnWidths {
			if v >= minColumnWidth(i) && v <= maxColumnWidth(i) {
				a.columnWidths[i] = v
			}
		}
	}
	if s.LanguageLocale != "" {
		a.languageLocale = s.LanguageLocale
	}
	a.languageFile = s.CustomLanguageFile
}

func validColumnOrder(order []string) bool {
	seen := [3]bool{}
	for _, v := range order {
		switch v {
		case "name", "type", "path":
		default:
			return false
		}
		idx := columnIndexFromKey(v)
		if seen[idx] {
			return false
		}
		seen[idx] = true
	}
	return seen[0] && seen[1] && seen[2]
}

func columnIndexFromKey(key string) int {
	switch key {
	case "name":
		return 0
	case "type":
		return 1
	case "path":
		return 2
	default:
		return -1
	}
}

func columnKey(idx int) string {
	switch idx {
	case 0:
		return "name"
	case 1:
		return "type"
	case 2:
		return "path"
	default:
		return ""
	}
}

func minColumnWidth(idx int) int32 {
	switch idx {
	case 0:
		return 105
	case 1:
		return 78
	case 2:
		return 150
	default:
		return 100
	}
}

func maxColumnWidth(idx int) int32 {
	switch idx {
	case 0:
		return 600
	case 1:
		return 260
	case 2:
		return 900
	default:
		return 900
	}
}

func (a *App) saveSettings() {
	path := settingsPath()
	if path == "" {
		return
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return
	}
	order := make([]string, 0, 3)
	for _, idx := range a.columnOrder {
		order = append(order, columnKey(idx))
	}
	s := mkLINKSettings{
		Version:            1,
		LanguageLocale:     a.languageLocale,
		CustomLanguageFile: a.languageFile,
		ColumnOrder:        order,
		ColumnWidths:       []int32{a.columnWidths[0], a.columnWidths[1], a.columnWidths[2]},
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return
	}
	_ = os.Rename(tmp, path)
}
