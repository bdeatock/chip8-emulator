package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bdeatock/chip8-emulator/chip8"
)

func main() {
	romPath := flag.String("rom", "", "Path to the ROM")
	flag.Parse()

	if *romPath == "" {
		fmt.Println("Please provide a ROM path using the -rom flag")
		os.Exit(1)
	}

	emu := chip8.New()
	fmt.Println("=== CHIP-8 Emulator initialized ===")

	if err := emu.LoadROM(*romPath); err != nil {
		fmt.Printf("Error loading ROM: %v\n", err)
		os.Exit(1)
	}
	emu.Print()
	for {
		// Wait for user to press Enter before continuing
		fmt.Println("\nPress Enter to continue to next cycle...")
		fmt.Scanln()
		emu.RunCycle()
		emu.Print()
	}
}
