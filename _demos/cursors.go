//go:build ignore
// +build ignore

// Copyright 2022 The TCell Authors
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

// beep makes a beep every second until you press ESC
package main

import (
	"fmt"
	"github.com/ernestrc/tcell/v3"
	"os"
)

func main() {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault)
	s.Clear()

	text := "This demonstrates cursor styles.  Press 0 through 6 to change the style."
	x := 1
	for _, r := range text {
		s.SetContent(x, 1, r, nil, 1, tcell.StyleDefault)
		x++
	}
	s.SetContent(2, 2, '0', nil, 1, tcell.StyleDefault)
	s.SetCursorStyle(tcell.CursorStyleDefault)
	s.ShowCursor(3, 2)
	quit := make(chan struct{})
	style := tcell.StyleDefault
	go func() {
		for {
			ev := <-s.Poll()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyRune:
					switch ev.Rune() {
					case '0':
						s.SetContent(2, 2, '0', nil, 1, style)
						s.SetCursorStyle(tcell.CursorStyleDefault)
					case '1':
						s.SetContent(2, 2, '1', nil, 1, style)
						s.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
					case '2':
						s.SetContent(2, 2, '2', nil, 1, tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleSteadyBlock)
					case '3':
						s.SetContent(2, 2, '3', nil, 1, tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleBlinkingUnderline)
					case '4':
						s.SetContent(2, 2, '4', nil, 1, tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleSteadyUnderline)
					case '5':
						s.SetContent(2, 2, '5', nil, 1, tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleBlinkingBar)
					case '6':
						s.SetContent(2, 2, '6', nil, 1, tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleSteadyBar)
					}
					s.Show()

				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
					close(quit)
					return
				}
			}
		}
	}()
	<-quit
	s.Fini()
}
