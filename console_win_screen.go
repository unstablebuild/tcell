//go:build windows

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

// These methods exist on tScreen and CellBuffer but were never wired up on
// the Windows console backend after the fork's Screen interface grew extra
// parameters (width on SetContent, dirty on GetContent, UnionStyle).
// Implementing them against the existing CellBuffer lets GOOS=windows compile.

func (s *cScreen) Clear() {
	s.Fill(' ', s.style)
}

func (s *cScreen) Fill(r rune, style Style) {
	s.Lock()
	defer s.Unlock()
	s.cells.Fill(r, style)
}

func (s *cScreen) GetContent(x, y int) (primary rune, combining []rune, style Style, width uint8, dirty bool) {
	s.Lock()
	defer s.Unlock()
	return s.cells.GetContent(x, y)
}

func (s *cScreen) SetContent(x int, y int, primary rune, combining []rune, width uint8, style Style) {
	s.Lock()
	defer s.Unlock()
	s.cells.SetContentWidth(x, y, primary, combining, width, style)
}

func (s *cScreen) UnionStyle(x int, y int, style Style) {
	s.Lock()
	defer s.Unlock()
	s.cells.UnionStyle(x, y, style)
}

func (s *cScreen) Bell() {
	_ = s.Beep()
}

func (s *cScreen) Sync() {
	s.Show()
}
