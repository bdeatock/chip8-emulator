//go:build js

package main

import "syscall/js"

type jsEnvironment struct{}

func (j *jsEnvironment) setupWasm(game *Game) {
	js.Global().Set("loadROM", js.FuncOf(createLoadROMHandler(game)))
	js.Global().Set("switchMode", js.FuncOf(createModeSwitchHandler(game)))
	js.Global().Set("updateCycleRate", js.FuncOf(createSetCycleRateHandler(game)))
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

func createModeSwitchHandler(game *Game) func(js.Value, []js.Value) any {
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
