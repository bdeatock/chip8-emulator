package chip8

import "fmt"

// Returns current opcode at program counter as string, with optional description
func (e *Emulator) GetCurrentOpcode(addDescription bool) string {
	opcode := uint16(e.Memory[e.PC])<<8 | uint16(e.Memory[e.PC+1])

	if !addDescription {
		return fmt.Sprintf("0x%04X", opcode)
	}

	x := byte((opcode & 0x0F00) >> 8)
	y := byte((opcode & 0x00F0) >> 4)
	n := byte(opcode & 0x000F)
	nn := byte(opcode & 0x00FF)
	nnn := opcode & 0x0FFF

	description := ""

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			// 00E0: Clear screen
			description = "Clear screen"
		case 0x00EE:
			// Return from subroutine
			description = "Return from subroutine"
		}
	case 0x1000:
		// 1NNN: Jump
		description = fmt.Sprintf("Jump to address 0x%03X", nnn)
	case 0x2000:
		// 2NNN: call subroutine at NNN
		description = fmt.Sprintf("Call subroutine at 0x%03X", nnn)
	case 0x3000:
		// 3XNN: Skip next instruction if VX equals NN
		description = fmt.Sprintf("Skip next instruction if V%X (0x%02X) == 0x%02X", x, e.Registers[x], nn)
	case 0x4000:
		// 4XNN: Skip next instruction if VX not equal to NN
		description = fmt.Sprintf("Skip next instruction if V%X (0x%02X) != 0x%02X", x, e.Registers[x], nn)
	case 0x5000:
		// 5XY0: Skip next instruction if VX equal to VY
		if n == 0 {
			description = fmt.Sprintf("Skip next instruction if V%X (0x%02X) == V%X (0x%02X)", x, e.Registers[x], y, e.Registers[y])
		}
	case 0x6000:
		// 6XNN: Set
		description = fmt.Sprintf("Set V%X = 0x%02X", x, nn)
	case 0x7000:
		// 7XNN: Add
		description = fmt.Sprintf("Add V%X += 0x%02X (becomes 0x%02X)", x, nn, e.Registers[x]+nn)
	case 0x8000:
		switch n {
		case 0x0:
			// 8XY0: Set VX to value of VY
			description = fmt.Sprintf("Set V%X = V%X (0x%02X)", x, y, e.Registers[y])
		case 0x1:
			// 8XY1: Set VX to bitwise VX OR VY
			description = fmt.Sprintf("Set V%X = V%X | V%X (0x%02X | 0x%02X)", x, x, y, e.Registers[x], e.Registers[y])
		case 0x2:
			// 8XY2: Set VX to bitwise VX AND VY
			description = fmt.Sprintf("Set V%X = V%X & V%X (0x%02X & 0x%02X)", x, x, y, e.Registers[x], e.Registers[y])
		case 0x3:
			// 8XY3: Set VX to bitwise VX XOR VY
			description = fmt.Sprintf("Set V%X = V%X ^ V%X (0x%02X ^ 0x%02X)", x, x, y, e.Registers[x], e.Registers[y])
		case 0x4:
			// 8XY4: Add VY to VX with carry
			description = fmt.Sprintf("Add V%X += V%X (0x%02X + 0x%02X) with carry", x, y, e.Registers[x], e.Registers[y])
		case 0x5:
			// 8XY5: Subtract VY from VX with borrow
			description = fmt.Sprintf("Subtract V%X -= V%X (0x%02X - 0x%02X) with borrow", x, y, e.Registers[x], e.Registers[y])
		case 0x6:
			// 8XY6: Shift right
			if e.Config.LegacyShift {
				description = fmt.Sprintf("Set V%X = V%X (0x%02X) >> 1 with VF = LSB", x, y, e.Registers[y])
			} else {
				description = fmt.Sprintf("Shift V%X (0x%02X) >> 1 with VF = LSB", x, e.Registers[x])
			}
		case 0x7:
			// 8XY7: Set VX to VY - VX with borrow
			description = fmt.Sprintf("Set V%X = V%X - V%X (0x%02X - 0x%02X) with borrow", x, y, x, e.Registers[y], e.Registers[x])
		case 0xE:
			// 8XYE: Shift left
			if e.Config.LegacyShift {
				description = fmt.Sprintf("Set V%X = V%X (0x%02X) << 1 with VF = MSB", x, y, e.Registers[y])
			} else {
				description = fmt.Sprintf("Shift V%X (0x%02X) << 1 with VF = MSB", x, e.Registers[x])
			}
		}
	case 0x9000:
		// 9XY0: Skip next instruction if VX not equal to VY
		if n == 0 {
			description = fmt.Sprintf("Skip next instruction if V%X (0x%02X) != V%X (0x%02X)", x, e.Registers[x], y, e.Registers[y])
		}
	case 0xA000:
		// ANNN: Set index
		description = fmt.Sprintf("Set I = 0x%03X", nnn)
	case 0xB000:
		// BNNN: Jump with offset
		if e.Config.LegacyJump {
			description = fmt.Sprintf("Jump to address 0x%03X + V0 (0x%02X)", nnn, e.Registers[0])
		} else {
			description = fmt.Sprintf("Jump to address 0x%03X + V%X (0x%02X)", nnn, x, e.Registers[x])
		}
	case 0xC000:
		// CXNN: Random
		description = fmt.Sprintf("Set V%X = random & 0x%02X", x, nn)
	case 0xD000:
		// DXYN: Display
		description = fmt.Sprintf("Draw sprite at (V%X,V%X) = (%d,%d) with height %d", x, y, e.Registers[x], e.Registers[y], n)
	case 0xE000:
		switch nn {
		case 0x9E:
			// EX9E Skip if key X is pressed
			description = fmt.Sprintf("Skip next instruction if key V%X (0x%02X) is pressed", x, e.Registers[x])
		case 0xA1:
			// EXA1 Skip if key X is not pressed
			description = fmt.Sprintf("Skip next instruction if key V%X (0x%02X) is not pressed", x, e.Registers[x])
		}
	case 0xF000:
		switch nn {
		case 0x07:
			// FX07 Set VX to current value of delay timer
			description = fmt.Sprintf("Set V%X = delay timer (0x%02X)", x, e.DelayTimer)
		case 0x0A:
			// 0xFX0A Wait for a key press and set VX to key
			description = fmt.Sprintf("Wait for key press and store in V%X", x)
		case 0x15:
			// 0xFX15 Set delay timer to VX
			description = fmt.Sprintf("Set delay timer = V%X (0x%02X)", x, e.Registers[x])
		case 0x18:
			// 0xFX18 Set sound timer to VX
			description = fmt.Sprintf("Set sound timer = V%X (0x%02X)", x, e.Registers[x])
		case 0x1E:
			// 0xFX1E Add VX to I
			description = fmt.Sprintf("Add I += V%X (0x%04X + 0x%02X)", x, e.I, e.Registers[x])
		case 0x29:
			// 0xFX29 Set I to address of font for hex char in VX
			description = fmt.Sprintf("Set I to font address for hex digit V%X (0x%02X)", x, e.Registers[x])
		case 0x33:
			// 0xFX33 Take number in VX, convert to three decimal digits, and store at address in I, I+1, I+2
			description = fmt.Sprintf("Store BCD of V%X (0x%02X) at I, I+1, I+2", x, e.Registers[x])
		case 0x55:
			// 0xFX55 Store V0-VX at address I
			description = fmt.Sprintf("Store registers V0-V%X at address I (0x%04X)", x, e.I)
		case 0x65:
			// 0xFX65 Load memory from address I into V0-VX
			description = fmt.Sprintf("Load registers V0-V%X from address I (0x%04X)", x, e.I)
		}
	}
	return fmt.Sprintf("0x%04X - %s", opcode, description)
}
