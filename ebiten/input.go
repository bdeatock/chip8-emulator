package main

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Handles input and returns true if a cycle should happen
func (g *Game) handleInput() bool {
	for i, key := range keyArray {
		if inpututil.IsKeyJustPressed(key) {
			g.emulator.PressKey(byte(i))
		} else if inpututil.IsKeyJustReleased(key) {
			g.emulator.ReleaseKey(byte(i))
		}
	}

	return !g.stepMode || g.inputForStepCycle()
}

func (g *Game) inputForStepCycle() bool {
	for _, key := range stepKeys {
		if inpututil.IsKeyJustPressed(key) {
			return true
		}
	}
	return false
}
