//go:build !js

package main

type defaultEnvironment struct{}

func (d *defaultEnvironment) setupWasm(game *Game) {
	// Do nothing in non-WASM builds
}

func newEnvironment() environment {
	return &defaultEnvironment{}
}
