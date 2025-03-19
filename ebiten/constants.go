package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Display constants
const (
	// CHIP-8 display constants
	chip8PixelSize     = 8  // Size of each CHIP-8 pixel in screen pixels
	chip8DisplayWidth  = 64 // CHIP-8 display width in pixels
	chip8DisplayHeight = 32 // CHIP-8 display height in pixels

	// Ebiten window display constants
	marginX      = 11                                                                        // Horizontal margin for display elements
	marginY      = 5                                                                         // Vertical margin for display elements
	screenWidth  = chip8DisplayWidth*chip8PixelSize + marginX*3 + rightSpacing + 40          // Total width of the application window
	screenHeight = chip8DisplayHeight*chip8PixelSize + marginY*2 + lineHeight*(memNumRows+4) // Total height of the application window
	rightSpacing = marginY + 80

	// Memory display constants
	memWidth            = 16  // Number of bytes per memory view row
	memNumRows          = 6   // Number of rows in memory view
	lineHeight          = 20  // Height of each text line
	memoryHeaderSpacing = 70  // X Spacing for memory header
	topRowSpacing       = 120 // X Spacing for elements displayed above memory (PC, I, delaytimer, soundtimer)
)

// Sound constants
const (
	sampleRate = 48000
	frequency  = 440 // 440Hz = A4 Note
)

// Input key mapping
var keyArray = [16]ebiten.Key{
	ebiten.KeyX, // 0x0
	ebiten.Key1, // 0x1
	ebiten.Key2, // 0x2
	ebiten.Key3, // 0x3
	ebiten.KeyQ, // 0x4
	ebiten.KeyW, // 0x5
	ebiten.KeyE, // 0x6
	ebiten.KeyA, // 0x7
	ebiten.KeyS, // 0x8
	ebiten.KeyD, // 0x9
	ebiten.KeyZ, // 0xA
	ebiten.KeyC, // 0xB
	ebiten.Key4, // 0xC
	ebiten.KeyR, // 0xD
	ebiten.KeyF, // 0xE
	ebiten.KeyV, // 0xF
}

// Colours
var colorBackground = color.RGBA{
	51, 51, 51, 255,
}
var colorPrimary = color.RGBA{
	245, 245, 245, 255,
}
var colorAccent = color.RGBA{
	0, 255, 0, 255,
}
