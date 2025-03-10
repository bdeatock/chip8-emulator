package main

import (
	"fmt"

	"github.com/bdeatock/chip8-emulator/chip8"
)

func main() {
	emu := chip8.New()

	fmt.Println("=== CHIP-8 Emulator initialized ===")
	emu.Print()
}