//go:build windows

package main

import "unsafe"

const (
	columnName = iota
	columnType
	columnPath

	sourceTableFixedW  int32 = 94
	sourceHeaderH      int32 = 46
	sourceScrollbarH   int32 = 17
	sourceScrollbarGap int32 = 6

	columnDragNone    = 0
	columnDragResize  = 1
	columnDragReorder = 2
)

func (a *App) scaledColumnWidth(idx int) int32 {
	return a.sp(a.columnWidths[idx])
}

func (a *App) sourceDataViewportWidth() int32 {
	w := a.sourceListArea.Right - a.sourceListArea.Left - a.sp(sourceTableFixedW)
	if w < a.sp(40) {
		return a.sp(40)
	}
	return w
}

func (a *App) sourceColumnsContentWidth() int32 {
	total := int32(0)
	for i := 0; i < 3; i++ {
		total += a.scaledColumnWidth(i)
	}
	return total
}

func (a *App) sourceHorizontalMaxScroll() int32 {
	max := a.sourceColumnsContentWidth() - a.sourceDataViewportWidth()
	if max < 0 {
		return 0
	}
	return max
}

func (a *App) clampSourceHScroll() {
	max := a.sourceHorizontalMaxScroll()
	if a.hScroll < 0 {
		a.hScroll = 0
	}
	if a.hScroll > max {
		a.hScroll = max
	}
}

func (a *App) sourceColumnRects() [3]RECT {
	var out [3]RECT
	x := a.sourceListArea.Left + a.sp(sourceTableFixedW) - a.hScroll
	for pos := 0; pos < len(a.columnOrder); pos++ {
		idx := a.columnOrder[pos]
		w := a.scaledColumnWidth(idx)
		out[idx] = RECT{Left: x, Top: a.sourceListArea.Top, Right: x + w, Bottom: a.sourceListArea.Bottom}
		x += w
	}
	return out
}

func (a *App) sourceHeaderColumnAt(x, y int32) (col int, resize bool, ok bool) {
	if !pointIn(a.sourceListArea, x, y) {
		return -1, false, false
	}
	headerBottom := a.sourceListArea.Top + a.sp(sourceHeaderH)
	if y < a.sourceListArea.Top || y > headerBottom {
		return -1, false, false
	}
	rects := a.sourceColumnRects()
	const resizeSlop int32 = 6
	for pos := 0; pos < len(a.columnOrder); pos++ {
		idx := a.columnOrder[pos]
		r := rects[idx]
		if x >= r.Right-resizeSlop && x <= r.Right+resizeSlop {
			return idx, true, true
		}
		if pointIn(RECT{r.Left, r.Top, r.Right, headerBottom}, x, y) {
			return idx, false, true
		}
	}
	return -1, false, false
}

func (a *App) sourceHeaderPositionAtX(x int32) int {
	rects := a.sourceColumnRects()
	fixedRight := a.sourceListArea.Left + a.sp(sourceTableFixedW)
	if x <= fixedRight {
		return 0
	}
	for pos, idx := range a.columnOrder {
		r := rects[idx]
		mid := (r.Left + r.Right) / 2
		if x < mid {
			return pos
		}
	}
	return len(a.columnOrder) - 1
}

func (a *App) isSourceHeader(x, y int32) bool {
	if !pointIn(a.sourceListArea, x, y) {
		return false
	}
	headerBottom := a.sourceListArea.Top + a.sp(sourceHeaderH)
	if y < a.sourceListArea.Top || y > headerBottom {
		return false
	}
	return x >= a.sourceListArea.Left+a.sp(sourceTableFixedW) && x <= a.sourceListArea.Right
}

func (a *App) sourceHeaderResizeAt(x, y int32) bool {
	_, resize, ok := a.sourceHeaderColumnAt(x, y)
	return ok && resize
}

func (a *App) beginColumnDrag(x, y int32) bool {
	col, resize, ok := a.sourceHeaderColumnAt(x, y)
	if !ok {
		return false
	}
	a.columnDragCol = col
	a.columnDragStartX = x
	a.columnDragOrigWidth = a.columnWidths[col]
	a.columnDragFrom = a.columnPosition(col)
	a.columnDragTo = a.columnDragFrom
	if resize {
		a.columnDragMode = columnDragResize
	} else {
		a.columnDragMode = columnDragReorder
	}
	pSetCapture.Call(uintptr(a.hwnd))
	return true
}

func (a *App) columnPosition(col int) int {
	for i, v := range a.columnOrder {
		if v == col {
			return i
		}
	}
	return -1
}

