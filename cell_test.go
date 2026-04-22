package tcell

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyleUnion(t *testing.T) {
	tsuite := []struct {
		inA, inB Style
		wantOut  Style
	}{
		{},

		{Style{Fg: ColorRed}, Style{Fg: ColorDefault}, Style{Fg: ColorRed}},
		{Style{Bg: ColorRed}, Style{Bg: ColorDefault}, Style{Bg: ColorRed}},

		{Style{Fg: ColorWhite}, Style{Fg: ColorRed}, Style{Fg: ColorRed}},
		{Style{Bg: ColorWhite}, Style{Bg: ColorRed}, Style{Bg: ColorRed}},

		{Style{Fg: ColorRed, Attrs: AttrBold | AttrUnderline | AttrReverse},
			Style{Fg: ColorRed, Attrs: AttrBold | AttrUnderline | AttrReverse},
			Style{Fg: ColorRed, Attrs: AttrBold | AttrUnderline | AttrReverse}},
		{Style{Bg: ColorRed, Attrs: AttrBold | AttrUnderline | AttrReverse},
			Style{Bg: ColorRed, Attrs: AttrBold | AttrUnderline | AttrReverse},
			Style{Bg: ColorRed, Attrs: AttrBold | AttrUnderline | AttrReverse}},

		{Style{Attrs: AttrBold}, Style{Fg: ColorRed}, Style{Fg: ColorRed, Attrs: AttrBold}},
		{Style{Attrs: AttrBold}, Style{Bg: ColorRed}, Style{Bg: ColorRed, Attrs: AttrBold}},
		{Style{Attrs: AttrUnderline}, Style{Fg: ColorRed}, Style{Fg: ColorRed, Attrs: AttrUnderline}},
		{Style{Attrs: AttrUnderline}, Style{Bg: ColorRed}, Style{Bg: ColorRed, Attrs: AttrUnderline}},
		{Style{Attrs: AttrReverse}, Style{Fg: ColorRed}, Style{Fg: ColorRed, Attrs: AttrReverse}},
		{Style{Attrs: AttrReverse}, Style{Bg: ColorRed}, Style{Bg: ColorRed, Attrs: AttrReverse}},

		{Style{Fg: ColorRed}, Style{Attrs: AttrBold}, Style{Fg: ColorRed, Attrs: AttrBold}},
		{Style{Bg: ColorRed}, Style{Attrs: AttrBold}, Style{Bg: ColorRed, Attrs: AttrBold}},
		{Style{Fg: ColorRed}, Style{Attrs: AttrUnderline}, Style{Fg: ColorRed, Attrs: AttrUnderline}},
		{Style{Bg: ColorRed}, Style{Attrs: AttrUnderline}, Style{Bg: ColorRed, Attrs: AttrUnderline}},
		{Style{Fg: ColorRed}, Style{Attrs: AttrReverse}, Style{Fg: ColorRed, Attrs: AttrReverse}},
		{Style{Bg: ColorRed}, Style{Attrs: AttrReverse}, Style{Bg: ColorRed, Attrs: AttrReverse}},

		{Style{Attrs: AttrReverse | AttrBold | AttrUnderline},
			Style{Bg: ColorRed},
			Style{Bg: ColorRed, Attrs: AttrReverse | AttrBold | AttrUnderline}},
		{Style{Bg: ColorRed},
			Style{Attrs: AttrReverse | AttrBold | AttrUnderline},
			Style{Bg: ColorRed, Attrs: AttrReverse | AttrBold | AttrUnderline}},

		{Style{Fg: ColorDefault}, Style{Fg: ColorRed}, Style{Fg: ColorRed}},
		{Style{Bg: ColorDefault}, Style{Bg: ColorRed}, Style{Bg: ColorRed}},
	}

	for i, tcase := range tsuite {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actualOut := unionStyle(tcase.inA, tcase.inB)
			assert.Equal(t, tcase.wantOut, actualOut)
		})
	}
}

// TestCellBufferResizePreservesContent verifies that CellBuffer.Resize
// honors its documented contract: "resize the cells array, with
// different dimensions, while preserving the original contents."
//
// This is a regression test for a bug where Resize allocated a fresh
// zeroed buffer and discarded every cell previously written via
// SetContentWidth. The symptom visible to applications built on top
// of tScreen is a blank frame after every terminal resize, because
// Show() absorbs a pending resize AFTER Clear+SetContent have run for
// that frame — and the absorption wiped every cell that was just
// written.
func TestCellBufferResizePreservesContent(t *testing.T) {
	t.Run("content preserved within overlap when enlarging", func(t *testing.T) {
		var cb CellBuffer
		cb.Resize(10, 4)
		st := Style{Fg: ColorRed}
		cb.SetContentWidth(1, 2, 'a', nil, 1, st)
		cb.SetContentWidth(9, 3, 'b', nil, 1, st)

		cb.Resize(20, 8)

		w, h := cb.Size()
		assert.Equal(t, 20, w)
		assert.Equal(t, 8, h)

		mainA, _, styleA, _, _ := cb.GetContent(1, 2)
		assert.Equal(t, 'a', mainA,
			"cell (1,2) should be preserved when enlarging")
		assert.Equal(t, st, styleA)

		mainB, _, styleB, _, _ := cb.GetContent(9, 3)
		assert.Equal(t, 'b', mainB,
			"cell (9,3) should be preserved when enlarging")
		assert.Equal(t, st, styleB)
	})

	t.Run("content preserved within overlap when shrinking", func(t *testing.T) {
		var cb CellBuffer
		cb.Resize(20, 8)
		st := Style{Bg: ColorBlue}
		cb.SetContentWidth(3, 3, 'x', nil, 1, st)

		cb.Resize(10, 4)

		w, h := cb.Size()
		assert.Equal(t, 10, w)
		assert.Equal(t, 4, h)

		mainX, _, styleX, _, _ := cb.GetContent(3, 3)
		assert.Equal(t, 'x', mainX,
			"cell (3,3) should be preserved when shrinking")
		assert.Equal(t, st, styleX)
	})

	t.Run("resized cells are marked dirty so they redraw", func(t *testing.T) {
		var cb CellBuffer
		cb.Resize(4, 2)
		cb.SetContentWidth(1, 1, 'z', nil, 1, Style{Fg: ColorGreen})
		// Mark clean so we can detect the re-dirty after Resize.
		cb.ClearDirty(1, 1)
		assert.False(t, cb.DirtyAt(1, 1))

		cb.Resize(6, 3)

		assert.True(t, cb.DirtyAt(1, 1),
			"Resize should invalidate preserved cells so they repaint")
	})

	t.Run("newly exposed cells are zero-valued", func(t *testing.T) {
		var cb CellBuffer
		cb.Resize(4, 2)
		cb.SetContentWidth(1, 1, 'z', nil, 1, Style{Fg: ColorGreen})

		cb.Resize(6, 3)

		// Read the underlying cell directly: GetContent normalises
		// zero-valued mainc to a space via ProcessCell.
		c := cb.GetCell(5, 2)
		assert.Equal(t, rune(0), c.currMain,
			"freshly-allocated cells should start zero-valued")
		assert.Equal(t, Style{}, c.currStyle)
	})
}
