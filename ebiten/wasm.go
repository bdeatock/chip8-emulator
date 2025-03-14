//go:build js

package main

import "syscall/js"

type jsEnvironment struct{}

func (j *jsEnvironment) setupWasm(game *Game) {
	js.Global().Set("loadROM", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return nil
		}

		romData := make([]byte, args[0].Length())
		js.CopyBytesToGo(romData, args[0])

		game.emulator.LoadROMFromData(romData)
		game.isRunning = true

		return nil
	}))
}

func newEnvironment() environment {
	return &jsEnvironment{}
}
