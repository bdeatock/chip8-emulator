package main

import "github.com/hajimehoshi/ebiten/v2"

// Display constants
const (
	// CHIP-8 display constants
	pixelSize  = 8  // Size of each CHIP-8 pixel in screen pixels
	gridWidth  = 64 // CHIP-8 display width in pixels
	gridHeight = 32 // CHIP-8 display height in pixels

	// Ebiten window display constants
	marginX      = 11                                   // Horizontal margin for display elements
	marginY      = 5                                    // Vertical margin for display elements
	screenWidth  = gridWidth*pixelSize + marginX*3 + 65 // Total width of the application window
	screenHeight = 480                                  // Total height of the application window

	// Memory display constants
	memWidth   = 16 // Number of bytes per memory view row
	memNumRows = 8  // Number of rows in memory view
	lineHeight = 20 // Height of each text line
)

// Input key mapping
var keyArray = [16]ebiten.Key{
	ebiten.Key1,
	ebiten.Key2,
	ebiten.Key3,
	ebiten.Key4,
	ebiten.KeyQ,
	ebiten.KeyW,
	ebiten.KeyE,
	ebiten.KeyR,
	ebiten.KeyA,
	ebiten.KeyS,
	ebiten.KeyD,
	ebiten.KeyF,
	ebiten.KeyZ,
	ebiten.KeyX,
	ebiten.KeyC,
	ebiten.KeyV,
}
