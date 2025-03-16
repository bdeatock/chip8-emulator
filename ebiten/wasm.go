//go:build js

package main

import (
	"syscall/js"
)

type jsEnvironment struct{}

func (je *jsEnvironment) setupWasm(game *Game) {
	js.Global().Set("loadROM", js.FuncOf(createLoadROMHandler(game)))
	js.Global().Set("switchMode", js.FuncOf(createSwitchModeHandler(game)))
	js.Global().Set("updateCycleRate", js.FuncOf(createSetCycleRateHandler(game)))
	js.Global().Set("toggleLegacyShift", js.FuncOf(createToggleLegacyShiftHandler(game)))
	js.Global().Set("toggleLegacyJump", js.FuncOf(createToggleLegacyJumpHandler(game)))
	js.Global().Set("toggleLegacyStoreLoad", js.FuncOf(createToggleLegacyStoreLoadHandler(game)))
	js.Global().Set("resetEmulator", js.FuncOf(createResetEmulatorHandler(game)))
}

func createLoadROMHandler(g *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return js.ValueOf(map[string]any{
				"error": "No ROM data provided",
			})
		}

		romData := make([]byte, args[0].Length())
		js.CopyBytesToGo(romData, args[0])

		g.emulator.Reset()
		if err := g.emulator.LoadROMFromData(romData); err != nil {
			return js.ValueOf(map[string]any{
				"error": err.Error(),
			})
		}
		g.currentRom = romData
		g.isRunning = true

		return nil
	}
}

// Legacy mode: Set VX to VY then shift, Modern mode: Shift VX directly
func createToggleLegacyShiftHandler(g *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		g.emulator.Config.LegacyShift = !g.emulator.Config.LegacyShift
		return g.emulator.Config.LegacyShift
	}
}

// Legacy mode: Jump to NNN + V0, Modern mode: Jump to NNN + VX
func createToggleLegacyJumpHandler(g *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		g.emulator.Config.LegacyJump = !g.emulator.Config.LegacyJump
		return g.emulator.Config.LegacyJump
	}
}

// Legacy mode: Increment I after store/load, Modern mode: I remains unchanged
func createToggleLegacyStoreLoadHandler(g *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		g.emulator.Config.LegacyStoreLoad = !g.emulator.Config.LegacyStoreLoad
		return g.emulator.Config.LegacyStoreLoad
	}
}

func createSwitchModeHandler(g *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		g.ToggleStepMode()
		return nil
	}
}

func createSetCycleRateHandler(g *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return js.ValueOf(map[string]any{
				"error": "No cycle rate provided",
			})
		}

		cycleRate := args[0].Int()
		cycleRate = max(1, cycleRate)

		g.SetCyclesPerSecond(cycleRate)
		return nil
	}
}

func createResetEmulatorHandler(g *Game) func(js.Value, []js.Value) any {
	return func(this js.Value, args []js.Value) any {
		g.emulator.Reset()

		if g.currentRom != nil {
			if err := g.emulator.LoadROMFromData(g.currentRom); err != nil {
				return js.ValueOf(map[string]any{
					"error": err.Error(),
				})
			}
		}
		return nil
	}
}

func newEnvironment() environment {
	return &jsEnvironment{}
}
