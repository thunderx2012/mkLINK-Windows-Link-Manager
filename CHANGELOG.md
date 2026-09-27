# Changelog

All notable user-visible changes to mkLINK are recorded in this file.

## [1.5.22] - 2026-09-26

### Fixed

- Corrected the Traditional Chinese link-method terminology from `資料夾連節點` to `資料夾連接點`.
- Synchronized the corrected terminology across the relevant interface text, permission messages, and project documentation.

### Preserved

- English link-method titles can wrap to a maximum of two lines within the existing cards and remain vertically centered beside their icons.
- External language-pack architecture with four bundled locales: `zh-TW`, `en-US`, `ja-JP`, and `zh-CN`.
- Standard Windows UAC relaunch through **Run as Administrator**.
- Source-table column reordering, column resizing, vertical scrolling, horizontal scrolling, and persisted layout settings.
- Destination write-access probing and temporary-file cleanup.

## [1.5.21]

### Fixed

- Link-method titles no longer become unnecessarily truncated in English when the fixed card width is too narrow.
- Titles may wrap naturally to two lines and remain vertically centered in the existing title area.

### Preserved

- Existing card dimensions, icon assignments, and the two-line usage descriptions below the titles.

## [1.5.20]

### Added

- External language-pack loading and language selection.
- Bundled Traditional Chinese, English, Japanese, and Simplified Chinese language packs.
- **Run as Administrator** with standard Windows UAC elevation.
- Configurable source-table column order and column widths.
- Horizontal source-table scrolling when the configured table width exceeds the visible area.
- Persistence for language and source-table settings.

### Removed

- The toolbar **Recheck** command was removed because preflight updates automatically when the relevant inputs change.
- **Paste Path** and **Remove Selected** were removed from the toolbar and remain available from the source area through keyboard shortcuts and the context menu.

## [1.5.19]

### Fixed

- Improved cleanup of the temporary file used by destination write-access probing by registering removal with `defer` immediately after the temporary file is created.

## [1.5.18]

### Added

- Destination write-access probing before link creation for Symbolic Link, Hard Link, and Directory Junction operations.
- Permission warnings for destinations that may require elevated privileges or additional NTFS/share permissions.

## [1.5.17]

### Added

- Pre-creation elevation warning for Symbolic Link operations.

### Changed

- Source-area guidance changed to `支援滑鼠拖曳`.
- Expanded the creator-information area so `@thunderx2012 · Threads` can be displayed completely.

## [1.5.16]

### Improved

- Standardized GDI resource management in drawing helpers.
- Replaced an internal placeholder radius helper with a named constant without changing the visual geometry.

## [1.5.15]

### Fixed

- Cached the `SetFocus` procedure lookup.
- Improved checkmark rendering for high-DPI layouts.
- Hardened GDI resource cleanup in rounded panels.
- Used decoded PNG dimensions for Octicon bitmap metadata.
- Improved Shift+Click range selection by introducing a persistent selection anchor.
- Preserved a reasonable source-list scroll position after removing selected entries.

---

Copyright (C) 2026 thunderx2012.