func (a *App) updateColumnDrag(x int32) {
	switch a.columnDragMode {
	case columnDragResize:
		logicalDelta := int32(float64(x-a.columnDragStartX) * 96.0 / float64(a.dpi) / a.uiScale)
		width := a.columnDragOrigWidth + logicalDelta
		minW := minColumnWidth(a.columnDragCol)
		maxW := maxColumnWidth(a.columnDragCol)
		if width < minW {
			width = minW
		}
		if width > maxW {
			width = maxW
		}
		a.columnWidths[a.columnDragCol] = width
		a.clampSourceHScroll()
		pInvalidateRect.Call(uintptr(a.hwnd), 0, 0)
	case columnDragReorder:
		a.columnDragTo = a.sourceHeaderPositionAtX(x)
		if a.columnDragTo < 0 {
			a.columnDragTo = 0
		}
		if a.columnDragTo >= len(a.columnOrder) {
			a.columnDragTo = len(a.columnOrder) - 1
		}
		pInvalidateRect.Call(uintptr(a.hwnd), 0, 0)
	}
}

func (a *App) finishColumnDrag() {
	mode := a.columnDragMode
	if mode == columnDragReorder && a.columnDragFrom >= 0 && a.columnDragTo >= 0 && a.columnDragFrom != a.columnDragTo {
		moved := a.columnOrder[a.columnDragFrom]
		if a.columnDragTo > a.columnDragFrom {
			copy(a.columnOrder[a.columnDragFrom:a.columnDragTo], a.columnOrder[a.columnDragFrom+1:a.columnDragTo+1])
		} else {
			copy(a.columnOrder[a.columnDragTo+1:a.columnDragFrom+1], a.columnOrder[a.columnDragTo:a.columnDragFrom])
		}
		a.columnOrder[a.columnDragTo] = moved
	}
	if mode != columnDragNone {
		pReleaseCapture.Call()
		a.saveSettings()
	}
	if mode == columnDragResize {
		// If the new total width crossed the viewport threshold, layout the
		// list once so the native horizontal scrollbar appears/disappears.
		a.layoutSourceTable()
	}
	a.columnDragMode = columnDragNone
	a.columnDragCol = -1
	a.columnDragFrom = -1
	a.columnDragTo = -1
	pInvalidateRect.Call(uintptr(a.hwnd), 0, 1)
}

func (a *App) layoutSourceTable() {
	baseBottom := a.sourceRect.Bottom - a.sp(46)
	listTop := a.sourceRect.Top + a.sp(70)
	listRight := a.sourceRect.Right - a.sp(12)
	scrollbarW := a.sp(17)
	scrollbarLeft := listRight - scrollbarW
	if scrollbarLeft < a.sourceRect.Left+a.sp(120) {
		scrollbarLeft = a.sourceRect.Left + a.sp(120)
	}
	a.sourceListArea = RECT{a.sourceRect.Left + a.sp(12), listTop, scrollbarLeft - a.sp(6), baseBottom}
	if a.sourceListArea.Right < a.sourceListArea.Left+a.sp(80) {
		a.sourceListArea.Right = a.sourceListArea.Left + a.sp(80)
	}
	needH := a.sourceColumnsContentWidth() > a.sourceDataViewportWidth()
	hbarH := a.sp(sourceScrollbarH)
	if needH {
		gap := a.sp(sourceScrollbarGap)
		a.sourceListArea.Bottom -= hbarH + gap
		if a.sourceListArea.Bottom < a.sourceListArea.Top+a.sp(100) {
			a.sourceListArea.Bottom = a.sourceListArea.Top + a.sp(100)
		}
	}
	if a.sourceScrollbar != 0 {
		pMoveWindow.Call(uintptr(a.sourceScrollbar), uintptr(scrollbarLeft), uintptr(listTop), uintptr(scrollbarW), uintptr(a.sourceListArea.Bottom-listTop), 1)
	}
	if a.sourceHScrollbar != 0 {
		if needH {
			dataLeft := a.sourceListArea.Left + a.sp(sourceTableFixedW)
			width := a.sourceListArea.Right - dataLeft
			y := a.sourceListArea.Bottom + a.sp(sourceScrollbarGap)
			pMoveWindow.Call(uintptr(a.sourceHScrollbar), uintptr(dataLeft), uintptr(y), uintptr(width), uintptr(hbarH), 1)
		}
	}
	a.updateSourceScrollbars()
}

