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

package terminfo

import (
	"bytes"
	"testing"
	"time"
)

// This terminfo entry is a stripped down version from
// xterm-256color, but I've added some of my own entries.
var testTerminfo = &Terminfo{
	Name:      "simulation_test",
	Columns:   80,
	Lines:     24,
	Colors:    256,
	Bell:      "\a",
	Blink:     "\x1b2ms$<20>something",
	Reverse:   "\x1b[7m",
	SetFg:     "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m",
	SetBg:     "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m",
	AltChars:  "``aaffggiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
	Mouse:     "\x1b[M",
	SetCursor: "\x1b[%i%p1%d;%p2%dH",
	PadChar:   "\x00",
	EnterUrl:  "\x1b]8;%p2%s;%p1%s\x1b\\",
}

func TestTerminfoExpansion(t *testing.T) {
	ti := testTerminfo

	// Tests %i and basic parameter strings too
	if string(ti.TGoto(7, 9)) != "\x1b[10;8H" {
		t.Error("TGoto expansion failed")
	}

	// This tests some conditionals
	if string(ti.TParm("A[%p1%2.2X]B", []int{47})) != "A[2F]B" {
		t.Error("TParm conditionals failed")
	}

	// Color tests.
	if string(ti.TParm(ti.SetFg, []int{7})) != "\x1b[37m" {
		t.Error("SetFg(7) failed")
	}
	if string(ti.TParm(ti.SetFg, []int{15})) != "\x1b[97m" {
		t.Error("SetFg(15) failed")
	}
	if string(ti.TParm(ti.SetFg, []int{200})) != "\x1b[38;5;200m" {
		t.Error("SetFg(200) failed")
	}

	type testCase struct {
		expect string
		format string
		params []int
	}

	cases := []testCase{
		{expect: "0a", format: "%p1%02x", params: []int{10}},
		{expect: "0A", format: "%p1%02X", params: []int{10}},
		{expect: "A", format: "%p1%c", params: []int{65}},
		{expect: "A", format: "%'A'%c", params: []int{}},
		{expect: "65", format: "%'A'%d", params: []int{}},
		{expect: "7", format: "%i%p1%p2%+%d", params: []int{2, 3}},
		{expect: "1%d", format: "1%%d", params: []int{}},
		{expect: "15", format: "%{3}%p1%*%d", params: []int{5}},
		{expect: " A", format: "%p1%2c", params: []int{65}},
		{expect: "0", format: "%pA%d", params: []int{}}, // missing/invalid parameter
		{expect: "5", format: "%p1%p2%/%d", params: []int{15, 3}},
		{expect: "0", format: "%p1%p2%/%d", params: []int{3, 15}},
		{expect: "0", format: "%p1%p2%/%d", params: []int{3, 0}},
		{expect: "3", format: "%p1%p2%m%d", params: []int{15, 4}},
		{expect: "0", format: "%p1%p2%m%d", params: []int{3, 0}},
		{expect: "0", format: "%d", params: []int{}}, // underflow
		{expect: "", format: "%s", params: []int{}},  // underflow
		{expect: "1", format: "%p1%p2%=%d", params: []int{3, 3}},
		{expect: "0", format: "%p1%p2%=%d", params: []int{3, 4}},
		{expect: "1", format: "%p1%p2%=%!%d", params: []int{3, 4}},
		{expect: "1", format: "%p1%p2%>%d", params: []int{4, 3}},
		{expect: "3", format: "%p1%p2%|%d", params: []int{1, 2}},
		{expect: "2", format: "%p1%p2%&%d", params: []int{2, 3}},
		{expect: "1", format: "%p1%p2%^%d", params: []int{2, 3}},
		{expect: "f", format: "%p1%~%{255}%&%x", params: []int{0xf0}},
		{expect: "%Z", format: "%Z", params: []int{2, 3}}, // unknown sequence
		{expect: "2", format: "%p1%Pa%{4}%{3}%ga%d", params: []int{2}},
		{expect: "2", format: "%p1%PA%{4}%{3}%gA%d", params: []int{2}},
		{expect: "0", format: "%p1%PA%{4}%{3}%ga%d", params: []int{2}},
		{expect: "0", format: "%p1%Pz%{4}%{3}%gZ%d", params: []int{2}},
	}

	for i := range cases {
		if res := ti.TParm(cases[i].format, cases[i].params); string(res) != cases[i].expect {
			t.Errorf("Format case %d failed: Format %q got %q", i, cases[i].expect, res)
		}
	}
	t.Logf("Tested %d cases", len(cases))
}

func TestTerminfoDelay(t *testing.T) {
	ti := testTerminfo
	buf := bytes.NewBuffer(nil)
	now := time.Now()
	ti.TPuts(buf, []byte(ti.Blink))
	then := time.Now()
	s := string(buf.Bytes())
	if s != "\x1b2mssomething" {
		t.Errorf("Terminfo delay failed: %s", s)
	}
	if then.Sub(now) < time.Millisecond*20 {
		t.Error("Too short delay")
	}
	if then.Sub(now) > time.Millisecond*50 {
		t.Error("Too late delay")
	}
}

func BenchmarkSetFgBg(b *testing.B) {
	ti := testTerminfo

	for i := 0; i < b.N; i++ {
		ti.TParm(ti.SetFg, []int{100, 200})
		ti.TParm(ti.SetBg, []int{100, 200})
	}
}
