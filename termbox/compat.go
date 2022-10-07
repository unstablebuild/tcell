// Copyright 2020 The TCell Authors
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

// Package termbox is a compatibility layer to allow tcell to emulate
// the github.com/nsf/termbox package.
package termbox

import (
	"errors"

	"github.com/ernestrc/tcell/v2"
)

var screen tcell.Screen
var outMode OutputMode

// Init initializes the screen for use.
func Init() error {
	outMode = OutputNormal
	if s, e := tcell.NewScreen(); e != nil {
		return e
	} else if e = s.Init(); e != nil {
		return e
	} else {
		screen = s
		return nil
	}
}

// Close cleans up the terminal, restoring terminal modes, etc.
func Close() {
	screen.Fini()
}

// Flush updates the screen.
func Flush() error {
	screen.Show()
	return nil
}

// SetCursor displays the terminal cursor at the given location.
func SetCursor(x, y int) {
	screen.ShowCursor(x, y)
}

// HideCursor hides the terminal cursor.
func HideCursor() {
	SetCursor(-1, -1)
}

// Size returns the screen size as width, height in character cells.
func Size() (int, int) {
	return screen.Size()
}

// Attribute affects the presentation of characters, such as color, boldness,
// and so forth.
type Attribute uint16

// Colors first.  The order here is significant.
const (
	ColorDefault Attribute = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// Other attributes.
const (
	AttrBold Attribute = 1 << (9 + iota)
	AttrUnderline
	AttrReverse
)

func fixColor(c tcell.Color) tcell.Color {
	if c == tcell.ColorDefault {
		return c
	}
	switch outMode {
	case OutputNormal:
		c = tcell.PaletteColor(int(c) & 0xf)
	case Output256:
		c = tcell.PaletteColor(int(c) & 0xff)
	case Output216:
		c = tcell.PaletteColor(int(c)%216 + 16)
	case OutputGrayscale:
		c = tcell.PaletteColor(int(c)%24 + 232)
	default:
		c = tcell.ColorDefault
	}
	return c
}

func mkStyle(fg, bg Attribute) tcell.Style {
	st := tcell.StyleDefault

	var f, b tcell.Color

	fColor := int(fg) & 0x1ff
	bColor := int(bg) & 0x1ff

	if Attribute(fColor) != ColorDefault {
		f = tcell.PaletteColor(fColor - 1)
	}
	if Attribute(bColor) != ColorDefault {
		b = tcell.PaletteColor(bColor - 1)
	}

	f = fixColor(f)
	b = fixColor(b)
	st = st.Foreground(f).Background(b)
	if (fg|bg)&AttrBold != 0 {
		st = st.Bold(true)
	}
	if (fg|bg)&AttrUnderline != 0 {
		st = st.Underline(true)
	}
	if (fg|bg)&AttrReverse != 0 {
		st = st.Reverse(true)
	}
	return st
}

// Clear clears the screen with the given attributes.
func Clear(fg, bg Attribute) {
	st := mkStyle(fg, bg)
	w, h := screen.Size()
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			screen.SetContent(col, row, ' ', nil, st)
		}
	}
}

// InputMode is not used.
type InputMode int

// Input mode. See SetInputMode function.
const (
	InputEsc InputMode = 1 << iota
	InputAlt
	InputMouse
	InputCurrent InputMode = 0
)

// SetInputMode does not do anything in this version.
func SetInputMode(mode InputMode) InputMode {
	if mode == InputCurrent {
		return InputEsc
	}
	if mode&InputMouse != 0 {
		screen.EnableMouse()
	} else {
		screen.DisableMouse()
	}
	return InputEsc
}

// OutputMode represents an output mode, which determines how colors
// are used.  See the termbox documentation for an explanation.
type OutputMode int

// OutputMode values.
const (
	OutputCurrent OutputMode = iota
	OutputNormal
	Output256
	Output216
	OutputGrayscale
)

// SetOutputMode is used to set the color palette used.
func SetOutputMode(mode OutputMode) OutputMode {
	if screen == nil || screen.Colors() < 256 {
		mode = OutputNormal
	}
	switch mode {
	case OutputCurrent:
		return outMode
	case OutputNormal, Output256, Output216, OutputGrayscale:
		outMode = mode
		return mode
	default:
		return outMode
	}
}

// Sync forces a resync of the screen.
func Sync() error {
	screen.Sync()
	return nil
}

// SetCell sets the character cell at a given location to the given
// content (rune) and attributes.
func SetCell(x, y int, ch rune, fg, bg Attribute) {
	st := mkStyle(fg, bg)
	screen.SetContent(x, y, ch, nil, st)
}

// EventType represents the type of event.
type EventType uint8

// Modifier represents the possible modifier keys.
type Modifier tcell.ModMask

// Key is a key press.
type Key tcell.Key

// Event represents an event like a key press, mouse action, or window resize.
type Event struct {
	Type   EventType
	Mod    Modifier
	Key    Key
	Ch     rune
	Width  int
	Height int
	Err    error
	MouseX int
	MouseY int
	N      int
	Raw    []byte
}

