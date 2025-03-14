package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

func (g *Game) drawChip8Display(screen *ebiten.Image) {
	// Draw a border around screen
	borderColor := color.RGBA{100, 100, 100, 255}
	borderWidth := float32(1.0)
	displayWidth := float32(chip8DisplayWidth * chip8PixelSize)
	displayHeight := float32(chip8DisplayHeight * chip8PixelSize)
	// Top
	vector.DrawFilledRect(screen, marginX, marginY, displayWidth, borderWidth, borderColor, false)
	// Bottom
	vector.DrawFilledRect(screen, marginX, displayHeight+marginY, displayWidth, borderWidth, borderColor, false)
	// Left
	vector.DrawFilledRect(screen, marginX, marginY, borderWidth, displayHeight, borderColor, false)
	// Right
	vector.DrawFilledRect(screen, displayWidth+marginX, marginY, borderWidth, displayHeight, borderColor, false)

	// Draw emulator display (pixel grid)
	for x := range chip8DisplayWidth {
		for y := range chip8DisplayHeight {
			if g.emulator.Display[y*64+x] {
				vector.DrawFilledRect(
					screen,
					float32(x*chip8PixelSize)+marginX,
					float32(y*chip8PixelSize)+marginY,
					float32(chip8PixelSize),
					float32(chip8PixelSize),
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
	registersX := float64(chip8DisplayWidth*chip8PixelSize + marginX*2)
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
	memoryViewY := chip8DisplayHeight*chip8PixelSize + marginY*2
	byteWidth := chip8DisplayWidth*chip8PixelSize/memWidth - 4

	textOptions := &text.DrawOptions{}
	textOptions.GeoM.Translate(float64(memoryViewX), float64(memoryViewY))
	textOptions.ColorScale.ScaleWithColor(color.White)

	currentOpcode := g.emulator.GetCurrentOpcode()
	currentAddress := g.emulator.PC

	text.Draw(screen, fmt.Sprintf("PC: 0x%04X", g.emulator.PC), face, textOptions)
	textOptions.GeoM.Translate(topRowSpacing, 0)
	text.Draw(screen, fmt.Sprintf("I : 0x%04X", g.emulator.I), face, textOptions)
	textOptions.GeoM.Translate(topRowSpacing, 0)
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
		g.memViewStart = 4096 - memViewSize // Adjust start address to keep endAddress within memory bounds
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

			byteX := memoryViewX + memoryHeaderSpacing + (offset * byteWidth)
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
