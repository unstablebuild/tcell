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