func (a *App) updateSourceScrollbars() {
	if a.sourceScrollbar != 0 {
		_, visible, _ := a.sourceMetrics()
		a.clampSourceScroll()
		if len(a.items) <= visible {
			pShowWindow.Call(uintptr(a.sourceScrollbar), SW_HIDE)
		} else {
			info := SCROLLINFO{CbSize: uint32(unsafe.Sizeof(SCROLLINFO{})), FMask: SIF_RANGE | SIF_PAGE | SIF_POS, NMin: 0, NMax: int32(len(a.items) - 1), NPage: uint32(visible), NPos: int32(a.scroll)}
			pSetScrollInfo.Call(uintptr(a.sourceScrollbar), SB_CTL, uintptr(unsafe.Pointer(&info)), 1)
			pShowWindow.Call(uintptr(a.sourceScrollbar), SW_SHOW)
		}
	}
	if a.sourceHScrollbar != 0 {
		a.clampSourceHScroll()
		max := a.sourceHorizontalMaxScroll()
		if max <= 0 {
			pShowWindow.Call(uintptr(a.sourceHScrollbar), SW_HIDE)
		} else {
			page := uint32(a.sourceDataViewportWidth())
			info := SCROLLINFO{CbSize: uint32(unsafe.Sizeof(SCROLLINFO{})), FMask: SIF_RANGE | SIF_PAGE | SIF_POS, NMin: 0, NMax: a.sourceColumnsContentWidth(), NPage: page, NPos: a.hScroll}
			pSetScrollInfo.Call(uintptr(a.sourceHScrollbar), SB_CTL, uintptr(unsafe.Pointer(&info)), 1)
			pShowWindow.Call(uintptr(a.sourceHScrollbar), SW_SHOW)
		}
	}
}

func (a *App) drawSourceTable(hdc uintptr, list RECT) {
	fillRect(hdc, list, cSurface)
	frameRect(hdc, list, cBorder2)
	header := RECT{list.Left, list.Top, list.Right, list.Top + a.sp(sourceHeaderH)}
	fillRect(hdc, header, cSurface2)

	fixedRight := list.Left + a.sp(sourceTableFixedW)
	textRect(hdc, RECT{list.Left, list.Top, fixedRight, header.Bottom}, "", cText2, a.smallFont, DT_CENTER|DT_SINGLELINE|DT_VCENTER|DT_NOPREFIX)

	rects := a.sourceColumnRects()
	for _, idx := range a.columnOrder {
		r := rects[idx]
		vr := r
		if vr.Left < fixedRight {
			vr.Left = fixedRight
		}
		if vr.Right > list.Right {
			vr.Right = list.Right
		}
		if vr.Right <= vr.Left {
			continue
		}
		title := a.tr(columnTitleKey(idx))
		textRect(hdc, RECT{vr.Left + a.sp(8), header.Top, vr.Right - a.sp(8), header.Bottom}, title, cText2, a.smallFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
		// Subtle header separators make the draggable column boundaries discoverable.
		dividerX := r.Right
		if dividerX > fixedRight && dividerX < list.Right {
			fillRect(hdc, RECT{dividerX, header.Top + a.sp(8), dividerX + a.sp(1), header.Bottom - a.sp(8)}, cBorder2)
		}
	}

	rowH, visible, _ := a.sourceMetrics()
	start := a.scroll
	if start > len(a.items)-visible {
		start = len(a.items) - visible
	}
	if start < 0 {
		start = 0
	}
	end := start + visible
	if end > len(a.items) {
		end = len(a.items)
	}
	for i := start; i < end; i++ {
		y := list.Top + a.sp(sourceHeaderH) + int32(i-start)*rowH
		row := RECT{list.Left + 1, y, list.Right - 1, y + rowH}
		if a.items[i].Selected {
			fillRect(hdc, row, cSelected)
		}
		a.drawCheckbox(hdc, RECT{row.Left + a.sp(14), row.Top + a.sp(15), row.Left + a.sp(38), row.Top + a.sp(39)}, a.items[i].Selected)
		drawItemIcon(hdc, RECT{row.Left + a.sp(50), row.Top + a.sp(11), row.Left + a.sp(80), row.Top + a.sp(43)}, a.items[i].Kind)
		for _, idx := range a.columnOrder {
			r := rects[idx]
			vr := RECT{r.Left + a.sp(8), row.Top, r.Right - a.sp(8), row.Bottom}
			if vr.Left < fixedRight {
				vr.Left = fixedRight + a.sp(8)
			}
			if vr.Right > list.Right-a.sp(2) {
				vr.Right = list.Right - a.sp(2)
			}
			if vr.Right <= vr.Left {
				continue
			}
			var value, color string
			switch idx {
			case columnName:
				value = displayName(a.items[i].Path)
				color = "text"
			case columnType:
				if a.items[i].Kind == "folder" {
					value = a.tr("source.kind.folder")
				} else {
					value = a.tr("source.kind.file")
				}
				color = "text2"
			case columnPath:
				value = a.items[i].Path
				color = "text2"
			}
			c := cText2
			if color == "text" {
				c = cText
			}
			textRect(hdc, vr, value, c, a.bodyFont, DT_LEFT|DT_SINGLELINE|DT_VCENTER|DT_END_ELLIPSIS|DT_NOPREFIX)
		}
		frameRect(hdc, RECT{row.Left, row.Bottom - 1, row.Right, row.Bottom}, cBorder)
	}
}

func columnTitleKey(idx int) string {
	switch idx {
	case columnName:
		return "source.col.name"
	case columnType:
		return "source.col.type"
	default:
		return "source.col.path"
	}
}
