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

	if !g.stepMode || g.inputForStepCycle() {
		g.cycleCount++
		deltaTime := time.Second / time.Duration(g.cyclesPerSecond)
		return g.emulator.Step(deltaTime)
	}

	return nil
}

func (g *Game) inputForStepCycle() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace)
}
