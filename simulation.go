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

import (
	"unicode/utf8"
)

// NewSimulationScreen returns a SimulationScreen.  Note that
// SimulationScreen is also a Screen.
func NewSimulationScreen() SimulationScreen {
	ss := &simscreen{}
	ss.Screen = &baseScreen{screenImpl: ss, cb: ss.GetCells()}
	return ss
}

// SimulationScreen represents a screen simulation.  This is intended to
// be a superset of normal Screens, but also adds some important interfaces
// for testing.
type SimulationScreen interface {
	Screen

	// InjectKeyBytes injects a stream of bytes corresponding to
	// the native encoding (see charset).  It turns true if the entire
	// set of bytes were processed and delivered as KeyEvents, false
	// if any bytes were not fully understood.  Any bytes that are not
	// fully converted are discarded.
	InjectKeyBytes(buf []byte) bool

	// InjectKey injects a key event.  The rune is a UTF-8 rune, post
	// any translation.
	InjectKey(key Key, r rune, mod ModMask)

	// InjectMouse injects a mouse event.
	InjectMouse(x, y int, buttons ButtonMask, mod ModMask)

	// GetContents returns screen contents as an array of
	// cells, along with the physical width & height.   Note that the
	// physical contents will be used until the next time SetSize()
	// is called.
	GetContents() (cells []SimCell, width int, height int)

	// GetCursor returns the cursor details.
	GetCursor() (x int, y int, visible bool)

	SetSize(int, int)
}

// SimCell represents a simulated screen cell.  The purpose of this
// is to track on screen content.
type SimCell struct {
	// Bytes is the actual character bytes.  Normally this is
	// rune data, but it could be be data in another encoding system.
	Bytes []byte

	// Style is the style used to display the data.
	Style Style

	// Runes is the list of runes, unadulterated, in UTF-8.
	Runes []rune
}

type simscreen struct {
	dirty bool
	physw int
	physh int
	fini  bool
	style Style
	evch  chan Event
	quit  chan struct{}

	front     []SimCell
	back      CellBuffer
	cursorx   int
	cursory   int
	cursorvis bool
	mouse     bool
	paste     bool
	fillchar  rune
	fillstyle Style

	Screen
}

func (s *simscreen) Init() error {
	s.evch = make(chan Event, 10)
	s.quit = make(chan struct{})
	s.fillchar = 'X'
	s.fillstyle = StyleDefault
	s.mouse = false
	s.physw = 80
	s.physh = 25
	s.cursorx = -1
	s.cursory = -1
	s.style = StyleDefault

	s.front = make([]SimCell, s.physw*s.physh)
	s.back.Resize(80, 25)
	return nil
}

func (s *simscreen) Fini() {
	s.fini = true
	s.back.Resize(0, 0)
	if s.quit != nil {
		close(s.quit)
	}
	s.physw = 0
	s.physh = 0
	s.front = nil
}

func (s *simscreen) SetStyle(style Style) {
	s.style = style
}

func (s *simscreen) drawCell(x, y int) uint8 {

	mainc, combc, style, width, dirty := s.back.GetContent(x, y)
	if !dirty || x >= s.physw || y >= s.physh {
		return width
	}

	simc := &s.front[(y*s.physw)+x]
	if style == StyleDefault {
		style = s.style
	}
	simc.Style = style

	if x > s.physw-int(width) {
		simc.Runes = []rune{' '}
		simc.Bytes = []byte{' '}
		s.back.ClearDirty(x, y)
		return width
	}

	simc.Runes = append([]rune{mainc}, combc...)
	simc.Bytes = make([]byte, 0, len(simc.Runes))
	simc.Bytes = append(simc.Bytes, []byte(string(simc.Runes))...)
	s.back.ClearDirty(x, y)
	return width
}

func (s *simscreen) ShowCursor(x, y int) {
	s.cursorx, s.cursory = x, y
	s.showCursor()
}

func (s *simscreen) HideCursor() {
	s.ShowCursor(-1, -1)
}

func (s *simscreen) showCursor() {

	x, y := s.cursorx, s.cursory
	if x < 0 || y < 0 || x >= s.physw || y >= s.physh {
		s.cursorvis = false
	} else {
		s.cursorvis = true
	}
}

func (s *simscreen) hideCursor() {
	// does not update cursor position
	s.cursorvis = false
}

func (s *simscreen) SetCursorStyle(CursorStyle) {}

func (s *simscreen) Show() {
	if s.dirty {
		s.resize()
	}
	s.draw()
}

