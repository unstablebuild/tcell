package tcell

import (
	"testing"
)

func BenchmarkIntegrationLarge(b *testing.B) {
	benchmarkIntegration(b, false, 800, 600)
}

func BenchmarkIntegrationMedium(b *testing.B) {
	benchmarkIntegration(b, false, 200, 100)
}

func BenchmarkIntegrationSmall(b *testing.B) {
	benchmarkIntegration(b, false, 80, 30)
}


func BenchmarkIntegrationLargeResize(b *testing.B) {
	benchmarkIntegration(b, true, 800, 600)
}

func BenchmarkIntegrationMediumResize(b *testing.B) {
	benchmarkIntegration(b, true, 200, 100)
}

func BenchmarkIntegrationSmallResize(b *testing.B) {
	benchmarkIntegration(b, true, 80, 30)
}

func benchmarkIntegration(b *testing.B, resize bool, width, height int) {
	s := NewSimulationScreen()
	err := s.Init()
	if err != nil {
		b.Fatalf("Init screen: %v", err)
	}
	s.SetSize(width, height)
	st := StyleDefault.Foreground(ColorBlack).Background(ColorWhite)
	s.Fill(' ', st)

	b.ResetTimer()

	// simulate termbox event loop
	if !resize {
		for i := 0; i < b.N; i++ {
			clear(s)
			makeBox(s, i)
			s.Show()
		}
	} else {
		for i := 0; i < b.N; i++ {
			if i%10 == 0 {
				s.SetSize(width*1, height+1)
			} else if i%5 == 0 {
				s.SetSize(width, height)
			}
			clear(s)
			makeBox(s, i)
			s.Show()
		}
	}
}

var glyphs = []rune{'@', '#', '&', '*', '=', '%', 'Z', 'A'}

func makeBox(s Screen, i int) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	lx := i % w
	ly := i % h
	lw := i % (w - lx)
	lh := i % (h - ly)
	st := StyleDefault
	gl := ' '
	st = st.Reverse(i%2 == 0)
	gl = glyphs[i%len(glyphs)]

	for row := 0; row < lh; row++ {
		for col := 0; col < lw; col++ {
			s.SetContent(lx+col, ly+row, gl, nil, 1, st)
		}
	}
}

func clear(s Screen) {
	w, h := s.Size()
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			s.SetContent(col, row, ' ', nil, 1, StyleDefault)
		}
	}
}
