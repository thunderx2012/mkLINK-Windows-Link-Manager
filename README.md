# mkLINK

**Windows Link Manager** — a lightweight Windows GUI for creating symbolic links, hard links, and directory junctions.

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)

## Overview

mkLINK provides a focused graphical workflow for selecting one or more files and folders, choosing a link type, reviewing the preflight result, and creating the requested links.

The project is designed for Windows users who want direct access to standard NTFS link operations without relying on a command prompt for each operation.

## Features

- Add files and folders as link sources.
- Add multiple sources in one operation.
- Drag files or folders into the source area with the mouse.
- Select and deselect sources with the left mouse button.
- Use `Ctrl+V`, the context menu, `Delete`, `Ctrl+A`, and `Shift+Click` for source-list operations.
- Create three Windows link types:
  - **Symbolic Link**
  - **Hard Link**
  - **Directory Junction**
- Preview the planned link operations before creation.
- Automatically check source compatibility and the destination before creation.
- Detect destination write-access problems before starting a batch.
- Relaunch the application through the standard Windows UAC flow when elevated privileges are required.
- Resize and reorder the source-table columns **Name**, **Type**, and **Path**.
- Use vertical scrolling and horizontal scrolling when the source list exceeds the visible area.
- Persist column layout and language selection in the user's configuration directory.

## Link types

| Type | Supported source | Typical purpose |
|---|---|---|
| Symbolic Link | Files and folders | Creates a symbolic reference to another path. |
| Hard Link | Files | Provides another directory entry for the same file data on the same volume. |
| Directory Junction | Folders | Provides a directory-level link to another folder. |

Compatibility checks are performed before creation. The application reports conditions such as an unsupported source type, an existing target, a cross-volume hard-link attempt, or a recursive directory relationship when applicable.

## Permissions

Creating a link may require permissions that are not available to the current process.

Before creation, mkLINK checks whether the destination can be written by the current process. The check applies to all three link types. When required, mkLINK displays a permission warning before the batch starts.

Symbolic-link creation also has a dedicated elevation check. The exact requirement depends on the Windows configuration and user privileges.

The **Run as Administrator** command uses the standard Windows elevation mechanism. mkLINK does not attempt to bypass UAC or modify NTFS permissions automatically.

## Source table

The source area contains three configurable columns:

- **Name**
- **Type**
- **Path**

Column headers can be dragged to change their order. Column separators can be dragged to change widths. The layout is retained between sessions.

The source list supports both vertical and horizontal scrolling when the visible area is insufficient.

When the list is empty, the source area displays the drag-and-drop guidance. The guidance disappears as soon as a source is added and returns when the list is cleared.

## Languages

mkLINK separates interface text and user-facing messages from the application logic through language packs.

The distribution includes these languages:

| Locale | Language |
|---|---|
| `zh-TW` | Traditional Chinese (Taiwan) |
| `en-US` | English |
| `ja-JP` | Japanese |
| `zh-CN` | Simplified Chinese |

Language files are JSON documents stored in `languages/`. Users can select a bundled language or load another compatible language file. The executable also contains embedded copies of the four bundled languages as a fallback.

Language-pack schema validation prevents an incomplete or malformed file from disabling the application. Missing language files or invalid custom files fall back safely to the built-in Traditional Chinese resources when necessary.

## UI and terminology

The three link-method names use the following terminology in the Traditional Chinese interface:

- `符號連結` — Symbolic Link
- `硬連結` — Hard Link
- `資料夾連接點` — Directory Junction

The source-area guidance is:

> 支援滑鼠拖曳

The application retains the creator information `@thunderx2012 · Threads` and its associated profile link.

## Icons

The interface uses Primer Octicons assets supplied with this project. The original SVG files are retained under `assets/octicons/`, together with the generated PNG assets used by the Win32/GDI renderer.

See [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) for the applicable third-party license information.

## Requirements

- Windows x64
- Go 1.23 or later for source builds
- No third-party Go modules are required by the project

## Building on Windows

Open Windows PowerShell in the project directory and run:

```powershell
.\在Windows建立EXE.ps1
```

The script configures a Windows x64 build and produces:

```text
mkLINK_v1.5.22.exe
```

For a manual build, the equivalent command is:

```powershell
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go test ./...
go build -trimpath -ldflags='-s -w -H=windowsgui' -o mkLINK_v1.5.22.exe .
```

## Project layout

```text
.
├── assets/
│   └── octicons/
├── languages/
│   ├── zh-TW.json
│   ├── en-US.json
│   ├── ja-JP.json
│   └── zh-CN.json
├── main.go
├── language_windows.go
├── source_table_windows.go
├── settings_windows.go
├── octicons_windows.go
├── createEXE.ps1
├── README.md
├── CHANGELOG.md
├── THIRD_PARTY_NOTICES.md
└── LICENSE
```

## License

The mkLINK source code and original project documentation are licensed under the **GNU General Public License, version 3.0**.

See the complete license text in [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

Third-party assets retain their own applicable licenses. In particular, the bundled Primer Octicons assets are distributed under the MIT License; see [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

---

Copyright (C) 2026 thunderx2012.
