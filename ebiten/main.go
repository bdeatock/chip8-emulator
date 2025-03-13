package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/bdeatock/chip8-emulator/chip8"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenHeight = 480
	pixelSize    = 8
	gridWidth    = 64
	gridHeight   = 32
	marginX      = 11
	marginY      = 5
	screenWidth  = gridWidth*pixelSize + marginX*3 + 65
	memWidth     = 16
	memNumRows   = 8
	lineHeight   = 20
)

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

type Game struct {
	cycleCount      int
	emulator        *chip8.Emulator
	memViewStart    uint16
	stepMode        bool
	cyclesPerSecond int
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	options := parseCommandLineOptions()

	emu := chip8.New()
	fmt.Println("=== CHIP-8 Emulator initialized ===")

	if err := emu.LoadROM(options.romPath); err != nil {
		fmt.Printf("Error loading ROM: %v\n", err)
		os.Exit(1)
	}

	cyclesPerSecond := 4
	if options.cycleMode != "step" {
		cyclesPerSecond = options.cyclesPerSecond
	}

	initEbiten(emu, cyclesPerSecond, options.cycleMode == "step")
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

func initEbiten(emu *chip8.Emulator, cyclesPerSecond int, stepMode bool) {
	game := &Game{
		emulator:        emu,
		stepMode:        stepMode,
		cyclesPerSecond: cyclesPerSecond,
	}

	ebiten.SetWindowSize(screenWidth*1.5, screenHeight*1.5)
	ebiten.SetWindowTitle("Emulator Display")
	if !stepMode {
		ebiten.SetTPS(cyclesPerSecond)
	}

	if err := ebiten.RunGame(game); err != nil {
		fmt.Printf("Error while running: %v\n", err)
		os.Exit(1)
	}
}

func (g *Game) Update() error {
	for i, key := range keyArray {
		if inpututil.IsKeyJustPressed(key) {
			g.emulator.PressKey(byte(i))
		} else if inpututil.IsKeyJustReleased(key) {
			g.emulator.ReleaseKey(byte(i))
		}
	}

	if !g.stepMode || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.cycleCount++
		deltaTime := time.Second / time.Duration(g.cyclesPerSecond)
		return g.emulator.RunCycle(deltaTime)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{40, 40, 40, 255})

	g.drawDisplay(screen)
	g.drawRegisters(screen)
	g.drawMemoryView(screen)
}

func (g *Game) drawDisplay(screen *ebiten.Image) {
	// Draw a border around screen
	borderColor := color.RGBA{100, 100, 100, 255}
	borderWidth := float32(1.0)
	displayWidth := float32(gridWidth * pixelSize)
	displayHeight := float32(gridHeight * pixelSize)
	// Top
	vector.DrawFilledRect(screen, marginX, marginY, displayWidth, borderWidth, borderColor, false)
	// Bottom
	vector.DrawFilledRect(screen, marginX, displayHeight+marginY, displayWidth, borderWidth, borderColor, false)
	// Left
	vector.DrawFilledRect(screen, marginX, marginY, borderWidth, displayHeight, borderColor, false)
	// Right
	vector.DrawFilledRect(screen, displayWidth+marginX, marginY, borderWidth, displayHeight, borderColor, false)

	// Draw emulator display (pixel grid)
	for x := range gridWidth {
		for y := range gridHeight {
			if g.emulator.Display[y*64+x] {
				vector.DrawFilledRect(
					screen,
					float32(x*pixelSize)+marginX,
					float32(y*pixelSize)+marginY,
					float32(pixelSize),
					float32(pixelSize),
					color.RGBA{200, 200, 100, 255},
					false,
				)
			}
		}
	}
}

func (g *Game) drawRegisters(screen *ebiten.Image) {
	face := text.NewGoXFace(basicfont.Face7x13)
	textOptions := &text.DrawOptions{}
	textOptions.ColorScale.ScaleWithColor(color.White)

	// Draw registers on right of the display
	registersX := float64(gridWidth*pixelSize + marginX*2)
	registersY := float64(marginY)
	textOptions.GeoM.Translate(registersX, registersY)

	for i := range byte(0x10) {
		text.Draw(screen, fmt.Sprintf("V%X:  0x%02X", i, g.emulator.Registers[i]), face, textOptions)
		textOptions.GeoM.Translate(0, lineHeight)

	}

}

func (g *Game) drawMemoryView(screen *ebiten.Image) {
	face := text.NewGoXFace(basicfont.Face7x13)

	memoryViewX := marginX
	memoryViewY := gridHeight*pixelSize + marginY*2
	byteWidth := gridWidth*pixelSize/memWidth - 4

	textOptions := &text.DrawOptions{}
	textOptions.GeoM.Translate(float64(memoryViewX), float64(memoryViewY))
	textOptions.ColorScale.ScaleWithColor(color.White)

	currentOpcode := g.emulator.GetCurrentOpcode()
	currentAddress := g.emulator.PC

	text.Draw(screen, fmt.Sprintf("PC: 0x%04X", g.emulator.PC), face, textOptions)
	textOptions.GeoM.Translate(100, 0)
	text.Draw(screen, fmt.Sprintf("I : 0x%04X", g.emulator.I), face, textOptions)
	textOptions.GeoM.Translate(100, 0)
	textOptions.ColorScale.ScaleWithColor(color.RGBA{0, 255, 0, 255})
	text.Draw(screen, fmt.Sprintf("Opcode: 0x%04X", currentOpcode), face, textOptions)
	textOptions.ColorScale.Reset()

	memViewSize := uint16(memWidth * memNumRows)

	// Only move memory view if currentAddress falls out of visible range
	if g.memViewStart == 0 || currentAddress < g.memViewStart || currentAddress >= g.memViewStart+memViewSize {
		g.memViewStart = (currentAddress / 16) * 16 // Round down to nearest 16-byte boundary

	}

	endAddress := g.memViewStart + memViewSize

	if endAddress > 4096 {
		g.memViewStart -= endAddress - 4096 // Adjust start address to keep endAddress within memory bounds
		endAddress = 4096
	}

	// Draw memory rows
	for addr := g.memViewStart; addr < endAddress; addr += uint16(memWidth) {
		textOptions.GeoM.Translate(0, float64(lineHeight))

		// Draw row address
		rowText := fmt.Sprintf("0x%04X:", addr)
		textOptions.GeoM.SetElement(0, 2, float64(memoryViewX)) // set x to memoryViewX
		text.Draw(screen, rowText, face, textOptions)

		// Draw bytes in this row
		for offset := 0; offset < memWidth; offset++ {
			byteAddr := addr + uint16(offset)
			byteValue := g.emulator.Memory[byteAddr]

			byteX := memoryViewX + 70 + (offset * byteWidth)
			byteY := int(textOptions.GeoM.Element(1, 2))

			// Highlight the current opcode bytes (2 bytes)
			if byteAddr == currentAddress || byteAddr == currentAddress+1 {
				vector.DrawFilledRect(
					screen,
					float32(byteX-(byteWidth/5)),
					float32(byteY-(lineHeight/5)),
					float32(byteWidth),
					float32(lineHeight),
					color.RGBA{100, 100, 200, 255},
					false,
				)
			}

			textOptions.GeoM.SetElement(0, 2, float64(byteX)) // set x to byteX
			text.Draw(screen, fmt.Sprintf("%02X", byteValue), face, textOptions)
		}
	}
}
