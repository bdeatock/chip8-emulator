//go:build !js

package main

// defaultEnvironment is a dummy implementation that satisfies the environment interface
// when not in WebAssembly
type defaultEnvironment struct{}

func (de *defaultEnvironment) setupWasm(game *Game) {
	// Do nothing in non-WASM builds
}

func newEnvironment() environment {
	return &defaultEnvironment{}
}
