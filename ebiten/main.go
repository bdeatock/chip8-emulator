package main

import (
	"fmt"
	"image/color"
	"os"
	"runtime"
	"time"

	"github.com/bdeatock/chip8-emulator/chip8"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// environment abstracts platform-specific functionality
type environment interface {
	setupWasm(game *Game)
}

type Game struct {
	cycleCount      int
	emulator        *chip8.Emulator
	memViewStart    uint16
	stepMode        bool
	cyclesPerSecond int
	isRunning       bool
	isWasm          bool
	audioContext    *audio.Context
	audioPlayer     *audio.Player
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	options := parseCommandLineOptions()

	emu := chip8.New()
	fmt.Println("=== CHIP-8 Emulator initialized ===")

	initEbiten(emu, options)
}

func initEbiten(emu *chip8.Emulator, options *Options) {
	cyclesPerSecond := 4
	if options.cycleMode != "step" {
		cyclesPerSecond = options.cyclesPerSecond
	}

	game := &Game{
		emulator:        emu,
		stepMode:        options.cycleMode == "step",
		cyclesPerSecond: cyclesPerSecond,
		isWasm:          runtime.GOOS == "js",
	}

	if err := game.initSound(); err != nil {
		fmt.Printf("Error loading sound: %v\n", err)
		os.Exit(1)
	}

	env := newEnvironment()
	if game.isWasm {
		env.setupWasm(game)
	}

	if options.romPath != "" {
		if err := emu.LoadROMFromPath(options.romPath); err != nil {
			fmt.Printf("Error loading ROM: %v\n", err)
			os.Exit(1)
		}
		game.isRunning = true
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Emulator Display")
	if !game.stepMode {
		ebiten.SetTPS(cyclesPerSecond)
	}

	if err := ebiten.RunGame(game); err != nil {
		fmt.Printf("Error while running: %v\n", err)
		os.Exit(1)
	}
}

func (g *Game) Update() error {
	if !g.isRunning {
		return nil
	}

	if g.handleInput() {
		// time to run a cycle
		g.cycleCount++
		deltaTime := time.Second / time.Duration(g.cyclesPerSecond)
		if err := g.emulator.Step(deltaTime); err != nil {
			return err
		}
	}

	g.handleSound()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{40, 40, 40, 255})

	g.drawChip8Display(screen)
	g.drawUI(screen)
}

func (g *Game) ToggleStepMode() {
	if g.stepMode {
		g.stepMode = false
		ebiten.SetTPS(g.cyclesPerSecond)
	} else {
		g.stepMode = true
		ebiten.SetTPS(60)
	}
}

func (g *Game) SetCyclesPerSecond(cycles int) {
	g.cyclesPerSecond = cycles
	if !g.stepMode {
		ebiten.SetTPS(g.cyclesPerSecond)
	}
}
