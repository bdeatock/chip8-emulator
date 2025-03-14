package main

import (
	"flag"
	"fmt"
	"os"
)

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
