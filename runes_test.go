// Copyright 2018 The TCell Authors
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
	"testing"
)

func TestCanDisplayUTF8(t *testing.T) {
	s := mkTestScreen(t)
	defer s.Fini()

	if !s.CanDisplay('a', true) {
		t.Errorf("Should be able to display 'a'")
	}
	if !s.CanDisplay(RuneHLine, true) {
		t.Errorf("Should be able to display hline (with fallback)")
	}
	if !s.CanDisplay(RuneHLine, false) {
		t.Errorf("Should be able to display hline (no fallback)")
	}
	if !s.CanDisplay('⌀', false) {
		t.Errorf("Should be able to display null")
	}
}
