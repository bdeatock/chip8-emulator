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
