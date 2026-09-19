//go:build windows

package tcell

import "errors"

// NewConsoleScreen is not implemented for this fork yet. Rune's GPU
// frontend uses Ebiten, not the Windows Console API.
func NewConsoleScreen() (Screen, error) {
	return nil, errors.New("tcell: windows console screen is not implemented in this fork")
}
