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

// Screen represents the physical (or emulated) screen.
// This can be a terminal window or a physical console.  Platforms implement
// this differently.
type Screen interface {
	// Init initializes the screen for use.
	Init() error

	// Fini finalizes the screen also releasing resources.
	Fini()

	// Clear logically erases the screen.
	// This is effectively a short-cut for Fill(' ', StyleDefault).
	Clear()

	// Fill fills the screen with the given character and style.
	// The effect of filling the screen is not visible until Show
	// is called (or Sync).
	Fill(rune, Style)

	// GetContent returns the contents at the given location.  If the
	// coordinates are out of range, then the values will be 0, nil,
	// StyleDefault.  Note that the contents returned are logical contents
	// and may not actually be what is displayed, but rather are what will
	// be displayed if Show() or Sync() is called.  The width is the width
	// in screen cells; most often this will be 1, but some East Asian
	// characters and emoji require two cells.
	GetContent(x, y int) (primary rune, combining []rune, style Style, width uint8, dirty bool)

	// SetContent sets the contents of the given cell location.  If
	// the coordinates are out of range, then the operation is ignored.
	//
	// The first rune is the primary non-zero width rune.  The array
	// that follows is a possible list of combining characters to append,
	// and will usually be nil (no combining characters.)
	//
	// The results are not displayed until Show() or Sync() is called.
	//
	// Note that wide (East Asian full width and emoji) runes occupy two cells,
	// and attempts to place character at next cell to the right will have
	// undefined effects.  Wide runes that are printed in the
	// last column will be replaced with a single width space on output.
	SetContent(x int, y int, primary rune, combining []rune, width uint8, style Style)

	// UnionStyle computes the set union between a and b,
	// that is overrides a set of attributes that contain
	// all the bit flags set in a, b or both, and uses the color
	// defined in b or if not set, uses the color in a.
	UnionStyle(x int, y int, style Style)

	// SetStyle sets the default style to use when clearing the screen
	// or when StyleDefault is specified.  If it is also StyleDefault,
	// then whatever system/terminal default is relevant will be used.
	SetStyle(style Style)

	// ShowCursor is used to display the cursor at a given location.
	// If the coordinates -1, -1 are given or are otherwise outside the
	// dimensions of the screen, the cursor will be hidden.
	ShowCursor(x int, y int)

	// HideCursor is used to hide the cursor.  It's an alias for
	// ShowCursor(-1, -1).sim
	HideCursor()

	// SetCursorStyle is used to set the cursor style.  If the style
	// is not supported (or cursor styles are not supported at all),
	// then this will have no effect.
	SetCursorStyle(CursorStyle)

	// Size returns the screen size as width, height.  This changes in
	// response to a call to Clear or Flush.
	Size() (width, height int)

	// Poll returns the underlying event channel.
	Poll() <-chan Event

	// PostEvent tries to post an event into the event stream.  This
	// can fail if the event queue is full.  In that case, the event
	// is dropped, and ErrEventQFull is returned.
	PostEvent(ev Event) error

	// EnableMouse enables the mouse.  (If your terminal supports it.)
	// If no flags are specified, then all events are reported, if the
	// terminal supports them.
	EnableMouse(...MouseFlags)

	// DisableMouse disables the mouse.
	DisableMouse()

	// EnablePaste enables bracketed paste mode, if supported.
	EnablePaste()

	// DisablePaste disables bracketed paste mode.
	DisablePaste()

	// EnableFocus enables reporting of focus events, if your terminal supports it.
	EnableFocus()

	// DisableFocus disables reporting of focus events.
	DisableFocus()

	// Bell makes an audible noise.
	Bell()

	// HasMouse returns true if the terminal (apparently) supports a
	// mouse.  Note that the return value of true doesn't guarantee that
	// a mouse/pointing device is present; a false return definitely
	// indicates no mouse support is available.
	HasMouse() bool

	// Colors returns the number of colors.  All colors are assumed to
	// use the ANSI color map.  If a terminal is monochrome, it will
	// return 0.
	Colors() int

	// Show makes all the content changes made using SetContent() visible
	// on the display.
	//
	// It does so in the most efficient and least visually disruptive
	// manner possible.
	Show()

	// HasKey returns true if the keyboard is believed to have the
	// key.  In some cases a keyboard may have keys with this name
	// but no support for them, while in others a key may be reported
	// as supported but not actually be usable (such as some emulators
	// that hijack certain keys).  Its best not to depend to strictly
	// on this function, but it can be used for hinting when building
	// menus, displayed hot-keys, etc.  Note that KeyRune (literal
	// runes) is always true.
	HasKey(Key) bool

	// Beep attempts to sound an OS-dependent audible alert and returns an error
	// when unsuccessful.
	Beep() error

	// Tty returns the underlying Tty. If the screen is not a terminal, the
	// returned bool will be false
	Tty() (Tty, bool)
}

// NewScreen returns a default Screen suitable for the user's terminal
// environment.
func NewScreen() (Screen, error) {
	// Windows is happier if we try for a console screen first.
	if s, _ := NewConsoleScreen(); s != nil {
		return s, nil
	} else if s, e := NewTerminfoScreen(); s != nil {
		return s, nil
	} else {
		return nil, e
	}
}

// MouseFlags are options to modify the handling of mouse events.
// Actual events can be ORed together.
type MouseFlags int

const (
	MouseButtonEvents = MouseFlags(1) // Click events only
	MouseDragEvents   = MouseFlags(2) // Click-drag events (includes button events)
	MouseMotionEvents = MouseFlags(4) // All mouse events (includes click and drag events)
)

// CursorStyle represents a given cursor style, which can include the shape and
// whether the cursor blinks or is solid.  Support for changing this is not universal.
type CursorStyle int

const (
	CursorStyleDefault = CursorStyle(iota) // The default
	CursorStyleBlinkingBlock
	CursorStyleSteadyBlock
	CursorStyleBlinkingUnderline
	CursorStyleSteadyUnderline
	CursorStyleBlinkingBar
	CursorStyleSteadyBar
)

// screenImpl is a subset of Screen that can be used with baseScreen to formulate
// a complete implementation of Screen.  See Screen for doc comments about methods.
type screenImpl interface {
	Init() error
	Fini()
	SetStyle(style Style)
	ShowCursor(x int, y int)
	HideCursor()
	SetCursorStyle(CursorStyle)
	Size() (width, height int)
	EnableMouse(...MouseFlags)
	DisableMouse()
	EnablePaste()
	DisablePaste()
	EnableFocus()
	DisableFocus()
	Bell()
	HasMouse() bool
	Colors() int
	Show()
	HasKey(Key) bool
	Beep() error
	Tty() (Tty, bool)
	Poll() <-chan Event
	PostEvent(ev Event) error
}

type baseScreen struct {
	screenImpl
	cb *CellBuffer
}

func (b *baseScreen) Clear() {
	b.cb.Fill(' ', StyleDefault)
}

func (b *baseScreen) Fill(r rune, style Style) {
	b.cb.Fill(r, style)
}

func (b *baseScreen) SetContent(x, y int, mainc rune, combc []rune, width uint8, st Style) {
	b.cb.SetContentWidth(x, y, mainc, combc, width, st)
}

func (b *baseScreen) UnionStyle(x int, y int, style Style) {
	b.cb.UnionStyle(x, y, style)
}

func (b *baseScreen) GetContent(x, y int) (rune, []rune, Style, uint8, bool) {
	return b.cb.GetContent(x, y)
}
