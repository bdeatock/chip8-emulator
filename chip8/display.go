package chip8

import "fmt"

// printDisplay renders the current state of the CHIP-8 display to the console.
func (e *Emulator) printDisplay() {
	for y := range 32 {
		fmt.Print("|")
		for x := range DisplayWidth {
			if e.Display[y*DisplayWidth+x] {
				fmt.Print("██")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Print("|\n")
	}
}

// Draws sprite with specified height at specified coordinates.
// Sprite is read from address pointed to by Index register.
func (e *Emulator) drawSprite(xPos, yPos, height int) {
	// Wrap coordinates
	xPos = xPos % DisplayWidth
	yPos = yPos % DisplayHeight

	// Reset collision flag
	e.Registers[0xF] = 0

	for row := range height {
		if yPos+row >= DisplayHeight {
			break
		}

		sprite := e.Memory[e.I+uint16(row)]

		for col := range 8 {
			if sprite&(128>>col) > 0 && xPos+col < DisplayWidth {
				if !e.flipPixel(xPos+col, yPos+row) {
					// If pixel turned off, set collision flag
					e.Registers[0xF] = 1
				}
			}
		}
	}
}

// Flips a pixel at coord, and returns the resulting state of pixel
func (e *Emulator) flipPixel(x int, y int) bool {
	e.Display[y*DisplayWidth+x] = !e.Display[y*DisplayWidth+x]

	return e.Display[y*DisplayWidth+x]
}

func (e *Emulator) clearDisplay() {
	for i := range e.Display {
		e.Display[i] = false
	}
}
