//go:build js

package main

import (
	"syscall/js"
)

type jsEnvironment struct{}

func (j *jsEnvironment) setupWasm(game *Game) {
	js.Global().Set("loadROM", js.FuncOf(createLoadROMHandler(game)))
	js.Global().Set("switchMode", js.FuncOf(createSwitchModeHandler(game)))
	js.Global().Set("updateCycleRate", js.FuncOf(createSetCycleRateHandler(game)))
	js.Global().Set("toggleLegacyShift", js.FuncOf(createToggleLegacyShiftHandler(game)))
	js.Global().Set("toggleLegacyJump", js.FuncOf(createToggleLegacyJumpHandler(game)))
	js.Global().Set("toggleLegacyStoreLoad", js.FuncOf(createToggleLegacyStoreLoadHandler(game)))
}

func createLoadROMHandler(game *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return js.ValueOf(map[string]any{
				"error": "No ROM data provided",
			})
		}

		romData := make([]byte, args[0].Length())
		js.CopyBytesToGo(romData, args[0])

		game.emulator.Reset()
		if err := game.emulator.LoadROMFromData(romData); err != nil {
			return js.ValueOf(map[string]any{
				"error": err.Error(),
			})
		}
		game.isRunning = true

		return nil
	}
}

// Legacy mode: Set VX to VY then shift, Modern mode: Shift VX directly
func createToggleLegacyShiftHandler(game *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		game.emulator.Config.LegacyShift = !game.emulator.Config.LegacyShift
		return game.emulator.Config.LegacyShift
	}
}

// Legacy mode: Jump to NNN + V0, Modern mode: Jump to NNN + VX
func createToggleLegacyJumpHandler(game *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		game.emulator.Config.LegacyJump = !game.emulator.Config.LegacyJump
		return game.emulator.Config.LegacyJump
	}
}

// Legacy mode: Increment I after store/load, Modern mode: I remains unchanged
func createToggleLegacyStoreLoadHandler(game *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		game.emulator.Config.LegacyStoreLoad = !game.emulator.Config.LegacyStoreLoad
		return game.emulator.Config.LegacyStoreLoad
	}
}

func createSwitchModeHandler(game *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		game.ToggleStepMode()
		return nil
	}
}

func createSetCycleRateHandler(game *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return js.ValueOf(map[string]any{
				"error": "No cycle rate provided",
			})
		}

		cycleRate := args[0].Int()
		cycleRate = max(1, cycleRate)

		game.SetCyclesPerSecond(cycleRate)
		return nil
	}
}

func newEnvironment() environment {
	return &jsEnvironment{}
}