// Event types.
const (
	EventNone EventType = iota
	EventKey
	EventResize
	EventMouse
	EventInterrupt
	EventError
	EventRaw
)

// Keys codes.
const (
	KeyF1             = Key(tcell.KeyF1)
	KeyF2             = Key(tcell.KeyF2)
	KeyF3             = Key(tcell.KeyF3)
	KeyF4             = Key(tcell.KeyF4)
	KeyF5             = Key(tcell.KeyF5)
	KeyF6             = Key(tcell.KeyF6)
	KeyF7             = Key(tcell.KeyF7)
	KeyF8             = Key(tcell.KeyF8)
	KeyF9             = Key(tcell.KeyF9)
	KeyF10            = Key(tcell.KeyF10)
	KeyF11            = Key(tcell.KeyF11)
	KeyF12            = Key(tcell.KeyF12)
	KeyInsert         = Key(tcell.KeyInsert)
	KeyDelete         = Key(tcell.KeyDelete)
	KeyHome           = Key(tcell.KeyHome)
	KeyEnd            = Key(tcell.KeyEnd)
	KeyArrowUp        = Key(tcell.KeyUp)
	KeyArrowDown      = Key(tcell.KeyDown)
	KeyArrowRight     = Key(tcell.KeyRight)
	KeyArrowLeft      = Key(tcell.KeyLeft)
	KeyCtrlA          = Key(tcell.KeyCtrlA)
	KeyCtrlB          = Key(tcell.KeyCtrlB)
	KeyCtrlC          = Key(tcell.KeyCtrlC)
	KeyCtrlD          = Key(tcell.KeyCtrlD)
	KeyCtrlE          = Key(tcell.KeyCtrlE)
	KeyCtrlF          = Key(tcell.KeyCtrlF)
	KeyCtrlG          = Key(tcell.KeyCtrlG)
	KeyCtrlH          = Key(tcell.KeyCtrlH)
	KeyCtrlI          = Key(tcell.KeyCtrlI)
	KeyCtrlJ          = Key(tcell.KeyCtrlJ)
	KeyCtrlK          = Key(tcell.KeyCtrlK)
	KeyCtrlL          = Key(tcell.KeyCtrlL)
	KeyCtrlM          = Key(tcell.KeyCtrlM)
	KeyCtrlN          = Key(tcell.KeyCtrlN)
	KeyCtrlO          = Key(tcell.KeyCtrlO)
	KeyCtrlP          = Key(tcell.KeyCtrlP)
	KeyCtrlQ          = Key(tcell.KeyCtrlQ)
	KeyCtrlR          = Key(tcell.KeyCtrlR)
	KeyCtrlS          = Key(tcell.KeyCtrlS)
	KeyCtrlT          = Key(tcell.KeyCtrlT)
	KeyCtrlU          = Key(tcell.KeyCtrlU)
	KeyCtrlV          = Key(tcell.KeyCtrlV)
	KeyCtrlW          = Key(tcell.KeyCtrlW)
	KeyCtrlX          = Key(tcell.KeyCtrlX)
	KeyCtrlY          = Key(tcell.KeyCtrlY)
	KeyCtrlZ          = Key(tcell.KeyCtrlZ)
	KeyCtrlUnderscore = Key(tcell.KeyCtrlUnderscore)
	KeyBackspace      = Key(tcell.KeyBackspace)
	KeyBackspace2     = Key(tcell.KeyBackspace2)
	KeyTab            = Key(tcell.KeyTab)
	KeyEnter          = Key(tcell.KeyEnter)
	KeyEsc            = Key(tcell.KeyEscape)
	KeyPgdn           = Key(tcell.KeyPgDn)
	KeyPgup           = Key(tcell.KeyPgUp)
	KeySpace          = Key(tcell.Key(' '))
	KeyTilde          = Key(tcell.Key('~'))
	KeyCtrlSpace      = Key(tcell.KeyCtrlSpace)

	// The following assignments are provided for termbox
	// compatibility.  Their use in applications is discouraged.
	// The mouse keys are completely not supported as tcell uses
	// a separate mouse event instead of key strokes.
	MouseLeft         = Key(tcell.KeyF63) // arbitrary assignments
	MouseRight        = Key(tcell.KeyF62)
	MouseMiddle       = Key(tcell.KeyF61)
	MouseRelease      = Key(tcell.KeyF60)
	MouseWheelUp      = Key(tcell.KeyF59)
	MouseWheelDown    = Key(tcell.KeyF58)
	KeyCtrlTilde      = Key(tcell.KeyCtrlSpace) // termbox defines a bunch of weird ones, don't use them
	KeyCtrl2          = Key(tcell.KeyNUL)
	KeyCtrl3          = Key(tcell.KeyEscape)
	KeyCtrl4          = Key(tcell.KeyCtrlBackslash)
	KeyCtrl5          = Key(tcell.KeyCtrlRightSq)
	KeyCtrl6          = Key(tcell.KeyCtrlCarat)
	KeyCtrl7          = Key(tcell.KeyCtrlUnderscore)
	KeyCtrl8          = Key(tcell.KeyDEL)
	KeyCtrlSlash      = Key(tcell.KeyCtrlUnderscore)
	KeyCtrlRsqBracket = Key(tcell.KeyCtrlRightSq)
	KeyCtrlBackslash  = Key(tcell.KeyCtrlBackslash)
	KeyCtrlLsqBracket = Key(tcell.KeyCtrlLeftSq)
)

