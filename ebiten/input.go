package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) handleInput() error {
	for i, key := range keyArray {
		if inpututil.IsKeyJustPressed(key) {
			g.emulator.PressKey(byte(i))
		} else if inpututil.IsKeyJustReleased(key) {
			g.emulator.ReleaseKey(byte(i))
		}
	}

	if !g.stepMode || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.cycleCount++
		deltaTime := time.Second / time.Duration(g.cyclesPerSecond)
		return g.emulator.RunCycle(deltaTime)
	}

	return nil
}
