// Copyright 2023 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

// Cell represents a cell in a two-dimensional array of character cells.
type Cell struct {
	currMain  rune
	currComb  []rune
	currStyle Style
	lastMain  rune
	lastStyle Style
	width     uint8
}

// CellBuffer represents a two-dimensional array of character cells.
// This is primarily intended for use by Screen implementors; it
// contains much of the common code they need.  To create one, just
// declare a variable of its type; no explicit initialization is necessary.
//
// CellBuffer is not thread safe.
type CellBuffer struct {
	w     int
	h     int
	cells []Cell
}

// SetContentWidth behaves like SetContent, but allows clients
// to pass the grapheme width if known. The behaviour is undefined
// if the passed width is ever 0.
func (cb *CellBuffer) SetContentWidth(x int, y int,
	mainc rune, combc []rune, width uint8, style Style,
) {
	if x >= cb.w || y >= cb.h {
		return
	}

	c := &cb.cells[(y*cb.w)+x]
	c.currComb = combc
	c.width = width
	c.currMain = mainc
	c.currStyle = style

	iwidth := int(width)
	for i := 1; i < iwidth; i++ {
		cb.SetDirty(x+i, y)
	}
}

// UnionStyle computes the set union between a and b,
// that is it overrides the set of attributes at x, y that contain
// all the bit flags set in a, b or both, and uses the color
// defined in b or if not set, uses the color in a.
func (cb *CellBuffer) UnionStyle(x int, y int, style Style) {
	if x >= cb.w || y >= cb.h {
		return
	}

	c := &cb.cells[(y*cb.w)+x]
	c.currStyle = unionStyle(c.currStyle, style)
}

// GetContent returns the contents of a character cell, including the
// primary rune, any combining character runes (which will usually be
// nil), the style, and the display width in cells.
func (cb *CellBuffer) GetContent(x, y int) (
	mainc rune, combc []rune, style Style, width uint8, dirty bool,
) {
	c := cb.GetCell(x, y)
	mainc, combc, style, width = cb.ProcessCell(c)
	dirty = cb.Dirty(c)
	return
}

// GetCell gets the cell at x and y coordinates.
func (cb *CellBuffer) GetCell(x, y int) *Cell {
	return &cb.cells[(y*cb.w)+x]
}

// ProcessCell unwraps the given cell
func (cb *CellBuffer) ProcessCell(c *Cell) (
	mainc rune, combc []rune, style Style, width uint8,
) {
	mainc = c.currMain
	combc = c.currComb
	style = c.currStyle
	width = c.width
	if width == 0 || mainc < ' ' {
		width = 1
		mainc = ' '
		// no need to reset combc on mainc < ' '
		// as it is only useful for cells that were not set
		// and those will not have a combc set.
	}
	return
}

func (cb *CellBuffer) getCells() []Cell {
	return cb.cells
}

// Size returns the (width, height) in cells of the buffer.
func (cb *CellBuffer) Size() (int, int) {
	return cb.w, cb.h
}

// Dirty checks if a character at the given location needs to be
// refreshed on the physical display.  This returns true if the cell
// content is different since the last time it was marked clean.
func (cb *CellBuffer) DirtyAt(x, y int) bool {
	return cb.Dirty(&cb.cells[(y*cb.w)+x])
}

// Dirty checks if the given cell needs to be
// refreshed on the physical display.  This returns true if the cell
// content is different since the last time it was marked clean.
func (cb *CellBuffer) Dirty(c *Cell) bool {
	return c.lastMain == rune(0) ||
		c.lastMain != c.currMain ||
		c.lastStyle != c.currStyle
	// we don't clear dirty when width > 1
	// so there's no need to check for this
	// len(c.lastComb) != len(c.currComb)
}

// SetDirty is normally used to manually
// force a cell to be marked dirty.
func (cb *CellBuffer) SetDirty(x, y int) {
	cb.cells[(y*cb.w)+x].lastMain = 0
}

// SetDirty marks the given Cell as dirty.
func SetDirty(c *Cell) {
	c.lastMain = 0
}

// ClearDirty is normally used to indicate that a cell has
// been displayed.
func (cb *CellBuffer) ClearDirty(x, y int) {
	ClearDirty(&cb.cells[(y*cb.w)+x])
}

// ClearDirty marks the given Cell as not dirty.
func ClearDirty(c *Cell) {
	c.lastMain = c.currMain
	c.lastStyle = c.currStyle
	// no need to copy currComb as this
	// should never be called for multi-width cells.
}

// Resize is used to resize the cells array, with different dimensions,
// while preserving the original contents within the overlapping
// rectangle. Cells within that overlap are invalidated so that they
// are redrawn on the next Show/Sync. Cells outside the overlap
// (newly exposed area, when enlarging) are zero-valued.
func (cb *CellBuffer) Resize(w, h int) {
	if cb.h == h && cb.w == w {
		return
	}

	newc := make([]Cell, w*h)
	minW, minH := w, h
	if cb.w < minW {
		minW = cb.w
	}
	if cb.h < minH {
		minH = cb.h
	}
	for y := 0; y < minH; y++ {
		for x := 0; x < minW; x++ {
			c := cb.cells[(y*cb.w)+x]
			// Invalidate so the next draw repaints the cell at its
			// new position in the buffer.
			c.lastMain = 0
			newc[(y*w)+x] = c
		}
	}
	cb.cells = newc
	cb.h = h
	cb.w = w
}

// Fill fills the entire cell buffer array with the specified character
// and style.  Normally choose ' ' to clear the screen.  This API doesn't
// support combining characters, or characters with a width larger than one.
// If either the foreground or background are ColorNone, then the respective
// color is unchanged.
func (cb *CellBuffer) Fill(r rune, style Style) {
	if style.Fg == ColorNone && style.Bg == ColorNone {
		for i := range cb.cells {
			c := &cb.cells[i]
			c.currMain = r
			c.currComb = nil
			c.width = 1
		}
		return
	}
	for i := range cb.cells {
		c := &cb.cells[i]
		c.currMain = r
		c.currComb = nil
		c.currStyle = style
		c.width = 1
	}
}

func unionStyle(a, b Style) Style {
	retFgColor := b.Fg
	retBgColor := b.Bg
	if b.Fg == ColorDefault {
		retFgColor = a.Fg
	}
	if b.Bg == ColorDefault {
		retBgColor = a.Bg
	}
	return Style{Fg: retFgColor, Bg: retBgColor, Attrs: (a.Attrs | b.Attrs)}
}
