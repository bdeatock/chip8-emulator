package main

import (
	"fmt"
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
	stepMode        bool // True = paused (can manually step)
	cyclesPerSecond int
	isRunning       bool
	isWasm          bool
	audioContext    *audio.Context
	audioPlayer     *audio.Player
	currentRom      []byte // stores last loaded rom to re-load after reset
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	options := parseCommandLineOptions()

	emu := chip8.New()
	fmt.Println("=== CHIP-8 Emulator initialized ===")

	if err := initEbiten(emu, options); err != nil {
		fmt.Printf("Failed to initialize ebiten: %v\n", err)
		os.Exit(1)
	}
}

func initEbiten(emu *chip8.Emulator, options *Options) error {
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
		return fmt.Errorf("error loading sound: %w", err)
	}

	env := newEnvironment()
	if game.isWasm {
		env.setupWasm(game)
	}

	if options.romPath != "" {
		if err := emu.LoadROMFromPath(options.romPath); err != nil {
			return fmt.Errorf("error loading ROM: %w", err)
		}
		game.isRunning = true
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Emulator Display")
	if !game.stepMode {
		ebiten.SetTPS(cyclesPerSecond)
	}

	if err := ebiten.RunGame(game); err != nil {
		return fmt.Errorf("error while running: %w", err)
	}

	return nil
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
	screen.Fill(colorBackground)

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
