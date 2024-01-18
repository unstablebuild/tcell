package tcell

import (
	"context"
	"io"
	"os"
	"testing"
)

func BenchmarkIntegration(b *testing.B) {
	suite := []struct {
		description string
		resize      bool
		width       int
		height      int
		simulation  bool
	}{
		{"simulation large", false, 800, 600, true},
		{"simulation medium", false, 200, 100, true},
		{"simulation small", false, 80, 30, true},
		{"simulation large with resize", true, 800, 600, true},
		{"simulation medium with resize", true, 200, 100, true},
		{"simulation small with resize", true, 80, 30, true},

		{"tscreen large", false, 800, 600, false},
		{"tscreen medium", false, 200, 100, false},
		{"tscreen small", false, 80, 30, false},
		{"tscreen large with resize", true, 800, 600, false},
		{"tscreen medium with resize", true, 200, 100, false},
		{"tscreen small with resize", true, 80, 30, false},
	}
	for _, test := range suite {
		b.Run(test.description, func(b *testing.B) {
			if test.simulation {
				benchmarkSimulation(b, test.resize, test.width, test.height)
				return
			}
			benchmarkTscreen(b, test.resize, test.width, test.height)
		})
	}
}

func benchmarkSimulation(b *testing.B, resize bool, width, height int) {
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
			s.Clear()
			makeBox(s, i)
			s.Show()
		}
	} else {
		for i := 0; i < b.N; i++ {
			if i%10 == 0 {
				s.SetSize(width+1, height+1)
			} else if i%5 == 0 {
				s.SetSize(width, height)
			}
			s.Clear()
			makeBox(s, i)
			s.Show()
		}
	}
}

func benchmarkTscreen(b *testing.B, resize bool, width, height int) {
	var tty benchTty
	ti, err := LookupTerminfo(os.Getenv("TERM"))
	if err != nil {
		b.Fatalf("lookup term info: %v", err)
	}
	s := tScreen{tty: &tty, ti: ti}
	bs := baseScreen{&s, s.GetCells()}
	if err := s.Init(); err != nil {
		b.Fatalf("Init screen: %v", err)
	}
	tty.setWindowSize(width, height)
	s.resize()
	st := StyleDefault.Foreground(ColorBlack).Background(ColorWhite)
	s.GetCells().Fill(' ', st)

	ctx, cancel := context.WithCancel(context.Background())
	go func() { // drain events
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.eventQ:
			}
		}
	}()

	b.Cleanup(func() {
		bs.Fini()
		cancel()
	})

	b.ResetTimer()

	// simulate termbox event loop
	if !resize {
		for i := 0; i < b.N; i++ {
			bs.Clear()
			makeBox(&bs, i)
			s.Show()
		}
	} else {
		for i := 0; i < b.N; i++ {
			if i%10 == 0 {
				tty.setWindowSize(width+1, height+1)
			} else if i%5 == 0 {
				tty.setWindowSize(width, height)
			}
			s.resize()
			bs.Clear()
			makeBox(&bs, i)
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

	if i%5 != 0 {
		st.Foreground(Color(i)%(ColorYellowGreen-ColorValid) + ColorValid)
	}

	for row := 0; row < lh; row++ {
		for col := 0; col < lw; col++ {
			s.SetContent(lx+col, ly+row, gl, nil, 1, st)
		}
	}
}

type benchTty struct {
	notifyResize func()
	size         WindowSize
}

func (b *benchTty) Start() error {
	return nil
}

func (b *benchTty) Stop() error {
	return nil
}

func (b *benchTty) Drain() error {
	return nil
}

func (b *benchTty) NotifyResize(cb func()) {
	b.notifyResize = cb
}

func (b *benchTty) WindowSize() (WindowSize, error) {
	return b.size, nil
}

func (b *benchTty) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (b *benchTty) Write(p []byte) (n int, err error) {
	return io.Discard.Write(p)
}

func (b *benchTty) Close() error {
	return nil
}

func (b *benchTty) setWindowSize(width, height int) {
	b.size.Width = width
	b.size.Height = height
	if b.notifyResize != nil {
		b.notifyResize()
	}
}
