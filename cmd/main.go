package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/bdeatock/chip8-emulator/chip8"
)

func main() {
	options := parseCommandLineOptions()

	emu := chip8.New()
	fmt.Println("=== CHIP-8 Emulator initialized ===")

	if err := emu.LoadROM(options.romPath); err != nil {
		fmt.Printf("Error loading ROM: %v\n", err)
		os.Exit(1)
	}
	emu.Print()

	if options.cycleMode == "continuous" {
		runContinuousMode(emu, options.cyclesPerSecond, options.displayRate)
	} else {
		runStepMode(emu)
	}
}

type options struct {
	romPath         string
	cycleMode       string
	cyclesPerSecond int
	displayRate     int
}

func parseCommandLineOptions() *options {
	romPath := flag.String("rom", "", "Path to the ROM")
	cycleMode := flag.String("mode", "continuous", "Execution mode: 'step' for manual stepping or 'continuous' for continuous execution")
	cyclesPerSecond := flag.Int("speed", 700, "Number of cycles per second in continuous mode")
	displayRate := flag.Int("refresh", 60, "Display refresh rate in Hz")
	flag.Parse()

	if *romPath == "" {
		fmt.Println("Please provide a ROM path using the -rom flag")
		os.Exit(1)
	}

	if *cycleMode != "step" && *cycleMode != "continuous" {
		fmt.Println("Invalid mode. Use 'step' or 'continuous'")
		os.Exit(1)
	}

	if *cyclesPerSecond <= 0 {
		fmt.Println("Speed must be a positive number")
		os.Exit(1)
	}
	if *displayRate <= 0 {
		fmt.Println("Display rate must be a positive number")
		os.Exit(1)
	}

	return &options{
		romPath:         *romPath,
		cycleMode:       *cycleMode,
		cyclesPerSecond: *cyclesPerSecond,
		displayRate:     min(*displayRate, *cyclesPerSecond), // display rate has no reason to ever be above cycle rate
	}
}

func runContinuousMode(emu *chip8.Emulator, cyclesPerSecond int, displayRate int) {
	errCh := emu.Run(cyclesPerSecond)

	displayRefreshClock := time.NewTicker(time.Second / time.Duration(displayRate))

	go func() {
		for range displayRefreshClock.C {
			emu.Print()
		}
	}()

	if err := <-errCh; err != nil {
		fmt.Printf("\nEmulation stopped with error: %v\n", err)
	}
}

func runStepMode(emu *chip8.Emulator) {
	for {
		fmt.Println("\nPress Enter to continue to next cycle...")
		fmt.Scanln()
		fmt.Printf("Executing opcode: %s\n", emu.GetCurrentOpcode(false))
		if err := emu.Step(time.Second / 4); err != nil {
			fmt.Printf("\nEmulation stopped with error: %v\n", err)
			return
		}
		emu.Print()
	}
}
