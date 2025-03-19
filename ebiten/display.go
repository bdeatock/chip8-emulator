package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

func (g *Game) drawChip8Display(screen *ebiten.Image) {
	// Draw a border around screen
	borderColor := colorPrimary
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
					colorAccent,
					false,
				)
			}
		}
	}
}

func (g *Game) drawUI(screen *ebiten.Image) {
	face := text.NewGoXFace(basicfont.Face7x13)
	textOptions := &text.DrawOptions{}
	textOptions.ColorScale.ScaleWithColor(colorPrimary)

	rightX := float64(chip8DisplayWidth*chip8PixelSize + marginX*2)
	rightY := float64(marginY)
	textOptions.GeoM.Translate(rightX, rightY)

	g.drawRegisters(screen, face, textOptions)

	textOptions.GeoM.SetElement(0, 2, float64(rightX)+rightSpacing)
	textOptions.GeoM.SetElement(1, 2, float64(rightY))
	g.drawStack(screen, face, textOptions)

	bottomX := marginX
	bottomY := chip8DisplayHeight*chip8PixelSize + marginY*2
	textOptions.GeoM.SetElement(0, 2, float64(bottomX))
	textOptions.GeoM.SetElement(1, 2, float64(bottomY))

	g.drawStats(screen, face, textOptions)
	g.drawMemoryView(screen, face, textOptions)
}

func (g *Game) drawRegisters(screen *ebiten.Image, face *text.GoXFace, textOptions *text.DrawOptions) {
	text.Draw(screen, "Registers", face, textOptions)
	textOptions.GeoM.Translate(0, lineHeight*2)

	for i := range byte(0x10) {
		text.Draw(screen, fmt.Sprintf("V%X: 0x%02X", i, g.emulator.Registers[i]), face, textOptions)
		textOptions.GeoM.Translate(0, lineHeight)
	}
}

func (g *Game) drawStack(screen *ebiten.Image, face *text.GoXFace, textOptions *text.DrawOptions) {
	text.Draw(screen, "Stack", face, textOptions)
	textOptions.GeoM.Translate(0, lineHeight*2)

	if g.emulator.SP == 0 {
		// Stack empty
		text.Draw(screen, "Empty", face, textOptions)
	} else {
		for i := range g.emulator.SP {
			text.Draw(screen, fmt.Sprintf("0x%04X", g.emulator.Stack[i]), face, textOptions)
			textOptions.GeoM.Translate(0, lineHeight)
		}
	}
}

func (g *Game) drawStats(screen *ebiten.Image, face *text.GoXFace, textOptions *text.DrawOptions) {
	currentOpcode := g.emulator.GetCurrentOpcode(true)
	startX := textOptions.GeoM.Element(0, 2)

	// textOptions.ColorScale.ScaleWithColor(color.RGBA{52, 152, 219, 255})
	text.Draw(screen, fmt.Sprintf("Opcode: %s", currentOpcode), face, textOptions)
	// textOptions.ColorScale.Reset()
	textOptions.GeoM.Translate(0, lineHeight)
	text.Draw(screen, fmt.Sprintf("PC: 0x%04X", g.emulator.PC), face, textOptions)
	textOptions.GeoM.Translate(topRowSpacing, 0)
	text.Draw(screen, fmt.Sprintf("I: 0x%04X", g.emulator.I), face, textOptions)
	textOptions.GeoM.Translate(topRowSpacing, 0)
	text.Draw(screen, fmt.Sprintf("Timer: 0x%02X", g.emulator.DelayTimer), face, textOptions)
	textOptions.GeoM.Translate(topRowSpacing, 0)
	text.Draw(screen, fmt.Sprintf("SoundTimer: 0x%02X", g.emulator.SoundTimer), face, textOptions)
	textOptions.GeoM.Translate(0, lineHeight)
	textOptions.GeoM.SetElement(0, 2, float64(startX)) // reset x to starting value
}

func (g *Game) drawMemoryView(screen *ebiten.Image, face *text.GoXFace, textOptions *text.DrawOptions) {

	byteWidth := chip8DisplayWidth*chip8PixelSize/memWidth - 4
	currentAddress := g.emulator.PC
	startX := int(textOptions.GeoM.Element(0, 2))

	textOptions.GeoM.Translate(0, lineHeight) // extra line gap
	text.Draw(screen, "Memory section", face, textOptions)

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
		textOptions.GeoM.Translate(0, lineHeight)

		// Draw row address
		rowText := fmt.Sprintf("0x%04X:", addr)
		textOptions.GeoM.SetElement(0, 2, float64(startX)) // set x to memoryViewX
		text.Draw(screen, rowText, face, textOptions)

		// Draw bytes in this row
		for offset := 0; offset < memWidth; offset++ {
			byteAddr := addr + uint16(offset)
			byteValue := g.emulator.Memory[byteAddr]

			byteX := startX + memoryHeaderSpacing + (offset * byteWidth)
			byteY := int(textOptions.GeoM.Element(1, 2))

			// Highlight the current opcode bytes (2 bytes)
			if byteAddr == currentAddress || byteAddr == currentAddress+1 {
				vector.DrawFilledRect(
					screen,
					float32(byteX-(byteWidth/5)),
					float32(byteY-(lineHeight/5)),
					float32(byteWidth),
					float32(lineHeight),
					colorAccent,
					false,
				)
			}

			textOptions.GeoM.SetElement(0, 2, float64(byteX)) // set x to byteX
			text.Draw(screen, fmt.Sprintf("%02X", byteValue), face, textOptions)
		}
	}
}