func (s *simscreen) clearScreen() {
	// We emulate a hardware clear by filling with a specific pattern
	for i := range s.front {
		s.front[i].Style = s.fillstyle
		s.front[i].Runes = []rune{s.fillchar}
		s.front[i].Bytes = []byte{byte(s.fillchar)}
	}
}

func (s *simscreen) draw() {
	s.hideCursor()

	w, h := s.back.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			width := s.drawCell(x, y)
			if width > 1 {
				x += int(width) - 1
			}
		}
	}
	s.showCursor()
}

func (s *simscreen) EnableMouse(...MouseFlags) {
	s.mouse = true
}

func (s *simscreen) DisableMouse() {
	s.mouse = false
}

func (s *simscreen) EnablePaste() {
	s.paste = true
}

func (s *simscreen) DisablePaste() {
	s.paste = false
}

func (s *simscreen) EnableFocus() {
}

func (s *simscreen) DisableFocus() {
}

func (s *simscreen) Size() (int, int) {
	w, h := s.back.Size()
	return w, h
}

func (s *simscreen) resize() {
	w, h := s.physw, s.physh
	ow, oh := s.back.Size()
	if w != ow || h != oh {
		s.back.Resize(w, h)
	}
	s.dirty = false
}

func (s *simscreen) Colors() int {
	return 256
}

func (s *simscreen) PollEvent() Event {
	select {
	case <-s.quit:
		return nil
	case ev := <-s.evch:
		return ev
	}
}

func (s *simscreen) Poll() <-chan Event {
	return s.evch
}

func (s *simscreen) PostEventWait(ev Event) {
	s.evch <- ev
}

func (s *simscreen) postEvent(ev Event) {
	select {
	case s.evch <- ev:
	case <-s.quit:
	}
}

func (s *simscreen) InjectMouse(x, y int, buttons ButtonMask, mod ModMask) {
	ev := NewEventMouse(x, y, buttons, mod, nil)
	s.postEvent(ev)
}

func (s *simscreen) InjectKey(key Key, r rune, mod ModMask) {
	ev := NewEventKey(key, r, mod, nil)
	s.postEvent(ev)
}

func (s *simscreen) InjectKeyBytes(b []byte) bool {
	failed := false

outer:
	for len(b) > 0 {
		if b[0] >= ' ' && b[0] <= 0x7F {
			// printable ASCII easy to deal with -- no encodings
			ev := NewEventKey(KeyRune, rune(b[0]), ModNone, nil)
			s.postEvent(ev)
			b = b[1:]
			continue
		}

		if b[0] < 0x80 {
			mod := ModNone
			// No encodings start with low numbered values
			if Key(b[0]) >= KeyCtrlA && Key(b[0]) <= KeyCtrlZ {
				mod = ModCtrl
			}
			ev := NewEventKey(Key(b[0]), 0, mod, nil)
			s.postEvent(ev)
			b = b[1:]
			continue
		}

		for l := 1; l < len(b); l++ {
			r, nin := utf8.DecodeRune(b[:l])
			if r != utf8.RuneError {
				ev := NewEventKey(KeyRune, r, ModNone, nil)
				s.postEvent(ev)
			}
			b = b[nin:]
			continue outer
		}
		failed = true
		b = b[1:]
		continue
	}

	return !failed
}

func (s *simscreen) SetSize(w, h int) {
	newc := make([]SimCell, w*h)
	for row := 0; row < h && row < s.physh; row++ {
		for col := 0; col < w && col < s.physw; col++ {
			newc[(row*w)+col] = s.front[(row*s.physw)+col]
		}
	}
	s.cursorx, s.cursory = -1, -1
	s.physw, s.physh = w, h
	s.front = newc
	s.back.Resize(w, h)
	s.dirty = true
}

func (s *simscreen) GetContents() ([]SimCell, int, int) {
	cells, w, h := s.front, s.physw, s.physh
	return cells, w, h
}

func (s *simscreen) GetCursor() (int, int, bool) {
	x, y, vis := s.cursorx, s.cursory, s.cursorvis
	return x, y, vis
}

func (s *simscreen) HasMouse() bool {
	return false
}

func (s *simscreen) Resize(int, int, int, int) {}

func (s *simscreen) HasKey(Key) bool {
	return true
}

func (s *simscreen) Beep() error {
	return nil
}

func (s *simscreen) Suspend() error {
	return nil
}

func (s *simscreen) Resume() error {
	return nil
}

func (s *simscreen) Tty() (Tty, bool) {
	return nil, false
}

func (s *simscreen) GetCells() *CellBuffer {
	return &s.back
}

func (s *simscreen) EventQ() chan Event {
	return s.evch
}

func (s *simscreen) StopQ() <-chan struct{} {
	return s.quit
}