// Modifiers.
const (
	ModAlt = Modifier(tcell.ModAlt)
)

var (
	// provide MouseRelease compat
	prevKey Key
)

// NewEvent makes a new event from a tcell.Event.
func NewEvent(tev tcell.Event) Event {
	switch tev := tev.(type) {
	case *tcell.EventInterrupt:
		if tev.Data() == EventNone {
			return Event{Type: EventNone}
		}
		return Event{Type: EventInterrupt}
	case *tcell.EventResize:
		w, h := tev.Size()
		return Event{Type: EventResize, Width: w, Height: h}
	case *tcell.EventMouse:
		x, y := tev.Position()
		var mod Modifier
		var key Key
		if tev.Modifiers()&tcell.ModAlt != 0 {
			mod = ModAlt
		}
		buttons := tev.Buttons()
		if buttons&tcell.ButtonPrimary != 0 {
			key = MouseLeft
		} else if buttons&tcell.ButtonSecondary != 0 {
			key = MouseRight
		} else if buttons&tcell.ButtonMiddle != 0 {
			key = MouseMiddle
		} else if buttons&tcell.WheelUp != 0 {
			key = MouseWheelUp
		} else if buttons&tcell.WheelDown != 0 {
			key = MouseWheelDown
		} else if prevKey == MouseLeft ||
			prevKey == MouseRight ||
			prevKey == MouseMiddle {
			key = MouseRelease
		}
		prevKey = key
		return Event{
			Type:   EventMouse,
			MouseX: x,
			MouseY: y,
			Mod:    mod,
			Key:    key,
			Raw:    tev.Raw(),
			N:      len(tev.Raw()),
		}
	// TODO
	// case *tcell.EventPaste:
	case *tcell.EventKey:
		k := tev.Key()
		ch := rune(0)
		if k == tcell.KeyRune {
			ch = tev.Rune()
			if ch == ' ' {
				k = tcell.Key(' ')
			} else {
				k = tcell.Key(0)
			}
		}
		var mod Modifier
		if tev.Modifiers()&tcell.ModAlt != 0 {
			mod = ModAlt
		}

		return Event{
			Type: EventKey,
			Key:  Key(k),
			Ch:   ch,
			Mod:  mod,
			N:    len(tev.Raw()),
			Raw:  tev.Raw(),
		}
	default:
		return Event{Type: EventNone}
	}
}

// ParseEvent is not supported.
func ParseEvent(data []byte) Event {
	// Not supported
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

// PollRawEvent is not supported.
func PollRawEvent(data []byte) Event {
	// Not supported
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

// PollEvent blocks until an event is ready, and then returns it.
func PollEvent() Event {
	ev := screen.PollEvent()
	return NewEvent(ev)
}

// Interrupt posts an interrupt event.
func Interrupt() {
	screen.PostEvent(tcell.NewEventInterrupt(nil))
}

// PublishEvent tries to publish an event into the event stream.
// This can fail if the event queue is full in which case false
// is returned. It doesn't support EventMouse and EventRaw events.
func PublishEvent(ev Event) bool {
	var tev tcell.Event
	switch ev.Type {
	case EventNone:
		tev = tcell.NewEventInterrupt(EventNone)
	case EventKey:
		var mod tcell.ModMask
		if ev.Mod&ModAlt != 0 {
			mod = tcell.ModAlt
		}
		tev = tcell.NewEventKey(tcell.Key(ev.Key), ev.Ch, mod, ev.Raw)
	case EventResize:
		tev = tcell.NewEventResize(ev.Width, ev.Height)
	case EventInterrupt:
		tev = tcell.NewEventInterrupt(nil)
	case EventError:
		tev = tcell.NewEventError(ev.Err)
	default /* + EventRaw + EventMouse */ :
		return false
	}
	err := screen.PostEvent(tev)
	return err != tcell.ErrEventQFull
}

// HasPendingEvent returns true if PollEvent would return an event
// without blocking.  If the screen is stopped and PollEvent would
// return nil, then the return value from this function is unspecified.
// The purpose of this function is to allow multiple events to be collected
// at once, to minimize screen redraws.
func HasPendingEvent() bool {
	return screen.HasPendingEvent()
}

// Screen exposes the underlying tcell.Screen.
func Screen() tcell.Screen {
	return screen
}

// Cell represents a single character cell on screen.
type Cell struct {
	Ch rune
	Fg Attribute
	Bg Attribute
}
